package main

import (
	"log"
	"net/http"
	"os"
	"storage/part8/apiServer/heartbeat"
	"storage/part8/apiServer/locate"
	"storage/part8/apiServer/objects"
	"storage/part8/apiServer/temp"
	"storage/part8/apiServer/version"
)

func main() {
	go heartbeat.ListenHeartBeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", version.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
