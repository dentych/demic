package pyramid

import (
	"gitlab.com/dentych/demic/pakage/cardDeck"
	"gitlab.com/dentych/demic/pakage/room"
	"log"
)

//deals 4 cardDeck.Card to room.Player
func NewHand(roomID string, playerName string) {
	element, exist := room.Rooms[roomID]
	if exist {
		element.Players[playerName].Hand, element.Deck = cardDeck.Deal(element.Deck, 4)
		log.Println("Room", roomID, ": New hand added to", playerName, element.Players[playerName].Hand)
	}
	room.Rooms[roomID] = element
}
