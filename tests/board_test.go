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
