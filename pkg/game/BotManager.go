package game

// MaxBotCycleRounds limits number of rounds in one BotManager.Cycle call.
//
// Games played by bots only can repeat the same position forever
const MaxBotCycleRounds = 1000

// BotManager handles bot interaction
type BotManager struct {
	bots []*Bot
}

func NewBotManager(bots []*Bot) *BotManager {
	botManager := &BotManager{}
	botManager.setBots(bots)
	return botManager
}

func (b *BotManager) setBots(bots []*Bot) {
	b.bots = bots
}

// Cycle processes actions of all bots until no more actions can be made
//
// Should be called on every players action including game start
func (b *BotManager) Cycle(g *Game) {
	if !g.IsStarted() || g.IsOver() {
		return
	}
	for round := 0; !g.IsOver() && round < MaxBotCycleRounds; round++ {
		acted := false

		for _, bot := range b.bots {
			if bot.Player.quitGame {
				continue
			}
			acted = bot.Act(g) || acted
		}

		if !acted {
			break
		}
	}
}
