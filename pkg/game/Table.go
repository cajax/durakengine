package game

// Table represent the list of card pairs on table
type Table struct {
	pairs []*TablePair
}

// TablePair contains details of pair of attack card with its owner and optionally defense card with defender
type TablePair struct {
	Attack   *Card
	Defense  *Card
	Attacker *Player
	Defender *Player
}

func (t *Table) GetCardsOnTable() []*Card {
	var cards []*Card
	for _, pair := range t.pairs {
		if nil != pair.Attack {
			cards = append(cards, pair.Attack)
		}

		if nil != pair.Defense {
			cards = append(cards, pair.Defense)
		}
	}

	return cards
}

func (t *Table) clear() {
	t.pairs = []*TablePair{}
}

func (t *Table) IsEmpty() bool {
	if t.pairs == nil || len(t.pairs) == 0 {
		return true
	}
	return false
}

func (t *Table) cardMatchesSomeRanks(c *Card) bool {
	for _, tc := range t.GetCardsOnTable() {
		if tc.Rank == c.Rank {
			return true
		}
	}

	return false
}

func (t *Table) attack(p *Player, c *Card) {
	t.pairs = append(t.pairs, t.createPair(p, c))
}

func (t *Table) createPair(attacker *Player, card *Card) *TablePair {
	return &TablePair{Attacker: attacker, Attack: card}
}

// defend put defense card on pair
func (t *Table) defend(p *Player, tp *TablePair, c *Card) {
	tp.Defense = c
	tp.Defender = p
}

// defenseStarted returns true if at least one card on table is covered
func (t *Table) defenseStarted() bool {
	for _, pair := range t.pairs {
		if nil != pair.Defense {
			return true
		}
	}
	return false
}

// cardMatchesAttackRank returns true if all cards on table are of the same rank as given card and none of them covered
//
// Use it to check whether cards can be redirected
func (t *Table) cardMatchesAttackRank(c *Card) bool {
	if t.defenseStarted() {
		return false
	}

	for _, p := range t.pairs {
		if p.Attack.Rank != c.Rank {
			return false
		}
	}

	return true
}

// HasUnbeatenCards checks if there it at least one card to beat
func (t *Table) HasUnbeatenCards() bool {
	for _, p := range t.pairs {
		if p.Defense == nil {

			return true
		}
	}

	return false
}

// GetCardsToBeat returns cards that are not covered yet
func (t *Table) GetCardsToBeat() []*Card {
	var cards []*Card
	for _, p := range t.pairs {
		if p.Defense == nil {
			cards = append(cards, p.Attack)
		}
	}

	return cards
}

// GetPairs returns list of pairs currently on table
func (t *Table) GetPairs() []*TablePair {
	return t.pairs
}

func (t *Table) HasPair(i int) bool {
	return i > -1 && i < len(t.pairs)
}
func (t *Table) GetPair(i int) *TablePair {
	return t.pairs[i]
}
