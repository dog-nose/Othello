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
