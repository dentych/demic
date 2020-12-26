package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab/dentych/demic/pakage/pyramid"
	"gitlab/dentych/demic/pakage/room"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
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
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func reader(conn *websocket.Conn) {
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
		var message message

		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println(err)
			return
		}
		if message.Action == "new-room" {
			room.CreateRoom()
		}
		if message.Action == "join-room" {
			room.AddPlayerToRoom(message.PlayerName, message.RoomID)
			fmt.Println(room.Rooms)
		}
		if message.Action == "start-game" {
			pyramid.StartGame(message.RoomID)
			log.Println("Game started in room ID:", message.RoomID)
		}
		if message.Action == "new-card" {
			pyramid.NewCard(message.RoomID, message.PlayerName, 2)
		}
		if message.Action == "end-game" {
			pyramid.EndGame(message.RoomID)
			log.Println("Game ended in room ID:", message.RoomID)
		}
		if message.Action == "print-room" {
			printRoom(message.RoomID)
		}
	}
}
func printRoom(roomID string) {
	fmt.Println(room.Rooms)
}

//wb:
type message struct {
	Action     string `json:"action"`
	PlayerName string `json:"name"`
	RoomID     string `json:"room"`
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8081", nil))

	playerName := "Noer"
	fmt.Println("hello")
	roomID := room.CreateRoom()
	room.AddPlayerToRoom(playerName, roomID)
	pyramid.StartGame(roomID)
	fmt.Println(room.Rooms[roomID].Players[playerName].Hand)
	pyramid.NewCard(roomID, playerName, 0)
	fmt.Println(room.Rooms[roomID].Players[playerName].Hand)
}
