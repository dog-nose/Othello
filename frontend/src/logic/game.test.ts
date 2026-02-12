import { describe, it, expect } from 'vitest';
import { createInitialGameState, handleMove, applyOpponentMove } from './game';

describe('createInitialGameState', () => {
  it('should start with black as current player', () => {
    const state = createInitialGameState();
    expect(state.currentPlayer).toBe('black');
  });

  it('should not be game over', () => {
    const state = createInitialGameState();
    expect(state.isGameOver).toBe(false);
  });

  it('should have 2 black and 2 white stones', () => {
    const state = createInitialGameState();
    expect(state.blackCount).toBe(2);
    expect(state.whiteCount).toBe(2);
  });

  it('should have 4 valid moves', () => {
    const state = createInitialGameState();
    expect(state.validMoves.length).toBe(4);
  });

  it('should have no winner', () => {
    const state = createInitialGameState();
    expect(state.winner).toBeNull();
  });
});

describe('handleMove', () => {
  it('should switch to white after black moves', () => {
    const state = createInitialGameState();
    const next = handleMove(state, 2, 3);
    expect(next.currentPlayer).toBe('white');
  });

  it('should update stone counts after a move', () => {
    const state = createInitialGameState();
    const next = handleMove(state, 2, 3);
    expect(next.blackCount).toBe(4);
    expect(next.whiteCount).toBe(1);
  });

  it('should record the last move', () => {
    const state = createInitialGameState();
    const next = handleMove(state, 2, 3);
    expect(next.lastMove).toEqual({ row: 2, col: 3 });
  });

  it('should not change state for an invalid move', () => {
    const state = createInitialGameState();
    const next = handleMove(state, 0, 0);
    expect(next).toBe(state);
  });

  it('should update valid moves for the next player', () => {
    const state = createInitialGameState();
    const next = handleMove(state, 2, 3);
    expect(next.validMoves.length).toBeGreaterThan(0);
    // White should have valid moves after black plays (2,3)
    expect(next.currentPlayer).toBe('white');
  });
});

describe('applyOpponentMove', () => {
  it('should apply opponent move to the board', () => {
    const state = createInitialGameState();
    // Black plays at (2,3) - this is a valid move for black
    const next = applyOpponentMove(state, 2, 3, 'black');
    expect(next.board[2][3]).toBe('black');
    expect(next.blackCount).toBe(4);
    expect(next.whiteCount).toBe(1);
  });

  it('should switch to the correct next player', () => {
    const state = createInitialGameState();
    const next = applyOpponentMove(state, 2, 3, 'black');
    expect(next.currentPlayer).toBe('white');
  });

  it('should compute valid moves for the next player', () => {
    const state = createInitialGameState();
    const next = applyOpponentMove(state, 2, 3, 'black');
    expect(next.validMoves.length).toBeGreaterThan(0);
  });
});
