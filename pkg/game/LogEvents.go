package game

type StartEvent struct {
	*Event
	CardsCount    int               `json:"cards_count"`
	MinRank       Rank              `json:"min_rank"`
	FirstAttacker Player            `json:"first_attacker"`
	FirstDefender Player            `json:"first_defender"`
	TrumpSuit     Suit              `json:"trump_suit"`
	Options       map[string]Option `json:"Options"`
}

func NewStartEvent(cardsCount int, minRank Rank, firstAttacker Player, firstDefender Player, trumpSuit Suit, options map[string]Option) *StartEvent {
	return &StartEvent{
		CardsCount:    cardsCount,
		MinRank:       minRank,
		FirstAttacker: firstAttacker,
		FirstDefender: firstDefender,
		TrumpSuit:     trumpSuit,
		Options:       options,
		Event:         &Event{Type: StartEventType},
	}

}

type AttackEvent struct {
	*Event
	Player Player  `json:"player"`
	Cards  []*Card `json:"cards"`
}

func NewAttackEvent(player Player, cards []*Card) *AttackEvent {
	return &AttackEvent{Player: player, Cards: cards, Event: &Event{Type: AttackEventType}}
}

func (a *AttackEvent) GetCards() []*Card {
	return a.Cards
}

func (a *AttackEvent) GetPlayer() Player {
	return a.Player
}

type PickupEvent struct {
	*Event
	Player Player  `json:"player"`
	Cards  []*Card `json:"cards"`
}

func NewPickupEvent(player Player, cards []*Card) *PickupEvent {
	return &PickupEvent{Player: player, Cards: cards, Event: &Event{Type: PickupEventType}}
}

func (e *PickupEvent) GetCards() []*Card {
	return e.Cards
}

func (e *PickupEvent) GetPlayer() Player {
	return e.Player
}

type DefenseEvent struct {
	*Event
	Player Player    `json:"player"`
	Pair   TablePair `json:"pair"`
	Card   Card      `json:"card"`
}

func NewDefenseEvent(player Player, pair TablePair, card Card) *DefenseEvent {
	return &DefenseEvent{Player: player, Pair: pair, Card: card, Event: &Event{Type: DefenseEventType}}
}

func (e *DefenseEvent) GetPlayer() Player {
	return e.Player
}

func (e *DefenseEvent) GetCard() Card {
	return e.Card
}

func (e *DefenseEvent) GetPair() TablePair {
	return e.Pair
}

type RefillEvent struct {
	*Event
	Player Player `json:"player"`
	Cards  []Card `json:"cards"`
}

func NewRefillEvent(player Player, cards []Card) *RefillEvent {
	return &RefillEvent{Player: player, Cards: cards, Event: &Event{Type: RefillEventType}}
}

func (e *RefillEvent) GetPlayer() Player {
	return e.Player
}

type EndTurnEvent struct {
	*Event
	NextAttacker Player `json:"next_attacker"`
	NextDefender Player `json:"next_defender"`
}

func NewEndTurnEvent(attacker Player, defender Player) *EndTurnEvent {
	return &EndTurnEvent{NextAttacker: attacker, NextDefender: defender, Event: &Event{Type: EndTurnEventType}}
}

type OverEvent struct {
	*Event
	LastPlayers []Player `json:"last_players"`
}

func NewGameOverEvent(lastPlayers []Player) *OverEvent {
	return &OverEvent{LastPlayers: lastPlayers, Event: &Event{Type: GameOverEventType}}
}

type AbandonEvent struct {
	*Event
	Player Player  `json:"player"`
	Cards  []*Card `json:"cards"`
}

func NewAbandonEvent(player Player, cards []*Card) *AbandonEvent {
	return &AbandonEvent{Player: player, Cards: cards, Event: &Event{Type: AbandonEventType}}
}

func (e *AbandonEvent) GetCards() []*Card {
	return e.Cards
}

func (e *AbandonEvent) GetPlayer() Player {
	return e.Player
}
