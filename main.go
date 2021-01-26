package main

import (
	"fmt"
	"gitlab.com/dentych/demic/ws"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Application started")
	SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SetupRoutes() {
	http.HandleFunc("/ws", ws.WebsocketEndpoint)
}
