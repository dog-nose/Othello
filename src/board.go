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
