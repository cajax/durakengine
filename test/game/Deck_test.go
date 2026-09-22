package game

import (
	"github.com/cajax/durakengine/pkg/game"
	"testing"
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
