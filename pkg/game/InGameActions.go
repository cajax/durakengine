package game

import (
	"errors"
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
)

func (g *Game) StartGame() error {
	if len(g.players) < 2 {
		return errors.New(ErrorTooFewPlayers)
	}

	if len(g.players) > 6 {
		return errors.New(ErrorTooManyPlayers)
	}

	minRank := Six
	if value := g.GetOption(OptionMinRank).Value; value != "" {
		minRank = RankFromString(value)
		if minRank < Two || minRank > Ace {
			return errors.New(ErrorInvalidMinRank)
		}
	}

	if deckSize(minRank) < len(g.players)*handSize {
		return errors.New(ErrorDeckTooSmall)
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
	g.SpreadCards()
	g.SelectFirstAttacker()
	g.selectDefender()

	g.started = true

	_, attacker := g.getAttacker()
	_, defender := g.GetDefender()
	g.Log.Add(NewStartEvent(g.deck.GetCount(), minRank, attacker.snapshot(), defender.snapshot(), *g.deck.trumpSuit, maps.Clone(g.options)))
}

// Attack player if is one of attacker
func (g *Game) Attack(p *Player, cards []*Card) error {
	if !g.inProgress() {
		return errors.New(ErrorNotInPlayingState)
	}
	isNeighbor := g.IsOneOfAttackers(p)

	if p.IsAttacker() {
		if g.table.IsEmpty() && !g.CardsOfSameRank(cards) {
			return errors.New(ErrorFirstAttackWithDifferentRanks)
		}
	} else if isNeighbor {
		if g.table.IsEmpty() {
			return errors.New(ErrorFirstAttackByNeighbor)
		}
	} else {
		return errors.New(ErrorAttackByWrongPlayer)
	}

	if !g.table.IsEmpty() {
		for _, c := range cards {
			if !g.table.cardMatchesSomeRanks(c) {
				return errors.New(ErrorAttackRankNotOnTable)
			}
		}
	}

	_, defender := g.GetDefender()

	if len(defender.cards) < len(cards)+len(g.table.GetCardsToBeat()) {
		return errors.New(ErrorAttackIsTooBig)
	}

	if !g.withinAttackLimit(len(cards)) {
		return errors.New(ErrorAttackLimitReached)
	}

	if !p.hasCards(cards) {
		return errors.New(ErrorAttackerHasNoCard)
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
		return errors.New(ErrorNotInPlayingState)
	}

	if !p.IsDefender() {
		return errors.New(ErrorNotDefender)
	}

	if !p.hasCard(c) {
		return errors.New(ErrorDefenderHasNoCard)
	}

	if !g.table.HasPair(i) {
		return errors.New(ErrorPairIndexOutOfRange)
	}
	tp := g.table.GetPair(i)

	if tp.Defender != nil {
		return errors.New(ErrorAlreadyDefended)
	}

	if !g.CardCanBeatOther(c, tp.Attack) {
		return errors.New(ErrorCardCantBeat)
	}

	// From here on we no longer expect errors
	g.advanceSequence()
	g.table.defend(p, tp, c)
	p.removeCard(c)
	g.Log.Add(NewDefenseEvent(p.snapshot(), *tp, *c))
	return nil
}

// Pickup collects all cards from table if defender
func (g *Game) Pickup(p *Player) error {
	if !g.inProgress() {
		return errors.New(ErrorNotInPlayingState)
	}

	if !p.IsDefender() {
		return errors.New(ErrorDefenseByNotDefender)
	}

	//From here on we no longer expect errors
	g.advanceSequence()

	p.addCards(g.table.GetCardsOnTable())
	g.Log.Add(NewPickupEvent(p.snapshot(), g.table.GetCardsOnTable()))
	g.table.clear()
	// defender loses the turn
	g.endTurn(g.GetPlayerIndex(p)+1, p)
	return nil
}

// EndAttack ends turn by attacker
func (g *Game) EndAttack(p *Player) error {
	if !g.inProgress() {
		return errors.New(ErrorNotInPlayingState)
	}

	if !p.IsAttacker() {
		return errors.New(ErrorGameEndTurnByNotAttacker)
	}
	if g.table.IsEmpty() {
		return errors.New(ErrorGameEndTurnWithoutAttack)
	}

	if g.table.HasUnbeatenCards() {
		return errors.New(ErrorGameEndTurnWithoutUnbeatenCards)
	}

	// From here on we no longer expect errors
	g.advanceSequence()
	g.Log.Add(NewDiscardEvent(p.snapshot(), g.table.GetCardsOnTable()))

	g.table.clear()

	defenderIndex, _ := g.GetDefender()
	g.endTurn(defenderIndex, nil)

	return nil
}

// Redirect to the left with laying on table card(s) of the same rank
func (g *Game) Redirect(defender *Player, cards []*Card) error {
	if !g.inProgress() {
		return errors.New(ErrorNotInPlayingState)
	}

	if g.GetOption(OptionRedirect).Value != "1" {
		return errors.New(ErrorNoRedirectsAllowed)
	}

	if !defender.IsDefender() {
		return errors.New(ErrorNotDefender)
	}

	if len(cards) < 1 || !g.CardsOfSameRank(cards) {
		return errors.New(ErrorRedirectWithNoOrMixedCards)
	}

	if !defender.hasCards(cards) {
		return errors.New(ErrorDefenderHasNoCard)
	}

	if g.table.defenseStarted() {
		return errors.New(ErrorAlreadyDefending)
	}

	if !g.table.cardMatchesAttackRank(cards[0]) {
		return errors.New(ErrorRedirectRankMismatch)
	}

	defenderIndex := g.GetPlayerIndex(defender)
	_, nextDefender := g.getActivePlayerToTheLeft(defenderIndex)

	if nextDefender == nil || nextDefender == defender {
		return errors.New(ErrorNoPlayerToRedirect)
	}

	// shown cards stay in defender's hand
	keepCards := g.GetOption(OptionRedirectKeepCard).Value == "1"
	addedCards := len(cards)
	if keepCards {
		addedCards = 0
		for _, card := range cards {
			if g.table.wasShown(card) {
				return errors.New(ErrorRedirectCardAlreadyShown)
			}
		}
	}

	if len(nextDefender.cards) < len(g.table.GetCardsOnTable())+addedCards {
		return errors.New(ErrorAttackIsTooBig)
	}

	if !g.withinAttackLimit(addedCards) {
		return errors.New(ErrorAttackLimitReached)
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
		return errors.New(ErrorNotInPlayingState)
	}

	if player.quitGame {
		return errors.New(ErrorPlayerAlreadyQuit)
	}

	g.advanceSequence()

	cards := player.cards
	player.quitGame = true
	player.abandonedGame = true
	for _, card := range cards {
		g.deck.AddCard(card)
	}
	player.cards = nil
	g.Log.Add(NewAbandonEvent(player.snapshot(), cards))
	if player.IsAttacker() || player.IsDefender() {
		g.cancelTurn()
		g.endTurn(g.GetPlayerIndex(player)+1, nil)
	}

	return nil
}
