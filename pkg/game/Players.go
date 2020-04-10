package game

import (
	"errors"
)

// GetPlayers returns point to list of game players
func (g *Game) GetPlayers() []*Player {
	return g.players
}

// AddPlayer adds new player to game
func (g *Game) AddPlayer(player *Player) ([]*Player, error) {
	if g.IsStarted() || g.over {
		return nil, errors.New("cannot_add_player_while_in_game")
	}
	g.players = append(g.players, player)

	return g.players, nil
}

// SetBots sets bots to game
func (g *Game) SetBots(bots []*Bot) error {
	if g.IsStarted() {
		return errors.New(ErrorGameAlreadyStarted)
	}
	g.botManager = &BotManager{}
	g.botManager.setBots(bots)
	return nil
}

// CycleBots Perform bot actions
func (g *Game) CycleBots() {
	if g.botManager != nil {
		g.botManager.Cycle(g)
	}
}
