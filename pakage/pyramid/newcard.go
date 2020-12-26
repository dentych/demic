package pyramid

import (
	"gitlab/dentych/demic/pakage/cardDeck"
	"gitlab/dentych/demic/pakage/room"
)

func NewCard(roomID string, playerName string, handIDX int) {
	element, exist := room.Rooms[roomID]
	if exist {
		var newCardList []cardDeck.Card
		oldCard := element.Players[playerName].Hand[handIDX]
		newCardList, element.Deck = cardDeck.Deal(element.Deck, 1)
		element.Players[playerName].Hand[handIDX] = newCardList[0]
		room.AddCardToDeck(roomID, oldCard)
	}
	room.Rooms[roomID] = element
}
