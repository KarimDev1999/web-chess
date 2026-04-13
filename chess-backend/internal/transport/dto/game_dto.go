package dto

import "time"

type PlayerInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type GameResponse struct {
	ID             string              `json:"id"`
	WhitePlayer    PlayerInfo          `json:"white_player"`
	BlackPlayer    *PlayerInfo         `json:"black_player"`
	Status         string              `json:"status"`
	Result         *string             `json:"result,omitempty"`
	EndReason      *string             `json:"end_reason,omitempty"`
	FEN            string              `json:"fen"`
	Turn           string              `json:"turn"`
	Moves          []MoveResponse      `json:"moves"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	DrawOffer      *DrawOfferResponse  `json:"draw_offer,omitempty"`
	TimeControl    TimeControlResponse `json:"time_control"`
	WhiteRemaining *int64              `json:"white_remaining,omitempty"`
	BlackRemaining *int64              `json:"black_remaining,omitempty"`
	LastMoveAt     *time.Time          `json:"last_move_at,omitempty"`
}

type TimeControlResponse struct {
	Base      int `json:"base"`
	Increment int `json:"increment"`
}

type TimeControlPresetResponse struct {
	Label     string `json:"label"`
	Base      int    `json:"base"`
	Increment int    `json:"increment"`
}

type CreateGameRequest struct {
	TimeBase      int    `json:"time_base"`
	TimeIncrement int    `json:"time_increment"`
	ColorPref     string `json:"color_pref"`
}

type DrawOfferResponse struct {
	OfferedBy string    `json:"offered_by"`
	OfferedAt time.Time `json:"offered_at"`
}

type MoveResponse struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Promotion *string   `json:"promotion"`
	MadeAt    time.Time `json:"made_at"`
	Castle    bool      `json:"castle"`
	EnPassant bool      `json:"en_passant"`
}

type MoveRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
}
