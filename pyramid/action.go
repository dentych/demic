package pyramid

const (
	ActionCreateGame   = "create-game"
	ActionPlayerJoin   = "player-join"
	ActionStartGame    = "start-game"
	ActionDealHand     = "player-deal-hand"
	ActionPlayerQuit   = "player-quit"
	ActionAttack       = "player-attack"
	ActionAcceptAttack = "player-accept-attack"
	ActionRejectAttack = "player-reject-attack"
	ActionPickCard     = "player-pick-card"
	ActionAttackState  = "attack-state"
	ActionContinue     = "continue"
	ActionHost         = "host"
	ActionNewCard      = "new-card"
	ActionNewRound     = "new-round"
	ActionShowCards    = "show-cards"
)

// Action struct that is used to send/receive messages to/from clients.
//
// Types of actions: attack, accept attack, reject attack, next card gets turned, attack mode start, attack mode end.
type Action struct {
	ActionType string `json:"action_type"`
	Origin     string `json:"origin"`
	Target     string `json:"target"`
}
