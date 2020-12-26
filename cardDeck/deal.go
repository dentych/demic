package cardDeck

// Deal a specified amount of cards
func Deal(deck []Card, n int) ([]Card, Deck) {
	var cardList []Card
	for _, card := range deck[0:n] {
		cardList = append(cardList, card)
		deck = deck[1:]
	}
	return cardList, deck
}
