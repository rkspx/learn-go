package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rikisan1993/learn-go/yt-sub-monitor/websocket"
	"github.com/rikisan1993/learn-go/yt-sub-monitor/youtube"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func stats(w http.ResponseWriter, r *http.Request) {
	ws := websocket.New(youtube.GetStatistics)

	ws.Upgrade(w, r)
	go ws.Write()
}

func setupRouter() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/stat", stats)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	fmt.Println("YouTube Subscriber Monitor")
	setupRouter()
}
