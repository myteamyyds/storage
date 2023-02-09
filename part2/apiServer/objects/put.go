package objects

import (
	"log"
	"net/http"
	"strings"
)

//func put(w http.ResponseWriter, r *http.Request) {
//	file, err := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" +
//		strings.Split(r.URL.EscapedPath(), "/")[2])
//	if err != nil {
//		fmt.Println(err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	defer file.Close()
//	io.Copy(file, r.Body)
//
//}

func put(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	c, err := storeObject(r.Body, object)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(c)
}
