package chess

import "time"

type Move struct {
	From      Position
	To        Position
	PlayerID  string
	Timestamp time.Time
	Promotion *PieceType
	Castle    bool
	EnPassant bool
}

func (m Move) String() string {
	return m.From.String() + m.To.String()
}
