package tests

import (
	"othello/src"
	"testing"
)

func TestIsValidMove(t *testing.T) {
	board := src.NewBoard()
	board.SetupInitialPosition()

	// 初期配置: d4: 白、e4: 黒、d5: 黒、e5: 白
	// 黒の有効な手: d3, c4, f5, e6

	tests := []struct {
		pos    src.Position
		player src.CellState
		valid  bool
	}{
		{src.Position{Col: 3, Row: 2}, src.Black, true},  // d3 (黒の有効手)
		{src.Position{Col: 2, Row: 3}, src.Black, true},  // c4 (黒の有効手)
		{src.Position{Col: 5, Row: 4}, src.Black, true},  // f5 (黒の有効手)
		{src.Position{Col: 4, Row: 5}, src.Black, true},  // e6 (黒の有効手)
		{src.Position{Col: 0, Row: 0}, src.Black, false}, // a1 (無効)
		{src.Position{Col: 3, Row: 3}, src.Black, false}, // d4 (既に石がある)
	}

	for _, tt := range tests {
		got := board.IsValidMove(tt.pos, tt.player)
		if got != tt.valid {
			t.Errorf("IsValidMove(%v, %v) = %v, want %v",
				tt.pos, tt.player, got, tt.valid)
		}
	}
}

func TestPlaceStone(t *testing.T) {
	board := src.NewBoard()
	board.SetupInitialPosition()

	// 黒がd3に置く
	pos := src.Position{Col: 3, Row: 2} // d3
	board.PlaceStone(pos, src.Black)

	// d3に黒石があることを確認
	if board.Get(pos) != src.Black {
		t.Errorf("Expected Black at %v, got %v", pos, board.Get(pos))
	}

	// d4の白石が黒にひっくり返されたことを確認
	d4 := src.Position{Col: 3, Row: 3}
	if board.Get(d4) != src.Black {
		t.Errorf("Expected Black at d4 after flipping, got %v", board.Get(d4))
	}

	// d5は元々黒なのでそのまま
	d5 := src.Position{Col: 3, Row: 4}
	if board.Get(d5) != src.Black {
		t.Errorf("Expected Black at d5, got %v", board.Get(d5))
	}
}

func TestGetValidMoves(t *testing.T) {
	board := src.NewBoard()
	board.SetupInitialPosition()

	// 黒の有効な手を取得
	moves := board.GetValidMoves(src.Black)

	// 初期配置では黒は4つの有効な手がある
	if len(moves) != 4 {
		t.Errorf("Expected 4 valid moves for Black, got %d", len(moves))
	}

	// 有効な手が含まれていることを確認
	expectedMoves := []src.Position{
		{Col: 3, Row: 2}, // d3
		{Col: 2, Row: 3}, // c4
		{Col: 5, Row: 4}, // f5
		{Col: 4, Row: 5}, // e6
	}

	for _, expected := range expectedMoves {
		found := false
		for _, move := range moves {
			if move == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected move %v not found in valid moves", expected)
		}
	}
}
