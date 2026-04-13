package chess

import (
	"fmt"
	"strconv"
	"strings"
)

func (g *Game) buildFEN(board *Board) string {
	boardPart := board.ToFEN()

	activeColor := "w"
	if g.Turn == Black {
		activeColor = "b"
	}

	castling := g.castlingString()
	if castling == "" {
		castling = "-"
	}

	enPassant := "-"
	if g.EnPassantTarget != nil {
		enPassant = g.EnPassantTarget.String()
	}

	halfmove := 0
	fullmove := len(g.Moves)/2 + 1

	return fmt.Sprintf("%s %s %s %s %d %d", boardPart, activeColor, castling, enPassant, halfmove, fullmove)
}

func (g *Game) ToFEN() string {
	var board *Board
	var err error
	if g.CurrentFEN != "" {
		board, err = BoardFromFEN(g.CurrentFEN)
		if err != nil {
			board = NewBoard()
		}
	} else {
		board = NewBoard()
	}

	boardPart := board.ToFEN()

	activeColor := "w"
	if g.Turn == Black {
		activeColor = "b"
	}

	castling := g.castlingString()
	if castling == "" {
		castling = "-"
	}

	enPassant := "-"
	if g.EnPassantTarget != nil {
		enPassant = g.EnPassantTarget.String()
	}

	halfmove := 0
	fullmove := len(g.Moves)/2 + 1

	return fmt.Sprintf("%s %s %s %s %d %d", boardPart, activeColor, castling, enPassant, halfmove, fullmove)
}

func (g *Game) castlingString() string {
	var sb strings.Builder
	if g.CanCastleWhiteKingside {
		sb.WriteByte('K')
	}
	if g.CanCastleWhiteQueenside {
		sb.WriteByte('Q')
	}
	if g.CanCastleBlackKingside {
		sb.WriteByte('k')
	}
	if g.CanCastleBlackQueenside {
		sb.WriteByte('q')
	}
	return sb.String()
}

func (b *Board) ToFEN() string {
	var fen strings.Builder
	for row := 0; row < BoardSize; row++ {
		emptyCount := 0
		for col := 0; col < BoardSize; col++ {
			piece := b[row][col]
			if piece.IsEmpty() {
				emptyCount++
			} else {
				if emptyCount > 0 {
					fen.WriteString(strconv.Itoa(emptyCount))
					emptyCount = 0
				}
				fen.WriteString(pieceToFEN(piece))
			}
		}
		if emptyCount > 0 {
			fen.WriteString(strconv.Itoa(emptyCount))
		}
		if row < BoardMax {
			fen.WriteString("/")
		}
	}
	return fen.String()
}

func pieceToFEN(p Piece) string {
	var c byte
	switch p.Type {
	case Pawn:
		c = 'p'
	case Rook:
		c = 'r'
	case Knight:
		c = 'n'
	case Bishop:
		c = 'b'
	case Queen:
		c = 'q'
	case King:
		c = 'k'
	default:
		return ""
	}
	if p.Color == White {
		c = byte(c) - 'a' + 'A'
	}
	return string(c)
}

func BoardFromFEN(fen string) (*Board, error) {
	board, _, _, _, err := ParseFullFEN(fen)
	return board, err
}

func ParseFullFEN(fen string) (*Board, Color, CastlingRights, *Position, error) {
	parts := strings.Fields(fen)
	if len(parts) < 2 {
		return nil, White, CastlingRights{}, nil, fmt.Errorf("invalid FEN string: expected at least board and active color")
	}

	board, err := parseBoard(parts[0])
	if err != nil {
		return nil, White, CastlingRights{}, nil, err
	}

	turn := White
	if parts[1] == "b" {
		turn = Black
	}

	castling := CastlingRights{}
	if len(parts) >= 3 {
		castling = parseCastling(parts[2])
	}

	var enPassant *Position
	if len(parts) >= 4 && parts[3] != "-" {
		pos, err := ParseAlgebraic(parts[3])
		if err == nil {
			enPassant = &pos
		}
	}

	return board, turn, castling, enPassant, nil
}

type CastlingRights struct {
	WhiteKingside  bool
	WhiteQueenside bool
	BlackKingside  bool
	BlackQueenside bool
}

func parseCastling(s string) CastlingRights {
	c := CastlingRights{}
	for _, ch := range s {
		switch ch {
		case 'K':
			c.WhiteKingside = true
		case 'Q':
			c.WhiteQueenside = true
		case 'k':
			c.BlackKingside = true
		case 'q':
			c.BlackQueenside = true
		}
	}
	return c
}

func parseBoard(boardPart string) (*Board, error) {
	ranks := strings.Split(boardPart, "/")
	if len(ranks) != BoardSize {
		return nil, fmt.Errorf("FEN board part must have %d ranks", BoardSize)
	}

	board := &Board{}
	for row, rank := range ranks {
		col := 0
		for _, ch := range rank {
			if col >= BoardSize {
				return nil, fmt.Errorf("too many columns in rank %d", row)
			}
			if ch >= '1' && ch <= '8' {

				empty := int(ch - '0')
				for i := 0; i < empty; i++ {
					if col >= BoardSize {
						return nil, fmt.Errorf("too many empty squares in rank %d", row)
					}
					board[row][col] = Piece{}
					col++
				}
			} else {

				piece, err := fenToPiece(ch)
				if err != nil {
					return nil, fmt.Errorf("invalid piece char '%c' at rank %d, col %d: %w", ch, row, col, err)
				}
				board[row][col] = piece
				col++
			}
		}
		if col != BoardSize {
			return nil, fmt.Errorf("rank %d has %d squares, expected %d", row, col, BoardSize)
		}
	}
	return board, nil
}

func fenToPiece(ch rune) (Piece, error) {
	var p Piece
	switch ch {
	case 'P':
		p = Piece{Pawn, White}
	case 'N':
		p = Piece{Knight, White}
	case 'B':
		p = Piece{Bishop, White}
	case 'R':
		p = Piece{Rook, White}
	case 'Q':
		p = Piece{Queen, White}
	case 'K':
		p = Piece{King, White}
	case 'p':
		p = Piece{Pawn, Black}
	case 'n':
		p = Piece{Knight, Black}
	case 'b':
		p = Piece{Bishop, Black}
	case 'r':
		p = Piece{Rook, Black}
	case 'q':
		p = Piece{Queen, Black}
	case 'k':
		p = Piece{King, Black}
	default:
		return Piece{}, fmt.Errorf("unknown piece character: %c", ch)
	}
	return p, nil
}
