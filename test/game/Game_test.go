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

func TestStartGameWhenDealExhaustsDeck(t *testing.T) {
	var players []*game.Player
	for i := 0; i < 6; i++ {
		players = append(players, &game.Player{Name: "Player"})
	}
	options := map[string]game.Option{"min_rank": {Value: "6"}}
	g := game.NewGame(&game.Deck{}, players, options, false, false, game.Table{}, nil)

	if err := g.StartGame(); err != nil {
		t.Fatal(err)
	}

	deck := g.GetDeck()
	if deck.GetCount() != 0 {
		t.Error("Expected all 36 cards to be dealt")
	}

	if _, d := g.GetDefender(); d == nil {
		t.Error("Expected defender to be selected")
	}
}

func TestDefendFailWithInvalidPairIndex(t *testing.T) {
	players := []*game.Player{
		game.NewPlayer("1", false, false, false, "Attacker", []*game.Card{
			{Suit: game.Clubs, Rank: game.Six},
		}, false, true),
		game.NewPlayer("2", false, false, false, "Defender", []*game.Card{
			{Suit: game.Clubs, Rank: game.Ace},
		}, true, false),
	}
	g := game.NewGame(game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Hearts, Rank: game.King}), players, map[string]game.Option{}, true, false, game.Table{}, &game.BotManager{})

	if err := g.Attack(players[0], players[0].GetCards()); err != nil {
		t.Fatal(err)
	}

	for _, i := range []int{-1, 1} {
		err := g.Defend(players[1], i, players[1].GetCards()[0])
		if err == nil || err.Error() != game.ErrorPairIndexOutOfRange {
			t.Errorf("Defend with pair index %d should fail with %q, got %v", i, game.ErrorPairIndexOutOfRange, err)
		}
	}
}

func TestAttackFailWithDuplicateCard(t *testing.T) {
	players := []*game.Player{
		game.NewPlayer("1", false, false, false, "Attacker", []*game.Card{
			{Suit: game.Clubs, Rank: game.Six},
		}, false, true),
		game.NewPlayer("2", false, false, false, "Defender", []*game.Card{
			{Suit: game.Clubs, Rank: game.Seven},
			{Suit: game.Clubs, Rank: game.Eight},
		}, true, false),
	}
	g := game.NewGame(game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Hearts, Rank: game.King}), players, map[string]game.Option{}, true, false, game.Table{}, &game.BotManager{})

	c := players[0].GetCards()[0]
	err := g.Attack(players[0], []*game.Card{c, c})

	if err == nil || err.Error() != game.ErrorAttackerHasNoCard {
		t.Errorf("Attack with the same card twice should fail with %q, got %v", game.ErrorAttackerHasNoCard, err)
	}
	if !g.GetTable().IsEmpty() {
		t.Error("Table should stay empty after a rejected attack")
	}
}

func newRedirectGame() (*game.Game, []*game.Player) {
	players := []*game.Player{
		game.NewPlayer("1", false, false, false, "Attacker", []*game.Card{
			{Suit: game.Clubs, Rank: game.Six},
		}, false, true),
		game.NewPlayer("2", false, false, false, "Defender", []*game.Card{
			{Suit: game.Hearts, Rank: game.Six},
			{Suit: game.Clubs, Rank: game.Ace},
		}, true, false),
		game.NewPlayer("3", false, false, false, "Next defender", []*game.Card{
			{Suit: game.Spades, Rank: game.Ten},
			{Suit: game.Spades, Rank: game.Jack},
		}, false, false),
	}
	options := map[string]game.Option{"with_redirect": {Value: "1"}}
	g := game.NewGame(game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Spades, Rank: game.King}), players, options, true, false, game.Table{}, &game.BotManager{})
	return g, players
}

func TestRedirectRemovesCardFromHand(t *testing.T) {
	g, players := newRedirectGame()
	if err := g.Attack(players[0], players[0].GetCards()); err != nil {
		t.Fatal(err)
	}

	sixHearts := players[1].GetCards()[0]
	if err := g.Redirect(players[1], []*game.Card{sixHearts}); err != nil {
		t.Fatal(err)
	}

	for _, c := range players[1].GetCards() {
		if c.Rank == sixHearts.Rank && c.Suit == sixHearts.Suit {
			t.Error("Redirected card should be removed from defender's hand")
		}
	}
	if len(g.GetTable().GetPairs()) != 2 {
		t.Error("Expected 2 cards on table after redirect")
	}
	if !players[2].IsDefender() {
		t.Error("Expected next player to become defender")
	}
}

func TestRedirectFailWithCardNotInHand(t *testing.T) {
	g, players := newRedirectGame()
	if err := g.Attack(players[0], players[0].GetCards()); err != nil {
		t.Fatal(err)
	}

	err := g.Redirect(players[1], []*game.Card{{Suit: game.Diamonds, Rank: game.Six}})
	if err == nil || err.Error() != game.ErrorDefenderHasNoCard {
		t.Errorf("Redirect with a card not in hand should fail with %q, got %v", game.ErrorDefenderHasNoCard, err)
	}
}
