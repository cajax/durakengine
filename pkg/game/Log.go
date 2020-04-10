package game

const (
	StartEventType    = "started"
	AttackEventType   = "attack"
	DefenseEventType  = "defense"
	TransferEventType = "transfer"
	PickupEventType   = "pickup"
	RefillEventType   = "refill"
	EndTurnEventType  = "end_turn"
	GameOverEventType = "game_over"
	AbandonEventType  = "abandon"
	//todo add other Events like Card spread and refills
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
	firstSeq--
	lastSeq--
	if firstSeq < 0 {
		firstSeq = 0
	}
	if firstSeq >= lastSeq {
		lastSeq = firstSeq + 1
	}

	if lastSeq > len(l.Events) {
		lastSeq = len(l.Events)
	}

	for _, eventGroup := range l.Events[firstSeq:lastSeq] {
		for _, event := range eventGroup {
			events = append(events, event)
		}
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
