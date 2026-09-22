package game

import (
	"errors"
	"math/rand/v2"
)

// Deck with cards and trump
type Deck struct {
	cards     []*Card
	trump     *Card
	trumpSuit *Suit
	rng       *rand.Rand
}

const ErrorDeckNoTrump = "No trump card in deck"

// NewDeck makes deck with predefined values for testing. Cards are drawn in given order
func NewDeck(cards []*Card, trump *Card) *Deck {
	return &Deck{cards: cards, trump: trump, trumpSuit: &trump.Suit}
}

// SetRandom sets source of randomness for shuffling. Global source is used if nil
func (d *Deck) SetRandom(r *rand.Rand) {
	d.rng = r
}

// GetCard pops card from top of deck
func (d *Deck) GetCard() *Card {
	if len(d.cards) > 0 {
		card := d.cards[0]
		d.cards = d.cards[1:]
		return card
	}
	// last card
	if d.trump != nil {
		card := d.trump
		d.trump = nil
		return card
	}
	return nil
}

// GetCount returns number of cards in deck (inc. trump)
func (d *Deck) GetCount() int {
	c := len(d.cards)

	if d.trump != nil {
		c++
	}
	return c
}

// ResetDeck refills deck with shuffled cards
func (d *Deck) ResetDeck(leastRank Rank) {
	d.cards = make([]*Card, 0, (Ace+1-leastRank)*4)

	for rank := leastRank; rank <= Ace; rank++ {
		for suit := Hearts; suit <= Diamonds; suit++ {
			d.cards = append(d.cards, &Card{Rank: rank, Suit: suit})
		}
	}

	shuffle := rand.Shuffle
	if d.rng != nil {
		shuffle = d.rng.Shuffle
	}
	shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})

	d.trump = d.GetCard()
	d.trumpSuit = &d.trump.Suit
}

// AddCard puts card into random position in deck
func (d *Deck) AddCard(card *Card) {
	i := intN(d.rng, len(d.cards)+1)
	d.cards = append(d.cards[:i], append([]*Card{card}, d.cards[i:]...)...)
}

// intN returns random number in [0, n) from given source, or from global source if nil
func intN(r *rand.Rand, n int) int {
	if r == nil {
		return rand.IntN(n)
	}
	return r.IntN(n)
}

// GetTrumpSuit returns suit of trump card
func (d *Deck) GetTrumpSuit() *Suit {
	return d.trumpSuit
}

// GetTrump returns trump card for rendering
func (d *Deck) GetTrump() (Card, error) {
	if d.trump != nil {
		return *d.trump, nil
	}
	return Card{}, errors.New(ErrorDeckNoTrump)
}
