package chess

import "errors"

var (
	ErrGameNotFound         = errors.New("game not found")
	ErrGameNotActive        = errors.New("game is not active")
	ErrNotYourTurn          = errors.New("not your turn")
	ErrInvalidMove          = errors.New("invalid move")
	ErrGameFull             = errors.New("game is full")
	ErrPlayerNotInGame      = errors.New("player not in this game")
	ErrGameNotWaiting       = errors.New("game is not waiting for players")
	ErrConcurrentUpdate     = errors.New("game state changed, please retry")
	ErrOwnGame              = errors.New("cannot join your own game")
	ErrDrawAlreadyOffered   = errors.New("you already have a pending draw offer")
	ErrNoDrawOffer          = errors.New("no draw offer to respond to")
	ErrCannotAcceptOwnOffer = errors.New("cannot accept your own draw offer")
	ErrTimeout              = errors.New("you ran out of time")
)
