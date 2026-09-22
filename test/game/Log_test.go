package game

import (
	"testing"

	"github.com/cajax/durakengine/pkg/game"
)

func TestLogSequence(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Nine, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	mustSucceed(t, g.Defend(p[1], 0, card(game.Eight, game.Clubs)))

	if g.GetSequence() != 2 {
		t.Fatalf("Expected sequence 2, got %d", g.GetSequence())
	}

	events := g.Log.GetEvents(1, 3)
	if len(events) != 2 || events[0].GetType() != game.AttackEventType || events[1].GetType() != game.DefenseEventType {
		t.Fatalf("Expected attack and defense events, got %v", events)
	}
	for i, e := range events {
		if e.GetSequence() != i+1 {
			t.Errorf("Expected event %d to have sequence %d, got %d", i, i+1, e.GetSequence())
		}
	}

	events = g.Log.GetEvents(2, 3)
	if len(events) != 1 || events[0].GetType() != game.DefenseEventType {
		t.Errorf("Expected only defense event, got %v", events)
	}

	if len(g.Log.GetEvents(0, 100)) != 2 {
		t.Error("Expected out of range bounds to be clamped")
	}
}
