package game

import "fmt"
// Suit family
type Suit int
//Rank level
type Rank int8

// Card suits
const (
	Hearts Suit = iota + 1
	Spades
	Clubs
	Diamonds
)

// Card ranks
const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

//Card representation
type Card struct {
	Suit Suit `json:"suit"`
	Rank Rank `json:"rank"`
}

// ToString renders suit to readable form
func (s *Suit) ToString() string {
	suits := map[Suit]string{
		Hearts:   "♥",
		Spades:   "♠",
		Clubs:    "♣",
		Diamonds: "♦",
	}
	return suits[*s]
}

// ToString Render rank to readable form
func (r *Rank) ToString() string {
	ranks := map[Rank]string{
		Two:   "2",
		Three: "3",
		Four:  "4",
		Five:  "5",
		Six:   "6",
		Seven: "7",
		Eight: "8",
		Nine:  "9",
		Ten:   "10",
		Jack:  "J",
		Queen: "Q",
		King:  "K",
		Ace:   "A",
	}

	return ranks[*r]
}

// RankFromString convert rank from string to internal format
func RankFromString(str string) Rank {
	ranks := map[string]Rank{
		"2":  Two,
		"3":  Three,
		"4":  Four,
		"5":  Five,
		"6":  Six,
		"7":  Seven,
		"8":  Eight,
		"9":  Nine,
		"10": Ten,
		"J":  Jack,
		"Q":  Queen,
		"K":  King,
		"A":  Ace,
	}

	return ranks[str]
}

// ToString renders card to readable format
func (c *Card) ToString() string {
	return fmt.Sprintf("%s%s", c.Rank.ToString(), c.Suit.ToString())
}
