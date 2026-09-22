package game

// OptionMaxAttackCards limits number of attack cards per turn. Value is an integer, 0 or empty means unlimited
const OptionMaxAttackCards = "max_attack_cards"

// OptionMinRank sets least card rank in deck, e.g. "6" for 36 cards deck
const OptionMinRank = "min_rank"

// OptionRedirect allows defender to redirect attack to the next player when set to "1"
const OptionRedirect = "with_redirect"

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
