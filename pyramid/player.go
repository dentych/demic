package pyramid

import (
	"gitlab.com/dentych/demic/card"
	"gitlab.com/dentych/demic/models"
)

//Player holds the Players hand
type Player struct {
	ClientID string
	Name     string
	Hand     []card.Card
	Sips     int
	Output   chan models.OutgoingMessage `json:"-"`
}

func NewPlayer(clientID, name string) *Player {
	p := Player{}
	p.ClientID = clientID
	p.Name = name
	p.Output = make(chan models.OutgoingMessage, 25)
	return &p
}
