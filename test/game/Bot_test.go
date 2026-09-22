package game

import (
	"github.com/cajax/durakengine/pkg/game"
	"testing"
)

func TestSelectLeastCardThatBeatOtherWithTrump(t *testing.T) {
	d := game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Spades, Rank: game.Seven})
	g := game.NewGame(d, []*game.Player{}, map[string]game.Option{}, true, false, game.Table{}, &game.BotManager{})
	b := game.Bot{}

	c := []*game.Card{}
	c = append(c, &game.Card{Suit: game.Clubs, Rank: game.Ten})
	c = append(c, &game.Card{Suit: game.Clubs, Rank: game.Seven})
	c = append(c, &game.Card{Suit: game.Spades, Rank: game.Six})
	c = append(c, &game.Card{Suit: game.Diamonds, Rank: game.Ace})

	o := &game.Card{Suit: game.Clubs, Rank: game.King}
	card := b.SelectLeastCardThatBeatOther(g, c, o)
	if card != c[2] {
		t.Error("Trump card supposed to beat the regular card")
	}
}

func TestSelectLeastCardThatBeatOtherWithoutTrump(t *testing.T) {
	d := game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Spades, Rank: game.Seven})
	g := game.NewGame(d, []*game.Player{}, map[string]game.Option{}, true, false, game.Table{}, &game.BotManager{})
	b := game.Bot{}

	c := []*game.Card{}
	c = append(c, &game.Card{Suit: game.Spades, Rank: game.Six})
	c = append(c, &game.Card{Suit: game.Clubs, Rank: game.Ten})
	c = append(c, &game.Card{Suit: game.Clubs, Rank: game.Seven})
	c = append(c, &game.Card{Suit: game.Clubs, Rank: game.Ace})
	c = append(c, &game.Card{Suit: game.Clubs, Rank: game.Ten})
	c = append(c, &game.Card{Suit: game.Spades, Rank: game.Ace})

	o := &game.Card{Suit: game.Clubs, Rank: game.King}
	card := b.SelectLeastCardThatBeatOther(g, c, o)
	if card != c[3] {
		t.Error("Regular card of the same suit supposed to beat the card")
	}
}

func TestSelectLeastCardThatBeatOtherFail(t *testing.T) {
	d := game.NewDeck([]*game.Card{}, &game.Card{Suit: game.Spades, Rank: game.Seven})
	g := game.NewGame(d, []*game.Player{}, map[string]game.Option{}, true, false, game.Table{}, &game.BotManager{})
	b := game.Bot{}

	c := []*game.Card{}
	c = append(c, &game.Card{Suit: game.Clubs, Rank: game.Ten})
	c = append(c, &game.Card{Suit: game.Clubs, Rank: game.Seven})
	c = append(c, &game.Card{Suit: game.Clubs, Rank: game.Ten})
	c = append(c, &game.Card{Suit: game.Hearts, Rank: game.Ace})

	o := &game.Card{Suit: game.Clubs, Rank: game.King}
	card := b.SelectLeastCardThatBeatOther(g, c, o)
	if card != nil {
		t.Error("There is no card that can beat the other")
	}
}

func TestBotAttacksWithLeastNonTrump(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Hearts), card(game.Nine, game.Clubs), card(game.Seven, game.Spades)},
		[]*game.Card{card(game.Eight, game.Clubs), card(game.Nine, game.Spades)},
	)
	b := game.Bot{Player: p[0]}

	if !b.Act(g) {
		t.Fatal("Expected bot to attack")
	}

	cards := g.GetTable().GetCardsToBeat()
	if len(cards) != 1 || *cards[0] != *card(game.Seven, game.Spades) {
		t.Errorf("Expected bot to attack with 7♠, got %v", cards)
	}
}

func TestBotDefends(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Six, game.Hearts), card(game.Ace, game.Clubs), card(game.Nine, game.Clubs)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	b := game.Bot{Player: p[1]}

	if !b.Act(g) {
		t.Fatal("Expected bot to defend")
	}

	pairs := g.GetTable().GetPairs()
	if len(pairs) != 1 || pairs[0].Defense == nil || *pairs[0].Defense != *card(game.Nine, game.Clubs) {
		t.Error("Expected bot to defend with 9♣")
	}
}

func TestBotPicksUpWhenCannotDefend(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Six, game.Clubs), card(game.Ace, game.Spades)},
		[]*game.Card{card(game.Ten, game.Diamonds)},
	)
	mustSucceed(t, g.Attack(p[0], p[0].GetCards()))
	b := game.Bot{Player: p[1]}

	if !b.Act(g) {
		t.Fatal("Expected bot to pick up")
	}
	if !g.GetTable().IsEmpty() || !hasCard(p[1].GetCards(), card(game.Seven, game.Clubs)) {
		t.Error("Expected bot to pick up cards from table")
	}
}

func TestBotEndsAttackWhenAllCardsBeaten(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Seven, game.Clubs), card(game.Ace, game.Spades)},
		[]*game.Card{card(game.Nine, game.Clubs), card(game.Six, game.Diamonds)},
		[]*game.Card{card(game.Ten, game.Diamonds)},
	)
	mustSucceed(t, g.Attack(p[0], []*game.Card{card(game.Seven, game.Clubs)}))
	mustSucceed(t, g.Defend(p[1], 0, card(game.Nine, game.Clubs)))
	b := game.Bot{Player: p[0]}

	if !b.Act(g) {
		t.Fatal("Expected bot to end attack")
	}
	if !g.GetTable().IsEmpty() || !p[1].IsAttacker() {
		t.Error("Expected turn to pass to defender")
	}
}

func TestBotsPlayFullGames(t *testing.T) {
	for n := 2; n <= 6; n++ {
		for _, redirect := range []string{"0", "1"} {
			finished := 0
			games := 100
			for i := 0; i < games; i++ {
				var players []*game.Player
				var bots []*game.Bot
				for j := 0; j < n; j++ {
					p := &game.Player{Name: "Bot"}
					players = append(players, p)
					bots = append(bots, &game.Bot{Player: p})
				}
				options := map[string]game.Option{"min_rank": {Value: "6"}, "with_redirect": {Value: redirect}}
				g := game.NewGame(&game.Deck{}, players, options, false, false, game.Table{}, game.NewBotManager(bots))
				mustSucceed(t, g.StartGame())

				g.CycleBots()

				if g.IsOver() {
					finished++
				}
			}
			// bots can rarely end up repeating the same position
			if finished < games*95/100 {
				t.Errorf("%d players, redirect %s: only %d of %d games finished", n, redirect, finished, games)
			}
		}
	}
}

func TestBotDefendsWithTrumpAce(t *testing.T) {
	g, p := newTestGame(nil,
		[]*game.Card{card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Ace, game.Hearts), card(game.Six, game.Clubs)},
	)
	b := game.Bot{}

	if c := b.SelectLeastCardThatBeatOther(g, p[1].GetCards(), card(game.Seven, game.Clubs)); c == nil || *c != *card(game.Ace, game.Hearts) {
		t.Errorf("Expected trump ace to beat the card, got %v", c)
	}
}
