package tests

import (
	"othello/src"
	"strings"
	"testing"
)

func TestBoardString(t *testing.T) {
	board := src.NewBoard()
	board.SetupInitialPosition()

	// 盤面を文字列に変換
	str := board.String()

	// 文字列が空でないことを確認
	if str == "" {
		t.Error("Board.String() should not be empty")
	}

	// 列のラベル (a-h) が含まれていることを確認
	if !strings.Contains(str, "a") || !strings.Contains(str, "h") {
		t.Error("Board.String() should contain column labels a-h")
	}

	// 行のラベル (1-8) が含まれていることを確認
	if !strings.Contains(str, "1") || !strings.Contains(str, "8") {
		t.Error("Board.String() should contain row labels 1-8")
	}

	// 黒石（●）と白石（○）が含まれていることを確認
	if !strings.Contains(str, "●") || !strings.Contains(str, "○") {
		t.Error("Board.String() should contain black (●) and white (○) stones")
	}
}

func TestFormatPosition(t *testing.T) {
	tests := []struct {
		pos  src.Position
		want string
	}{
		{src.Position{Col: 0, Row: 0}, "a1"},
		{src.Position{Col: 3, Row: 3}, "d4"},
		{src.Position{Col: 7, Row: 7}, "h8"},
		{src.Position{Col: 4, Row: 4}, "e5"},
	}

	for _, tt := range tests {
		got := src.FormatPosition(tt.pos)
		if got != tt.want {
			t.Errorf("FormatPosition(%v) = %s, want %s", tt.pos, got, tt.want)
		}
	}
}
