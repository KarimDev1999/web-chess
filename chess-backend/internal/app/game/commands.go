package game

import "chess-backend/internal/domain/chess"

type CreateGameCommand struct {
	PlayerID    string
	TimeControl chess.TimeControl
	ColorPref   chess.ColorPreference
}

type JoinGameCommand struct {
	GameID   string
	PlayerID string
}

type MakeMoveCommand struct {
	GameID   string
	PlayerID string
	From     string
	To       string
}

type ResignGameCommand struct {
	GameID   string
	PlayerID string
}

type OfferDrawCommand struct {
	GameID   string
	PlayerID string
}

type AcceptDrawCommand struct {
	GameID   string
	PlayerID string
}

type DeclineDrawCommand struct {
	GameID   string
	PlayerID string
}

type GetWaitingGamesQuery struct{}
