package main

import (
	"encoding/json"
	"gitlab.com/dentych/demic/pyramid"
	"gitlab.com/dentych/demic/ws"
	"log"
	"net/http"
	"os"
	"time"
)

var exit chan string

func main() {
	exit = make(chan string, 1)
	go func() {
		<-exit
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()
	log.Println("Starting demic")
	log.Println("Setting up routes")
	SetupRoutes()
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SetupRoutes() {
	http.HandleFunc("/ws", ws.WebsocketEndpoint)

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		output, err := json.Marshal(pyramid.Rooms)
		if err != nil {
			log.Println("Failed to marshal pyramid rooms", err)
			return
		}
		_, err = w.Write(output)
		if err != nil {
			log.Println("Failed to write output", err)
		}
	})

	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		output, _ := json.Marshal("shutting down")
		_, _ = w.Write(output)
		exit <- ""
	})
}
