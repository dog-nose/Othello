package src

import "fmt"

// Position は盤面上の座標を表す型
type Position struct {
	Col int // 列 (0-7: a-h)
	Row int // 行 (0-7: 1-8)
}

// IsValid は座標が有効範囲内かどうかを判定する
func (p Position) IsValid() bool {
	return p.Col >= 0 && p.Col < 8 && p.Row >= 0 && p.Row < 8
}

// ParsePosition は文字列（例: "d4"）を座標に変換する
func ParsePosition(s string) (Position, error) {
	if len(s) != 2 {
		return Position{}, fmt.Errorf("invalid position format: %s", s)
	}

	col := int(s[0] - 'a')
	row := int(s[1] - '1')

	pos := Position{Col: col, Row: row}
	if !pos.IsValid() {
		return Position{}, fmt.Errorf("position out of range: %s", s)
	}

	return pos, nil
}
