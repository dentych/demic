package pyramid

import (
	"fmt"
	"gitlab.com/dentych/demic/card"
)

var (
	ErrGameStarted      = fmt.Errorf("game is started")
	ErrGameNotStarted   = fmt.Errorf("game is not started")
	ErrPlayerNameExists = fmt.Errorf("player name already exists")
	ErrNoMoreCards      = fmt.Errorf("no more cards to turn")
)

type Pyramid struct {
	RoomId         string
	Players        []Player
	BoardCardIndex int
	Board          []card.Card
	Deck           []card.Card
	Locked         bool // Locks as soon as cards are dealt, means the game is started
}

func NewPyramidGame() *Pyramid {
	p := Pyramid{}
	p.RoomId = GenerateId(4)
	p.Players = make([]Player, 0)
	p.Deck = card.NewDeck()
	p.BoardCardIndex = 0

	return &p
}

func (p *Pyramid) AddPlayer(player *Player) error {
	if p.Locked {
		return ErrGameStarted
	}

	for _, v := range p.Players {
		if player.Name == v.Name {
			return ErrPlayerNameExists
		}
	}

	p.Players = append(p.Players, *player)
	return nil
}

func (p *Pyramid) DealCards() {
	p.Locked = true
	p.Deck = card.Shuffle(p.Deck)
	p.Board, p.Deck = card.Deal(p.Deck, 10)
	for k := range p.Players {
		p.Players[k].Hand, p.Deck = card.Deal(p.Deck, 4)
	}
}

func (p *Pyramid) TurnNextCard() (*card.Card, error) {
	if !p.Locked {
		return nil, ErrGameNotStarted
	}

	if p.BoardCardIndex < len(p.Board) {
		c := p.Board[p.BoardCardIndex]
		p.BoardCardIndex++
		return &c, nil
	} else {
		return nil, ErrNoMoreCards
	}
}
