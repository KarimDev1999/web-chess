package events

import (
	"encoding/json"
	"fmt"
	"time"

	"chess-backend/internal/transport/wsmsg"
)

type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

type BaseEvent struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func (e BaseEvent) EventName() string     { return e.Type }
func (e BaseEvent) OccurredAt() time.Time { return e.Timestamp }

type GameCreated struct {
	BaseEvent
	GameID        string `json:"game_id"`
	WhitePlayerID string `json:"white_player_id"`
}

type GameJoined struct {
	BaseEvent
	GameID        string `json:"game_id"`
	WhitePlayerID string `json:"white_player_id"`
	BlackPlayerID string `json:"black_player_id"`
}

type MoveMade struct {
	BaseEvent
	GameID        string  `json:"game_id"`
	WhitePlayerID string  `json:"white_player_id"`
	BlackPlayerID string  `json:"black_player_id"`
	PlayerID      string  `json:"player_id"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	Result        *string `json:"result,omitempty"`
	EndReason     *string `json:"end_reason,omitempty"`
}

type DrawOffered struct {
	BaseEvent
	GameID        string `json:"game_id"`
	WhitePlayerID string `json:"white_player_id"`
	BlackPlayerID string `json:"black_player_id"`
	OfferedBy     string `json:"offered_by"`
}

type DrawDeclined struct {
	BaseEvent
	GameID        string `json:"game_id"`
	WhitePlayerID string `json:"white_player_id"`
	BlackPlayerID string `json:"black_player_id"`
	DeclinedBy    string `json:"declined_by"`
}

type DrawOfferExpired struct {
	BaseEvent
	GameID        string `json:"game_id"`
	WhitePlayerID string `json:"white_player_id"`
	BlackPlayerID string `json:"black_player_id"`
}

type GameResigned struct {
	BaseEvent
	GameID        string  `json:"game_id"`
	WhitePlayerID string  `json:"white_player_id"`
	BlackPlayerID string  `json:"black_player_id"`
	ResignedBy    string  `json:"resigned_by"`
	Result        *string `json:"result,omitempty"`
	EndReason     *string `json:"end_reason,omitempty"`
}

type DrawAccepted struct {
	BaseEvent
	GameID        string  `json:"game_id"`
	WhitePlayerID string  `json:"white_player_id"`
	BlackPlayerID string  `json:"black_player_id"`
	AcceptedBy    string  `json:"accepted_by"`
	Result        *string `json:"result,omitempty"`
	EndReason     *string `json:"end_reason,omitempty"`
}

type GameTimedOut struct {
	BaseEvent
	GameID        string  `json:"game_id"`
	WhitePlayerID string  `json:"white_player_id"`
	BlackPlayerID string  `json:"black_player_id"`
	PlayerID      string  `json:"player_id"`
	Result        *string `json:"result,omitempty"`
	EndReason     *string `json:"end_reason,omitempty"`
}

const (
	EventGameCreated      = wsmsg.TypeGameCreated
	EventGameJoined       = wsmsg.TypeGameJoined
	EventMoveMade         = wsmsg.TypeMoveMade
	EventDrawOffered      = wsmsg.TypeDrawOffered
	EventDrawDeclined     = wsmsg.TypeDrawDeclined
	EventDrawOfferExpired = wsmsg.TypeDrawOfferExpired
	EventGameResigned     = wsmsg.TypeGameResigned
	EventDrawAccepted     = wsmsg.TypeDrawAccepted
	EventGameTimedOut     = wsmsg.TypeGameTimedOut
)

func UnmarshalEvent(eventName string, data []byte) (DomainEvent, error) {
	switch eventName {
	case EventGameCreated:
		var e GameCreated
		err := json.Unmarshal(data, &e)
		return e, err
	case EventGameJoined:
		var e GameJoined
		err := json.Unmarshal(data, &e)
		return e, err
	case EventMoveMade:
		var e MoveMade
		err := json.Unmarshal(data, &e)
		return e, err
	case EventDrawOffered:
		var e DrawOffered
		err := json.Unmarshal(data, &e)
		return e, err
	case EventDrawDeclined:
		var e DrawDeclined
		err := json.Unmarshal(data, &e)
		return e, err
	case EventDrawOfferExpired:
		var e DrawOfferExpired
		err := json.Unmarshal(data, &e)
		return e, err
	case EventGameResigned:
		var e GameResigned
		err := json.Unmarshal(data, &e)
		return e, err
	case EventDrawAccepted:
		var e DrawAccepted
		err := json.Unmarshal(data, &e)
		return e, err
	case EventGameTimedOut:
		var e GameTimedOut
		err := json.Unmarshal(data, &e)
		return e, err
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventName)
	}
}
