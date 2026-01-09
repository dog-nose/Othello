package tests

import (
	"othello/src"
	"testing"
)

func TestNewGame(t *testing.T) {
	game := src.NewGame()

	// ゲームが正しく初期化されていることを確認
	if game == nil {
		t.Error("NewGame() should not return nil")
	}

	// 黒が先手
	if game.CurrentPlayer() != src.Black {
		t.Error("Black should be the first player")
	}
}

func TestSwitchPlayer(t *testing.T) {
	game := src.NewGame()

	// 初期状態は黒
	if game.CurrentPlayer() != src.Black {
		t.Error("Initial player should be Black")
	}

	// 手番を交代
	game.SwitchPlayer()

	// 白に変わる
	if game.CurrentPlayer() != src.White {
		t.Error("Player should switch to White")
	}

	// もう一度交代
	game.SwitchPlayer()

	// 黒に戻る
	if game.CurrentPlayer() != src.Black {
		t.Error("Player should switch back to Black")
	}
}

func TestCountStones(t *testing.T) {
	game := src.NewGame()

	// 初期配置では黒2、白2
	blackCount, whiteCount := game.CountStones()
	if blackCount != 2 || whiteCount != 2 {
		t.Errorf("Initial count should be Black:2, White:2, got Black:%d, White:%d",
			blackCount, whiteCount)
	}

	// 黒がd3に置く
	pos := src.Position{Col: 3, Row: 2}
	game.PlaceStone(pos)

	// 黒4、白1になる
	blackCount, whiteCount = game.CountStones()
	if blackCount != 4 || whiteCount != 1 {
		t.Errorf("After move, count should be Black:4, White:1, got Black:%d, White:%d",
			blackCount, whiteCount)
	}
}

func TestIsGameOver(t *testing.T) {
	game := src.NewGame()

	// 初期状態ではゲームは終了していない
	if game.IsGameOver() {
		t.Error("Game should not be over at start")
	}
}

func TestGetWinner(t *testing.T) {
	// ゲームを新規作成して、特定の状態を作る
	game := src.NewGame()

	// 初期状態では引き分け（2:2）
	winner := game.GetWinner()
	if winner != src.Empty {
		t.Errorf("Initial state should be a draw, got %v", winner)
	}

	// 黒がd3に置く
	pos := src.Position{Col: 3, Row: 2}
	game.PlaceStone(pos)

	// 黒4、白1なので黒の勝ち
	winner = game.GetWinner()
	if winner != src.Black {
		t.Errorf("Black should be winning, got %v", winner)
	}
}

func TestCanPlayerMove(t *testing.T) {
	game := src.NewGame()

	// 初期状態では黒は手を打てる
	if !game.CanPlayerMove(src.Black) {
		t.Error("Black should be able to move at start")
	}

	// 白も手を打てる
	if !game.CanPlayerMove(src.White) {
		t.Error("White should be able to move at start")
	}
}
