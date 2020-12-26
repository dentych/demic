package pyramid

import (
	"gitlab.com/dentych/demic/pakage/cardDeck"
	"gitlab.com/dentych/demic/pakage/room"
	"log"
)

//starts game by preparing room.Player and room.Board
func StartGame(roomID string) {
	element, exist := room.Rooms[roomID]
	if exist {
		element.Deck = cardDeck.NewDeck()
		element.Deck = cardDeck.Shuffle(element.Deck)
		room.Rooms[roomID] = element
		NewBoard(roomID)
		for playerName, _ := range element.Players {
			NewHand(roomID, playerName)
		}
		log.Println("Room", roomID, ": Game started")
	}
}
