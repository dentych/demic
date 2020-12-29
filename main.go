package main

import (
	"fmt"
	"gitlab.com/dentych/demic/ws"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Application Started")
	ws.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8081", nil))
}
