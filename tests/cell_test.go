package tests

import (
	"othello/src"
	"testing"
)

func TestCellState(t *testing.T) {
	// マスの状態を表す型が定義されていることをテスト
	var empty src.CellState = src.Empty
	var black src.CellState = src.Black
	var white src.CellState = src.White

	if empty == black || empty == white || black == white {
		t.Error("CellState values should be distinct")
	}
}
