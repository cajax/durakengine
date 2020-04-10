package game

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
	for {
		if g.IsOver() {
			break
		}
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
