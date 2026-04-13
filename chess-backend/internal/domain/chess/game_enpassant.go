package chess

func (g *Game) isEnPassantCapture(piece Piece, to Position) bool {
	if piece.Type != Pawn {
		return false
	}
	if g.EnPassantTarget == nil {
		return false
	}
	return to == *g.EnPassantTarget
}
