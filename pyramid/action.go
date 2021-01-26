package pyramid

const (
	ActionCreateGame  = "create-game"
	ActionGameCreated = "game-created"
)

// Action struct that is used to send/receive messages to/from clients.
//
// Types of actions: attack, reject attack, next card gets turned, attack mode start, attack mode end.
type Action struct {
	ActionType string `json:"action_type"`
	Origin     string `json:"origin"`
	Target     string `json:"target"`
}
