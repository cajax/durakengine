package game

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/cajax/durakengine/pkg/game"
)

// expectPlaces checks finishing place of each player, in seat order
func expectPlaces(t *testing.T, players []*game.Player, places ...int) {
	t.Helper()
	for i, p := range players {
		if p.Place() != places[i] {
			t.Errorf("Expected player %d to finish %d, got %d", i+1, places[i], p.Place())
		}
	}
}

func TestPlacesAfterNormalEnd(t *testing.T) {
	g, p := newEmptyDeckGame(
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Seven, game.Spades), card(game.Eight, game.Spades)},
	)
	expectPlaces(t, p, 0, 0)

	// defender picks up the attacker's last card
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	mustSucceed(t, g.Pickup(p[1]))

	if !g.IsOver() || !p[1].IsLoser() {
		t.Fatal("Expected defender to lose")
	}
	expectPlaces(t, p, 1, 2)
}

func TestPlacesSharedByPlayersGoingOutTogether(t *testing.T) {
	g, p := newEmptyDeckGame(
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Ten, game.Spades)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	mustSucceed(t, g.Defend(p[1], 0, card(game.Seven, game.Clubs)))
	mustSucceed(t, g.EndAttack(p[0]))

	if !g.IsOver() || !p[2].IsLoser() {
		t.Fatal("Expected third player to lose")
	}
	expectPlaces(t, p, 1, 1, 3)
	if !p[0].IsFirstWinner() || !p[1].IsFirstWinner() {
		t.Error("Expected both players sharing first place to be first winners")
	}
}

func TestPlacesInDraw(t *testing.T) {
	g, p := newEmptyDeckGame(
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Seven, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	mustSucceed(t, g.Defend(p[1], 0, card(game.Seven, game.Clubs)))
	mustSucceed(t, g.EndAttack(p[0]))

	if !g.IsOver() || p[0].IsLoser() || p[1].IsLoser() {
		t.Fatal("Expected game to end in a draw")
	}
	expectPlaces(t, p, 1, 1)
}

func TestPlacesWithEarlierAbandon(t *testing.T) {
	// seats: C attacks, A defends, B, D
	g, p := newEmptyDeckGame(
		[]*game.Card{card(game.Six, game.Clubs), card(game.Nine, game.Diamonds)},
		[]*game.Card{card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Spades)},
		[]*game.Card{card(game.Ten, game.Spades)},
	)
	c, a, b, d := p[0], p[1], p[2], p[3]

	mustSucceed(t, g.Abandon(d))
	if d.Place() != 4 {
		t.Errorf("Expected first abandoner to get the last place, got %d", d.Place())
	}

	// A beats the last card and goes out, C draws the abandoned card
	mustSucceed(t, g.Attack(c, []*game.Card{card(game.Six, game.Clubs)}))
	mustSucceed(t, g.Defend(a, 0, card(game.Seven, game.Clubs)))
	mustSucceed(t, g.EndAttack(c))
	if a.Place() != 1 || !b.IsAttacker() || !c.IsDefender() {
		t.Fatal("Expected A to go out first and B to attack C")
	}

	// C picks up B's last card
	mustSucceed(t, g.Attack(b, b.GetCards()))
	mustSucceed(t, g.Pickup(c))

	if !g.IsOver() || !c.IsLoser() {
		t.Fatal("Expected C to lose")
	}
	expectPlaces(t, []*game.Player{a, b, c, d}, 1, 2, 3, 4)
	if b.IsFirstWinner() {
		t.Error("Expected only A to be first winner")
	}
}

func TestPlacesWithAbandonAndPlayersGoingOutTogether(t *testing.T) {
	// seats: C attacks, D defends, A, B
	g, p := newEmptyDeckGame(
		[]*game.Card{card(game.Nine, game.Diamonds)},
		[]*game.Card{card(game.Ten, game.Spades)},
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Seven, game.Clubs)},
	)
	c, d, a, b := p[0], p[1], p[2], p[3]

	// C draws the abandoned card, A attacks B next
	mustSucceed(t, g.Abandon(d))
	if !a.IsAttacker() || !b.IsDefender() {
		t.Fatal("Expected A to attack B after the abandon")
	}

	mustSucceed(t, g.Attack(a, a.GetCards()))
	mustSucceed(t, g.Defend(b, 0, card(game.Seven, game.Clubs)))
	mustSucceed(t, g.EndAttack(a))

	if !g.IsOver() || !c.IsLoser() {
		t.Fatal("Expected C to lose")
	}
	expectPlaces(t, []*game.Player{a, b, c, d}, 1, 1, 3, 4)
}

func TestPlacesAfterAbandonEndsGame(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Seven, game.Clubs)},
	)
	mustSucceed(t, g.Abandon(p[1]))

	expectPlaces(t, p, 1, 2)
}

func TestGameOverEventCarriesPlaces(t *testing.T) {
	g, p := newEmptyDeckGame(
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Seven, game.Spades), card(game.Eight, game.Spades)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	mustSucceed(t, g.Pickup(p[1]))

	events := g.Log.Events[len(g.Log.Events)-1]
	over, ok := events[len(events)-1].(*game.OverEvent)
	if !ok {
		t.Fatal("Expected game over event to be logged")
	}
	if len(over.LastPlayers) != 1 || over.LastPlayers[0].Place() != 2 {
		t.Fatal("Expected game over event to carry the loser's place")
	}
	data, err := json.Marshal(over)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"place":2`) {
		t.Errorf("Expected place in JSON, got %s", data)
	}
}
