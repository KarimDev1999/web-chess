package chess

func (g *Game) findKing(board *Board, color Color) Position {
	for r := 0; r < BoardSize; r++ {
		for c := 0; c < BoardSize; c++ {
			p := board[r][c]
			if p.Type == King && p.Color == color {
				return Position{Row: r, Col: c}
			}
		}
	}
	return Position{}
}

func (g *Game) isSquareAttacked(board *Board, pos Position, defenderColor Color) bool {
	var attackerColor Color
	if defenderColor == White {
		attackerColor = Black
	} else {
		attackerColor = White
	}
	for r := 0; r < BoardSize; r++ {
		for c := 0; c < BoardSize; c++ {
			piece := board[r][c]
			if piece.IsEmpty() || piece.Color != attackerColor {
				continue
			}
			from := Position{Row: r, Col: c}
			if g.isValidMove(board, piece, from, pos) {
				return true
			}
		}
	}
	return false
}

func (g *Game) hasAnyLegalMove(board *Board, color Color) bool {
	for r := 0; r < BoardSize; r++ {
		for c := 0; c < BoardSize; c++ {
			piece := board[r][c]
			if piece.IsEmpty() || piece.Color != color {
				continue
			}
			from := Position{Row: r, Col: c}
			for tr := 0; tr < BoardSize; tr++ {
				for tc := 0; tc < BoardSize; tc++ {
					to := Position{Row: tr, Col: tc}
					if from == to {
						continue
					}
					if g.isValidMove(board, piece, from, to) && g.moveWouldBeSafe(board, color, from, to) {
						return true
					}
				}
			}
		}
	}
	return false
}

func isPawnPromotion(piece Piece, to Position) bool {
	if piece.Type != Pawn {
		return false
	}
	return to.Row == BoardMin && piece.Color == White || to.Row == BoardMax && piece.Color == Black
}

func (b *Board) promotePawn(pos Position, color Color) {
	b[pos.Row][pos.Col] = Piece{Queen, color}
}
