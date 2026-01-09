package tests

import (
	"othello/src"
	"testing"
)

func TestPosition(t *testing.T) {
	// 座標を表す型が定義されていることをテスト
	// 列: a-h (0-7), 行: 1-8 (0-7)
	pos := src.Position{Col: 3, Row: 3} // d4

	if pos.Col < 0 || pos.Col > 7 {
		t.Error("Col should be in range 0-7")
	}

	if pos.Row < 0 || pos.Row > 7 {
		t.Error("Row should be in range 0-7")
	}
}

func TestPositionIsValid(t *testing.T) {
	tests := []struct {
		pos   src.Position
		valid bool
	}{
		{src.Position{Col: 0, Row: 0}, true},   // a1
		{src.Position{Col: 7, Row: 7}, true},   // h8
		{src.Position{Col: 3, Row: 3}, true},   // d4
		{src.Position{Col: -1, Row: 0}, false}, // 範囲外
		{src.Position{Col: 0, Row: -1}, false}, // 範囲外
		{src.Position{Col: 8, Row: 0}, false},  // 範囲外
		{src.Position{Col: 0, Row: 8}, false},  // 範囲外
	}

	for _, tt := range tests {
		if got := tt.pos.IsValid(); got != tt.valid {
			t.Errorf("Position{%d, %d}.IsValid() = %v, want %v",
				tt.pos.Col, tt.pos.Row, got, tt.valid)
		}
	}
}

func TestParsePosition(t *testing.T) {
	tests := []struct {
		input   string
		want    src.Position
		wantErr bool
	}{
		{"a1", src.Position{Col: 0, Row: 0}, false},
		{"d4", src.Position{Col: 3, Row: 3}, false},
		{"h8", src.Position{Col: 7, Row: 7}, false},
		{"e5", src.Position{Col: 4, Row: 4}, false},
		{"a9", src.Position{}, true},  // 範囲外
		{"i1", src.Position{}, true},  // 範囲外
		{"a0", src.Position{}, true},  // 範囲外
		{"", src.Position{}, true},    // 空文字列
		{"abc", src.Position{}, true}, // 不正な形式
		{"1a", src.Position{}, true},  // 順序が逆
	}

	for _, tt := range tests {
		got, err := src.ParsePosition(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParsePosition(%q) error = %v, wantErr %v",
				tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("ParsePosition(%q) = %v, want %v",
				tt.input, got, tt.want)
		}
	}
}
