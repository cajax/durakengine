package game

import (
	"encoding/json"
	"slices"
)

// Player state and details
type Player struct {
	ID            string
	skipTurn      bool
	quitGame      bool
	wonGame       bool
	firstWinner   bool
	lostGame      bool
	abandonedGame bool
	Name          string
	cards         []*Card
	defender      bool
	attacker      bool
}

type playerJSON struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Cards         []*Card `json:"cards"`
	Attacker      bool    `json:"attacker"`
	Defender      bool    `json:"defender"`
	SkipTurn      bool    `json:"skip_turn"`
	QuitGame      bool    `json:"quit_game"`
	WonGame       bool    `json:"won_game"`
	FirstWinner   bool    `json:"first_winner"`
	LostGame      bool    `json:"lost_game"`
	AbandonedGame bool    `json:"abandoned_game"`
}

// MarshalJSON serializes player including hand and state
func (p Player) MarshalJSON() ([]byte, error) {
	return json.Marshal(playerJSON{
		ID:            p.ID,
		Name:          p.Name,
		Cards:         p.cards,
		Attacker:      p.attacker,
		Defender:      p.defender,
		SkipTurn:      p.skipTurn,
		QuitGame:      p.quitGame,
		WonGame:       p.wonGame,
		FirstWinner:   p.firstWinner,
		LostGame:      p.lostGame,
		AbandonedGame: p.abandonedGame,
	})
}

// snapshot returns copy of player that is not affected by further changes
func (p *Player) snapshot() Player {
	s := *p
	s.cards = slices.Clone(p.cards)
	return s
}

// NewPlayer Creates new instance of in-game player
func NewPlayer(ID string, skipTurn bool, quitGame bool, abandonedGame bool, Name string, cards []*Card, defender bool, attacker bool) *Player {
	return &Player{ID: ID, skipTurn: skipTurn, quitGame: quitGame, abandonedGame: abandonedGame, Name: Name, cards: cards, defender: defender, attacker: attacker}
}

// addCards adds card to players hand
func (p *Player) addCards(cards []*Card) {
	p.cards = append(p.cards, cards...)
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

// hasCards checks if all cards are in players hand, each listed only once
func (p *Player) hasCards(cards []*Card) bool {
	for i, c := range cards {
		if !p.hasCard(c) {
			return false
		}
		for _, other := range cards[:i] {
			if other.Rank == c.Rank && other.Suit == c.Suit {
				return false
			}
		}
	}
	return true
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

// HasWon is true when player got rid of all cards
func (p *Player) HasWon() bool {
	return p.wonGame
}

// IsFirstWinner is true when player was the first to get rid of all cards
func (p *Player) IsFirstWinner() bool {
	return p.firstWinner
}

// IsLoser is true when player was the last one left in game
func (p *Player) IsLoser() bool {
	return p.lostGame
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
