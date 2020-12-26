package cardDeck

// Deal a specified amount of cards
func Deal(d Deck, n int) ([]Card, Deck) {
	var list []Card
	for i := 0; i < n; i++ {
		list = append(list, d[i])
		d = d[1:]
	}
	return list, d
}
