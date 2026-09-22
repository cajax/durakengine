package game

import (
	"errors"
	"maps"
	"math/rand/v2"
)

// Game state
type Game struct {
	players    []*Player
	deck       Deck
	options    map[string]Option
	started    bool
	over       bool
	table      Table
	botManager *BotManager
	rng        *rand.Rand
	Log        Log
}

const ErrorGameAlreadyStarted = "Game is already started"

// NewGame makes new game state
func NewGame(deck *Deck, players []*Player, options map[string]Option, started bool, over bool, table Table, botManager *BotManager) *Game {
	return &Game{
		deck:       *deck,
		players:    players,
		options:    options,
		started:    started,
		over:       over,
		table:      table,
		botManager: botManager,
		Log:        NewLog(),
	}
}

// advanceSequence increment internal action counter
func (g *Game) advanceSequence() {
	g.Log.Advance()
}

// GetSequence returns sequence number for render
func (g *Game) GetSequence() int {
	return len(g.Log.Events)
}

// GAME ACTIONS

// SetRandom sets source of randomness for shuffling and choosing first attacker.
// Global source is used if nil. Use seeded source to make games reproducible
func (g *Game) SetRandom(r *rand.Rand) error {
	if g.IsStarted() {
		return errors.New(ErrorGameAlreadyStarted)
	}
	g.rng = r
	return nil
}

// SetPlayers sets list of players to game
func (g *Game) SetPlayers(players []*Player) error {
	if g.IsStarted() {
		return errors.New(ErrorGameAlreadyStarted)
	}

	g.players = players
	return nil
}

// endGame marks game as over and the last player in game, if any, as loser
func (g *Game) endGame(loser *Player) {
	g.over = true
	if loser != nil {
		loser.lostGame = true
	}
}

func (g *Game) IsOver() bool {
	return g.over
}

// inProgress is true when game is started and not over yet
func (g *Game) inProgress() bool {
	return g.started && !g.over
}

// winPlayer marks player as winner
//
// also marks player as first to win if applicable
func (g *Game) winPlayer(p *Player) {
	p.quitGame = true
	firstToWin := true
	for _, p := range g.players {
		if p.wonGame {
			firstToWin = false
			break
		}
	}
	p.firstWinner = firstToWin
	p.wonGame = true
}

// checkGameOver checks if only one active player left
func (g *Game) checkGameOver() {
	// count number of active players. if <2 game is over
	var activePlayers []*Player
	for _, player := range g.players {
		if !player.quitGame {
			activePlayers = append(activePlayers, player)
		}
	}
	if len(activePlayers) >= 2 {
		return
	}

	var loser *Player
	if len(activePlayers) == 1 {
		loser = activePlayers[0]
	}
	g.endGame(loser)

	lastPlayers := make([]Player, 0, len(activePlayers))
	for _, player := range activePlayers {
		lastPlayers = append(lastPlayers, player.snapshot())
	}
	g.Log.Add(NewGameOverEvent(lastPlayers))
}

// endTurn refills hands, detects winners and prepares next turn
//
// Next attacker is the first active player starting from the given seat index.
// Skipping player, if any, is marked as missing the next turn
func (g *Game) endTurn(nextAttackerIndex int, skipping *Player) {
	for _, player := range g.players {
		player.skipTurn = player == skipping
	}
	g.RefillUsers()
	g.detectWinners()
	g.checkGameOver()
	if g.over {
		return
	}

	// getActivePlayerToTheLeft starts looking at the next seat
	index, nextAttacker := g.getActivePlayerToTheLeft((nextAttackerIndex + len(g.players) - 1) % len(g.players))
	_, nextDefender := g.getActivePlayerToTheLeft(index)
	g.setAttacker(nextAttacker)
	g.setDefender(nextDefender)
	g.Log.Add(NewEndTurnEvent(nextAttacker.snapshot(), nextDefender.snapshot()))
}

// GetOptions returns copy of game options
func (g *Game) GetOptions() map[string]Option {
	return maps.Clone(g.options)
}

// GetTable returns pointer to game table
func (g *Game) GetTable() *Table {
	return &g.table
}

// IsStarted is true when game is started
func (g *Game) IsStarted() bool {
	return g.started
}

func (g *Game) GetDeck() Deck {
	return g.deck
}

func (g *Game) GetPairs() []TablePair {
	pairs := make([]TablePair, 0, len(g.table.pairs))
	for _, pair := range g.table.pairs {
		pairs = append(pairs, *pair)
	}
	return pairs
}
