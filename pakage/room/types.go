package room

import "gitlab/dentych/demic/pakage/cardDeck"

type Room struct {
	Players map[string]*Player
	Board   []cardDeck.Card
	Deck    []cardDeck.Card
}

type Player struct {
	PlayerName string
	Hand       []cardDeck.Card
}

var (
	Rooms = make(map[string]Room)
)

//
//type Card struct {
//	Type string
//	Suit string
//}
//
//// Deck holds the cards in the deck to be shuffled
//type Deck []Card
