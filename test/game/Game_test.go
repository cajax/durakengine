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

	if err := g.Attack(players[0], c1); err != nil {
		t.Fatal(err)
	}
	if err := g.Defend(players[1], 0, players[1].GetCards()[0]); err != nil {
		t.Fatal(err)
	}

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

func hasCard(cards []*game.Card, c *game.Card) bool {
	for _, card := range cards {
		if card.Rank == c.Rank && card.Suit == c.Suit {
			return true
		}
	}
	return false
}

func newAbandonGame() (*game.Game, []*game.Player) {
	players := []*game.Player{
		game.NewPlayer("1", false, false, false, "Attacker", []*game.Card{
			{Suit: game.Clubs, Rank: game.Six},
			{Suit: game.Hearts, Rank: game.Nine},
		}, false, true),
		game.NewPlayer("2", false, false, false, "Defender", []*game.Card{
			{Suit: game.Clubs, Rank: game.Seven},
			{Suit: game.Hearts, Rank: game.Ten},
		}, true, false),
		game.NewPlayer("3", false, false, false, "Player 3", []*game.Card{
			{Suit: game.Spades, Rank: game.Ten},
			{Suit: game.Spades, Rank: game.Jack},
		}, false, false),
	}
	g := game.NewGame(game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Diamonds, Rank: game.King}), players, map[string]game.Option{}, true, false, game.Table{}, &game.BotManager{})
	return g, players
}

func TestAbandonByDefenderReturnsCardsToAttacker(t *testing.T) {
	g, players := newAbandonGame()
	attack := players[0].GetCards()[0]
	if err := g.Attack(players[0], []*game.Card{attack}); err != nil {
		t.Fatal(err)
	}
	if err := g.Defend(players[1], 0, players[1].GetCards()[0]); err != nil {
		t.Fatal(err)
	}

	if err := g.Abandon(players[1]); err != nil {
		t.Fatal(err)
	}

	if !g.GetTable().IsEmpty() {
		t.Error("Table should be cleared when defender abandons")
	}
	if !hasCard(players[0].GetCards(), attack) {
		t.Error("Attack card should be returned to attacker")
	}
	if len(players[1].GetCards()) != 0 {
		t.Error("Abandoning player should have no cards")
	}
	if n := countCards(g); n != 7 {
		t.Errorf("Expected 7 cards in game after abandon, got %d", n)
	}
}

func TestAbandonByAttackerClearsTable(t *testing.T) {
	g, players := newAbandonGame()
	if err := g.Attack(players[0], []*game.Card{players[0].GetCards()[0]}); err != nil {
		t.Fatal(err)
	}

	if err := g.Abandon(players[0]); err != nil {
		t.Fatal(err)
	}

	if !g.GetTable().IsEmpty() {
		t.Error("Table should be cleared when attacker abandons")
	}
	if len(players[0].GetCards()) != 0 {
		t.Error("Abandoning player should have no cards")
	}
	if n := countCards(g); n != 7 {
		t.Errorf("Expected 7 cards in game after abandon, got %d", n)
	}
}

// countCards returns number of cards in hands, on table and in deck
func countCards(g *game.Game) int {
	deck := g.GetDeck()
	n := deck.GetCount() + len(g.GetTable().GetCardsOnTable())
	for _, p := range g.GetPlayers() {
		n += len(p.GetCards())
	}
	return n
}

func TestGetPairs(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs), card(game.Six, game.Spades)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Nine, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	mustSucceed(t, g.Defend(p[1], 0, card(game.Eight, game.Clubs)))

	pairs := g.GetPairs()
	if len(pairs) != 2 {
		t.Fatalf("Expected 2 pairs, got %d", len(pairs))
	}
	if *pairs[0].Attack != *card(game.Six, game.Clubs) || *pairs[0].Defense != *card(game.Eight, game.Clubs) {
		t.Error("Expected first pair to be 6♣ beaten by 8♣")
	}
	if *pairs[1].Attack != *card(game.Six, game.Spades) || pairs[1].Defense != nil {
		t.Error("Expected second pair to be unbeaten 6♠")
	}
}

func TestRefillAfterDefenseStartsWithAttacker(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Seven, game.Clubs), card(game.Eight, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	mustSucceed(t, g.Defend(p[1], 0, card(game.Seven, game.Clubs)))
	mustSucceed(t, g.EndAttack(p[0]))

	if g.IsOver() || p[0].HasQuit() {
		t.Fatal("Attacker should draw the last card from deck and stay in game")
	}
	if !hasCard(p[0].GetCards(), card(game.King, game.Hearts)) || len(p[1].GetCards()) != 1 {
		t.Error("Expected attacker to draw before defender")
	}
	if !p[1].IsAttacker() || !p[0].IsDefender() {
		t.Error("Expected defender to attack next")
	}
}

func TestRefillAfterPickupEndsWithDefender(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs), card(game.Seven, game.Clubs), card(game.Eight, game.Clubs), card(game.Nine, game.Clubs), card(game.Ten, game.Clubs), card(game.Jack, game.Clubs)},
		[]*game.Card{card(game.Six, game.Spades)},
		[]*game.Card{card(game.Six, game.Diamonds), card(game.Seven, game.Diamonds), card(game.Eight, game.Diamonds), card(game.Nine, game.Diamonds), card(game.Ten, game.Diamonds)},
	)
	mustSucceed(t, g.Attack(p[0], []*game.Card{card(game.Six, game.Clubs)}))
	mustSucceed(t, g.Pickup(p[1]))

	if !hasCard(p[0].GetCards(), card(game.King, game.Hearts)) {
		t.Error("Expected attacker to draw first")
	}
	if len(p[2].GetCards()) != 5 || len(p[1].GetCards()) != 2 {
		t.Error("Expected other players to draw only after attacker")
	}
}
