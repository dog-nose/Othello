package src

// Game はオセロゲームの状態を管理する型
type Game struct {
	board         *Board
	currentPlayer CellState
}

// NewGame は新しいゲームを作成する
func NewGame() *Game {
	board := NewBoard()
	board.SetupInitialPosition()
	return &Game{
		board:         board,
		currentPlayer: Black, // 黒が先手
	}
}

// CurrentPlayer は現在の手番のプレイヤーを返す
func (g *Game) CurrentPlayer() CellState {
	return g.currentPlayer
}

// SwitchPlayer は手番を交代する
func (g *Game) SwitchPlayer() {
	if g.currentPlayer == Black {
		g.currentPlayer = White
	} else {
		g.currentPlayer = Black
	}
}

// CountStones は黒石と白石の数をカウントする
func (g *Game) CountStones() (blackCount, whiteCount int) {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			pos := Position{Col: col, Row: row}
			cell := g.board.Get(pos)
			if cell == Black {
				blackCount++
			} else if cell == White {
				whiteCount++
			}
		}
	}
	return blackCount, whiteCount
}

// IsGameOver はゲームが終了したかどうかを判定する
func (g *Game) IsGameOver() bool {
	// 両者とも手を打てない場合、ゲーム終了
	return !g.CanPlayerMove(Black) && !g.CanPlayerMove(White)
}

// GetWinner は勝者を返す（引き分けの場合は Empty）
func (g *Game) GetWinner() CellState {
	blackCount, whiteCount := g.CountStones()
	if blackCount > whiteCount {
		return Black
	} else if whiteCount > blackCount {
		return White
	}
	return Empty // 引き分け
}

// CanPlayerMove は指定プレイヤーが手を打てるかどうかを判定する
func (g *Game) CanPlayerMove(player CellState) bool {
	moves := g.board.GetValidMoves(player)
	return len(moves) > 0
}

// PlaceStone は現在のプレイヤーが指定位置に石を置く
func (g *Game) PlaceStone(pos Position) {
	g.board.PlaceStone(pos, g.currentPlayer)
}

// GetBoard は盤面を返す（テスト用）
func (g *Game) GetBoard() *Board {
	return g.board
}
