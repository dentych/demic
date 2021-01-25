package pyramid

import (
	"gitlab.com/dentych/demic/card"
)

//Player holds the players hand
type Player struct {
	Name string
	Hand []card.Card
}

func NewPlayer(name string) *Player {
	p := Player{}
	p.Name = name
	return &Player{}
}
