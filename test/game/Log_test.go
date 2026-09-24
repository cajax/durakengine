package game

import (
	"encoding/json"
	"math"
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

func TestLogGetEventsOutOfRange(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Nine, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	mustSucceed(t, g.Defend(p[1], 0, card(game.Eight, game.Clubs)))
	n := len(g.Log.Events)

	cases := []struct {
		first, last, expected int
	}{
		{n + 1, n + 1, 0},
		{n + 2, n + 10, 0},
		{math.MaxInt, math.MinInt, 0},
		{math.MinInt, math.MaxInt, 2},
		{n, n, 1},
	}
	for _, c := range cases {
		if events := g.Log.GetEvents(c.first, c.last); len(events) != c.expected {
			t.Errorf("GetEvents(%d, %d): expected %d events, got %d", c.first, c.last, c.expected, len(events))
		}
	}

	empty := game.NewLog()
	if events := empty.GetEvents(2, 1); len(events) != 0 {
		t.Errorf("Expected no events from empty log, got %v", events)
	}
}

func TestLogDefensePairIndex(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs), card(game.Six, game.Spades)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Eight, game.Spades)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	mustSucceed(t, g.Defend(p[1], 1, card(game.Eight, game.Spades)))
	mustSucceed(t, g.Defend(p[1], 0, card(game.Eight, game.Clubs)))
	mustSucceed(t, g.EndAttack(p[0]))

	var indices []int
	for _, e := range g.Log.GetEvents(1, g.GetSequence()+1) {
		if d, ok := e.(*game.DefenseEvent); ok {
			indices = append(indices, d.GetPairIndex())
		}
	}
	if len(indices) != 2 || indices[0] != 1 || indices[1] != 0 {
		t.Errorf("Expected defense pair indices [1 0], got %v", indices)
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

func TestLogSerializesPlayerState(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs), card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Nine, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], []*game.Card{card(game.Six, game.Clubs)}))
	mustSucceed(t, g.Defend(p[1], 0, card(game.Eight, game.Clubs)))

	data, err := json.Marshal(g.Log.GetEvents(1, 3))
	if err != nil {
		t.Fatal(err)
	}

	var events []struct {
		Type   string `json:"type"`
		Player struct {
			ID       string `json:"id"`
			Attacker bool   `json:"attacker"`
			Defender bool   `json:"defender"`
			Cards    []game.Card
		} `json:"player"`
		Pair struct {
			AttackerID string `json:"attacker_id"`
			DefenderID string `json:"defender_id"`
		} `json:"pair"`
	}
	if err := json.Unmarshal(data, &events); err != nil {
		t.Fatal(err)
	}

	attack := events[0]
	if !attack.Player.Attacker || len(attack.Player.Cards) != 1 || attack.Player.Cards[0] != *card(game.Seven, game.Clubs) {
		t.Errorf("Expected attacker with remaining hand in attack event, got %+v", attack.Player)
	}

	defense := events[1]
	if !defense.Player.Defender || len(defense.Player.Cards) != 1 {
		t.Errorf("Expected defender with remaining hand in defense event, got %+v", defense.Player)
	}
	if defense.Pair.AttackerID != p[0].ID || defense.Pair.DefenderID != p[1].ID {
		t.Errorf("Expected pair to reference players by ID, got %+v", defense.Pair)
	}
}

func TestLogKeepsStateAtTimeOfEvent(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs), card(game.Six, game.Spades)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Nine, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], []*game.Card{card(game.Six, game.Clubs)}))
	attack := g.Log.GetEvents(1, 2)[0].(*game.AttackEvent)
	before, _ := json.Marshal(attack)

	mustSucceed(t, g.Attack(p[0], []*game.Card{card(game.Six, game.Spades)}))
	mustSucceed(t, g.Pickup(p[1]))

	after, _ := json.Marshal(attack)
	if string(before) != string(after) {
		t.Errorf("Expected logged event not to change, before %s, after %s", before, after)
	}
}
