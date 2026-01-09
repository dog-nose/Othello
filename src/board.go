package src

// Board は8×8のオセロ盤面を表す型
type Board struct {
	cells [8][8]CellState
}

// NewBoard は新しい盤面を作成する（すべて空）
func NewBoard() *Board {
	return &Board{}
}

// Get は指定座標のマスの状態を取得する
func (b *Board) Get(pos Position) CellState {
	return b.cells[pos.Row][pos.Col]
}

// Set は指定座標のマスに石を設定する
func (b *Board) Set(pos Position, state CellState) {
	b.cells[pos.Row][pos.Col] = state
}

// SetupInitialPosition はゲーム開始時の初期配置を設定する
// d4: 白、e4: 黒、d5: 黒、e5: 白
func (b *Board) SetupInitialPosition() {
	b.Set(Position{Col: 3, Row: 3}, White) // d4
	b.Set(Position{Col: 4, Row: 3}, Black) // e4
	b.Set(Position{Col: 3, Row: 4}, Black) // d5
	b.Set(Position{Col: 4, Row: 4}, White) // e5
}

// IsValidMove は指定位置が有効な手かどうかを判定する
func (b *Board) IsValidMove(pos Position, player CellState) bool {
	// 範囲外チェック
	if !pos.IsValid() {
		return false
	}

	// 既に石がある場合は無効
	if b.Get(pos) != Empty {
		return false
	}

	// いずれかの方向で挟めるか確認
	for _, dir := range GetDirections() {
		if b.CheckDirection(pos, dir, player) {
			return true
		}
	}

	return false
}

// PlaceStone は指定位置に石を置き、挟んだ石をひっくり返す
func (b *Board) PlaceStone(pos Position, player CellState) {
	// 石を置く
	b.Set(pos, player)

	// 各方向で挟める石をひっくり返す
	for _, dir := range GetDirections() {
		if b.CheckDirection(pos, dir, player) {
			b.flipStones(pos, dir, player)
		}
	}
}

// flipStones は指定方向の石をひっくり返す（内部用）
func (b *Board) flipStones(pos Position, dir Direction, player CellState) {
	// 現在位置から指定方向に1つ進む
	nextPos := Position{
		Col: pos.Col + dir.DCol,
		Row: pos.Row + dir.DRow,
	}

	// 相手の石をひっくり返していく
	for {
		cell := b.Get(nextPos)

		// 自分の石に到達したら終了
		if cell == player {
			break
		}

		// 相手の石をひっくり返す
		b.Set(nextPos, player)

		// 次のマスへ
		nextPos.Col += dir.DCol
		nextPos.Row += dir.DRow
	}
}

// GetValidMoves は現在のプレイヤーの有効な手一覧を取得する
func (b *Board) GetValidMoves(player CellState) []Position {
	var moves []Position

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pos := Position{Col: col, Row: row}
			if b.IsValidMove(pos, player) {
				moves = append(moves, pos)
			}
		}
	}

	return moves
}
