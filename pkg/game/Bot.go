package game

import (
	"sort"
)
// Simple bot
type Bot struct {
	Player           *Player
	waitForNextCycle bool
}

type CardChoice struct {
	PairIndex int
	Card *Card
}

func (b *Bot) newCardChoice(card *Card, i int) *CardChoice {
	return &CardChoice{Card: card, PairIndex: i}
}

func (b *Bot) Act(g *Game) bool {
	if b.Player.defender {
		return b.defend(g)
	}

	if g.IsOneOfAttackers(b.Player) {
		return b.attack(g)
	}
	return false
}

func (b *Bot) defend(g *Game) bool {
	cardsCount := len(g.table.GetCardsToBeat())

	if cardsCount == 0 {
		//nothing to defend
		return false
	}

	var cardChoices []*CardChoice

	tmpCards := append(b.Player.cards[:0:0], b.Player.cards...)

	for i, pair := range g.table.pairs {
		if pair.Defense != nil {
			continue
		}

		defenseCard := b.SelectLeastCardThatBeatOther(g, tmpCards, pair.Attack)
		if defenseCard == nil {
			break
		}

		cardChoices = append(cardChoices, b.newCardChoice(defenseCard, i))
		for i, card := range tmpCards {
			if card == defenseCard {
				tmpCards = append(tmpCards[:i], tmpCards[i+1:]...)
				break
			}
		}
	}

	if len(cardChoices) < cardsCount {
		g.Pickup(b.Player)
		return true
	}

	for _, choice := range cardChoices {
		g.Defend(b.Player, choice.PairIndex, choice.Card)
	}

	return true
}

func (b *Bot) attack(g *Game) bool {
	//simple strategy

	if b.Player.attacker && !g.table.HasUnbeatenCards() && !g.table.IsEmpty() {
		//Attacker has nothing to do
		g.EndAttack(b.Player)

		return true
	}

	if b.Player.attacker && g.table.IsEmpty() {
		cards := b.getLeastValuedCardsForAttack(g)
		if len(cards) == 0 {
			panic("WTF")
		}
		err := g.Attack(b.Player, cards)
		return err == nil
	}
	if g.IsOneOfAttackers(b.Player) && !g.table.IsEmpty() {
		cards := b.getLeastValuedCardsToAdd(g)
		if len(cards) > 0 {
			err := g.Attack(b.Player, cards)
			return err == nil
		}

		return false
	}

	g.table.HasUnbeatenCards()

	return false
}

func (b *Bot) getLeastValuedCardsForAttack(g *Game) []*Card {
	_, d := g.GetDefender()
	maxCards := len(d.cards)
	//regular Cards
	cards := b.selectLeastValuedRank(b.groupCardsByRank(b.selectNonTrumps(b.Player.cards, g.deck.trumpSuit)), maxCards)

	if cards != nil {
		return cards
	}

	cards = b.selectLeastValuedRank(b.groupCardsByRank(b.selectTrumps(b.Player.cards, g.deck.trumpSuit)), maxCards)
	if cards != nil {
		return cards
	}

	return []*Card{}
}

func (b *Bot) getLeastValuedCardsToAdd(g *Game) []*Card {
	_, d := g.GetDefender()
	maxCards := len(d.cards)
	return b.selectLeastValuedRank(b.groupCardsByRank(b.selectMatchingTable(b.selectNonTrumps(b.Player.cards, g.deck.trumpSuit), g.table.GetCardsOnTable())), maxCards)
}

func (b *Bot) selectTrumps(c []*Card, trumpSuit *Suit) []*Card {
	var cards []*Card
	for _, card := range c {
		if card.Suit == *trumpSuit {
			cards = append(cards, card)
		}
	}
	return cards
}

func (b *Bot) selectNonTrumps(c []*Card, trumpSuit *Suit) []*Card {
	var cards []*Card
	for _, card := range c {
		if card.Suit != *trumpSuit {
			cards = append(cards, card)
		}
	}
	return cards
}

func (b *Bot) groupCardsByRank(c []*Card) map[Rank][]*Card {
	ranks := map[Rank][]*Card{}
	for _, card := range c {
		ranks[card.Rank] = append(ranks[card.Rank], card)
	}

	return ranks
}

func (b *Bot) selectLeastValuedRank(ranks map[Rank][]*Card, maxCards int) []*Card {
	var keys []int
	for k := range ranks {
		keys = append(keys, int(k))
	}

	sort.Ints(keys)

	bestValue := 0
	bestRank := Two

	for i, k := range keys {
		// init least value
		if i == 0 {
			bestRank = Rank(k)
		}

		n := len(ranks[Rank(k)])

		if n > bestValue && n <= maxCards {
			bestValue = n
			bestRank = Rank(k)
		}

		if i > 2 {
			break
		}
	}

	if len(ranks[bestRank]) > maxCards {
		ranks[bestRank] = ranks[bestRank][0:maxCards]
	}
	return ranks[bestRank]
}

func (b *Bot) selectMatchingTable(c []*Card, otherCards []*Card) []*Card {
	var cards []*Card
	for _, card := range c {
		for _, otherCard := range otherCards {
			if card.Rank == otherCard.Rank {
				cards = append(cards, card)
				break
			}
		}
	}
	return cards
}

func (b *Bot) SelectLeastCardThatBeatOther(g *Game, c []*Card, other *Card) *Card {

	var leastCard *Card
	leastRank := Ace
	leastTrump := true

	for _, card := range c {
		if g.CardCanBeatOther(card, other) {
			if leastTrump && !g.IsTrump(card) {
				leastTrump = false
				leastRank = card.Rank
				leastCard = card
				continue
			}
			if !leastTrump && g.IsTrump(card) {
				continue
			}
			if leastRank > card.Rank {
				leastTrump = g.IsTrump(card)
				leastRank = card.Rank
				leastCard = card
			}
		}
	}

	return leastCard
}
