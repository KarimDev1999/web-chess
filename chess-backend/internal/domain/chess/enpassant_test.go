package chess

import "testing"

func TestEnPassant_TargetSetAfterPawnDoublePush(t *testing.T) {
	game := NewGame("white-player-id", TimeControl{}, PreferenceWhite)
	game.Join("black-player-id")

	game.MakeMove("white-player-id", fromAlg("e2"), toAlg("e4"))

	if game.EnPassantTarget == nil {
		t.Fatal("expected en-passant target after e2-e4")
	}
	if *game.EnPassantTarget != toAlg("e3") {
		t.Errorf("expected en-passant target e3, got %s", game.EnPassantTarget.String())
	}
}

func TestEnPassant_CaptureBlackPawn(t *testing.T) {
	// White pawn on e4, black pawn on d4. White just pushed e2-e4, ep target e3.
	// Black captures: d4 → e3
	wpid := "white-player-id"
	bpid := "black-player-id"
	game := &Game{
		ID:                      "test",
		WhitePlayerID:           &wpid,
		BlackPlayerID:           &bpid,
		Status:                  StatusActive,
		Turn:                    Black,
		Moves:                   []Move{},
		CurrentFEN:              "rnbqkbnr/pppp1ppp/8/8/3pP3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
		CanCastleWhiteKingside:  true,
		CanCastleWhiteQueenside: true,
		CanCastleBlackKingside:  true,
		CanCastleBlackQueenside: true,
	}
	game.loadFullState()

	err := game.MakeMove("black-player-id", fromAlg("d4"), toAlg("e3"))
	if err != nil {
		t.Fatalf("en-passant capture failed: %v", err)
	}

	board, _ := game.GetBoard()
	// Black pawn should be on e3
	pawn := board.PieceAt(toAlg("e3"))
	if pawn.Type != Pawn || pawn.Color != Black {
		t.Errorf("expected black pawn on e3, got %v", pawn)
	}
	// White pawn on e4 should be REMOVED (en-passant capture)
	e4Pawn := board.PieceAt(toAlg("e4"))
	if !e4Pawn.IsEmpty() {
		t.Errorf("expected e4 to be empty after en-passant capture, got %v", e4Pawn)
	}
	// d4 should be empty (the capturing pawn moved)
	d4Pawn := board.PieceAt(toAlg("d4"))
	if !d4Pawn.IsEmpty() {
		t.Errorf("expected d4 to be empty, got %v", d4Pawn)
	}

	lastMove := game.Moves[len(game.Moves)-1]
	if !lastMove.EnPassant {
		t.Error("expected last move to be marked as en-passant")
	}
}

func TestEnPassant_CaptureWhitePawn(t *testing.T) {
	// Black pawn on d5, white pawn on e5. Black just pushed d7-d5, ep target d6.
	// White captures: e5 → d6

	wpid := "white-player-id"
	bpid := "black-player-id"
	game := &Game{
		ID:                      "test",
		WhitePlayerID:           &wpid,
		BlackPlayerID:           &bpid,
		Status:                  StatusActive,
		Turn:                    White,
		Moves:                   []Move{},
		CurrentFEN:              "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 1",
		CanCastleWhiteKingside:  true,
		CanCastleWhiteQueenside: true,
		CanCastleBlackKingside:  true,
		CanCastleBlackQueenside: true,
	}
	game.loadFullState()

	err := game.MakeMove("white-player-id", fromAlg("e5"), toAlg("d6"))
	if err != nil {
		t.Fatalf("en-passant capture failed: %v", err)
	}

	board, _ := game.GetBoard()
	pawn := board.PieceAt(toAlg("d6"))
	if pawn.Type != Pawn || pawn.Color != White {
		t.Errorf("expected white pawn on d6, got %v", pawn)
	}

	d5Pawn := board.PieceAt(toAlg("d5"))
	if !d5Pawn.IsEmpty() {
		t.Errorf("expected d5 to be empty after en-passant capture, got %v", d5Pawn)
	}

	lastMove := game.Moves[len(game.Moves)-1]
	if !lastMove.EnPassant {
		t.Error("expected last move to be marked as en-passant")
	}
}

func TestEnPassant_OnlyAvailableImmediately(t *testing.T) {
	game := NewGame("white-player-id", TimeControl{}, PreferenceWhite)
	game.Join("black-player-id")

	game.MakeMove("white-player-id", fromAlg("d2"), toAlg("d4"))
	if game.EnPassantTarget == nil {
		t.Fatal("expected en-passant target")
	}

	// Black makes a different move
	game.MakeMove("black-player-id", fromAlg("e7"), toAlg("e6"))

	if game.EnPassantTarget != nil {
		t.Error("en-passant target should be cleared after one ply")
	}
}

func TestFENRoundTrip(t *testing.T) {
	game := NewGame("white-player-id", TimeControl{}, PreferenceWhite)
	game.Join("black-player-id")

	game.MakeMove("white-player-id", fromAlg("e2"), toAlg("e4"))

	fen := game.CurrentFEN
	board, turn, castling, ep, err := ParseFullFEN(fen)
	if err != nil {
		t.Fatalf("failed to parse FEN: %v", err)
	}

	if turn != Black {
		t.Errorf("expected turn Black, got %v", turn)
	}
	if !castling.WhiteKingside || !castling.WhiteQueenside || !castling.BlackKingside || !castling.BlackQueenside {
		t.Errorf("expected all castling rights, got %+v", castling)
	}
	if ep == nil {
		t.Error("expected en-passant target")
	}
	if ep != nil && *ep != toAlg("e3") {
		t.Errorf("expected en-passant e3, got %s", ep.String())
	}

	pawn := board.PieceAt(toAlg("e4"))
	if pawn.Type != Pawn || pawn.Color != White {
		t.Errorf("expected white pawn on e4, got %v", pawn)
	}
}
