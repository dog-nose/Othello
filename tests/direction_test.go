package tests

import (
	"othello/src"
	"testing"
)

func TestDirections(t *testing.T) {
	// 8方向が定義されていることを確認
	directions := src.GetDirections()

	if len(directions) != 8 {
		t.Errorf("Expected 8 directions, got %d", len(directions))
	}

	// 各方向が有効な変位であることを確認
	for _, dir := range directions {
		if dir.DCol == 0 && dir.DRow == 0 {
			t.Error("Direction should not be (0, 0)")
		}
		if dir.DCol < -1 || dir.DCol > 1 {
			t.Errorf("DCol should be -1, 0, or 1, got %d", dir.DCol)
		}
		if dir.DRow < -1 || dir.DRow > 1 {
			t.Errorf("DRow should be -1, 0, or 1, got %d", dir.DRow)
		}
	}
}

func TestCheckDirection(t *testing.T) {
	board := src.NewBoard()
	board.SetupInitialPosition()

	// d4: 白、e4: 黒、d5: 黒、e5: 白

	// 黒がd3に置く場合、下方向（d4の白）を挟める
	pos := src.Position{Col: 3, Row: 2}    // d3
	dir := src.Direction{DCol: 0, DRow: 1} // 下

	canFlip := board.CheckDirection(pos, dir, src.Black)
	if !canFlip {
		t.Error("Should be able to flip in down direction")
	}

	// 黒がa1に置く場合、どの方向にも挟めない
	pos2 := src.Position{Col: 0, Row: 0} // a1
	canFlip2 := board.CheckDirection(pos2, dir, src.Black)
	if canFlip2 {
		t.Error("Should not be able to flip from a1")
	}
}
