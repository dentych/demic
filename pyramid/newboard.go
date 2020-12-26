package pyramid

import (
	"gitlab.com/dentych/demic/cardDeck"
	"gitlab.com/dentych/demic/pakage/room"
	"log"
)

//deals 15 cardDeck.Card to room.Room.Board
func NewBoard(roomID string) {
	element, exist := room.Rooms[roomID]
	if exist {
		element.Board, element.Deck = cardDeck.Deal(element.Deck, 15)
		log.Println("Room", roomID, ": New board added", element.Board)
	}
	room.Rooms[roomID] = element
}
