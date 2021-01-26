package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader websocket.Upgrader

func WebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to websocket: ", err)
		return
	}

	handler := PyramidSocketHandler{webSocket: ws}
	handler.HandleMessages()
}
