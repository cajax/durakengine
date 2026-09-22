package game

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/cajax/durakengine/pkg/game"
)

func TestGetCard(t *testing.T) {
	deck := game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Spades, Rank: game.Seven})
	if deck.GetCount() != 1 {
		t.Error("Expected last card")
	}

	deck.GetCard()
	if deck.GetCount() != 0 {
		t.Error("Expected empty deck")
	}
}

func TestResetDeck(t *testing.T) {
	deck := game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Spades, Rank: game.Seven})

	deck.ResetDeck(game.Ace)
	if deck.GetCount() != 4 {
		t.Error("Expected deck with exactly 4 aces")
	}

	deck.ResetDeck(game.King)
	if deck.GetCount() != 8 {
		t.Error("Expected deck with exactly 4 aces and 4 kings")
	}
}

func TestGetTrumpSuit(t *testing.T) {
	deck := game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Spades, Rank: game.Seven})

	if *deck.GetTrumpSuit() != game.Spades {
		t.Error("Unexpected trump suit")
	}

	deck = game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Hearts, Rank: game.Seven})

	if *deck.GetTrumpSuit() != game.Hearts {
		t.Error("Unexpected trump suit")
	}
}

func TestDeckDealsEveryCardOnceWithTrumpLast(t *testing.T) {
	deck := game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Spades, Rank: game.Seven})
	deck.ResetDeck(game.Six)

	trump, err := deck.GetTrump()
	if err != nil {
		t.Fatal(err)
	}

	seen := map[game.Card]bool{}
	var last *game.Card
	for c := deck.GetCard(); c != nil; c = deck.GetCard() {
		if seen[*c] {
			t.Fatalf("Card %s dealt twice", c.ToString())
		}
		seen[*c] = true
		last = c
	}

	if len(seen) != 36 {
		t.Errorf("Expected 36 cards, got %d", len(seen))
	}
	if *last != trump {
		t.Error("Expected trump card to be dealt last")
	}
	if _, err := deck.GetTrump(); err == nil || err.Error() != game.ErrorDeckNoTrump {
		t.Error("Expected no trump after it was dealt")
	}
}

func TestDeckDrawsInGivenOrder(t *testing.T) {
	cards := []*game.Card{{Suit: game.Clubs, Rank: game.Six}, {Suit: game.Spades, Rank: game.Ten}}
	trump := &game.Card{Suit: game.Hearts, Rank: game.Ace}
	deck := game.NewDeck(cards, trump)

	for _, expected := range []*game.Card{cards[0], cards[1], trump} {
		if c := deck.GetCard(); c != expected {
			t.Errorf("Expected %s, got %v", expected.ToString(), c)
		}
	}
	if deck.GetCard() != nil {
		t.Error("Expected empty deck")
	}
}

func TestDeckShuffleIsReproducibleWithSeed(t *testing.T) {
	deal := func(seed uint64) []game.Card {
		deck := game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Spades, Rank: game.Seven})
		deck.SetRandom(rand.New(rand.NewPCG(seed, 0)))
		deck.ResetDeck(game.Six)
		var cards []game.Card
		for c := deck.GetCard(); c != nil; c = deck.GetCard() {
			cards = append(cards, *c)
		}
		return cards
	}

	if !slices.Equal(deal(1), deal(1)) {
		t.Error("Expected same order with same seed")
	}
	if slices.Equal(deal(1), deal(2)) {
		t.Error("Expected different order with different seed")
	}
}

func TestAddCard(t *testing.T) {
	deck := game.NewDeck([]*game.Card{{Suit: game.Clubs, Rank: game.Six}}, &game.Card{Suit: game.Spades, Rank: game.Seven})
	deck.AddCard(&game.Card{Suit: game.Clubs, Rank: game.Ten})

	if deck.GetCount() != 3 {
		t.Errorf("Expected 3 cards, got %d", deck.GetCount())
	}
	if c := deck.GetCard(); c.Suit != game.Clubs {
		t.Error("Expected trump card to stay at the bottom")
	}
}
