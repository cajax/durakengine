package game

// OptionMaxAttackCards limits number of attack cards per turn. Value is an integer, 0 or empty means unlimited
const OptionMaxAttackCards = "max_attack_cards"

// OptionThrowIn allows neighbors of defender to throw in cards when set to "1"
const OptionThrowIn = "with_throw_in"

// Option represents game option
type Option struct {
	Value   string `json:"value"`
	Exposed bool
}
