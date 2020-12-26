package room

import "gitlab/dentych/demic/pakage/cardDeck"

//Room holds the players for the game
type Room struct {
	Players map[string]*Player
	Board   []cardDeck.Card
	Deck    []cardDeck.Card
}

//Player holds the players hand
type Player struct {
	Hand []cardDeck.Card
}

//Rooms holds map for all rooms, roomID's is used as key values
var (
	Rooms = make(map[string]Room)
)
