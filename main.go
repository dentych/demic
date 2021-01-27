package main

import (
	"encoding/json"
	"gitlab.com/dentych/demic/pyramid"
	"gitlab.com/dentych/demic/ws"
	"log"
	"net/http"
)


func main() {
	log.Println("Starting demic")
	log.Println("Setting up routes")
	SetupRoutes()
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SetupRoutes() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.WebsocketEndpoint(w, r, pyramid.PyramidRooms)
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		output, err := json.Marshal(pyramid.PyramidRooms)
		if err != nil {
			log.Println("Failed to marshal pyramid rooms", err)
			return
		}
		_, err = w.Write(output)
		if err != nil {
			log.Println("Failed to write output", err)
		}
	})
}
