package pyramid

import (
	"gitlab.com/dentych/demic/pakage/cardDeck"
	"gitlab.com/dentych/demic/pakage/room"
	"log"
)

//adds new cardDeck.Card to room.Player.Hand[handIDX], shuffles old cardDeck.Card into cardDeck.Deck
func NewCard(roomID string, playerName string, handIDX int) {
	element, exist := room.Rooms[roomID]
	if exist {
		var newCardList []cardDeck.Card
		oldCard := element.Players[playerName].Hand[handIDX]
		newCardList, element.Deck = cardDeck.Deal(element.Deck, 1)
		element.Players[playerName].Hand[handIDX] = newCardList[0]
		room.AddCardToDeck(roomID, oldCard)
		log.Println("Room", roomID, ":", playerName, "Swapped", oldCard, "for", newCardList[0], "at hand index", handIDX)
		log.Println(len(element.Deck))
	}
	room.Rooms[roomID] = element
}
