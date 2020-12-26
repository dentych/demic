package cardDeck

// Card holds the card suits and types in the deck
type Card struct {
	Type string
	Suit string
}

// Deck holds the cards in the deck to be shuffled
type Deck []Card
