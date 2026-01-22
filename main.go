package main

import (
	"bufio"
	"fmt"
	"os"
	"othello/src"
	"strings"
)

func main() {
	game := src.NewGame()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("オセロゲームを開始します！")
	fmt.Println()

	for !game.IsGameOver() {
		// 盤面を表示
		fmt.Println(game.GetBoard().String())

		// 石の数を表示
		blackCount, whiteCount := game.CountStones()
		fmt.Printf("黒: %d, 白: %d\n", blackCount, whiteCount)
		fmt.Println()

		// 現在のプレイヤー
		currentPlayer := game.CurrentPlayer()
		playerName := "黒"
		if currentPlayer == src.White {
			playerName = "白"
		}

		// 有効な手をチェック
		if !game.CanPlayerMove(currentPlayer) {
			fmt.Printf("%sはパスします。\n", playerName)
			game.SwitchPlayer()
			continue
		}

		// 有効な手を表示
		validMoves := game.GetBoard().GetValidMoves(currentPlayer)
		fmt.Print("有効な手: ")
		for i, move := range validMoves {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(src.FormatPosition(move))
		}
		fmt.Println()

		// 入力を受け付ける
		fmt.Printf("%sの手番です。位置を入力してください (例: d3): ", playerName)
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// 入力をパース
		pos, err := src.ParsePosition(input)
		if err != nil {
			fmt.Println("無効な入力です。もう一度入力してください。")
			continue
		}

		// 有効な手かチェック
		if !game.GetBoard().IsValidMove(pos, currentPlayer) {
			fmt.Println("その位置には置けません。もう一度入力してください。")
			continue
		}

		// 石を置く
		game.PlaceStone(pos)

		// 手番を交代
		game.SwitchPlayer()

		fmt.Println()
	}

	// ゲーム終了
	fmt.Println("ゲーム終了！")
	fmt.Println(game.GetBoard().String())

	// 最終結果を表示
	blackCount, whiteCount := game.CountStones()
	fmt.Printf("黒: %d, 白: %d\n", blackCount, whiteCount)

	winner := game.GetWinner()
	switch winner {
	case src.Black:
		fmt.Println("黒の勝ちです！")
	case src.White:
		fmt.Println("白の勝ちです！")
	case src.Empty:
		fmt.Println("引き分けです！")
	}
}
