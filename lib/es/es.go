package es

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
)

type Metadata struct {
	Dname   string `json:"dname"`
	Version int    `json:"version"`
	Size    int64  `json:"size"`
	Hash    string `json:"hash"`
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total struct {
			Value    int
			Relation string
		}
		Hits []hit
	}
}

func getMetadata(dname string, versionId int) (meta Metadata, err error) {
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d/_source", os.Getenv("ES_SERVER"), dname, versionId)
	fmt.Println(url)
	result, err := http.Get(url)
	if err != nil {
		return
	}
	if result.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to get %s_%d:%d", dname, versionId, result.StatusCode)
		return
	}
	result2, _ := io.ReadAll(result.Body)
	json.Unmarshal(result2, &meta)
	return
}

func SearchLatestVersion(dname string) (meta Metadata, err error) {
	fmt.Println("进入SearchLatestVersion")
	url := fmt.Sprintf("http://%s/metadata/_search?q=dname:%s&size=1&sort=version:desc", os.Getenv("ES_SERVER"), url2.PathEscape(dname))
	fmt.Println(url)
	result, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	if result.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to search latest metadata:%s", result.StatusCode)
		return
	}
	result2, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println(err)
	}
	var sr searchResult
	json.Unmarshal(result2, &sr)
	fmt.Println("sr = ", sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

func GetMetadata(name string, version int) (Metadata, error) {
	fmt.Println("进入GetMetadata")
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

func PutMetadata(dname string, version int, size int64, hash string) error {
	document := fmt.Sprintf(`{"dname":"%s","version":%d,"size":%d,"hash":"%s"}`, dname, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d?op_type=create", os.Getenv("ES_SERVER"), dname, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(document))
	request.Header.Set("Content-Type", "application/json")
	result, err := client.Do(request)
	if err != nil {
		return err
	}
	if result.StatusCode == http.StatusConflict {
		return PutMetadata(dname, version+1, size, hash)
	}
	if result.StatusCode != http.StatusCreated {
		result2, _ := io.ReadAll(result.Body)
		return fmt.Errorf("fail to put metadata:%d %s", result.StatusCode, string(result2))
	}
	return nil
}

func AddVersion(dname, hash string, size int64) error {
	version, err := SearchLatestVersion(dname)
	fmt.Println("上个版本的hash = ", version.Hash)
	fmt.Println("上个版本的version = ", version.Version)
	if err != nil {
		return err
	}
	return PutMetadata(dname, version.Version+1, size, hash)
}

func SearchAllVersions(dname string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?sort=version&from=%d&size=%d", os.Getenv("ES_SERVER"), from, size)
	if dname != "" {
		url += "&q=dname:" + dname
	}
	result, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	metas := make([]Metadata, 0)
	result2, _ := io.ReadAll(result.Body)
	var sr searchResult
	json.Unmarshal(result2, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}

func DelMetadata(dname string, version int) {
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d", os.Getenv("ES_SERVER"), dname, version)
	client := http.Client{}
	request, _ := http.NewRequest("DELETE", url, nil)
	client.Do(request)
}

type Bucket struct {
	Key        string `json:"key"`
	DocCount   int    `json:"doc_count"`
	MinVersion struct {
		Value float32 `json:"value"`
	} `json:"min_version"`
}

type aggregateResult struct {
	Aggregations struct {
		Group_by_name struct {
			Buckets []Bucket `json:"buckets"`
		} `json:"group_by_name"`
	} `json:"aggregations"`
}

func SearchVersionStatus(minDocCount int) ([]Bucket, error) {
	url := fmt.Sprintf("http://%s/metadata/_search", os.Getenv("ES_SERVER"))
	body := fmt.Sprintf(`
		{
			"size": 0,
			"aggs": {
				"group_by_name": {
					"terms": {
						"field": "name",
						"min_doc_count": %d
					},
					"aggs": {
						"min_version": {
							"min": {
								"field": "version"
							}
						}
					}
				}
			}
		}`, minDocCount)
	client := http.Client{}
	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	result, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	responseBody, _ := io.ReadAll(result.Body)
	var ar aggregateResult
	err = json.Unmarshal(responseBody, &ar)
	return ar.Aggregations.Group_by_name.Buckets, nil
}

func HasHash(hash string) (bool, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=0", os.Getenv("ES_SERVER"), hash)
	result, err := http.Get(url)
	if err != nil {
		return false, err
	}
	body, _ := io.ReadAll(result.Body)
	var sr searchResult
	json.Unmarshal(body, &sr)
	return sr.Hits.Total.Value != 0, nil
}

func SearchHashSize(hash string) (size int64, err error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=1", os.Getenv("ES_SERVER"), hash)
	result, err := http.Get(url)
	if err != nil {
		return
	}
	if result.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to search hash size:%d", result.StatusCode)
		return
	}
	body, _ := io.ReadAll(result.Body)
	var sr searchResult
	json.Unmarshal(body, &sr)
	if len(sr.Hits.Hits) != 0 {
		size = sr.Hits.Hits[0].Source.Size
	}
	return
}
