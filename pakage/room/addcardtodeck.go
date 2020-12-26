package room

import "gitlab/dentych/demic/pakage/cardDeck"

func AddCardToDeck(roomID string, card cardDeck.Card) {
	element, exist := Rooms[roomID]
	if exist {
		element.Deck = append(element.Deck, card)
		//element.Deck = cardDeck.Shuffle(element.Deck)
	}
	Rooms[roomID] = element
}
