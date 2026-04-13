package chess

import (
	"fmt"
	"time"

	"chess-backend/pkg/mathutil"

	"github.com/google/uuid"
)

type GameID string

type GameStatus string

const (
	StatusWaiting  GameStatus = "waiting"
	StatusActive   GameStatus = "active"
	StatusFinished GameStatus = "finished"
)

type GameResult string

const (
	ResultWhiteWins GameResult = "white_wins"
	ResultBlackWins GameResult = "black_wins"
	ResultDraw      GameResult = "draw"
)

type GameEndReason string

const (
	ReasonCheckmate GameEndReason = "checkmate"
	ReasonStalemate GameEndReason = "stalemate"
	ReasonResign    GameEndReason = "resign"
	ReasonTimeout   GameEndReason = "timeout"
	ReasonDrawAgree GameEndReason = "draw_agreement"
)

type DrawOffer struct {
	OfferedBy string    `json:"offered_by"`
	OfferedAt time.Time `json:"offered_at"`
}

type Game struct {
	ID            GameID
	WhitePlayerID *string
	BlackPlayerID *string
	Status        GameStatus
	Result        *GameResult
	EndReason     *GameEndReason
	CurrentFEN    string
	Turn          Color
	Moves         []Move
	CreatedAt     time.Time
	UpdatedAt     time.Time

	CanCastleWhiteKingside  bool
	CanCastleWhiteQueenside bool
	CanCastleBlackKingside  bool
	CanCastleBlackQueenside bool

	EnPassantTarget *Position

	DrawOffer *DrawOffer

	TimeControl TimeControl

	WhiteRemaining int64
	BlackRemaining int64

	LastMoveAt time.Time
}

func NewGame(creatorID string, tc TimeControl, colorPref ColorPreference) *Game {
	assignedColor := ResolveColor(colorPref)
	g := &Game{
		ID:          GameID(uuid.New().String()),
		Status:      StatusWaiting,
		Turn:        White,
		Moves:       []Move{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		LastMoveAt:  time.Time{},
		TimeControl: tc,

		CanCastleWhiteKingside:  true,
		CanCastleWhiteQueenside: true,
		CanCastleBlackKingside:  true,
		CanCastleBlackQueenside: true,
		EnPassantTarget:         nil,
	}

	if assignedColor == White {
		g.WhitePlayerID = &creatorID
	} else {
		g.BlackPlayerID = &creatorID
	}

	if tc.IsTimed() {
		ms := int64(tc.Base) * 1000
		g.WhiteRemaining = ms
		g.BlackRemaining = ms
	}

	g.CurrentFEN = g.ToFEN()
	return g
}

func (g *Game) GetBoard() (*Board, error) {
	board, err := BoardFromFEN(g.CurrentFEN)
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (g *Game) loadFullState() (*Board, error) {
	board, turn, castling, enPassant, err := ParseFullFEN(g.CurrentFEN)
	if err != nil {
		return nil, err
	}
	g.Turn = turn
	g.CanCastleWhiteKingside = castling.WhiteKingside
	g.CanCastleWhiteQueenside = castling.WhiteQueenside
	g.CanCastleBlackKingside = castling.BlackKingside
	g.CanCastleBlackQueenside = castling.BlackQueenside
	g.EnPassantTarget = enPassant
	return board, nil
}

func (g *Game) Join(opponentID string) error {
	if g.Status != StatusWaiting {
		return ErrGameNotWaiting
	}
	if g.WhitePlayerID != nil && opponentID == *g.WhitePlayerID {
		return ErrOwnGame
	}
	if g.BlackPlayerID != nil && opponentID == *g.BlackPlayerID {
		return ErrOwnGame
	}

	if g.WhitePlayerID == nil {
		g.WhitePlayerID = &opponentID
	} else if g.BlackPlayerID == nil {
		g.BlackPlayerID = &opponentID
	} else {
		return ErrGameFull
	}
	g.Status = StatusActive
	g.UpdatedAt = time.Now()
	return nil
}

func (g *Game) MakeMove(playerID string, from, to Position) error {
	if g.Status != StatusActive {
		return ErrGameNotActive
	}

	playerColor, ok := g.colorOf(playerID)
	if !ok {
		return ErrPlayerNotInGame
	}

	if playerColor != g.Turn {
		return ErrNotYourTurn
	}

	if g.TimeControl.IsTimed() {
		var elapsed int64
		if !g.LastMoveAt.IsZero() {
			elapsed = time.Since(g.LastMoveAt).Milliseconds()
		}
		remaining := g.WhiteRemaining
		if playerColor == Black {
			remaining = g.BlackRemaining
		}
		remaining -= elapsed
		if remaining <= 0 {
			g.WhiteRemaining = 0
			g.BlackRemaining = 0
			g.Status = StatusFinished
			winner := ResultWhiteWins
			if playerColor == White {
				winner = ResultBlackWins
			}
			g.Result = &winner
			reason := ReasonTimeout
			g.EndReason = &reason
			g.DrawOffer = nil
			g.UpdatedAt = time.Now()
			g.LastMoveAt = time.Now().UTC()
			return ErrTimeout
		}
		remaining += int64(g.TimeControl.Increment) * 1000
		if playerColor == White {
			g.WhiteRemaining = remaining
		} else {
			g.BlackRemaining = remaining
		}
		g.LastMoveAt = time.Now().UTC()
	}

	board, err := g.loadFullState()
	if err != nil {
		return fmt.Errorf("invalid board state: %w", err)
	}

	piece := board.PieceAt(from)
	if piece.IsEmpty() {
		return ErrInvalidMove
	}
	if piece.Color != playerColor {
		return ErrInvalidMove
	}

	isCastle := g.isCastlingMove(piece, from, to)
	if isCastle {
		if err := g.validateCastling(board, playerColor, from, to); err != nil {
			return err
		}
	} else {
		if !g.isValidMove(board, piece, from, to) {
			return ErrInvalidMove
		}
		if !g.moveWouldBeSafe(board, playerColor, from, to) {
			return ErrInvalidMove
		}
	}

	isEP := g.isEnPassantCapture(piece, to)

	if isCastle {
		g.executeCastling(board, from, to)
	} else {
		if err := board.MovePiece(from, to); err != nil {
			return err
		}
	}

	if isEP {
		capturedPos := Position{Row: from.Row, Col: to.Col}
		board.SetPiece(capturedPos, Piece{})
	}

	var promotion *PieceType
	if !isCastle && isPawnPromotion(piece, to) {
		board.promotePawn(to, playerColor)
		p := Queen
		promotion = &p
	}

	g.updateCastlingRights(piece, from, to)

	g.EnPassantTarget = nil
	if piece.Type == Pawn && mathutil.Abs(to.Row-from.Row) == 2 {
		epTarget := Position{Row: (from.Row + to.Row) / 2, Col: from.Col}
		g.EnPassantTarget = &epTarget
	}

	if g.Turn == White {
		g.Turn = Black
	} else {
		g.Turn = White
	}

	g.CurrentFEN = g.buildFEN(board)

	g.DrawOffer = nil
	g.Moves = append(g.Moves, Move{
		From:      from,
		To:        to,
		PlayerID:  playerID,
		Timestamp: time.Now(),
		Promotion: promotion,
		Castle:    isCastle,
		EnPassant: isEP,
	})

	nextBoard, err := g.loadFullState()
	if err != nil {
		return fmt.Errorf("invalid board state after move: %w", err)
	}
	if !g.hasAnyLegalMove(nextBoard, g.Turn) {
		g.Status = StatusFinished
		if g.isSquareAttacked(nextBoard, g.findKing(nextBoard, g.Turn), g.Turn) {
			winner := ResultWhiteWins
			if g.Turn == White {
				winner = ResultBlackWins
			} else {
				winner = ResultWhiteWins
			}
			g.Result = &winner
			reason := ReasonCheckmate
			g.EndReason = &reason
		} else {
			draw := ResultDraw
			g.Result = &draw
			reason := ReasonStalemate
			g.EndReason = &reason
		}
	}

	g.UpdatedAt = time.Now()
	return nil
}

func (g *Game) Resign(playerID string) error {
	if g.Status != StatusActive {
		return ErrGameNotActive
	}
	if _, ok := g.colorOf(playerID); !ok {
		return ErrPlayerNotInGame
	}
	g.Status = StatusFinished
	winner := ResultWhiteWins
	if g.WhitePlayerID != nil && playerID == *g.WhitePlayerID {
		winner = ResultBlackWins
	}
	g.Result = &winner
	reason := ReasonResign
	g.EndReason = &reason
	g.DrawOffer = nil
	g.UpdatedAt = time.Now()
	return nil
}

func (g *Game) OfferDraw(playerID string) error {
	if g.Status != StatusActive {
		return ErrGameNotActive
	}
	if _, ok := g.colorOf(playerID); !ok {
		return ErrPlayerNotInGame
	}
	if g.DrawOffer != nil && g.DrawOffer.OfferedBy == playerID {
		return ErrDrawAlreadyOffered
	}
	g.DrawOffer = &DrawOffer{
		OfferedBy: playerID,
		OfferedAt: time.Now(),
	}
	g.UpdatedAt = time.Now()
	return nil
}

func (g *Game) AcceptDraw(playerID string) error {
	if g.Status != StatusActive {
		return ErrGameNotActive
	}
	if g.DrawOffer == nil {
		return ErrNoDrawOffer
	}

	if playerID == g.DrawOffer.OfferedBy {
		return ErrCannotAcceptOwnOffer
	}
	if _, ok := g.colorOf(playerID); !ok {
		return ErrPlayerNotInGame
	}
	g.Status = StatusFinished
	draw := ResultDraw
	g.Result = &draw
	reason := ReasonDrawAgree
	g.EndReason = &reason
	g.DrawOffer = nil
	g.UpdatedAt = time.Now()
	return nil
}

func (g *Game) DeclineDraw(playerID string) error {
	if g.DrawOffer == nil {
		return ErrNoDrawOffer
	}
	if _, ok := g.colorOf(playerID); !ok {
		return ErrPlayerNotInGame
	}

	g.DrawOffer = nil
	g.UpdatedAt = time.Now()
	return nil
}

func (g *Game) ExpireDrawOffer() {
	g.DrawOffer = nil
	g.UpdatedAt = time.Now()
}

func (g *Game) IsPlayerInGame(playerID string) bool {
	_, ok := g.colorOf(playerID)
	return ok
}

func (g *Game) RemainingFor(color Color) int64 {
	remaining := g.WhiteRemaining
	if color == Black {
		remaining = g.BlackRemaining
	}
	if !g.TimeControl.IsTimed() || g.Status != StatusActive {
		return remaining
	}

	if color == g.Turn && !g.LastMoveAt.IsZero() {
		elapsed := time.Since(g.LastMoveAt).Milliseconds()
		live := remaining - elapsed
		if live < 0 {
			live = 0
		}
		return live
	}
	return remaining
}

func (g *Game) CheckTimeout(playerID string) bool {
	if !g.TimeControl.IsTimed() || g.Status != StatusActive {
		return false
	}
	color, ok := g.colorOf(playerID)
	if !ok {
		return false
	}
	return g.RemainingFor(color) <= 0
}

func (g *Game) EndOnTimeout(playerID string) error {
	if g.Status != StatusActive {
		return ErrGameNotActive
	}
	if _, ok := g.colorOf(playerID); !ok {
		return ErrPlayerNotInGame
	}

	if g.WhitePlayerID != nil && playerID == *g.WhitePlayerID {
		g.WhiteRemaining = 0
	} else {
		g.BlackRemaining = 0
	}
	g.Status = StatusFinished
	winner := ResultWhiteWins
	if g.WhitePlayerID != nil && playerID == *g.WhitePlayerID {
		winner = ResultBlackWins
	}
	g.Result = &winner
	reason := ReasonTimeout
	g.EndReason = &reason
	g.DrawOffer = nil
	g.UpdatedAt = time.Now().UTC()
	g.LastMoveAt = time.Now().UTC()
	return nil
}

func (g *Game) colorOf(playerID string) (Color, bool) {
	if g.WhitePlayerID != nil && playerID == *g.WhitePlayerID {
		return White, true
	}
	if g.BlackPlayerID != nil && playerID == *g.BlackPlayerID {
		return Black, true
	}
	return "", false
}
