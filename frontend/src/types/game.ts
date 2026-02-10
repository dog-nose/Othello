export type Color = 'black' | 'white';
export type CellState = Color | null;
export type Board = CellState[][];

export interface Position {
  row: number;
  col: number;
}

export interface GameState {
  board: Board;
  currentPlayer: Color;
  isGameOver: boolean;
  playId: string | null;
  blackCount: number;
  whiteCount: number;
  winner: Color | 'draw' | null;
  validMoves: Position[];
  lastMove: Position | null;
}
