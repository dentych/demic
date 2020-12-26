package room

import "gitlab.com/dentych/demic/pakage/cardDeck"

//adds 1 cardDeck.Card to cardDeck.Deck and Shuffles the deck
func AddCardToDeck(roomID string, card cardDeck.Card) {
	element, exist := Rooms[roomID]
	if exist {
		element.Deck = append(element.Deck, card)
		element.Deck = cardDeck.Shuffle(element.Deck)
	}
	Rooms[roomID] = element
}
