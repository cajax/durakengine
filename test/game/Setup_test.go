package game

import (
	"testing"

	"github.com/cajax/durakengine/pkg/game"
)

func newUnstartedGame(n int) *game.Game {
	var players []*game.Player
	for i := 0; i < n; i++ {
		players = append(players, &game.Player{Name: "Player"})
	}
	options := map[string]game.Option{"min_rank": {Value: "6"}}
	return game.NewGame(&game.Deck{}, players, options, false, false, game.Table{}, nil)
}

func TestStartGameFailWithTooFewPlayers(t *testing.T) {
	expectError(t, newUnstartedGame(1).StartGame(), game.ErrorTooFewPlayers)
}

func TestStartGameFailWithTooManyPlayers(t *testing.T) {
	expectError(t, newUnstartedGame(7).StartGame(), game.ErrorTooManyPlayers)
}

func TestStartGame(t *testing.T) {
	g := newUnstartedGame(4)
	mustSucceed(t, g.StartGame())

	if !g.IsStarted() || g.IsOver() {
		t.Error("Expected game to be in progress")
	}

	deck := g.GetDeck()
	if deck.GetCount() != 36-4*6 {
		t.Errorf("Expected %d cards in deck, got %d", 36-4*6, deck.GetCount())
	}

	players := g.GetPlayers()
	attackerIndex := -1
	for i, p := range players {
		if len(p.GetCards()) != 6 {
			t.Errorf("Expected player %d to have 6 cards, got %d", i, len(p.GetCards()))
		}
		if p.IsAttacker() {
			attackerIndex = i
		}
	}
	if attackerIndex < 0 {
		t.Fatal("Expected attacker to be selected")
	}

	defenderIndex, _ := g.GetDefender()
	if defenderIndex != (attackerIndex+1)%len(players) {
		t.Error("Expected defender to sit to the left of attacker")
	}

	if g.GetSequence() != 1 || lastEventType(g) != game.StartEventType {
		t.Error("Expected start event in first sequence")
	}
}

func TestStartGameSelectsAttackerWithLeastTrump(t *testing.T) {
	for i := 0; i < 20; i++ {
		g := newUnstartedGame(4)
		mustSucceed(t, g.StartGame())

		deck := g.GetDeck()
		trump := *deck.GetTrumpSuit()
		leastRank := game.Ace + 1
		var expected *game.Player
		for _, p := range g.GetPlayers() {
			for _, c := range p.GetCards() {
				if c.Suit == trump && c.Rank < leastRank {
					leastRank = c.Rank
					expected = p
				}
			}
		}

		if expected != nil && !expected.IsAttacker() {
			t.Fatal("Expected player with the least trump to attack first")
		}
	}
}
