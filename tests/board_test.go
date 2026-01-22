package tests

import (
	"othello/src"
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := src.NewBoard()

	// ボードが nil でないことを確認
	if board == nil {
		t.Error("NewBoard() should not return nil")
	}
}

func TestBoardInitialState(t *testing.T) {
	board := src.NewBoard()

	// すべてのマスが空であることを確認
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pos := src.Position{Col: col, Row: row}
			cell := board.Get(pos)
			if cell != src.Empty {
				t.Errorf("Expected Empty at %v, got %v", pos, cell)
			}
		}
	}
}

func TestBoardSetupInitialPosition(t *testing.T) {
	board := src.NewBoard()
	board.SetupInitialPosition()

	// d4: 白、e4: 黒、d5: 黒、e5: 白
	tests := []struct {
		pos  src.Position
		want src.CellState
	}{
		{src.Position{Col: 3, Row: 3}, src.White}, // d4
		{src.Position{Col: 4, Row: 3}, src.Black}, // e4
		{src.Position{Col: 3, Row: 4}, src.Black}, // d5
		{src.Position{Col: 4, Row: 4}, src.White}, // e5
	}

	for _, tt := range tests {
		got := board.Get(tt.pos)
		if got != tt.want {
			t.Errorf("At %v: got %v, want %v", tt.pos, got, tt.want)
		}
	}

	// 他のマスは空であることを確認
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pos := src.Position{Col: col, Row: row}
			// 初期配置の4マスはスキップ
			if (col == 3 || col == 4) && (row == 3 || row == 4) {
				continue
			}
			if board.Get(pos) != src.Empty {
				t.Errorf("Expected Empty at %v, got %v", pos, board.Get(pos))
			}
		}
	}
}
