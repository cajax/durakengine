package game

import (
	"errors"
	"testing"

	"github.com/cajax/durakengine/pkg/game"
)

func card(r game.Rank, s game.Suit) *game.Card {
	return &game.Card{Rank: r, Suit: s}
}

// newTestGame creates a started game with hearts as trump and only the trump card left in deck.
// First player is attacker, second is defender.
func newTestGame(options map[string]game.Option, hands ...[]*game.Card) (*game.Game, []*game.Player) {
	var players []*game.Player
	for i, hand := range hands {
		players = append(players, game.NewPlayer(string(rune('1'+i)), false, false, false, "Player", hand, i == 1, i == 0))
	}
	if options == nil {
		options = map[string]game.Option{}
	}
	deck := game.NewDeck([]*game.Card{}, card(game.King, game.Hearts))
	return game.NewGame(deck, players, options, true, false, game.Table{}, &game.BotManager{}), players
}

// newEmptyDeckGame is newTestGame with nothing left in deck, so players who run out of cards go out
func newEmptyDeckGame(hands ...[]*game.Card) (*game.Game, []*game.Player) {
	_, players := newTestGame(nil, hands...)
	deck := game.NewDeck([]*game.Card{}, card(game.King, game.Hearts))
	deck.GetCard()
	return game.NewGame(deck, players, map[string]game.Option{}, true, false, game.Table{}, &game.BotManager{}), players
}

func expectError(t *testing.T, err error, expected string) {
	t.Helper()
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error %q, got %v", expected, err)
		return
	}
	var ruleErr *game.RuleError
	if !errors.As(err, &ruleErr) {
		t.Errorf("Expected a *game.RuleError, got %T", err)
	}
}

func mustSucceed(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func lastEventType(g *game.Game) game.EventType {
	events := g.Log.Events[len(g.Log.Events)-1]
	return events[len(events)-1].GetType()
}
