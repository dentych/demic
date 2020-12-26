package pyramid

import (
	"gitlab/dentych/demic/pakage/cardDeck"
	"gitlab/dentych/demic/pakage/room"
)

func StartGame(roomID string) {
	element, exist := room.Rooms[roomID]
	if exist {
		element.Deck = cardDeck.NewDeck()
		room.Rooms[roomID] = element
		NewBoard(roomID)
		for playerName, _ := range element.Players {
			NewHand(roomID, playerName)
		}
		room.Rooms[roomID] = element
	}
}
