package chess

import (
	"testing"
)

func TestCastling_KingsideWhite(t *testing.T) {
	// Clear path: no pieces on f1, g1
	board, _ := BoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1")
	game := &Game{
		Turn:                    White,
		CanCastleWhiteKingside:  true,
		CanCastleWhiteQueenside: true,
		CanCastleBlackKingside:  true,
		CanCastleBlackQueenside: true,
	}

	if !game.isCastlingMove(Piece{King, White}, fromAlg("e1"), toAlg("g1")) {
		t.Error("e1-g1 should be detected as castling move")
	}

	err := game.validateCastling(board, White, fromAlg("e1"), toAlg("g1"))
	if err != nil {
		t.Errorf("kingside castling should be valid: %v", err)
	}
}

func TestCastling_QueensideWhite(t *testing.T) {
	// Clear path: no pieces on b1, c1, d1
	board, _ := BoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1")
	game := &Game{
		Turn:                    White,
		CanCastleWhiteKingside:  true,
		CanCastleWhiteQueenside: true,
		CanCastleBlackKingside:  true,
		CanCastleBlackQueenside: true,
	}

	if !game.isCastlingMove(Piece{King, White}, fromAlg("e1"), toAlg("c1")) {
		t.Error("e1-c1 should be detected as castling move")
	}

	err := game.validateCastling(board, White, fromAlg("e1"), toAlg("c1"))
	if err != nil {
		t.Errorf("queenside castling should be valid: %v", err)
	}
}

func TestCastling_CannotWhenInCheck(t *testing.T) {
	board, _ := BoardFromFEN("4r3/8/8/8/8/8/8/R3K2R w KQk - 0 1")
	game := &Game{
		Turn:                   White,
		CanCastleWhiteKingside: true,
	}

	err := game.validateCastling(board, White, fromAlg("e1"), toAlg("g1"))
	if err == nil {
		t.Error("cannot castle when king is in check")
	}
}

func TestCastling_CannotPassThroughCheck(t *testing.T) {
	board, _ := BoardFromFEN("5r2/8/8/8/8/8/8/R3K2R w KQk - 0 1")
	game := &Game{
		Turn:                   White,
		CanCastleWhiteKingside: true,
	}

	err := game.validateCastling(board, White, fromAlg("e1"), toAlg("g1"))
	if err == nil {
		t.Error("cannot castle through a square attacked by opponent")
	}
}

func TestCastling_CannotLandOnAttackedSquare(t *testing.T) {
	board, _ := BoardFromFEN("6r1/8/8/8/8/8/8/R3K2R w KQk - 0 1")
	game := &Game{
		Turn:                   White,
		CanCastleWhiteKingside: true,
	}

	err := game.validateCastling(board, White, fromAlg("e1"), toAlg("g1"))
	if err == nil {
		t.Error("cannot castle to a square attacked by opponent")
	}
}

func TestCastling_CannotIfPiecesBetween(t *testing.T) {
	// Bishop on f1 blocks kingside castling: RNBQKB1R = R,N,B,Q,K,B,empty,R
	board, _ := BoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKB1R w KQkq - 0 1")
	game := &Game{
		Turn:                   White,
		CanCastleWhiteKingside: true,
	}

	err := game.validateCastling(board, White, fromAlg("e1"), toAlg("g1"))
	if err == nil {
		t.Error("cannot castle if there are pieces between king and rook")
	}
}

func TestCastling_CannotIfRookMissing(t *testing.T) {
	board, _ := BoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQK3 w Qkq - 0 1")
	game := &Game{
		Turn:                   White,
		CanCastleWhiteKingside: true,
	}

	err := game.validateCastling(board, White, fromAlg("e1"), toAlg("g1"))
	if err == nil {
		t.Error("cannot castle if the rook is missing")
	}
}

func TestCastling_CannotIfRightsLost(t *testing.T) {
	board, _ := BoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w Qkq - 0 1")
	game := &Game{
		Turn:                    White,
		CanCastleWhiteKingside:  false,
		CanCastleWhiteQueenside: true,
	}

	err := game.validateCastling(board, White, fromAlg("e1"), toAlg("g1"))
	if err == nil {
		t.Error("cannot castle kingside if rights were lost")
	}
}

func TestCastling_KingMovesRevokesRights(t *testing.T) {
	wpid := "white-player-id"
	bpid := "black-player-id"
	game := &Game{
		ID:                      "test",
		WhitePlayerID:           &wpid,
		BlackPlayerID:           &bpid,
		Status:                  StatusActive,
		Turn:                    White,
		CurrentFEN:              "rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR w KQkq - 0 1",
		CanCastleWhiteKingside:  true,
		CanCastleWhiteQueenside: true,
		CanCastleBlackKingside:  true,
		CanCastleBlackQueenside: true,
	}
	game.loadFullState()

	// King moves through the e2 gap
	game.MakeMove("white-player-id", fromAlg("e1"), toAlg("e2"))
	game.MakeMove("black-player-id", fromAlg("e7"), toAlg("e6"))
	game.MakeMove("white-player-id", fromAlg("e2"), toAlg("e1"))

	if game.CanCastleWhiteKingside || game.CanCastleWhiteQueenside {
		t.Error("moving the king should revoke all castling rights")
	}
}

func TestCastling_RookMovesRevokesThatSide(t *testing.T) {
	// Clear h2 so the rook can move
	wpid := "white-player-id"
	bpid := "black-player-id"
	game := &Game{
		ID:                      "test",
		WhitePlayerID:           &wpid,
		BlackPlayerID:           &bpid,
		Status:                  StatusActive,
		Turn:                    White,
		CurrentFEN:              "rnbqkbnr/pppppppp/8/8/8/8/PPPPP1P1/RNBQK2R w KQkq - 0 1",
		CanCastleWhiteKingside:  true,
		CanCastleWhiteQueenside: true,
		CanCastleBlackKingside:  true,
		CanCastleBlackQueenside: true,
	}
	game.loadFullState()

	game.MakeMove("white-player-id", fromAlg("h1"), toAlg("h2"))
	game.MakeMove("black-player-id", fromAlg("e7"), toAlg("e6"))

	if game.CanCastleWhiteKingside {
		t.Error("moving the kingside rook should revoke kingside castling")
	}
	if !game.CanCastleWhiteQueenside {
		t.Error("moving the kingside rook should NOT revoke queenside castling")
	}
}

func TestCastling_FullKingsideCastling(t *testing.T) {
	wpid := "white-player-id"
	bpid := "black-player-id"
	game := &Game{
		ID:                      "test",
		WhitePlayerID:           &wpid,
		BlackPlayerID:           &bpid,
		Status:                  StatusActive,
		Turn:                    White,
		CurrentFEN:              "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
		CanCastleWhiteKingside:  true,
		CanCastleWhiteQueenside: true,
		CanCastleBlackKingside:  true,
		CanCastleBlackQueenside: true,
	}
	game.loadFullState()

	err := game.MakeMove("white-player-id", fromAlg("e1"), toAlg("g1"))
	if err != nil {
		t.Fatalf("kingside castling failed: %v", err)
	}

	board, _ := game.GetBoard()
	king := board.PieceAt(toAlg("g1"))
	if king.Type != King || king.Color != White {
		t.Errorf("expected white king on g1, got %v", king)
	}

	rook := board.PieceAt(toAlg("f1"))
	if rook.Type != Rook || rook.Color != White {
		t.Errorf("expected white rook on f1, got %v", rook)
	}

	if !board.IsEmpty(toAlg("e1")) {
		t.Error("expected e1 to be empty after castling")
	}
	if !board.IsEmpty(toAlg("h1")) {
		t.Error("expected h1 to be empty after castling")
	}

	lastMove := game.Moves[len(game.Moves)-1]
	if !lastMove.Castle {
		t.Error("expected last move to be marked as castling")
	}

	// Kingside castling should revoke both white castling rights
	if game.CanCastleWhiteKingside || game.CanCastleWhiteQueenside {
		t.Error("castling should revoke all castling rights for that color")
	}
}

func TestCastling_FullQueensideCastling(t *testing.T) {
	wpid := "white-player-id"
	bpid := "black-player-id"
	game := &Game{
		ID:                      "test",
		WhitePlayerID:           &wpid,
		BlackPlayerID:           &bpid,
		Status:                  StatusActive,
		Turn:                    White,
		CurrentFEN:              "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1",
		CanCastleWhiteKingside:  true,
		CanCastleWhiteQueenside: true,
		CanCastleBlackKingside:  true,
		CanCastleBlackQueenside: true,
	}
	game.loadFullState()

	err := game.MakeMove("white-player-id", fromAlg("e1"), toAlg("c1"))
	if err != nil {
		t.Fatalf("queenside castling failed: %v", err)
	}

	board, _ := game.GetBoard()
	king := board.PieceAt(toAlg("c1"))
	if king.Type != King || king.Color != White {
		t.Errorf("expected white king on c1, got %v", king)
	}

	rook := board.PieceAt(toAlg("d1"))
	if rook.Type != Rook || rook.Color != White {
		t.Errorf("expected white rook on d1, got %v", rook)
	}

	lastMove := game.Moves[len(game.Moves)-1]
	if !lastMove.Castle {
		t.Error("expected last move to be marked as castling")
	}
}

func TestCastling_BlackKingside(t *testing.T) {
	wpid := "white-player-id"
	bpid := "black-player-id"
	game := &Game{
		ID:                      "test",
		WhitePlayerID:           &wpid,
		BlackPlayerID:           &bpid,
		Status:                  StatusActive,
		Turn:                    Black,
		CurrentFEN:              "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w kq - 0 1",
		CanCastleWhiteKingside:  false,
		CanCastleWhiteQueenside: false,
		CanCastleBlackKingside:  true,
		CanCastleBlackQueenside: true,
	}
	game.loadFullState()

	game.MakeMove("white-player-id", fromAlg("e2"), toAlg("e3"))
	game.MakeMove("black-player-id", fromAlg("e8"), toAlg("g8"))

	board, _ := game.GetBoard()
	king := board.PieceAt(toAlg("g8"))
	if king.Type != King || king.Color != Black {
		t.Errorf("expected black king on g8, got %v", king)
	}

	rook := board.PieceAt(toAlg("f8"))
	if rook.Type != Rook || rook.Color != Black {
		t.Errorf("expected black rook on f8, got %v", rook)
	}
}
