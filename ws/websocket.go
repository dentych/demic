package ws

import (
	"encoding/json"
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
	var clientType string

	for {
		if len(clientType) == 0 {
			t, m, err := ws.ReadMessage()
			if err != nil {
				log.Println("Failed to read message from websocket", err)
				return
			}
			if t == websocket.TextMessage {
				var message Message
				err = json.Unmarshal(m, &message)
				if err != nil {
					log.Println("Failed to unmarshal message from client", err)
					return
				}
				if message.Action.ActionType == pyramid.ActionCreateGame {
					roomId := pyramid.GenerateId(4)
					rooms[roomId] = pyramid.NewPyramidGame()
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
				} else {
					// Assume join received
					roomId := message.RoomId
					v, ok := rooms[roomId]
					if ok {
						err := v.AddPlayer(pyramid.NewPlayer(message.Action.Origin))
						if err != nil {
							log.Printf("Failed to add player '%s' to game: %s", message.Action.Origin, err)
							return
						}
					}
				}
			}
		} else {

		}
	}
}
