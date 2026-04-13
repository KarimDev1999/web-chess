package chess

import (
	"fmt"
	"regexp"
)

const (
	BoardSize = 8
	BoardMin  = 0
	BoardMax  = BoardSize - 1
)

type Position struct {
	Row int
	Col int
}

func (p Position) IsValid() bool {
	return p.Row >= BoardMin && p.Row < BoardSize && p.Col >= BoardMin && p.Col < BoardSize
}

func (p Position) String() string {
	if !p.IsValid() {
		return "invalid"
	}
	return fmt.Sprintf("%c%d", 'a'+rune(p.Col), 8-p.Row)
}

func ParseAlgebraic(s string) (Position, error) {
	if len(s) != 2 {
		return Position{}, fmt.Errorf("invalid position format")
	}
	re := regexp.MustCompile(`^[a-h][1-8]$`)
	if !re.MatchString(s) {
		return Position{}, fmt.Errorf("invalid position format")
	}
	col := int(s[0] - 'a')
	row := 8 - int(s[1]-'0')
	return Position{Row: row, Col: col}, nil
}
