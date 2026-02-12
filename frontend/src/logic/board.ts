import type { Board, CellState, Color, Position } from '../types/game';

const DIRECTIONS: [number, number][] = [
  [-1, -1], [-1, 0], [-1, 1],
  [0, -1],           [0, 1],
  [1, -1],  [1, 0],  [1, 1],
];

export function createInitialBoard(): Board {
  const board: Board = Array.from({ length: 8 }, () =>
    Array.from<CellState>({ length: 8 }).fill(null)
  );
  board[3][3] = 'white';
  board[3][4] = 'black';
  board[4][3] = 'black';
  board[4][4] = 'white';
  return board;
}

function opponent(color: Color): Color {
  return color === 'black' ? 'white' : 'black';
}

function getFlippableInDirection(
  board: Board,
  row: number,
  col: number,
  dr: number,
  dc: number,
  color: Color
): Position[] {
  const opp = opponent(color);
  const flippable: Position[] = [];
  let r = row + dr;
  let c = col + dc;

  while (r >= 0 && r < 8 && c >= 0 && c < 8 && board[r][c] === opp) {
    flippable.push({ row: r, col: c });
    r += dr;
    c += dc;
  }

  if (flippable.length > 0 && r >= 0 && r < 8 && c >= 0 && c < 8 && board[r][c] === color) {
    return flippable;
  }
  return [];
}

export function getFlippableStones(board: Board, row: number, col: number, color: Color): Position[] {
  if (board[row][col] !== null) {
    return [];
  }

  const allFlippable: Position[] = [];
  for (const [dr, dc] of DIRECTIONS) {
    allFlippable.push(...getFlippableInDirection(board, row, col, dr, dc, color));
  }
  return allFlippable;
}

export function isValidMove(board: Board, row: number, col: number, color: Color): boolean {
  return getFlippableStones(board, row, col, color).length > 0;
}

export function getValidMoves(board: Board, color: Color): Position[] {
  const moves: Position[] = [];
  for (let r = 0; r < 8; r++) {
    for (let c = 0; c < 8; c++) {
      if (isValidMove(board, r, c, color)) {
        moves.push({ row: r, col: c });
      }
    }
  }
  return moves;
}

export function placeStone(board: Board, row: number, col: number, color: Color): Board {
  const flippable = getFlippableStones(board, row, col, color);
  if (flippable.length === 0) {
    return board;
  }

  const newBoard = board.map(r => [...r]);
  newBoard[row][col] = color;
  for (const pos of flippable) {
    newBoard[pos.row][pos.col] = color;
  }
  return newBoard;
}

export function countStones(board: Board): { black: number; white: number } {
  let black = 0;
  let white = 0;
  for (let r = 0; r < 8; r++) {
    for (let c = 0; c < 8; c++) {
      if (board[r][c] === 'black') black++;
      else if (board[r][c] === 'white') white++;
    }
  }
  return { black, white };
}

export function isGameOver(board: Board): boolean {
  return getValidMoves(board, 'black').length === 0 && getValidMoves(board, 'white').length === 0;
}

export function getWinner(board: Board): Color | 'draw' {
  const { black, white } = countStones(board);
  if (black > white) return 'black';
  if (white > black) return 'white';
  return 'draw';
}
