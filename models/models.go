package models

import "encoding/json"

const (
	ActionHello       = "hello"
	ActionCreateGame  = "create-game"
	ActionGameCreated = "game-created"
	ActionJoinGame    = "join-game"
	ActionGameJoined  = "game-joined"
)

type IncomingMessage struct {
	ActionType string          `json:"action_type"`
	Payload    json.RawMessage `json:"payload"`
}

type OutgoingMessage struct {
	ActionType string      `json:"action_type"`
	Payload    interface{} `json:"payload"`
}

type PayloadCreateGame struct {
	Game string `json:"game"`
}

type PayloadGameCreated struct {
	Game   string `json:"game"`
	RoomID string `json:"room_id"`
}

type PayloadJoinGame struct {
	Game       string `json:"game"`
	RoomID     string `json:"room_id"`
	PlayerName string `json:"player_name"`
}

type PayloadHello struct {
	ClientID string `json:"client_id"`
}
