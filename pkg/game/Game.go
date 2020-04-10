package game

import (
	"errors"
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

//todo test win-lose case when attacker has some cards, but defender nothing more

// advanceSequence increment internal action counter
func (g *Game) advanceSequence() {
	g.Log.Advance()
}

// GetSequence returns sequence number for render
func (g *Game) GetSequence() int {
	return len(g.Log.Events)
}

// GAME ACTIONS

// SetPlayers sets list of players to game
func (g *Game) SetPlayers(players []*Player) error {
	if g.IsStarted() {
		return errors.New(ErrorGameAlreadyStarted)
	}

	g.players = players
	return nil
}

// TODO do rest of flag updates
func (g *Game) endGame() {
	g.over = true
}

func (g *Game) IsOver() bool {
	return g.over
}

// winPlayer marks player as winner
//
// also marks player as first to win if applicable
func (g *Game) winPlayer(p *Player) {
	p.quitGame = true
	// TODO refactor (first) winner detector. Instead of checking in the end of whole turn do a check after every valid action and count number of Cards in hand
	// Simulate refill to see if Player with empty hand now would have or not Cards from deck on refill. If no new Cards then user won.
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
	var activePlayers []Player
	for _, player := range g.players {
		if !player.quitGame {
			activePlayers = append(activePlayers, *player)
		}
	}
	if len(activePlayers) < 2 {
		g.endGame()
		g.Log.Add(NewGameOverEvent(activePlayers))
	}
}

// endTurn ends turn, does count, prepares next turn
func (g *Game) endTurn(nextAttacker *Player) {
	g.setAttacker(nextAttacker)
	_, nextDefender := g.getActivePlayerToTheLeft(g.GetPlayerIndex(nextAttacker))
	g.setDefender(nextDefender)
	g.Log.Add(NewEndTurnEvent(*nextAttacker, *nextDefender))
	g.RefillUsers()
	g.detectWinners()
	g.checkGameOver()
	//todo finalize game if game is over
}

// GetOptions returns game options set during creation
func (g *Game) GetOptions() {
	//todo
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
	pairs := make([]TablePair, len(g.table.pairs))
	for _, pair := range g.table.pairs {
		pairs = append(pairs, *pair)
	}
	return pairs
}
