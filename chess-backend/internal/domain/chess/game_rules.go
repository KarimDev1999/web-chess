package chess

func (g *Game) isValidMove(board *Board, piece Piece, from, to Position) bool {
	dr := to.Row - from.Row
	dc := to.Col - from.Col
	absDR := dr
	if absDR < 0 {
		absDR = -absDR
	}
	absDC := dc
	if absDC < 0 {
		absDC = -absDC
	}
	switch piece.Type {
	case Pawn:
		return g.isValidPawnMove(board, piece.Color, from, to, dr, dc, absDC)
	case Rook:
		return g.isValidRookMove(board, from, to, dr, dc) && !g.hasFriendlyPiece(board, piece.Color, to)
	case Knight:
		return ((absDR == 2 && absDC == 1) || (absDR == 1 && absDC == 2)) && !g.hasFriendlyPiece(board, piece.Color, to)
	case Bishop:
		return absDR == absDC && g.isPathClear(board, from, to) && !g.hasFriendlyPiece(board, piece.Color, to)
	case Queen:
		return (absDR == absDC || dr == 0 || dc == 0) && g.isPathClear(board, from, to) && !g.hasFriendlyPiece(board, piece.Color, to)
	case King:
		if absDR > 1 || absDC > 1 {
			return false
		}
		return !g.hasFriendlyPiece(board, piece.Color, to)
	}
	return false
}

func (g *Game) hasFriendlyPiece(board *Board, color Color, to Position) bool {
	target := board.PieceAt(to)
	return !target.IsEmpty() && target.Color == color
}

func (g *Game) isValidPawnMove(board *Board, color Color, from, to Position, dr, dc, absDC int) bool {
	if color == White {
		if dc == 0 {
			if dr == -1 && board.IsEmpty(to) {
				return true
			}
			if from.Row == 6 && dr == -2 && board.IsEmpty(to) && board.IsEmpty(Position{Row: from.Row - 1, Col: from.Col}) {
				return true
			}
		} else if absDC == 1 && dr == -1 {
			target := board.PieceAt(to)
			if !target.IsEmpty() && target.Color == Black {
				return true
			}
			if g.EnPassantTarget != nil && to == *g.EnPassantTarget {
				return true
			}
		}
	} else {
		if dc == 0 {
			if dr == 1 && board.IsEmpty(to) {
				return true
			}
			if from.Row == 1 && dr == 2 && board.IsEmpty(to) && board.IsEmpty(Position{Row: from.Row + 1, Col: from.Col}) {
				return true
			}
		} else if absDC == 1 && dr == 1 {
			target := board.PieceAt(to)
			if !target.IsEmpty() && target.Color == White {
				return true
			}
			if g.EnPassantTarget != nil && to == *g.EnPassantTarget {
				return true
			}
		}
	}
	return false
}

func (g *Game) isValidRookMove(board *Board, from, to Position, dr, dc int) bool {
	if dr != 0 && dc != 0 {
		return false
	}
	return g.isPathClear(board, from, to)
}

func (g *Game) isPathClear(board *Board, from, to Position) bool {
	dr := to.Row - from.Row
	dc := to.Col - from.Col
	stepRow := 0
	if dr > 0 {
		stepRow = 1
	} else if dr < 0 {
		stepRow = -1
	}
	stepCol := 0
	if dc > 0 {
		stepCol = 1
	} else if dc < 0 {
		stepCol = -1
	}
	r, c := from.Row+stepRow, from.Col+stepCol
	for r != to.Row || c != to.Col {
		if !board.IsEmpty(Position{Row: r, Col: c}) {
			return false
		}
		r += stepRow
		c += stepCol
	}
	return true
}

func (g *Game) moveWouldBeSafe(board *Board, playerColor Color, from, to Position) bool {
	originalFrom := board.PieceAt(from)
	originalTo := board.PieceAt(to)
	var savedEP Position
	var savedEPPiece Piece
	isEP := g.isEnPassantCapture(originalFrom, to)
	if isEP {
		capturedPos := Position{Row: from.Row, Col: to.Col}
		savedEP = capturedPos
		savedEPPiece = board.PieceAt(capturedPos)
	}

	board.SetPiece(to, originalFrom)
	board.SetPiece(from, Piece{})
	if isEP {
		board.SetPiece(savedEP, Piece{})
	}

	kingPos := g.findKing(board, playerColor)
	inCheck := g.isSquareAttacked(board, kingPos, playerColor)

	if isEP {
		board.SetPiece(savedEP, savedEPPiece)
	}
	board.SetPiece(from, originalFrom)
	board.SetPiece(to, originalTo)
	return !inCheck
}
