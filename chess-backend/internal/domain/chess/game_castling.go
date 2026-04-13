package chess

import "chess-backend/pkg/mathutil"

func (g *Game) isCastlingMove(piece Piece, from, to Position) bool {
	if piece.Type != King {
		return false
	}
	if mathutil.Abs(to.Col-from.Col) != 2 || to.Row != from.Row {
		return false
	}
	if piece.Color == White && from.Row != BoardMax {
		return false
	}
	if piece.Color == Black && from.Row != BoardMin {
		return false
	}
	return true
}

func (g *Game) validateCastling(board *Board, color Color, from, to Position) error {
	rookCol := BoardMax
	if to.Col < from.Col {
		rookCol = BoardMin
	}
	if !g.hasCastlingRight(color, rookCol == BoardMax) {
		return ErrInvalidMove
	}
	if g.isSquareAttacked(board, from, color) {
		return ErrInvalidMove
	}
	dir := 1
	if to.Col < from.Col {
		dir = -1
	}
	for c := from.Col + dir; c != rookCol; c += dir {
		sq := Position{Row: from.Row, Col: c}
		if !board.IsEmpty(sq) {
			return ErrInvalidMove
		}
		if c == from.Col+dir || c == to.Col {
			if g.isSquareAttacked(board, sq, color) {
				return ErrInvalidMove
			}
		}
	}
	rookPos := Position{Row: from.Row, Col: rookCol}
	rook := board.PieceAt(rookPos)
	if rook.Type != Rook || rook.Color != color {
		return ErrInvalidMove
	}
	return nil
}

func (g *Game) hasCastlingRight(color Color, kingside bool) bool {
	switch {
	case color == White && kingside:
		return g.CanCastleWhiteKingside
	case color == White:
		return g.CanCastleWhiteQueenside
	case color == Black && kingside:
		return g.CanCastleBlackKingside
	case color == Black:
		return g.CanCastleBlackQueenside
	}
	return false
}

func (g *Game) executeCastling(board *Board, from, to Position) {
	rookFromCol := BoardMax
	if to.Col < from.Col {
		rookFromCol = BoardMin
	}
	rookToCol := from.Col + mathutil.Sign(to.Col-from.Col)
	rook := board.PieceAt(Position{Row: from.Row, Col: rookFromCol})
	board.SetPiece(Position{Row: from.Row, Col: rookToCol}, rook)
	board.SetPiece(Position{Row: from.Row, Col: rookFromCol}, Piece{})
	board.SetPiece(to, board.PieceAt(from))
	board.SetPiece(from, Piece{})
}

func (g *Game) updateCastlingRights(piece Piece, from, to Position) {
	if piece.Type == King {
		if piece.Color == White {
			g.CanCastleWhiteKingside = false
			g.CanCastleWhiteQueenside = false
		} else {
			g.CanCastleBlackKingside = false
			g.CanCastleBlackQueenside = false
		}
	}
	if piece.Type == Rook {
		g.revokeRookRights(piece.Color, from)
	}
	g.checkRookCaptured(to)
}

func (g *Game) revokeRookRights(color Color, pos Position) {
	if color == White && pos.Row == BoardMax {
		if pos.Col == BoardMax {
			g.CanCastleWhiteKingside = false
		}
		if pos.Col == BoardMin {
			g.CanCastleWhiteQueenside = false
		}
	} else if color == Black && pos.Row == BoardMin {
		if pos.Col == BoardMax {
			g.CanCastleBlackKingside = false
		}
		if pos.Col == BoardMin {
			g.CanCastleBlackQueenside = false
		}
	}
}

func (g *Game) checkRookCaptured(pos Position) {
	if pos.Row == BoardMax {
		if pos.Col == BoardMax {
			g.CanCastleWhiteKingside = false
		}
		if pos.Col == BoardMin {
			g.CanCastleWhiteQueenside = false
		}
	}
	if pos.Row == BoardMin {
		if pos.Col == BoardMax {
			g.CanCastleBlackKingside = false
		}
		if pos.Col == BoardMin {
			g.CanCastleBlackQueenside = false
		}
	}
}
