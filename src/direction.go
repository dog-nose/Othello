package src

// Direction は盤面上の方向を表す型
type Direction struct {
	DCol int // 列の変位 (-1, 0, 1)
	DRow int // 行の変位 (-1, 0, 1)
}

// GetDirections は8方向のスライスを返す
func GetDirections() []Direction {
	return []Direction{
		{DCol: -1, DRow: -1}, // 左上
		{DCol: 0, DRow: -1},  // 上
		{DCol: 1, DRow: -1},  // 右上
		{DCol: -1, DRow: 0},  // 左
		{DCol: 1, DRow: 0},   // 右
		{DCol: -1, DRow: 1},  // 左下
		{DCol: 0, DRow: 1},   // 下
		{DCol: 1, DRow: 1},   // 右下
	}
}

// CheckDirection は指定方向に挟める石があるか判定する
func (b *Board) CheckDirection(pos Position, dir Direction, player CellState) bool {
	// 現在位置から指定方向に1つ進む
	nextPos := Position{
		Col: pos.Col + dir.DCol,
		Row: pos.Row + dir.DRow,
	}

	// 範囲外チェック
	if !nextPos.IsValid() {
		return false
	}

	// 最初のマスが相手の石でなければならない
	opponent := White
	if player == White {
		opponent = Black
	}

	if b.Get(nextPos) != opponent {
		return false
	}

	// 相手の石が続く間、進む
	for {
		nextPos.Col += dir.DCol
		nextPos.Row += dir.DRow

		// 範囲外チェック
		if !nextPos.IsValid() {
			return false
		}

		cell := b.Get(nextPos)

		// 空きマスなら挟めない
		if cell == Empty {
			return false
		}

		// 自分の石なら挟める
		if cell == player {
			return true
		}

		// 相手の石なら続ける
	}
}
