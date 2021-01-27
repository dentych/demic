package ws

import (
	"github.com/gorilla/websocket"
	"gitlab.com/dentych/demic/pyramid"
	"log"
	"net/http"
)

type Message struct {
	RoomId string         `json:"room_id"`
	Action pyramid.Action `json:"action"`
}

var upgrader websocket.Upgrader

func WebsocketEndpoint(w http.ResponseWriter, r *http.Request, rooms map[string]*pyramid.Pyramid) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to websocket: ", err)
		return
	}

	handleWsMessages(ws, rooms)
}

func handleWsMessages(ws *websocket.Conn, rooms map[string]*pyramid.Pyramid) {
	output := make(chan pyramid.Action, 5)
	var roomId string
	var message Message

	defer ws.Close()
	err := ws.ReadJSON(&message)
	if err != nil {
		log.Println("Failed to read initial message from websocket", err)
		return
	}
	switch message.Action.ActionType {
	case pyramid.ActionCreateGame:
		roomId = pyramid.GenerateId(4)
		rooms[roomId] = pyramid.NewPyramidGame()
		rooms[roomId].AddPlayer(&pyramid.Player{
			Name:   "HOST",
			Output: output,
		})
		err = ws.WriteJSON(Message{
			RoomId: roomId,
			Action: pyramid.Action{
				ActionType: pyramid.ActionGameCreated,
				Origin:     "",
				Target:     roomId,
			},
		})
		if err != nil {
			log.Println("Failed to write message back to client", err)
			return
		}
	case pyramid.ActionPlayerJoin:
		roomId = message.RoomId
		v, ok := rooms[roomId]
		if ok {
			player := pyramid.NewPlayer(message.Action.Origin)
			player.Output = output
			err := v.AddPlayer(player)
			if err != nil {
				log.Printf("Failed to add player '%s' to game: %s", message.Action.Origin, err)
				return
			}
		}
	default:
		log.Println("Incorrect initial message. Closing websocket")
		return
	}
	go handleOutput(ws, roomId, output)
	for {
		err = ws.ReadJSON(&message)
		if err != nil {
			log.Println("Failed to read JSON message from websocket", err)
			return
		}

	}
}

func handleOutput(ws *websocket.Conn, roomId string, output <-chan pyramid.Action) {
	for {
		action := <-output
		err := ws.WriteJSON(Message{
			RoomId: roomId,
			Action: action,
		})
		if err != nil {
			return
		}
	}
}
