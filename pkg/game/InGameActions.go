package game

import (
	"errors"
)

const (
	ErrorTooFewPlayers                   = "Too few players in game"
	ErrorTooManyPlayers                  = "Too many players in game"
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
)

func (g *Game) StartGame() error {
	if len(g.players) < 2 {
		return errors.New(ErrorTooFewPlayers)
	}

	if len(g.players) > 6 {
		return errors.New(ErrorTooManyPlayers)
	}
	// TODO move option name to constant
	minRank := RankFromString(g.GetOption("min_rank").Value)
	g.initTable(minRank)
	return nil
}

func (g *Game) initTable(minRank Rank) {
	g.advanceSequence()

	g.table = Table{}
	g.table.clear()
	g.deck = Deck{}

	g.deck.ResetDeck(minRank)
	g.SpreadCards()
	g.SelectFirstAttacker()
	g.selectDefender()

	g.started = true

	_, attacker := g.getAttacker()
	_, defender := g.GetDefender()
	g.Log.Add(NewStartEvent(g.deck.GetCount(), minRank, *attacker, *defender, *g.deck.trumpSuit, g.options))
}

// Attack player if is one of attacker
func (g *Game) Attack(p *Player, cards []*Card) error {
	if !g.IsStarted() || g.over {
		return errors.New(ErrorNotInPlayingState)
	}
	isNeighbor := g.IsOneOfAttackers(p)

	if p.IsAttacker() {
		if g.table.IsEmpty() && !g.CardsOfSameRank(cards) {
			return errors.New(ErrorFirstAttackWithDifferentRanks)
		}
	} else if isNeighbor { // TODO check if podkindnoy
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

	for _, c := range cards {
		if !p.hasCard(c) {
			return errors.New(ErrorAttackerHasNoCard)
		}
	}

	// From here on we no longer expect errors
	g.advanceSequence()
	g.Log.Add(NewAttackEvent(*p, cards))

	for _, c := range cards {
		g.table.attack(p, c)
		p.removeCard(c)
	}

	return nil
}

// Defend against cards on table if defender
// todo replace tp by index of tp
func (g *Game) Defend(p *Player, i int, c *Card) error {
	if !p.IsDefender() {
		return errors.New(ErrorNotDefender)
	}

	if !p.hasCard(c) {
		return errors.New(ErrorDefenderHasNoCard)
	}

	if !g.table.HasPair(i) {
		errors.New(ErrorPairIndexOutOfRange)
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
	g.Log.Add(NewDefenseEvent(*p, *tp, *c))
	g.table.defend(p, tp, c)
	p.removeCard(c)
	return nil
}

// Pickup collects all cards from table if defender
func (g *Game) Pickup(p *Player) error {
	if !p.IsDefender() {
		return errors.New(ErrorDefenseByNotDefender)
	}

	//From here on we no longer expect errors
	g.advanceSequence()
	g.Log.Add(NewPickupEvent(*p, g.table.GetCardsOnTable()))

	i := g.GetPlayerIndex(p)
	p.addCards(g.table.GetCardsOnTable())
	g.table.clear()
	_, attacker := g.getActivePlayerToTheLeft(i)
	g.endTurn(attacker)
	return nil
}

// EndAttack ends turn by attacker
func (g *Game) EndAttack(p *Player) error {

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

	g.table.clear()

	_, newAttacker := g.getActivePlayerToTheLeft(g.GetPlayerIndex(p))
	g.endTurn(newAttacker)

	return nil
}

//Redirect to the left with laying on table card(s) of the same rank
func (g *Game) Redirect(defender *Player, cards []*Card) error {
	// todo option name to constant
	if g.GetOption("with_redirect").Value != "1" {
		return errors.New(ErrorNoRedirectsAllowed)
	}

	if !defender.IsDefender() {
		return errors.New(ErrorNotDefender)
	}

	if len(cards) < 1 || !g.CardsOfSameRank(cards) {
		return errors.New(ErrorRedirectWithNoOrMixedCards)
	}

	if g.table.defenseStarted() {
		return errors.New(ErrorAlreadyDefending)
	}

	if !g.table.cardMatchesAttackRank(cards[0]) {
		return errors.New(ErrorRedirectRankMismatch)
	}

	defenderIndex := g.GetPlayerIndex(defender)
	_, nextDefender := g.getActivePlayerToTheLeft(defenderIndex)

	if len(nextDefender.cards) < len(g.table.GetCardsOnTable())+len(cards) {
		return errors.New(ErrorAttackIsTooBig)
	}

	//From here on we no longer expect errors
	g.advanceSequence()

	for _, card := range cards {
		g.table.attack(defender, card)
	}

	g.setAttacker(defender)
	g.setDefender(nextDefender)

	return nil
}

// Abandon the game and put crds from hand into deck
func (g *Game) Abandon(player *Player) error {
	if player.quitGame {
		return errors.New(ErrorPlayerAlreadyQuit)
	}

	g.advanceSequence()
	g.Log.Add(NewAbandonEvent(*player, player.cards))

	player.quitGame = true
	player.abandonedGame = true
	if len(player.cards) > 0 {
		for _, card := range player.cards {
			g.deck.AddCard(card)
		}
	}
	if player.IsAttacker() || player.IsDefender() {
		i := g.GetPlayerIndex(player)
		_, attacker := g.getActivePlayerToTheLeft(i)
		g.endTurn(attacker)
	}

	return nil
}
