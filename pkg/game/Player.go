package game

//Player state and details
type Player struct {
	ID            string  `json:"id"`
	skipTurn      bool    `json:"skip_turn"`
	quitGame      bool    `json:"quit_game"`
	wonGame       bool    `json:"won_game"`
	firstWinner   bool    `json:"first_winner"`
	abandonedGame bool    `json:"abandoned_game"`
	Name          string  `json:"name"`
	cards         []*Card `json:"Cards"`
	defender      bool    `json:"defender"`
	attacker      bool    `json:"attacker"`
}

// NewPlayer Creates new instance of in-game player
func NewPlayer(ID string, skipTurn bool, quitGame bool, abandonedGame bool, Name string, cards []*Card, defender bool, attacker bool) *Player {
	return &Player{ID: ID, skipTurn: skipTurn, quitGame: quitGame, abandonedGame: abandonedGame, Name: Name, cards: cards, defender: defender, attacker: attacker}
}

// addCards adds card to players hand
func (p *Player) addCards(cards []*Card) {
	for i := range cards {
		p.cards = append(p.cards, cards[i])
	}
}

// removeCard removes card from players hand
func (p *Player) removeCard(c *Card) {

	for i, card := range p.cards {
		if card.Rank == c.Rank && card.Suit == c.Suit {
			p.cards = append(p.cards[:i], p.cards[i+1:]...)
			return
		}
	}
}

// hasCard checks if card is in players card
func (p *Player) hasCard(c *Card) bool {
	for i := range p.cards {
		if p.cards[i].Rank == c.Rank && p.cards[i].Suit == c.Suit {
			return true
		}
	}
	return false
}

// GetCards returns pointer to player's cards
func (p *Player) GetCards() []*Card {
	return p.cards
}

// IsDefender returns true when player is defending
func (p *Player) IsDefender() bool {
	return p.defender
}

// IsAttacker returns true when player is main attacker (initiator of turn of last redirector)
func (p *Player) IsAttacker() bool {
	return p.attacker
}

// HasAbandoned is true when user quit game before game over
func (p *Player) HasAbandoned() bool {
	return p.abandonedGame
}

// HasQuit is true when user abandoned game or has no more cards
func (p *Player) HasQuit() bool {
	return p.quitGame
}

// IsSkipTurn is true when user picked up cards and miss turn
func (p *Player) IsSkipTurn() bool {
	return p.skipTurn
}
