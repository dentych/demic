package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Message struct {
	Action     string `json:"action"`
	PlayerName string `json:"name"`
	RoomID     string `json:"room"`
}

type PyramidSocketHandler struct {
	RoomId string
	PlayerName string
	IsHost bool

	webSocket *websocket.Conn
}

func (h *PyramidSocketHandler) HandleMessages() {
	msgType, msg, err := h.readMessage()
	var message Message
	err := h.webSocket.ReadJSON(&message)
	if err != nil {
		log.Println("Failed to read JSON message", err)
		return
	}

	err = h.handleInitialMessage()
	if err != nil {
		log.Println("unable to handle initial message", err)
		return
	}

	//for {
	//	// read in a message
	//	messageType, p, err := conn.ReadMessage()
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//	// print out that message for clarity
	//	fmt.Println(string(p))
	//
	//	if err := conn.WriteMessage(messageType, p); err != nil {
	//		log.Println(err)
	//		return
	//	}
	//	var message Message
	//
	//	err = json.Unmarshal(p, &message)
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//	if message.Action == "new-room" {
	//		room.CreateRoom()
	//	}
	//	if message.Action == "join-room" {
	//		room.AddPlayerToRoom(message.RoomID, message.PlayerName)
	//	}
	//	if message.Action == "start-game" {
	//		pyramid.StartGame(message.RoomID)
	//	}
	//	if message.Action == "new-card.go" {
	//		handIDX, _ := strconv.Atoi(message.HandIDX)
	//		pyramid.NewCard(message.RoomID, message.PlayerName, handIDX)
	//	}
	//	if message.Action == "end-game" {
	//		pyramid.EndGame(message.RoomID)
	//	}
	//	if message.Action == "print-room" {
	//		fmt.Println(room.Rooms)
	//	}
	//	if message.Action == "delete-room" {
	//		room.DeleteRoom(message.RoomID)
	//	}
	//	if message.Action == "remove-player" {
	//		room.RemovePlayer(message.RoomID, message.PlayerName)
	//	}
	//}
}

func (h *PyramidSocketHandler) handleInitialMessage() error {
	var message Message
	err := h.webSocket.ReadMessage()
	if err != nil {
		return err
	}

	switch message.Action {
	case "create-game":
		fmt.Println("Creating game")
	default:
		fmt.Println("Invalid message")
	}

	return nil
}

func (h *PyramidSocketHandler) readMessage() (interface{}, interface{}, interface{}) {
	msgType, msg, err := h.webSocket.ReadMessage()
}
