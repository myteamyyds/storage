package main

import (
	"log"
	"net/http"
	"os"
	"storage/part2/apiServer/objects"
	"storage/part2/dataServer/heartbeat"
	"storage/part2/dataServer/locate"
)

func main() {
	go heartbeat.StartHeartBeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
