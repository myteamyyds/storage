package objects

import (
	"io"
	"log"
	"net/http"
	"strings"
)

//func get(w http.ResponseWriter, r *http.Request) {
//	file, err := os.Open(os.Getenv("STORAGE_ROOT") + "/objects/" +
//		strings.Split(r.URL.EscapedPath(), "/")[2])
//	if err != nil {
//		fmt.Println(err)
//		w.WriteHeader(http.StatusFound)
//		return
//	}
//	defer file.Close()
//	io.Copy(w, file)
//}

func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := getStream(object)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}
