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
