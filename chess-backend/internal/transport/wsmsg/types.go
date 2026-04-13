package wsmsg

import (
	"time"

	"chess-backend/internal/transport/dto"
)

const (
	TypeGameCreated      = "game.created"
	TypeGameJoined       = "game.joined"
	TypeMoveMade         = "game.move_made"
	TypeDrawOffered      = "game.draw_offered"
	TypeDrawDeclined     = "game.draw_declined"
	TypeDrawOfferExpired = "game.draw_offer_expired"
	TypeGameResigned     = "game.resigned"
	TypeDrawAccepted     = "game.draw_accepted"
	TypeGameTimedOut     = "game.timed_out"
)

const (
	ClientTypePing        = "ping"
	ClientTypePong        = "pong"
	ClientTypeJoinGame    = "join_game"
	ClientTypeLeaveGame   = "leave_game"
	ClientTypeResign      = "resign"
	ClientTypeOfferDraw   = "offer_draw"
	ClientTypeAcceptDraw  = "accept_draw"
	ClientTypeDeclineDraw = "decline_draw"
)

const (
	TypePresence = "presence"
)

const (
	KeyGameID = "game_id"
	KeyType   = "type"
)

type Message struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

type GameMessage struct {
	Message
	GameID string           `json:"game_id"`
	Game   dto.GameResponse `json:"game"`
}

func NewGameMessage(msgType, gameID string, game dto.GameResponse) GameMessage {
	return GameMessage{
		Message: Message{Type: msgType, Timestamp: time.Now()},
		GameID:  gameID,
		Game:    game,
	}
}

type ClientMessage struct {
	Message
	Data map[string]any `json:"data,omitempty"`
}

type Presence struct {
	Type    string   `json:"type"`
	GameID  string   `json:"game_id"`
	Players []string `json:"players"`
}

func NewPresence(gameID string, players []string) Presence {
	return Presence{
		Type:    TypePresence,
		GameID:  gameID,
		Players: players,
	}
}
