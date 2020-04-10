package game

import (
	"errors"
	"math/rand"
	"time"
)

// Deck with cards and trump
type Deck struct {
	cards     []*Card
	trump     *Card
	trumpSuit *Suit
}

const ErrorDeckNoTrump = "No trump card in deck"

// NewDeck makes deck with predefined values for testing
func NewDeck(cards []*Card, trump *Card) *Deck {
	return &Deck{cards: cards, trump: trump, trumpSuit: &trump.Suit}
}

// GetCard pops random card from deck
func (d *Deck) GetCard() *Card {
	if len(d.cards) > 0 {
		rand.Seed(time.Now().Unix())
		idx := rand.Intn(len(d.cards))
		card := d.cards[idx]
		d.cards = append(d.cards[:idx], d.cards[idx+1:]...)
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

// ResetDeck refills deck with cards
func (d *Deck) ResetDeck(leastRank Rank) {
	d.cards = make([]*Card, 0, (Ace+1-leastRank)*4)

	for rank := leastRank; rank <= Ace; rank++ {
		for suit := Hearts; suit <= Diamonds; suit++ {
			d.AddCard(&Card{Rank: rank, Suit: suit})
		}
	}

	d.trump = d.GetCard()
	d.trumpSuit = &d.trump.Suit
}

// AddCard adds card to the end of deck
func (d *Deck) AddCard(card *Card) {
	d.cards = append(d.cards, card)
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
