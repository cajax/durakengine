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

func TestLogDiscard(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs), card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Nine, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], []*game.Card{card(game.Six, game.Clubs)}))
	mustSucceed(t, g.Defend(p[1], 0, card(game.Eight, game.Clubs)))
	mustSucceed(t, g.EndAttack(p[0]))

	events := g.Log.GetEvents(g.GetSequence(), g.GetSequence()+1)
	discard, ok := events[0].(*game.DiscardEvent)
	if !ok {
		t.Fatalf("Expected discard event, got %v", events[0].GetType())
	}
	if len(discard.GetCards()) != 2 || discard.GetPlayer().ID != p[0].ID {
		t.Error("Expected discard event with attacker and both cards from table")
	}
}

func TestLogTransfer(t *testing.T) {
	g, p := newRedirectGame()
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	mustSucceed(t, g.Redirect(p[1], []*game.Card{card(game.Six, game.Hearts)}))

	events := g.Log.GetEvents(g.GetSequence(), g.GetSequence()+1)
	transfer, ok := events[0].(*game.TransferEvent)
	if !ok {
		t.Fatalf("Expected transfer event, got %v", events[0].GetType())
	}
	if transfer.GetPlayer().ID != p[1].ID || transfer.NextDefender.ID != p[2].ID || len(transfer.GetCards()) != 1 {
		t.Error("Expected transfer event from defender to next defender with redirected card")
	}
}
