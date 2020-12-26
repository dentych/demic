package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.com/dentych/demic/pakage/pyramid"
	"gitlab.com/dentych/demic/pakage/room"
	"log"
	"strconv"
)

func Reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
		var message Message

		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println(err)
			return
		}
		if message.Action == "new-room" {
			room.CreateRoom()
		}
		if message.Action == "join-room" {
			room.AddPlayerToRoom(message.RoomID, message.PlayerName)
		}
		if message.Action == "start-game" {
			pyramid.StartGame(message.RoomID)
		}
		if message.Action == "new-card" {
			handIDX, _ := strconv.Atoi(message.HandIDX)
			pyramid.NewCard(message.RoomID, message.PlayerName, handIDX)
		}
		if message.Action == "end-game" {
			pyramid.EndGame(message.RoomID)
		}
		if message.Action == "print-room" {
			fmt.Println(room.Rooms)
		}
		if message.Action == "delete-room" {
			room.DeleteRoom(message.RoomID)
		}
		if message.Action == "remove-player" {
			room.RemovePlayer(message.RoomID, message.PlayerName)
		}
	}
}
