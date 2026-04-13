package chess

type PieceType string

const (
	Pawn   PieceType = "pawn"
	Rook   PieceType = "rook"
	Knight PieceType = "knight"
	Bishop PieceType = "bishop"
	Queen  PieceType = "queen"
	King   PieceType = "king"
)

type Piece struct {
	Type  PieceType
	Color Color
}

func (p Piece) IsEmpty() bool {
	return p.Type == ""
}

func (p Piece) String() string {
	if p.IsEmpty() {
		return "."
	}
	var c byte = ' '
	if p.Color == White {
		c = 'w'
	} else {
		c = 'b'
	}
	switch p.Type {
	case Pawn:
		return string(c) + "P"
	case Rook:
		return string(c) + "R"
	case Knight:
		return string(c) + "N"
	case Bishop:
		return string(c) + "B"
	case Queen:
		return string(c) + "Q"
	case King:
		return string(c) + "K"
	}
	return "?"
}
