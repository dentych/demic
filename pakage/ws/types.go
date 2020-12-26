package ws

import "github.com/gorilla/websocket"

type Message struct {
	Action     string `json:"action"`
	PlayerName string `json:"name"`
	RoomID     string `json:"room"`
	HandIDX    string `json:"card"`
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
