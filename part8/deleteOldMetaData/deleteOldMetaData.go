package main

import (
	"fmt"
	"log"
	"storage/lib/es"
)

const MinVersionCount = 5

func main() {
	fmt.Println("run")
	buckets, err := es.SearchVersionStatus(MinVersionCount + 1)
	fmt.Println(buckets)
	if err != nil {
		log.Println(err)
		return
	}
	for i := range buckets {
		bucket := buckets[i]
		for v := 0; v < bucket.DocCount-MinVersionCount; v++ {
			fmt.Println(bucket.Key, v+int(bucket.MinVersion.Value))
			es.DelMetadata(bucket.Key, v+int(bucket.MinVersion.Value))
		}
	}
}
