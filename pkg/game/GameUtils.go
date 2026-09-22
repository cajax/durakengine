package game

import (
	"strconv"
)

/**************************************************

	Game utility methods that not called directly

***************************************************/

// GetDefender returns player that must cover cards on table
func (g *Game) GetDefender() (int, *Player) {
	for i, player := range g.players {
		if player.defender {
			return i, player
		}
	}

	return 0, nil
}

// getAttacker returns player that is currently main attacker (to the right of defender)
func (g *Game) getAttacker() (int, *Player) {
	for i, player := range g.players {
		if player.attacker {
			return i, player
		}
	}

	return 0, nil
}

// GetNeighborAttackers returns up to 2 neighbors of defender. One each side
//
// Empty unless throw-ins are allowed by OptionThrowIn
func (g *Game) GetNeighborAttackers() []*Player {
	var neighbors []*Player
	if g.GetOption(OptionThrowIn).Value != "1" {
		return neighbors
	}

	defenderIndex, _ := g.GetDefender()
	attackerIndex, _ := g.getAttacker()

	rightIndex, rightPlayer := g.getActivePlayerToTheRight(attackerIndex)
	if rightPlayer != nil && rightIndex != defenderIndex {
		neighbors = append(neighbors, rightPlayer)
	}

	leftIndex, leftPlayer := g.getActivePlayerToTheLeft(defenderIndex)
	if leftPlayer != nil && leftIndex != attackerIndex {
		neighbors = append(neighbors, leftPlayer)
	}

	return neighbors
}

// IsOneOfAttackers is true when player is main attacker or if player is neighbor of defender
func (g *Game) IsOneOfAttackers(p *Player) bool {
	if p.attacker {
		return true
	}

	for _, np := range g.GetNeighborAttackers() {
		if np == p {
			return true
		}
	}

	return false
}

// RefillUsers fills hands of attacker, other players clockwise and then defender
func (g *Game) RefillUsers() {
	attackerIndex, _ := g.getAttacker()
	_, defender := g.GetDefender()
	for i := range g.players {
		player := g.players[(attackerIndex+i)%len(g.players)]
		if player != defender {
			g.refillUserCards(player)
		}
	}
	g.refillUserCards(defender)
}

// refillUserCards adds cards from deck to players hand until it's full
func (g *Game) refillUserCards(player *Player) {
	if player == nil || player.quitGame {
		return
	}
	var addedCards []Card
	for i := len(player.cards); i < handSize; i++ {
		card := g.deck.GetCard()
		if card == nil {
			break
		}
		player.cards = append(player.cards, card)
		addedCards = append(addedCards, *card)
	}
	if len(addedCards) > 0 {
		g.Log.Add(NewRefillEvent(player.snapshot(), addedCards))
	}
}

// SpreadCards fills each players hands with cards
func (g *Game) SpreadCards() {
	for i := range g.players {
		g.refillUserCards(g.players[i])
	}
}

// SetOption changes game options
func (g *Game) SetOption(id string, option Option) {
	if g.options == nil {
		g.options = map[string]Option{}
	}
	g.options[id] = option
}

// GetOption returns game option
func (g *Game) GetOption(id string) Option {
	return g.options[id]
}

// maxAttackCards returns number of attack cards allowed per turn, 0 if unlimited
func (g *Game) maxAttackCards() int {
	limit, err := strconv.Atoi(g.GetOption(OptionMaxAttackCards).Value)
	if err != nil || limit < 0 {
		return 0
	}
	return limit
}

// withinAttackLimit checks if n more attack cards can be put on table in current turn
func (g *Game) withinAttackLimit(n int) bool {
	limit := g.maxAttackCards()
	return limit == 0 || len(g.table.pairs)+n <= limit
}

// attackCapacity returns how many more cards can be added to attack in current turn
func (g *Game) attackCapacity() int {
	_, defender := g.GetDefender()
	capacity := len(defender.cards) - len(g.table.GetCardsToBeat())
	if limit := g.maxAttackCards(); limit > 0 {
		capacity = min(capacity, limit-len(g.table.pairs))
	}
	return max(capacity, 0)
}

// SelectFirstAttacker find Player with the least trump, or random Player if none has trumps
func (g *Game) SelectFirstAttacker() {
	leastRank := Ace + 1
	var candidate *Player = nil
	for i := range g.players {
		for _, card := range g.players[i].cards {
			if g.IsTrump(card) && card.Rank < leastRank {
				candidate = g.players[i]
				leastRank = card.Rank
			}
		}
	}

	if candidate == nil {
		candidate = g.players[intN(g.rng, len(g.players))]
	}

	g.setAttacker(candidate)
}

// selectDefender chooses user to the left of Attacker
func (g *Game) selectDefender() {
	index, _ := g.getAttacker()

	_, defender := g.getActivePlayerToTheLeft(index)
	g.setDefender(defender)
}

func (g *Game) getActivePlayerToTheLeft(index int) (int, *Player) {
	for i := 0; i < len(g.players); i++ {
		index++
		if index == len(g.players) {
			index = 0
		}

		//Player exists but not playing any more
		if !g.players[index].quitGame && len(g.players[index].cards) > 0 {
			return index, g.players[index]
		}
	}

	return 0, nil
}

func (g *Game) getActivePlayerToTheRight(index int) (int, *Player) {
	for i := 0; i < len(g.players); i++ {
		index--
		if index == -1 {
			index = len(g.players) - 1
		}

		//Player exists but not playing any more
		if !g.players[index].quitGame && len(g.players[index].cards) > 0 {
			return index, g.players[index]
		}
	}

	return 0, nil
}

func (g *Game) setAttacker(p *Player) {
	for _, player := range g.players {
		if player == p {
			player.attacker = true
		} else {
			player.attacker = false
		}
	}
}

func (g *Game) setDefender(p *Player) {
	for _, player := range g.players {
		if player == p {
			player.defender = true
		} else {
			player.defender = false
		}
	}
}

// cancelTurn clears the table, returning each card to the player who played it
//
// Cards of players who quit the game go back to the deck
func (g *Game) cancelTurn() {
	for _, pair := range g.table.pairs {
		g.returnCard(pair.Attacker, pair.Attack)
		if pair.Defense != nil {
			g.returnCard(pair.Defender, pair.Defense)
		}
	}
	g.table.clear()
}

func (g *Game) returnCard(p *Player, c *Card) {
	if p == nil || p.quitGame {
		g.deck.AddCard(c)
		return
	}
	p.addCards([]*Card{c})
}

// detectWinners finds players that have no Cards but still in game
func (g *Game) detectWinners() {
	for _, player := range g.players {
		if player.quitGame {
			continue
		}
		if len(player.cards) == 0 {
			g.winPlayer(player)
		}
	}
}

// GetPlayerIndex returns index of Player in players array
func (g *Game) GetPlayerIndex(p *Player) int {
	for i := range g.players {
		if g.players[i] == p {
			return i
		}
	}
	return -1
}

// CardCanBeatOther checks if first Card is in higher rank or is trump
func (g *Game) CardCanBeatOther(card *Card, other *Card) bool {
	if g.IsTrump(card) && !g.IsTrump(other) {
		return true
	}

	if card.Suit == other.Suit && card.Rank > other.Rank {
		return true
	}

	return false
}

// CardsOfSameRank verifies if all Cards in slice are of the same rank.
//
// Attack can be performed only with same rank Cards
func (g *Game) CardsOfSameRank(cards []*Card) bool {
	if len(cards) == 0 {
		return false
	}

	if len(cards) == 1 {
		return true
	}

	for _, c := range cards {
		if c.Rank != cards[0].Rank {
			return false
		}
	}

	return true
}

func (g *Game) IsTrump(c *Card) bool {
	return c.Suit == *g.deck.trumpSuit
}
