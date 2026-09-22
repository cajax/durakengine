package game

// OptionMaxAttackCards limits number of attack cards per turn. Value is an integer, 0 or empty means unlimited
const OptionMaxAttackCards = "max_attack_cards"

// OptionThrowIn allows neighbors of defender to throw in cards when set to "1"
const OptionThrowIn = "with_throw_in"

// OptionRedirectKeepCard lets defender redirect by showing cards instead of putting them on table when set to "1".
// Each card can be shown once per turn
const OptionRedirectKeepCard = "redirect_keep_card"

// Option represents game option
type Option struct {
	Value   string `json:"value"`
	Exposed bool
}
