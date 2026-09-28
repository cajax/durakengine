package game

import (
	"maps"
	"slices"
)

const (
	ErrorTooFewPlayers                   = "Too few players in game"
	ErrorTooManyPlayers                  = "Too many players in game"
	ErrorInvalidMinRank                  = "Invalid least card rank"
	ErrorDeckTooSmall                    = "Deck is too small to deal cards to all players"
	ErrorNotInPlayingState               = "Game is not in playing state"
	ErrorFirstAttackWithDifferentRanks   = "First attack in turn with different ranks"
	ErrorFirstAttackByNeighbor           = "First attack in turn can not be made by neighbor"
	ErrorAttackByWrongPlayer             = "Attack by someone who is not attacker"
	ErrorAttackRankNotOnTable            = "Attack by card that does not match any rank on table"
	ErrorAttackIsTooBig                  = "Attacking with more cards than defender can beat"
	ErrorAttackerHasNoCard               = "Attacker has no this card"
	ErrorNotDefender                     = "Player is not defender"
	ErrorDefenderHasNoCard               = "Defender has no this card"
	ErrorPairIndexOutOfRange             = "Invalid pair index"
	ErrorAlreadyDefended                 = "This card is already beaten"
	ErrorCardCantBeat                    = "The card is to weak to beat the other"
	ErrorDefenseByNotDefender            = "Defense by someone who is not defender"
	ErrorGameEndTurnByNotAttacker        = "Only attacker can end turn"
	ErrorGameEndTurnWithoutAttack        = "Can not end turn without at least one attack"
	ErrorGameEndTurnWithoutUnbeatenCards = "Can not end turn without unbeaten cards on table"
	ErrorNoRedirectsAllowed              = "Redirects are not allowed"
	ErrorRedirectWithNoOrMixedCards      = "Redirect with no cards or mixed ranks"
	ErrorAlreadyDefending                = "Some cards are already beaten"
	ErrorRedirectRankMismatch            = "Redirect card does not match rank on table"
	ErrorPlayerAlreadyQuit               = "Player already quit"
	ErrorNoPlayerToRedirect              = "No player to redirect to"
	ErrorAttackLimitReached              = "Attacking with more cards than allowed per turn"
	ErrorRedirectCardAlreadyShown        = "Card was already shown to redirect in this turn"
	ErrorPickupFromEmptyTable            = "Can not pick up from an empty table"
)

func (g *Game) StartGame() error {
	if g.IsStarted() || g.over {
		return ErrGameAlreadyStarted
	}

	if len(g.players) < 2 {
		return ErrTooFewPlayers
	}

	if len(g.players) > 6 {
		return ErrTooManyPlayers
	}

	minRank := Six
	if value := g.GetOption(OptionMinRank).Value; value != "" {
		minRank = RankFromString(value)
		if minRank < Two || minRank > Ace {
			return ErrInvalidMinRank
		}
	}

	if deckSize(minRank) < len(g.players)*handSize {
		return ErrDeckTooSmall
	}

	g.initTable(minRank)
	return nil
}

func (g *Game) initTable(minRank Rank) {
	g.advanceSequence()

	g.table = Table{}
	g.table.clear()
	g.deck = Deck{rng: g.rng}

	g.deck.ResetDeck(minRank)
	// the deal can take the trump card too when players share the whole deck
	trump := *g.deck.trump
	g.SpreadCards()
	g.SelectFirstAttacker()
	g.selectDefender()

	g.started = true

	_, attacker := g.getAttacker()
	_, defender := g.GetDefender()
	g.Log.Add(NewStartEvent(g.deck.GetCount(), minRank, attacker.snapshot(), defender.snapshot(), *g.deck.trumpSuit, trump, maps.Clone(g.options)))
}

// Attack player if is one of attacker
func (g *Game) Attack(p *Player, cards []*Card) error {
	if !g.inProgress() {
		return ErrNotInPlayingState
	}
	isNeighbor := g.IsOneOfAttackers(p)

	if p.IsAttacker() {
		if g.table.IsEmpty() && !g.CardsOfSameRank(cards) {
			return ErrFirstAttackWithDifferentRanks
		}
	} else if isNeighbor {
		if g.table.IsEmpty() {
			return ErrFirstAttackByNeighbor
		}
	} else {
		return ErrAttackByWrongPlayer
	}

	if !g.table.IsEmpty() {
		for _, c := range cards {
			if !g.table.cardMatchesSomeRanks(c) {
				return ErrAttackRankNotOnTable
			}
		}
	}

	_, defender := g.GetDefender()

	if len(defender.cards) < len(cards)+len(g.table.GetCardsToBeat()) {
		return ErrAttackIsTooBig
	}

	if !g.withinAttackLimit(len(cards)) {
		return ErrAttackLimitReached
	}

	if !p.hasCards(cards) {
		return ErrAttackerHasNoCard
	}

	// From here on we no longer expect errors
	g.advanceSequence()

	for _, c := range cards {
		g.table.attack(p, c)
		p.removeCard(c)
	}
	g.Log.Add(NewAttackEvent(p.snapshot(), slices.Clone(cards)))

	return nil
}

// Defend against cards on table if defender
func (g *Game) Defend(p *Player, i int, c *Card) error {
	if !g.inProgress() {
		return ErrNotInPlayingState
	}

	if !p.IsDefender() {
		return ErrNotDefender
	}

	if !p.hasCard(c) {
		return ErrDefenderHasNoCard
	}

	if !g.table.HasPair(i) {
		return ErrPairIndexOutOfRange
	}
	tp := g.table.GetPair(i)

	if tp.Defender != nil {
		return ErrAlreadyDefended
	}

	if !g.CardCanBeatOther(c, tp.Attack) {
		return ErrCardCantBeat
	}

	// From here on we no longer expect errors
	g.advanceSequence()
	g.table.defend(p, tp, c)
	p.removeCard(c)
	g.Log.Add(NewDefenseEvent(p.snapshot(), *tp, i, *c))
	return nil
}

// Pickup collects all cards from table if defender
func (g *Game) Pickup(p *Player) error {
	if !g.inProgress() {
		return ErrNotInPlayingState
	}

	if !p.IsDefender() {
		return ErrDefenseByNotDefender
	}

	if g.table.IsEmpty() {
		return ErrPickupFromEmptyTable
	}

	//From here on we no longer expect errors
	g.advanceSequence()

	p.addCards(g.table.GetCardsOnTable())
	g.Log.Add(NewPickupEvent(p.snapshot(), g.table.GetCardsOnTable()))
	g.table.clear()
	// defender loses the turn
	g.endTurn(g.GetPlayerIndex(p)+1, p, false)
	return nil
}

// EndAttack ends turn by attacker
func (g *Game) EndAttack(p *Player) error {
	if !g.inProgress() {
		return ErrNotInPlayingState
	}

	if !p.IsAttacker() {
		return ErrGameEndTurnByNotAttacker
	}
	if g.table.IsEmpty() {
		return ErrGameEndTurnWithoutAttack
	}

	if g.table.HasUnbeatenCards() {
		return ErrGameEndTurnWithoutUnbeatenCards
	}

	// From here on we no longer expect errors
	g.advanceSequence()
	g.Log.Add(NewDiscardEvent(p.snapshot(), g.table.GetCardsOnTable()))

	g.table.clear()

	defenderIndex, _ := g.GetDefender()
	g.endTurn(defenderIndex, nil, false)

	return nil
}

// Redirect to the left with laying on table card(s) of the same rank
func (g *Game) Redirect(defender *Player, cards []*Card) error {
	if !g.inProgress() {
		return ErrNotInPlayingState
	}

	if g.GetOption(OptionRedirect).Value != "1" {
		return ErrNoRedirectsAllowed
	}

	if !defender.IsDefender() {
		return ErrNotDefender
	}

	if len(cards) < 1 || !g.CardsOfSameRank(cards) {
		return ErrRedirectWithNoOrMixedCards
	}

	if !defender.hasCards(cards) {
		return ErrDefenderHasNoCard
	}

	if g.table.defenseStarted() {
		return ErrAlreadyDefending
	}

	if !g.table.cardMatchesAttackRank(cards[0]) {
		return ErrRedirectRankMismatch
	}

	defenderIndex := g.GetPlayerIndex(defender)
	_, nextDefender := g.getActivePlayerToTheLeft(defenderIndex)

	if nextDefender == nil || nextDefender == defender {
		return ErrNoPlayerToRedirect
	}

	// shown cards stay in defender's hand
	keepCards := g.GetOption(OptionRedirectKeepCard).Value == "1"
	addedCards := len(cards)
	if keepCards {
		addedCards = 0
		for _, card := range cards {
			if g.table.wasShown(card) {
				return ErrRedirectCardAlreadyShown
			}
		}
	}

	if len(nextDefender.cards) < len(g.table.GetCardsOnTable())+addedCards {
		return ErrAttackIsTooBig
	}

	if !g.withinAttackLimit(addedCards) {
		return ErrAttackLimitReached
	}

	//From here on we no longer expect errors
	g.advanceSequence()

	for _, card := range cards {
		if keepCards {
			g.table.shown = append(g.table.shown, *card)
			continue
		}
		g.table.attack(defender, card)
		defender.removeCard(card)
	}

	g.setAttacker(defender)
	g.setDefender(nextDefender)
	g.Log.Add(NewTransferEvent(defender.snapshot(), slices.Clone(cards), nextDefender.snapshot(), keepCards))

	return nil
}

// Abandon the game and put crds from hand into deck
func (g *Game) Abandon(player *Player) error {
	if !g.inProgress() {
		return ErrNotInPlayingState
	}

	if player.quitGame {
		return ErrPlayerAlreadyQuit
	}

	g.advanceSequence()

	// abandoners finish last, the earliest one at the very end
	place := len(g.players)
	for _, p := range g.players {
		if p.abandonedGame {
			place--
		}
	}

	cards := player.cards
	player.quitGame = true
	player.abandonedGame = true
	player.place = place
	for _, card := range cards {
		g.deck.AddCard(card)
	}
	player.cards = nil
	g.Log.Add(NewAbandonEvent(player.snapshot(), cards))
	if player.IsAttacker() || player.IsDefender() {
		g.cancelTurn()
		g.endTurn(g.GetPlayerIndex(player)+1, nil, true)
	}

	return nil
}
