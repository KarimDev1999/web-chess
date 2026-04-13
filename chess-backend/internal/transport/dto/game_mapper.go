package dto

import (
	"context"

	"chess-backend/internal/constants/appconst"
	"chess-backend/internal/domain/chess"
	"chess-backend/internal/domain/user"
)

type UserLookup func(ctx context.Context, userID string) (*user.User, error)

func ToGameResponse(ctx context.Context, lookup UserLookup, g *chess.Game, pendingDrawOffer ...*chess.DrawOffer) GameResponse {
	resp := GameResponse{
		ID:        string(g.ID),
		Status:    string(g.Status),
		FEN:       g.CurrentFEN,
		Turn:      string(g.Turn),
		Moves:     make([]MoveResponse, len(g.Moves)),
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
		TimeControl: TimeControlResponse{
			Base:      g.TimeControl.Base,
			Increment: g.TimeControl.Increment,
		},
	}

	if g.Result != nil {
		s := string(*g.Result)
		resp.Result = &s
	}
	if g.EndReason != nil {
		s := string(*g.EndReason)
		resp.EndReason = &s
	}

	if g.TimeControl.IsTimed() {
		resp.WhiteRemaining = &g.WhiteRemaining
		resp.BlackRemaining = &g.BlackRemaining
		if !g.LastMoveAt.IsZero() {
			resp.LastMoveAt = &g.LastMoveAt
		}
	}

	if g.WhitePlayerID != nil {
		if u, err := lookup(ctx, *g.WhitePlayerID); err == nil && u != nil {
			resp.WhitePlayer = PlayerInfo{ID: u.ID, Username: u.Username}
		} else {
			resp.WhitePlayer = PlayerInfo{ID: *g.WhitePlayerID, Username: appconst.UsernameUnknown}
		}
	}

	if g.BlackPlayerID != nil {
		if u, err := lookup(ctx, *g.BlackPlayerID); err == nil && u != nil {
			resp.BlackPlayer = &PlayerInfo{ID: u.ID, Username: u.Username}
		} else {
			resp.BlackPlayer = &PlayerInfo{ID: *g.BlackPlayerID, Username: appconst.UsernameUnknown}
		}
	}

	for i, m := range g.Moves {
		resp.Moves[i] = ToMoveResponse(m)
	}

	var drawOffer *chess.DrawOffer
	if len(pendingDrawOffer) > 0 && pendingDrawOffer[0] != nil {
		drawOffer = pendingDrawOffer[0]
	} else if g.DrawOffer != nil {
		drawOffer = g.DrawOffer
	}
	if drawOffer != nil {
		resp.DrawOffer = &DrawOfferResponse{
			OfferedBy: drawOffer.OfferedBy,
			OfferedAt: drawOffer.OfferedAt,
		}
	}

	return resp
}

func ToMoveResponse(m chess.Move) MoveResponse {
	resp := MoveResponse{
		From:      m.From.String(),
		To:        m.To.String(),
		MadeAt:    m.Timestamp,
		Castle:    m.Castle,
		EnPassant: m.EnPassant,
	}
	if m.Promotion != nil {
		s := string(*m.Promotion)
		resp.Promotion = &s
	}
	return resp
}

func ToGameResponses(ctx context.Context, lookup UserLookup, games []*chess.Game) []GameResponse {
	result := make([]GameResponse, len(games))
	for i, g := range games {
		result[i] = ToGameResponse(ctx, lookup, g)
	}
	return result
}
