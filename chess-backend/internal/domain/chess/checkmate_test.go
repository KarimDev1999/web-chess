package chess

import (
	"testing"
)

func fromAlg(s string) Position {
	p, _ := ParseAlgebraic(s)
	return p
}

func toAlg(s string) Position {
	p, _ := ParseAlgebraic(s)
	return p
}

// King + Queen vs King checkmate:
// White: Kg6, Qg7. Black: Kh8. Black to move.
// Q on g7 attacks h8, h7, g8. King can't go anywhere. CHECKMATE.
func TestCheckmateDetection_KingQueenMate(t *testing.T) {
	// 7k/6Q1/6K1/8/8/8/8/8 b - - 0 1
	// Row 0: 7k = h8=black king
	// Row 1: 6Q1 = g7=white queen
	// Row 2: 6K1 = g6=white king
	board, err := BoardFromFEN("7k/6Q1/6K1/8/8/8/8/8 b - - 0 1")
	if err != nil {
		t.Fatalf("failed to parse FEN: %v", err)
	}

	game := &Game{Turn: Black}

	hasLegal := game.hasAnyLegalMove(board, Black)
	if hasLegal {
		for r := 0; r < 8; r++ {
			for c := 0; c < 8; c++ {
				piece := board[r][c]
				if piece.IsEmpty() || piece.Color != Black {
					continue
				}
				from := Position{Row: r, Col: c}
				for tr := 0; tr < 8; tr++ {
					for tc := 0; tc < 8; tc++ {
						to := Position{Row: tr, Col: tc}
						if from == to {
							continue
						}
						if game.isValidMove(board, piece, from, to) {
							safe := game.moveWouldBeSafe(board, Black, from, to)
							t.Logf("Black %s at %s can move to %s (safe=%v)",
								piece.Type, from.String(), to.String(), safe)
						}
					}
				}
			}
		}
		t.Fatal("expected no legal moves for Black (K+Q vs K checkmate)")
	}

	kingPos := Position{Row: 0, Col: 7} // h8
	if !game.isSquareAttacked(board, kingPos, Black) {
		t.Error("expected Black king on h8 to be in check from queen on g7")
	}
}

// Back-rank mate: White rook on e8, Black king on g8, Black pawns on f7, g7, h7
// King can't go to f8/h8 (rook attacks), can't go to f7/g7/h7 (own pawns)
// No piece can block or capture the rook
func TestCheckmateDetection_BackRankMate(t *testing.T) {
	// Row 0: 4R1k1 = e8=R, g8=k  ← rook on e8, king on g8
	// Row 1: 5ppp = f7,p; g7,p; h7,p
	board, err := BoardFromFEN("4R1k1/5ppp/8/8/8/8/8/6K1 b - - 0 1")
	if err != nil {
		t.Fatalf("failed to parse FEN: %v", err)
	}

	game := &Game{Turn: Black}

	hasLegal := game.hasAnyLegalMove(board, Black)
	if hasLegal {
		for r := 0; r < 8; r++ {
			for c := 0; c < 8; c++ {
				piece := board[r][c]
				if piece.IsEmpty() || piece.Color != Black {
					continue
				}
				from := Position{Row: r, Col: c}
				for tr := 0; tr < 8; tr++ {
					for tc := 0; tc < 8; tc++ {
						to := Position{Row: tr, Col: tc}
						if from == to {
							continue
						}
						if game.isValidMove(board, piece, from, to) {
							safe := game.moveWouldBeSafe(board, Black, from, to)
							t.Logf("Black %s at %s can move to %s (safe=%v)",
								piece.Type, from.String(), to.String(), safe)
						}
					}
				}
			}
		}
		t.Fatal("expected no legal moves for Black (back-rank mate)")
	}

	kingPos := Position{Row: 0, Col: 6} // g8
	if !game.isSquareAttacked(board, kingPos, Black) {
		t.Error("expected Black king on g8 to be in check from rook on e8")
	}
}

// Simple stalemate position
// Black: Ka1. White: Ka3, Qb3. Black to move.
// King on a1 can go to: a2 (attacked by Ka3+Qb3), b1 (attacked by Qb3), b2 (attacked by Ka3+Qb3).
// King is NOT in check. STALEMATE.
func TestCheckmateDetection_Stalemate(t *testing.T) {
	// Row 5 (rank 3): Ka3, Qb3 → KQ6
	// Row 7 (rank 1): ka1 → k7
	board, err := BoardFromFEN("8/8/8/8/8/KQ6/8/k7 b - - 0 1")
	if err != nil {
		t.Fatalf("failed to parse FEN: %v", err)
	}

	game := &Game{Turn: Black}

	hasLegal := game.hasAnyLegalMove(board, Black)
	if hasLegal {
		for r := 0; r < 8; r++ {
			for c := 0; c < 8; c++ {
				piece := board[r][c]
				if piece.IsEmpty() || piece.Color != Black {
					continue
				}
				from := Position{Row: r, Col: c}
				for tr := 0; tr < 8; tr++ {
					for tc := 0; tc < 8; tc++ {
						to := Position{Row: tr, Col: tc}
						if from == to {
							continue
						}
						if game.isValidMove(board, piece, from, to) {
							safe := game.moveWouldBeSafe(board, Black, from, to)
							t.Logf("Black %s at %s can move to %s (safe=%v)",
								piece.Type, from.String(), to.String(), safe)
						}
					}
				}
			}
		}
		t.Fatal("expected no legal moves for Black (stalemate)")
	}

	kingPos := Position{Row: 7, Col: 0} // a1
	if game.isSquareAttacked(board, kingPos, Black) {
		t.Error("expected Black king on a1 to NOT be in check (stalemate)")
	}
}

func TestIsSquareAttacked_Pawn(t *testing.T) {
	// White pawn on e5, Black king on d6 (pawn attacks diagonally)
	board, _ := BoardFromFEN("8/8/3k4/4P3/8/8/8/8 w - - 0 1")
	game := &Game{}

	kingPos := Position{Row: 2, Col: 3} // d6
	if !game.isSquareAttacked(board, kingPos, Black) {
		t.Error("expected White pawn on e5 to attack Black king on d6")
	}
}

func TestIsSquareAttacked_Knight(t *testing.T) {
	// White knight on f7, Black king on d8
	board, _ := BoardFromFEN("3k4/5N2/8/8/8/8/8/8 w - - 0 1")
	game := &Game{}

	kingPos := Position{Row: 0, Col: 3} // d8
	if !game.isSquareAttacked(board, kingPos, Black) {
		t.Error("expected White knight on f7 to attack Black king on d8")
	}
}

func TestIsSquareAttacked_Queen(t *testing.T) {
	// White queen on f7, Black king on e8
	board, _ := BoardFromFEN("4k3/5Q2/8/8/8/8/8/8 w - - 0 1")
	game := &Game{}

	kingPos := Position{Row: 0, Col: 4} // e8
	if !game.isSquareAttacked(board, kingPos, Black) {
		t.Error("expected White queen on f7 to attack Black king on e8")
	}
}

func TestKingCannotCaptureOwnPieces(t *testing.T) {
	// Black king on e8, Black pawn on f7
	board, _ := BoardFromFEN("4k3/5p2/8/8/8/8/8/8 w - - 0 1")
	game := &Game{}

	kingPos := Position{Row: 0, Col: 4} // e8
	to := Position{Row: 1, Col: 5}      // f7 (has own pawn)

	if game.isValidMove(board, Piece{King, Black}, kingPos, to) {
		t.Error("king should not be able to move to square occupied by own piece")
	}
}

// Test full MakeMove checkmate detection through Game.MakeMove
// Fool's Mate: 1. f3 e5 2. g4 Qh4#
func TestMakeMove_DetectsCheckmate_FoolsMate(t *testing.T) {
	game := NewGame("white-player-id", TimeControl{}, PreferenceWhite)
	game.Join("black-player-id")

	// 1. f3
	err := game.MakeMove("white-player-id", fromAlg("f2"), toAlg("f3"))
	if err != nil {
		t.Fatalf("f2-f3 failed: %v", err)
	}

	// 1... e5
	err = game.MakeMove("black-player-id", fromAlg("e7"), toAlg("e5"))
	if err != nil {
		t.Fatalf("e7-e5 failed: %v", err)
	}

	// 2. g4
	err = game.MakeMove("white-player-id", fromAlg("g2"), toAlg("g4"))
	if err != nil {
		t.Fatalf("g2-g4 failed: %v", err)
	}

	// 2... Qh4# (checkmate!)
	err = game.MakeMove("black-player-id", fromAlg("d8"), toAlg("h4"))
	if err != nil {
		t.Fatalf("Qd8-h4 (checkmate) failed: %v", err)
	}

	if game.Status != StatusFinished {
		t.Errorf("expected game status 'finished', got '%s'", game.Status)
	}

	board, _ := game.GetBoard()
	whiteKingPos := game.findKing(board, White)
	if !game.isSquareAttacked(board, whiteKingPos, White) {
		t.Error("expected White king to be in check after Qh4#")
	}
}

// Test that checkmate detection also works for the losing side
func TestMakeMove_NoIllegalMovesAfterCheckmate(t *testing.T) {
	game := NewGame("white-player-id", TimeControl{}, PreferenceWhite)
	game.Join("black-player-id")

	game.MakeMove("white-player-id", fromAlg("f2"), toAlg("f3"))
	game.MakeMove("black-player-id", fromAlg("e7"), toAlg("e5"))
	game.MakeMove("white-player-id", fromAlg("g2"), toAlg("g4"))
	game.MakeMove("black-player-id", fromAlg("d8"), toAlg("h4"))

	err := game.MakeMove("white-player-id", fromAlg("e1"), toAlg("f2"))
	if err == nil {
		t.Error("expected error when trying to move after checkmate")
	}
}
