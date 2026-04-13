package chess

import (
	"testing"
)

// All pieces should NOT be able to move to squares occupied by friendly pieces.
func TestNoPieceCanCaptureFriendly(t *testing.T) {
	// White rook a1 (row7,col0), White pawn a2 (row6,col0)
	board, _ := BoardFromFEN("8/8/8/8/8/8/P7/R7 w - - 0 1")
	game := &Game{}

	// Rook a1 → a2 (has White pawn)
	if game.isValidMove(board, Piece{Rook, White}, Position{Row: 7, Col: 0}, Position{Row: 6, Col: 0}) {
		t.Error("rook at a1 SHOULD NOT be able to move to a2 (friendly piece)")
	}
}

// White rook on a1, White pawn on a3. Rook should NOT move to a3.
func TestRookCannotCaptureFriendly(t *testing.T) {
	// White rook a1 (row7,col0), White pawn a3 (row5,col0)
	// Row 5 = rank 3, but I want a3 = rank 3 = row 5. FEN row order: rank8,rank7,...,rank1
	// a3 = rank 3 = row 5. So P goes on row 5.
	board, _ := BoardFromFEN("8/8/8/8/8/P7/8/R7 w - - 0 1")
	game := &Game{}

	// Rook a1 (row7,col0) → a3 (row5,col0) has White pawn
	if game.isValidMove(board, Piece{Rook, White}, Position{Row: 7, Col: 0}, Position{Row: 5, Col: 0}) {
		t.Error("White rook should NOT be able to move to a3 (occupied by friendly pawn)")
	}
}

func TestKnightCannotCaptureFriendly(t *testing.T) {
	// White knight c3 (row5,col2), White pawn a4 (row4,col0)
	// a4 = rank 4 = row 4
	board, _ := BoardFromFEN("8/8/8/8/P7/2N5/8/8 w - - 0 1")
	game := &Game{}

	// Knight c3 (row5,col2) → a4 (row4,col0) has White pawn
	if game.isValidMove(board, Piece{Knight, White}, Position{Row: 5, Col: 2}, Position{Row: 4, Col: 0}) {
		t.Error("White knight should NOT be able to move to a4 (occupied by friendly pawn)")
	}
}

func TestBishopCannotCaptureFriendly(t *testing.T) {
	// White bishop c1 (row7,col2), White pawn a3 (row5,col0)
	// c1→a3: dr=-2,dc=-2 (diagonal, path=b2=row6,col1)
	board, _ := BoardFromFEN("8/8/8/8/8/P7/8/2B5 w - - 0 1")
	game := &Game{}

	// Bishop c1 (row7,col2) → a3 (row5,col0) has White pawn
	if game.isValidMove(board, Piece{Bishop, White}, Position{Row: 7, Col: 2}, Position{Row: 5, Col: 0}) {
		t.Error("White bishop should NOT be able to move to a3 (occupied by friendly pawn)")
	}
}

func TestQueenCannotCaptureFriendly(t *testing.T) {
	// White queen d1 (row7,col3), White pawn a4 (row4,col0)
	// d1→a4: dr=-3,dc=-3 (diagonal)
	board, _ := BoardFromFEN("8/8/8/8/P7/8/8/3Q4 w - - 0 1")
	game := &Game{}

	// Queen d1 (row7,col3) → a4 (row4,col0) has White pawn
	if game.isValidMove(board, Piece{Queen, White}, Position{Row: 7, Col: 3}, Position{Row: 4, Col: 0}) {
		t.Error("White queen should NOT be able to move to a4 (occupied by friendly pawn)")
	}
}

func TestPawnCannotCaptureFriendly(t *testing.T) {
	// White pawn e5, White pawn d6
	board, _ := BoardFromFEN("8/8/3P4/4P3/8/8/8/8 w - - 0 1")
	game := &Game{}

	// White pawn e5 → d6 (has White pawn): e5=row3,col4  d6=row2,col3
	// Pawn capture: dr=-1, dc=-1, absDC=1
	if game.isValidMove(board, Piece{Pawn, White}, Position{Row: 3, Col: 4}, Position{Row: 2, Col: 3}) {
		t.Error("White pawn should NOT be able to capture friendly pawn on d6")
	}
}
