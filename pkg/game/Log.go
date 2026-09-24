package game

const (
	StartEventType    = "started"
	AttackEventType   = "attack"
	DefenseEventType  = "defense"
	TransferEventType = "transfer"
	DiscardEventType  = "discard"
	PickupEventType   = "pickup"
	RefillEventType   = "refill"
	EndTurnEventType  = "end_turn"
	GameOverEventType = "game_over"
	AbandonEventType  = "abandon"
)

type Log struct {
	Events [][]EventInterface `json:"events"`
}

type EventType string

type EventInterface interface {
	GetType() EventType
	GetSequence() int
	SetSequence(seq int)
}

type Event struct {
	Sequence int       `json:"sequence"`
	Type     EventType `json:"type"`
}

func NewLog() Log {
	return Log{Events: [][]EventInterface{}}
}

func (l *Log) Advance() {
	l.Events = append(l.Events, []EventInterface{})
}

func (l *Log) Add(e EventInterface) {

	e.SetSequence(len(l.Events))

	l.Events[len(l.Events)-1] = append(l.Events[len(l.Events)-1], e)
}

func (l *Log) GetEvents(firstSeq int, lastSeq int) []EventInterface {
	events := []EventInterface{}

	// Convert to 0-based bounds without decrementing, so MinInt can't wrap
	lo, hi := 0, 0
	if firstSeq > 0 {
		lo = firstSeq - 1
	}
	if lastSeq > 0 {
		hi = lastSeq - 1
	}
	if lo >= hi {
		hi = lo + 1
	}

	if lo >= len(l.Events) {
		return events
	}
	if hi > len(l.Events) {
		hi = len(l.Events)
	}

	for _, eventGroup := range l.Events[lo:hi] {
		events = append(events, eventGroup...)
	}

	return events
}

func (e *Event) SetSequence(seq int) {
	e.Sequence = seq
}

func (e *Event) GetSequence() int {
	return e.Sequence
}

func (e *Event) GetType() EventType {
	return e.Type
}
