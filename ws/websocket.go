package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gitlab.com/dentych/demic/pyramid"
	"log"
	"net/http"
)

type Message struct {
	Action  string          `json:"action_type"`
	Payload json.RawMessage `json:"payload"`
}

var upgrader websocket.Upgrader

type Client struct {
	conn   *websocket.Conn
	game   *pyramid.Pyramid
	player *pyramid.Player
	output chan pyramid.Action
}

func WebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to websocket: ", err)
		return
	}

	client := Client{conn: ws, output: make(chan pyramid.Action, 10)}
	go client.handleIncoming()
	go client.handleOutgoing()
}

func (c *Client) handleIncoming() {
	var message Message
	for {
		err := c.conn.ReadJSON(&message)
		if err != nil {
			log.Println("Failed to read incoming message from websocket:", err)
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				if c.game != nil {
					c.game.PlayerLeave <- c.player
				}
			}
			return
		}
		err = c.handleIncomingMessage(&message)
		if err != nil {
			log.Println("Error handling incoming message:", err)
		}
	}
}

func (c *Client) handleOutgoing() {
	for {

		action := <-c.output

		myIn, err := json.Marshal(
			&pyramid.Action{
				ActionType: action.ActionType,
				Origin:     action.Origin,
				Target:     action.Target,
			},
		)
		myInRaw := json.RawMessage(myIn)

		err = c.conn.WriteJSON(Message{
			Action:  action.ActionType,
			Payload: myInRaw,
		})

		if err != nil {
			log.Println("Error writing to websocket:", err)
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				if c.game != nil {
					c.game.PlayerLeave <- c.player
				}
			}
			return
		}
	}
}

func (c *Client) handleIncomingMessage(msg *Message) error {
	var err error

	var payload pyramid.Action

	err = json.Unmarshal(msg.Payload, &payload)
	switch payload.ActionType {
	case pyramid.ActionCreateGame:
		err = c.createGame()
	case pyramid.ActionPlayerJoin:
		err = c.joinGame(&payload)
	default:
		if c.game != nil {
			c.game.Input <- payload
		}
	}
	return err
}

func (c *Client) createGame() error {
	roomId := pyramid.GenerateId(4)
	game := pyramid.NewPyramidGame()
	_, ok := pyramid.Rooms[roomId]
	if ok {
		return fmt.Errorf("tried to create pyramid game with existing ID: %s", roomId)
	} else {
		pyramid.Rooms[roomId] = game
		c.game = game
		c.output <- pyramid.Action{ActionType: pyramid.ActionCreateGame, Target: roomId}
		c.player = pyramid.NewPlayer("HOST")
		c.player.Output = c.output
		c.game.PlayerJoin <- c.player
	}
	return nil
}

func (c *Client) joinGame(action *pyramid.Action) error {
	if c.game == nil {
		var ok bool
		c.game, ok = pyramid.Rooms[action.Target]
		if !ok {
			return fmt.Errorf("player tried to join game with ID %s, which was not found", action.Target)
		}
	}
	if len(c.game.Players) > 6 {
		return fmt.Errorf("player tried to join game with ID %s, room is full", action.Target)
	}
	c.player = &pyramid.Player{Name: action.Origin, Output: c.output}
	c.game.PlayerJoin <- c.player
	return nil
}
