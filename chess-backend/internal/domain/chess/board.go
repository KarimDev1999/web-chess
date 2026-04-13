package chess

import (
	"errors"
)

type Board [BoardSize][BoardSize]Piece

func NewBoard() *Board {
	b := &Board{}

	b[BoardMax][0] = Piece{Rook, White}
	b[BoardMax][1] = Piece{Knight, White}
	b[BoardMax][2] = Piece{Bishop, White}
	b[BoardMax][3] = Piece{Queen, White}
	b[BoardMax][4] = Piece{King, White}
	b[BoardMax][5] = Piece{Bishop, White}
	b[BoardMax][6] = Piece{Knight, White}
	b[BoardMax][7] = Piece{Rook, White}
	for i := 0; i < BoardSize; i++ {
		b[BoardSize-2][i] = Piece{Pawn, White}
	}

	b[BoardMin][0] = Piece{Rook, Black}
	b[BoardMin][1] = Piece{Knight, Black}
	b[BoardMin][2] = Piece{Bishop, Black}
	b[BoardMin][3] = Piece{Queen, Black}
	b[BoardMin][4] = Piece{King, Black}
	b[BoardMin][5] = Piece{Bishop, Black}
	b[BoardMin][6] = Piece{Knight, Black}
	b[BoardMin][7] = Piece{Rook, Black}
	for i := 0; i < BoardSize; i++ {
		b[1][i] = Piece{Pawn, Black}
	}
	return b
}

func (b *Board) PieceAt(pos Position) Piece {
	return b[pos.Row][pos.Col]
}

func (b *Board) SetPiece(pos Position, p Piece) {
	b[pos.Row][pos.Col] = p
}

func (b *Board) MovePiece(from, to Position) error {
	if !from.IsValid() || !to.IsValid() {
		return errors.New("invalid position")
	}
	piece := b[from.Row][from.Col]
	if piece.IsEmpty() {
		return errors.New("no piece at source")
	}
	b[to.Row][to.Col] = piece
	b[from.Row][from.Col] = Piece{}
	return nil
}

func (b *Board) IsEmpty(pos Position) bool {
	return b[pos.Row][pos.Col].IsEmpty()
}

func (b *Board) String() string {
	s := ""
	for r := 0; r < BoardSize; r++ {
		for c := 0; c < BoardSize; c++ {
			s += b[r][c].String() + " "
		}
		s += "\n"
	}
	return s
}
