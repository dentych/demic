package ws

import "net/http"

func SetupRoutes() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/ws", WsEndpoint)
}
