package pyramid

import (
	"gitlab/dentych/demic/pakage/cardDeck"
	"gitlab/dentych/demic/pakage/room"
)

func NewHand(roomID string, playerName string) {
	element, exist := room.Rooms[roomID]
	if exist {
		element.Players[playerName].Hand, element.Deck = cardDeck.Deal(element.Deck, 4)
	}
	room.Rooms[roomID] = element
}
