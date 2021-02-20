package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gitlab.com/dentych/demic/models"
	"gitlab.com/dentych/demic/pyramid"
	"log"
	"net/http"
)

var upgrader websocket.Upgrader

type Client struct {
	clientID string
	conn     *websocket.Conn
	input    chan<- models.IncomingMessage
	output   <-chan models.OutgoingMessage
}

func WebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection to websocket:", err)
		return
	}
	client := Client{conn: ws}

	err = client.hello()
	if err != nil {
		log.Println("Hello procedure failed:", err)
		_ = ws.Close()
		return
	}

	go client.handleIncoming()
	client.handleOutgoing()

	_ = ws.Close()
}

func (c *Client) handleIncoming() {
	var message models.IncomingMessage
	for {
		err := c.conn.ReadJSON(&message)
		if err != nil {
			log.Println("Failed to read incoming message from websocket:", err)
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				// Do something
			}
			return
		}
		c.input <- message
	}
}

func (c *Client) handleOutgoing() {
	for action := range c.output {
		err := c.conn.WriteJSON(models.OutgoingMessage{
			ActionType: action.ActionType,
			Payload:    action,
		})
		if err != nil {
			log.Println("Error writing to websocket:", err)
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				// Do something
			}
			return
		}
	}
}

func (c *Client) createGame(payload models.PayloadCreateGame) error {
	var roomID string
	var err error
	switch payload.Game {
	case "pyramid":
		roomID, c.input, c.output, err = pyramid.Create(c.clientID)
		err = c.conn.WriteJSON(models.OutgoingMessage{
			ActionType: models.ActionGameCreated,
			Payload:    models.PayloadGameCreated{
				Game:   "pyramid",
				RoomID: roomID,
			},
		})
	default:
		log.Println("client tried to create invalid game:", payload.Game)
		return fmt.Errorf("client tried to create invalid game: %s", payload.Game)
	}

	return err
}

func (c *Client) joinGame(payload models.PayloadJoin) error {
	// TODO
	return nil
}

// hello exchanges clientID with the client. If the client provides one, it will be used, otherwise one will be
// generated. The server will always reply back with a hello message.
//
// After the initial hello, the client is expected to send a create or join command, to either create or join a game.
//
// Creating a game is achieved by sending a message with a PayloadCreateGame payload.
//
// Joining a game is achieved by sending a message with a PayloadJoin payload.
func (c *Client) hello() error {
	var msg models.IncomingMessage

	// exchange client ID procedure
	err := c.conn.ReadJSON(&msg)
	if err != nil {
		return fmt.Errorf("error reading first message: %w", err)
	}

	if msg.ActionType != models.ActionHello {
		return fmt.Errorf("incorrect action type: %s", msg.ActionType)
	}

	var payload models.PayloadHello
	err = json.Unmarshal(msg.Payload, &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal hello payload: %w", err)
	}

	if payload.ClientID == "" {
		payload.ClientID = uuid.NewString()
	}
	c.clientID = payload.ClientID

	err = c.conn.WriteJSON(models.OutgoingMessage{ActionType: models.ActionHello, Payload: models.PayloadHello{ClientID: c.clientID}})
	if err != nil {
		return fmt.Errorf("failed to write hello back to client: %w", err)
	}

	// join or create game
	err = c.conn.ReadJSON(&msg)
	if err != nil {
		return fmt.Errorf("error reading join/create message: %w", err)
	}

	switch msg.ActionType {
	case models.ActionCreateGame:
		var payload models.PayloadCreateGame
		err = json.Unmarshal(msg.Payload, &payload)
		if err != nil {
			log.Println("create game payload unmarshalling error:", err)
			return fmt.Errorf("create game payload unmarshalling error: %w", err)
		}
		return c.createGame(payload)
	case models.ActionJoinGame:
		var payload models.PayloadJoin
		err = json.Unmarshal(msg.Payload, &payload)
		if err != nil {
			log.Println("create game payload unmarshalling error:", err)
			return fmt.Errorf("create game payload unmarshalling error: %w", err)
		}
		return c.joinGame(payload)
	default:
		log.Println("Invalid action type during handshake. PayloadAction type:", msg.ActionType)
		return fmt.Errorf("handshake error due to invalid action type: %s", msg.ActionType)
	}

	return nil
}
