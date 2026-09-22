package game

// OptionMaxAttackCards limits number of attack cards per turn. Value is an integer, 0 or empty means unlimited
const OptionMaxAttackCards = "max_attack_cards"

// Option represents game option
type Option struct {
	Value   string `json:"value"`
	Exposed bool
}
