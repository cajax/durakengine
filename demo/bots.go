package main

import (
	"encoding/json"
	"github.com/cajax/durakengine/pkg/game"
	"log"
)

func main() {
	// set list of players
	players := []*game.Player{
		&game.Player{Name: "Player 1"},
		&game.Player{Name: "Player 2"},
		&game.Player{Name: "Player 3"},
		&game.Player{Name: "Player 4"},
		&game.Player{Name: "Player 5"},
	}

	// make bot for every player
	bots := make([]*game.Bot, 0, len(players))
	for _, p := range players {
		bots = append(bots, &game.Bot{Player: p})
	}

	// allow redirects
	options := map[string]game.Option{}
	options["with_redirect"] = game.Option{Value: "1", Exposed: true}

	// 36 cards, least rank is 6
	r := game.Six
	options["min_rank"] = game.Option{Value: r.ToString(), Exposed: true}

	// build the Bot Manager that takes care of bots actions
	botManager := game.NewBotManager(bots)

	// Create game
	g := game.NewGame(&game.Deck{},
		players,
		options,
		false,
		false,
		game.Table{},
		botManager)

	err := g.StartGame()
	if err != nil {
		log.Panic(err)
	}

	// Let bots do their actions. Bots will perform their actions until there is nothing to do.
	// In real game this should be called on every action of real player.
	// Since there is no real player in demo bots will act until game over
	g.CycleBots()

	//dumping actions
	l, _ := json.MarshalIndent(g.Log, "", "  ")
	log.Println(string(l))
}
