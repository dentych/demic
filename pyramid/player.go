package pyramid

import (
	"gitlab.com/dentych/demic/card"
)

//Player holds the Players hand
type Player struct {
	Name string
	Hand []card.Card
	Sips int
}

func NewPlayer(name string) *Player {
	p := Player{}
	p.Name = name
	return &p
}
