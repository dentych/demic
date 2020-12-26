package ws

import (
	"log"
	"net/http"
)

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	Upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	//log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hello frontend"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	Reader(ws)
}
