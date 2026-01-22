package src

import "fmt"

// String は盤面を文字列として表示する
func (b *Board) String() string {
	result := "    a b c d e f g h\n"
	result += "  ┌─┬─┬─┬─┬─┬─┬─┬─┐\n"

	for row := 0; row < 8; row++ {
		result += fmt.Sprintf("%d │", row+1)

		for col := 0; col < 8; col++ {
			pos := Position{Col: col, Row: row}
			cell := b.Get(pos)

			switch cell {
			case Black:
				result += "●"
			case White:
				result += "○"
			case Empty:
				result += " "
			}

			if col < 7 {
				result += "│"
			}
		}

		result += "│\n"

		if row < 7 {
			result += "  ├─┼─┼─┼─┼─┼─┼─┼─┤\n"
		}
	}

	result += "  └─┴─┴─┴─┴─┴─┴─┴─┘\n"
	return result
}

// FormatPosition は座標を文字列（例: "d4"）に変換する
func FormatPosition(pos Position) string {
	col := string(rune('a' + pos.Col))
	row := fmt.Sprintf("%d", pos.Row+1)
	return col + row
}
