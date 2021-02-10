package ws

import (
	"errors"
	"fmt"
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
		roomId := ""
		if c.game != nil {
			roomId = c.game.RoomId
		}
		err := c.conn.WriteJSON(Message{
			RoomId: roomId,
			Action: action,
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
	switch msg.Action.ActionType {
	case pyramid.ActionCreateGame:
		err = c.createGame(msg)
	case pyramid.ActionPlayerJoin:
		err = c.joinGame(msg)
	default:
		if c.game != nil {
			c.game.Input <- msg.Action
		}
	}
	return err
}

func (c *Client) createGame(msg *Message) error {
	roomId := pyramid.GenerateId(4)
	game := pyramid.NewPyramidGame()
	_, ok := pyramid.PyramidRooms[roomId]
	if ok {
		return fmt.Errorf("tried to create pyramid game with existing ID: %s", msg.RoomId)
	} else {
		pyramid.PyramidRooms[roomId] = game
		c.game = game
		c.output <- pyramid.Action{ActionType: pyramid.ActionCreateGame, Target: roomId}
		c.player = pyramid.NewPlayer("HOST")
		c.player.Output = c.output
		c.game.PlayerJoin <- c.player
	}
	return nil
}

func (c *Client) joinGame(msg *Message) error {
	if c.game == nil {
		var ok bool
		c.game, ok = pyramid.PyramidRooms[msg.RoomId]
		if !ok {
			return fmt.Errorf("player tried to join game with ID %s, which was not found", msg.RoomId)
		}
	}
	if len(c.game.Players) > 6 {
		return fmt.Errorf("player tried to join game with ID %s, room is full", msg.RoomId)
	}
	c.player = &pyramid.Player{Name: msg.Action.Origin, Output: c.output}
	c.game.PlayerJoin <- c.player
	return nil
}
