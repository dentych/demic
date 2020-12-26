package pyramid

import (
	"gitlab/dentych/demic/pakage/cardDeck"
	"gitlab/dentych/demic/pakage/room"
)

func NewBoard(roomID string) {
	element, exist := room.Rooms[roomID]
	if exist {
		element.Board, element.Deck = cardDeck.Deal(element.Deck, 15)
	}
	room.Rooms[roomID] = element
}
