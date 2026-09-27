package game

// RuleError is the error returned when an action breaks a game rule or does not fit the game state.
// Compare it with the Err* values using errors.Is, or get its code with errors.As.
type RuleError struct {
	code    string
	message string
}

func newRuleError(code, message string) *RuleError {
	return &RuleError{code: code, message: message}
}

// Error returns the message, one of the Error* constants.
func (e *RuleError) Error() string {
	return e.message
}

// Code returns a stable machine-readable code: the Error* constant's name without "Error", in UPPER_SNAKE.
func (e *RuleError) Code() string {
	return e.code
}

// Errors returned by the engine. Their messages are the Error* constants.
var (
	ErrGameAlreadyStarted              = newRuleError("GAME_ALREADY_STARTED", ErrorGameAlreadyStarted)
	ErrAddPlayerWhileInGame            = newRuleError("ADD_PLAYER_WHILE_IN_GAME", ErrorAddPlayerWhileInGame)
	ErrDeckNoTrump                     = newRuleError("DECK_NO_TRUMP", ErrorDeckNoTrump)
	ErrTooFewPlayers                   = newRuleError("TOO_FEW_PLAYERS", ErrorTooFewPlayers)
	ErrTooManyPlayers                  = newRuleError("TOO_MANY_PLAYERS", ErrorTooManyPlayers)
	ErrInvalidMinRank                  = newRuleError("INVALID_MIN_RANK", ErrorInvalidMinRank)
	ErrDeckTooSmall                    = newRuleError("DECK_TOO_SMALL", ErrorDeckTooSmall)
	ErrNotInPlayingState               = newRuleError("NOT_IN_PLAYING_STATE", ErrorNotInPlayingState)
	ErrFirstAttackWithDifferentRanks   = newRuleError("FIRST_ATTACK_WITH_DIFFERENT_RANKS", ErrorFirstAttackWithDifferentRanks)
	ErrFirstAttackByNeighbor           = newRuleError("FIRST_ATTACK_BY_NEIGHBOR", ErrorFirstAttackByNeighbor)
	ErrAttackByWrongPlayer             = newRuleError("ATTACK_BY_WRONG_PLAYER", ErrorAttackByWrongPlayer)
	ErrAttackRankNotOnTable            = newRuleError("ATTACK_RANK_NOT_ON_TABLE", ErrorAttackRankNotOnTable)
	ErrAttackIsTooBig                  = newRuleError("ATTACK_IS_TOO_BIG", ErrorAttackIsTooBig)
	ErrAttackerHasNoCard               = newRuleError("ATTACKER_HAS_NO_CARD", ErrorAttackerHasNoCard)
	ErrNotDefender                     = newRuleError("NOT_DEFENDER", ErrorNotDefender)
	ErrDefenderHasNoCard               = newRuleError("DEFENDER_HAS_NO_CARD", ErrorDefenderHasNoCard)
	ErrPairIndexOutOfRange             = newRuleError("PAIR_INDEX_OUT_OF_RANGE", ErrorPairIndexOutOfRange)
	ErrAlreadyDefended                 = newRuleError("ALREADY_DEFENDED", ErrorAlreadyDefended)
	ErrCardCantBeat                    = newRuleError("CARD_CANT_BEAT", ErrorCardCantBeat)
	ErrDefenseByNotDefender            = newRuleError("DEFENSE_BY_NOT_DEFENDER", ErrorDefenseByNotDefender)
	ErrGameEndTurnByNotAttacker        = newRuleError("GAME_END_TURN_BY_NOT_ATTACKER", ErrorGameEndTurnByNotAttacker)
	ErrGameEndTurnWithoutAttack        = newRuleError("GAME_END_TURN_WITHOUT_ATTACK", ErrorGameEndTurnWithoutAttack)
	ErrGameEndTurnWithoutUnbeatenCards = newRuleError("GAME_END_TURN_WITHOUT_UNBEATEN_CARDS", ErrorGameEndTurnWithoutUnbeatenCards)
	ErrNoRedirectsAllowed              = newRuleError("NO_REDIRECTS_ALLOWED", ErrorNoRedirectsAllowed)
	ErrRedirectWithNoOrMixedCards      = newRuleError("REDIRECT_WITH_NO_OR_MIXED_CARDS", ErrorRedirectWithNoOrMixedCards)
	ErrAlreadyDefending                = newRuleError("ALREADY_DEFENDING", ErrorAlreadyDefending)
	ErrRedirectRankMismatch            = newRuleError("REDIRECT_RANK_MISMATCH", ErrorRedirectRankMismatch)
	ErrPlayerAlreadyQuit               = newRuleError("PLAYER_ALREADY_QUIT", ErrorPlayerAlreadyQuit)
	ErrNoPlayerToRedirect              = newRuleError("NO_PLAYER_TO_REDIRECT", ErrorNoPlayerToRedirect)
	ErrAttackLimitReached              = newRuleError("ATTACK_LIMIT_REACHED", ErrorAttackLimitReached)
	ErrRedirectCardAlreadyShown        = newRuleError("REDIRECT_CARD_ALREADY_SHOWN", ErrorRedirectCardAlreadyShown)
	ErrPickupFromEmptyTable            = newRuleError("PICKUP_FROM_EMPTY_TABLE", ErrorPickupFromEmptyTable)
)
