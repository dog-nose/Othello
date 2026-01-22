package src

// CellState はマスの状態を表す型
type CellState int

const (
	Empty CellState = iota // 空
	Black                  // 黒
	White                  // 白
)
