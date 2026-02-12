import { describe, it, expect } from 'vitest';
import {
  createInitialBoard,
  getValidMoves,
  placeStone,
  countStones,
  isGameOver,
  getWinner,
  isValidMove,
  getFlippableStones,
} from './board';

describe('createInitialBoard', () => {
  it('should create an 8x8 board', () => {
    const board = createInitialBoard();
    expect(board.length).toBe(8);
    board.forEach(row => expect(row.length).toBe(8));
  });

  it('should have initial 4 stones in the center', () => {
    const board = createInitialBoard();
    expect(board[3][3]).toBe('white');
    expect(board[3][4]).toBe('black');
    expect(board[4][3]).toBe('black');
    expect(board[4][4]).toBe('white');
  });

  it('should have all other cells empty', () => {
    const board = createInitialBoard();
    const { black, white } = countStones(board);
    expect(black).toBe(2);
    expect(white).toBe(2);
  });
});

describe('getValidMoves', () => {
  it('should return 4 valid moves for black at game start', () => {
    const board = createInitialBoard();
    const moves = getValidMoves(board, 'black');
    expect(moves.length).toBe(4);

    const moveSet = new Set(moves.map(m => `${m.row},${m.col}`));
    expect(moveSet.has('2,3')).toBe(true);
    expect(moveSet.has('3,2')).toBe(true);
    expect(moveSet.has('4,5')).toBe(true);
    expect(moveSet.has('5,4')).toBe(true);
  });

  it('should return 4 valid moves for white at game start', () => {
    const board = createInitialBoard();
    const moves = getValidMoves(board, 'white');
    expect(moves.length).toBe(4);
  });
});

describe('isValidMove', () => {
  it('should return true for a valid move', () => {
    const board = createInitialBoard();
    expect(isValidMove(board, 2, 3, 'black')).toBe(true);
  });

  it('should return false for an occupied cell', () => {
    const board = createInitialBoard();
    expect(isValidMove(board, 3, 3, 'black')).toBe(false);
  });

  it('should return false for a cell that flips no stones', () => {
    const board = createInitialBoard();
    expect(isValidMove(board, 0, 0, 'black')).toBe(false);
  });
});

describe('getFlippableStones', () => {
  it('should return flippable stones for a valid move', () => {
    const board = createInitialBoard();
    const flippable = getFlippableStones(board, 2, 3, 'black');
    expect(flippable.length).toBe(1);
    expect(flippable[0]).toEqual({ row: 3, col: 3 });
  });
});

describe('placeStone', () => {
  it('should place a stone and flip opponent stones', () => {
    const board = createInitialBoard();
    const newBoard = placeStone(board, 2, 3, 'black');

    expect(newBoard[2][3]).toBe('black');
    expect(newBoard[3][3]).toBe('black');
    expect(newBoard[3][4]).toBe('black');

    const { black, white } = countStones(newBoard);
    expect(black).toBe(4);
    expect(white).toBe(1);
  });

  it('should not modify the original board', () => {
    const board = createInitialBoard();
    placeStone(board, 2, 3, 'black');
    expect(board[2][3]).toBeNull();
    expect(board[3][3]).toBe('white');
  });

  it('should return the same board if the move is invalid', () => {
    const board = createInitialBoard();
    const result = placeStone(board, 0, 0, 'black');
    expect(result).toBe(board);
  });
});

describe('countStones', () => {
  it('should count initial stones correctly', () => {
    const board = createInitialBoard();
    const { black, white } = countStones(board);
    expect(black).toBe(2);
    expect(white).toBe(2);
  });
});

describe('isGameOver', () => {
  it('should return false at game start', () => {
    const board = createInitialBoard();
    expect(isGameOver(board)).toBe(false);
  });

  it('should return true when no moves available for either player', () => {
    const board = createInitialBoard().map(row => row.map(() => null as ReturnType<typeof createInitialBoard>[0][0]));
    // Fill board completely with black except one white corner
    for (let r = 0; r < 8; r++) {
      for (let c = 0; c < 8; c++) {
        board[r][c] = 'black';
      }
    }
    expect(isGameOver(board)).toBe(true);
  });
});

describe('getWinner', () => {
  it('should return black when black has more stones', () => {
    const board = createInitialBoard();
    const newBoard = placeStone(board, 2, 3, 'black');
    expect(getWinner(newBoard)).toBe('black');
  });

  it('should return draw when counts are equal', () => {
    const board = createInitialBoard();
    expect(getWinner(board)).toBe('draw');
  });
});
