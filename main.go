package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var rooms []room

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

	log.Println("Client Connected")
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
			if len(rooms) > 0 {
				rooms = append(rooms, room{len(rooms) + 1})
			} else {
				rooms = append(rooms, room{1})
			}
			fmt.Println(rooms)
		}
		if message.Action == "join-room" {

			if len(rooms) > 0 {
				rooms = append(rooms, room{len(rooms) + 1})
			} else {
				rooms = append(rooms, room{1})
			}
			fmt.Println(rooms)
		}
	}
}

type message struct {
	Action string `json:"action"`
}

type room struct {
	id     int
	player []player
}

type player struct {
	name string
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8081", nil))
}
