package game

import "time"

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
	// clock gives the time of actions, time.Now if nil
	clock func() time.Time
	// at is the time of the current action
	at time.Time
}

type EventType string

type EventInterface interface {
	GetType() EventType
	GetSequence() int
	SetSequence(seq int)
	GetAt() time.Time
	SetAt(at time.Time)
}

type Event struct {
	Sequence int       `json:"sequence"`
	Type     EventType `json:"type"`
	// At is the time of the action that logged the event, in UTC
	At time.Time `json:"at"`
}

func NewLog() Log {
	return Log{Events: [][]EventInterface{}}
}

// Advance starts the group of events of the next action and reads its time from the clock
func (l *Log) Advance() {
	clock := l.clock
	if clock == nil {
		clock = time.Now
	}
	l.at = clock().UTC()
	l.Events = append(l.Events, []EventInterface{})
}

// Add appends event to the current action, stamped with the action's sequence and time
func (l *Log) Add(e EventInterface) {

	e.SetSequence(len(l.Events))
	e.SetAt(l.at)

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

func (e *Event) SetAt(at time.Time) {
	e.At = at
}

func (e *Event) GetAt() time.Time {
	return e.At
}

func (e *Event) GetType() EventType {
	return e.Type
}
