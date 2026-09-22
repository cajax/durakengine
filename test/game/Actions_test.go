package game

import (
	"testing"

	"github.com/cajax/durakengine/pkg/game"
)

var withRedirect = map[string]game.Option{"with_redirect": {Value: "1"}}

func TestActionsFailWhenGameNotInProgress(t *testing.T) {
	for _, state := range []struct{ started, over bool }{{false, false}, {true, true}} {
		players := []*game.Player{
			game.NewPlayer("1", false, false, false, "Attacker", []*game.Card{card(game.Six, game.Clubs)}, false, true),
			game.NewPlayer("2", false, false, false, "Defender", []*game.Card{card(game.Seven, game.Clubs)}, true, false),
		}
		deck := game.NewDeck([]*game.Card{}, card(game.King, game.Hearts))
		g := game.NewGame(deck, players, withRedirect, state.started, state.over, game.Table{}, nil)

		expectError(t, g.Attack(players[0], players[0].GetCards()), game.ErrorNotInPlayingState)
		expectError(t, g.Defend(players[1], 0, players[1].GetCards()[0]), game.ErrorNotInPlayingState)
		expectError(t, g.Pickup(players[1]), game.ErrorNotInPlayingState)
		expectError(t, g.EndAttack(players[0]), game.ErrorNotInPlayingState)
		expectError(t, g.Redirect(players[1], players[1].GetCards()), game.ErrorNotInPlayingState)
		expectError(t, g.Abandon(players[0]), game.ErrorNotInPlayingState)
	}
}

func TestAttack(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs), card(game.Six, game.Spades), card(game.Ten, game.Clubs)},
		[]*game.Card{card(game.Seven, game.Clubs), card(game.Eight, game.Clubs)},
	)

	mustSucceed(t, g.Attack(p[0], []*game.Card{card(game.Six, game.Clubs), card(game.Six, game.Spades)}))

	if len(g.GetTable().GetCardsToBeat()) != 2 {
		t.Error("Expected 2 cards to beat on table")
	}
	if len(p[0].GetCards()) != 1 {
		t.Error("Expected attack cards to be removed from attacker's hand")
	}
	if g.GetSequence() != 1 || lastEventType(g) != game.AttackEventType {
		t.Error("Expected attack event to be logged")
	}
}

func TestAttackFailWithDifferentRanks(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs), card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Nine, game.Clubs)},
	)
	expectError(t, g.Attack(p[0], p[0].GetCards()), game.ErrorFirstAttackWithDifferentRanks)
}

func TestAttackFailWithRankNotOnTable(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs), card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Nine, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], []*game.Card{card(game.Six, game.Clubs)}))
	expectError(t, g.Attack(p[0], []*game.Card{card(game.Seven, game.Clubs)}), game.ErrorAttackRankNotOnTable)
}

func TestAttackFailWithCardNotInHand(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Clubs)},
	)
	expectError(t, g.Attack(p[0], []*game.Card{card(game.Six, game.Spades)}), game.ErrorAttackerHasNoCard)
}

func TestAttackFailByWrongPlayer(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Clubs)},
		[]*game.Card{card(game.Six, game.Spades)},
		[]*game.Card{card(game.Six, game.Diamonds)},
	)
	expectError(t, g.Attack(p[2], p[2].GetCards()), game.ErrorAttackByWrongPlayer)
}

func TestAttackByNeighbor(t *testing.T) {
	g, p := newTestGame(withRedirect,
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Nine, game.Clubs)},
		[]*game.Card{card(game.Six, game.Spades)},
		[]*game.Card{card(game.Six, game.Diamonds)},
	)

	expectError(t, g.Attack(p[2], p[2].GetCards()), game.ErrorFirstAttackByNeighbor)

	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	mustSucceed(t, g.Attack(p[2], p[2].GetCards()))

	if len(g.GetTable().GetCardsToBeat()) != 2 {
		t.Error("Expected neighbor to add card to table")
	}
}

func TestDefend(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Ace, game.Clubs)},
		[]*game.Card{card(game.Six, game.Hearts), card(game.Eight, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))

	mustSucceed(t, g.Defend(p[1], 0, card(game.Six, game.Hearts)))

	if g.GetTable().HasUnbeatenCards() {
		t.Error("Expected trump to beat attack card")
	}
	if hasCard(p[1].GetCards(), card(game.Six, game.Hearts)) {
		t.Error("Expected defense card to be removed from defender's hand")
	}
	if lastEventType(g) != game.DefenseEventType {
		t.Error("Expected defense event to be logged")
	}
}

func TestDefendFails(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Six, game.Clubs), card(game.Ace, game.Spades), card(game.Eight, game.Clubs), card(game.Nine, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))

	expectError(t, g.Defend(p[0], 0, card(game.Seven, game.Clubs)), game.ErrorNotDefender)
	expectError(t, g.Defend(p[1], 0, card(game.Ten, game.Clubs)), game.ErrorDefenderHasNoCard)
	expectError(t, g.Defend(p[1], 0, card(game.Six, game.Clubs)), game.ErrorCardCantBeat)
	expectError(t, g.Defend(p[1], 0, card(game.Ace, game.Spades)), game.ErrorCardCantBeat)

	mustSucceed(t, g.Defend(p[1], 0, card(game.Eight, game.Clubs)))
	expectError(t, g.Defend(p[1], 0, card(game.Nine, game.Clubs)), game.ErrorAlreadyDefended)
}

func TestPickup(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs), card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Six, game.Spades)},
		[]*game.Card{card(game.Six, game.Diamonds)},
	)
	mustSucceed(t, g.Attack(p[0], []*game.Card{card(game.Six, game.Clubs)}))

	expectError(t, g.Pickup(p[0]), game.ErrorDefenseByNotDefender)
	mustSucceed(t, g.Pickup(p[1]))

	if !g.GetTable().IsEmpty() {
		t.Error("Expected table to be empty after pickup")
	}
	if !hasCard(p[1].GetCards(), card(game.Six, game.Clubs)) {
		t.Error("Expected defender to take cards from table")
	}
	if !p[2].IsAttacker() || !p[0].IsDefender() {
		t.Error("Expected player to the left of defender to attack next")
	}
}

func TestEndAttack(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs), card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Six, game.Spades)},
		[]*game.Card{card(game.Six, game.Diamonds)},
	)

	expectError(t, g.EndAttack(p[1]), game.ErrorGameEndTurnByNotAttacker)
	expectError(t, g.EndAttack(p[0]), game.ErrorGameEndTurnWithoutAttack)

	mustSucceed(t, g.Attack(p[0], []*game.Card{card(game.Six, game.Clubs)}))
	expectError(t, g.EndAttack(p[0]), game.ErrorGameEndTurnWithoutUnbeatenCards)

	mustSucceed(t, g.Defend(p[1], 0, card(game.Eight, game.Clubs)))
	mustSucceed(t, g.EndAttack(p[0]))

	if !g.GetTable().IsEmpty() {
		t.Error("Expected table to be cleared")
	}
	if !p[1].IsAttacker() || !p[2].IsDefender() {
		t.Error("Expected defender to attack next")
	}
	// 5 cards in hands and trump in deck, 2 cards discarded
	if n := countCards(g); n != 4 {
		t.Errorf("Expected 4 cards in game after discard, got %d", n)
	}
}

func TestRedirectFails(t *testing.T) {
	g, p := newRedirectGame()
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))

	expectError(t, g.Redirect(p[0], []*game.Card{card(game.Six, game.Hearts)}), game.ErrorNotDefender)
	expectError(t, g.Redirect(p[1], []*game.Card{}), game.ErrorRedirectWithNoOrMixedCards)
	expectError(t, g.Redirect(p[1], p[1].GetCards()), game.ErrorRedirectWithNoOrMixedCards)
	expectError(t, g.Redirect(p[1], []*game.Card{card(game.Ace, game.Clubs)}), game.ErrorRedirectRankMismatch)

	mustSucceed(t, g.Defend(p[1], 0, card(game.Ace, game.Clubs)))
	expectError(t, g.Redirect(p[1], []*game.Card{card(game.Six, game.Hearts)}), game.ErrorAlreadyDefending)
}

func TestRedirectFailWhenDisabled(t *testing.T) {
	g, p := newRedirectGame()
	g.SetOption("with_redirect", game.Option{Value: "0"})
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))

	expectError(t, g.Redirect(p[1], []*game.Card{card(game.Six, game.Hearts)}), game.ErrorNoRedirectsAllowed)
}

func TestRedirectFailWhenNextDefenderHasTooFewCards(t *testing.T) {
	g, p := newTestGame(withRedirect,
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Six, game.Spades), card(game.Ace, game.Clubs)},
		[]*game.Card{card(game.Ten, game.Spades)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))

	expectError(t, g.Redirect(p[1], []*game.Card{card(game.Six, game.Spades)}), game.ErrorAttackIsTooBig)
}

func TestAbandonByOtherPlayer(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Clubs)},
		[]*game.Card{card(game.Six, game.Diamonds)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))

	mustSucceed(t, g.Abandon(p[2]))
	expectError(t, g.Abandon(p[2]), game.ErrorPlayerAlreadyQuit)

	if !p[2].HasQuit() || !p[2].HasAbandoned() {
		t.Error("Expected player to be marked as abandoned")
	}
	if !p[0].IsAttacker() || !p[1].IsDefender() || len(g.GetTable().GetCardsToBeat()) != 1 {
		t.Error("Expected turn to continue when a player not involved in it abandons")
	}
}

func TestGameOver(t *testing.T) {
	players := []*game.Player{
		game.NewPlayer("1", false, false, false, "Attacker", []*game.Card{card(game.Six, game.Clubs)}, false, true),
		game.NewPlayer("2", false, false, false, "Defender", []*game.Card{card(game.Seven, game.Clubs), card(game.Eight, game.Clubs)}, true, false),
	}
	deck := game.NewDeck([]*game.Card{}, card(game.King, game.Hearts))
	deck.GetCard()
	g := game.NewGame(deck, players, map[string]game.Option{}, true, false, game.Table{}, nil)

	mustSucceed(t, g.Attack(players[0], players[0].GetCards()))
	mustSucceed(t, g.Defend(players[1], 0, card(game.Seven, game.Clubs)))
	mustSucceed(t, g.EndAttack(players[0]))

	if !g.IsOver() {
		t.Fatal("Expected game to be over")
	}
	if !players[0].HasQuit() || players[1].HasQuit() {
		t.Error("Expected only player with empty hand to quit")
	}
	if lastEventType(g) != game.GameOverEventType {
		t.Error("Expected game over event to be logged")
	}
}
