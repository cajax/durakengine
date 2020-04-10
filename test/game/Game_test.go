package game

import (
	"github.com/cajax/durakengine/pkg/game"
	"testing"
)

func TestAttackFailWithTooManyCards(t *testing.T) {
	players := []*game.Player{
		game.NewPlayer("1", false, false, false, "Attacker 1", []*game.Card{
			{Suit: game.Clubs, Rank: game.Ten},
			{Suit: game.Spades, Rank: game.Ten},
			{Suit: game.Diamonds, Rank: game.Ten},
		}, false, true),
		game.NewPlayer("2", false, false, false, "Defender", []*game.Card{
			{Suit: game.Clubs, Rank: game.Seven},
			//&game.Card{Suit: game.Spades, Rank: game.Seven},
			{Suit: game.Hearts, Rank: game.Ten},
		}, true, false),
		{Name: "Player 3"},
	}
	g := game.NewGame(&game.Deck{}, players, map[string]game.Option{}, true, false, game.Table{}, &game.BotManager{})

	c := []*game.Card{
		{Suit: game.Clubs, Rank: game.Ten},
		{Suit: game.Spades, Rank: game.Ten},
		{Suit: game.Diamonds, Rank: game.Ten},
	}
	err := g.Attack(players[0], c)

	if err.Error() != game.ErrorAttackIsTooBig {
		t.Error("Attack should fail with too many cards")
	}
}

func TestAttackFailWithTooManyAddedCards(t *testing.T) {
	players := []*game.Player{
		game.NewPlayer("1", false, false, false, "Attacker 1", []*game.Card{
			{Suit: game.Clubs, Rank: game.Six},
			{Suit: game.Spades, Rank: game.Ten},
			{Suit: game.Diamonds, Rank: game.Ten},
			{Suit: game.Hearts, Rank: game.Ten},
		}, false, true),
		game.NewPlayer("2", false, false, false, "Defender", []*game.Card{
			{Suit: game.Clubs, Rank: game.Ten},
			{Suit: game.Clubs, Rank: game.Jack},
		}, true, false),
		game.NewPlayer("3", false, false, false, "Attacker 2", []*game.Card{
			{Suit: game.Spades, Rank: game.Ten},
			{Suit: game.Diamonds, Rank: game.Ten},
			{Suit: game.Hearts, Rank: game.Ten},
		}, false, true),
	}
	g := game.NewGame(game.NewDeck([]*game.Card{}, &game.Card{Rank: game.King, Suit: game.Hearts}), players, map[string]game.Option{}, true, false, game.Table{}, &game.BotManager{})

	c1 := []*game.Card{
		players[0].GetCards()[0],
	}

	g.Attack(players[0], c1)
	g.Defend(players[1], 0, players[1].GetCards()[0])

	err := g.Attack(players[2], []*game.Card{
		players[2].GetCards()[0],
	})

	if err != nil {
		t.Error("Defender still have some cards")
	}

	err = g.Attack(players[2], []*game.Card{
		players[2].GetCards()[0],
	})
	if err.Error() != game.ErrorAttackIsTooBig {
		t.Error("Adding should fail with too many cards")
	}
}
