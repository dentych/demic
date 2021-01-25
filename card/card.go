package card

import "math/rand"

const (
	SuitHeart   = 'H'
	SuitDiamond = 'D'
	SuitClub    = 'C'
	SuitSpade   = 'S'
)

type Card struct {
	Rank string
	Suit rune
}

func Deal(deck []Card, n int) ([]Card, []Card) {
	var cardList []Card
	for _, card := range deck[0:n] {
		cardList = append(cardList, card)
		deck = deck[1:]
	}
	return cardList, deck
}

func NewDeck() (deck []Card) {
	// Valid ranks include Two, Three, Four, Five, Six
	// Seven, Eight, Nine, Ten, Jack, Queen, King & Ace
	ranks := []string{"2", "3", "4", "5", "6", "7",
		"8", "9", "10", "J", "Q", "K", "A"}

	// Valid suits include Heart, Diamond, Club & Spade
	suits := []rune{'H', 'D', 'C', 'S'}

	// Loop over each type and suit appending to the deck
	for i := 0; i < len(ranks); i++ {
		for n := 0; n < len(suits); n++ {
			card := Card{
				Rank: ranks[i],
				Suit: suits[n],
			}
			deck = append(deck, card)
		}
	}
	return
}

func Shuffle(d []Card) []Card {
	for i := 0; i < len(d); i++ {
		rand.Shuffle(len(d), func (i, j int) {
			d[i], d[j] = d[j], d[i]
		})
	}
	return d
}
