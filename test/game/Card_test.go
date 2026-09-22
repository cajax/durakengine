package game

import (
	"testing"

	"github.com/cajax/durakengine/pkg/game"
)

func TestCardToString(t *testing.T) {
	if s := card(game.Ten, game.Hearts).ToString(); s != "10♥" {
		t.Errorf("Expected 10♥, got %s", s)
	}
	if s := card(game.Ace, game.Spades).ToString(); s != "A♠" {
		t.Errorf("Expected A♠, got %s", s)
	}
}

func TestRankFromString(t *testing.T) {
	for r := game.Two; r <= game.Ace; r++ {
		if parsed := game.RankFromString(r.ToString()); parsed != r {
			t.Errorf("Expected %s to parse to %d, got %d", r.ToString(), r, parsed)
		}
	}
}
