package src

// Position は盤面上の座標を表す型
type Position struct {
	Col int // 列 (0-7: a-h)
	Row int // 行 (0-7: 1-8)
}

// IsValid は座標が有効範囲内かどうかを判定する
func (p Position) IsValid() bool {
	return p.Col >= 0 && p.Col < 8 && p.Row >= 0 && p.Row < 8
}
