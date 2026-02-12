import type { Board, Color, GameState } from '../types/game';
import { createInitialBoard, getValidMoves, placeStone, countStones, isGameOver, getWinner } from './board';

export function createInitialGameState(): GameState {
  const board = createInitialBoard();
  return {
    board,
    currentPlayer: 'black',
    isGameOver: false,
    playId: null,
    blackCount: 2,
    whiteCount: 2,
    winner: null,
    validMoves: getValidMoves(board, 'black'),
    lastMove: null,
  };
}

function nextPlayer(current: Color): Color {
  return current === 'black' ? 'white' : 'black';
}

export function applyOpponentMove(state: GameState, row: number, col: number, color: Color): GameState {
  const newBoard = placeStone(state.board, row, col, color);
  if (newBoard === state.board) {
    return state;
  }
  return resolveNextTurn(newBoard, nextPlayer(color), { row, col });
}

export function handleMove(state: GameState, row: number, col: number): GameState {
  const newBoard = placeStone(state.board, row, col, state.currentPlayer);
  if (newBoard === state.board) {
    return state;
  }

  return resolveNextTurn(newBoard, nextPlayer(state.currentPlayer), { row, col });
}

function resolveNextTurn(board: Board, next: Color, lastMove: { row: number; col: number }): GameState {
  const { black, white } = countStones(board);

  if (isGameOver(board)) {
    return {
      board,
      currentPlayer: next,
      isGameOver: true,
      playId: null,
      blackCount: black,
      whiteCount: white,
      winner: getWinner(board),
      validMoves: [],
      lastMove,
    };
  }

  const nextMoves = getValidMoves(board, next);
  if (nextMoves.length > 0) {
    return {
      board,
      currentPlayer: next,
      isGameOver: false,
      playId: null,
      blackCount: black,
      whiteCount: white,
      winner: null,
      validMoves: nextMoves,
      lastMove,
    };
  }

  // Pass: the next player has no valid moves, so skip back
  const passedPlayer = nextPlayer(next);
  const passedMoves = getValidMoves(board, passedPlayer);
  return {
    board,
    currentPlayer: passedPlayer,
    isGameOver: false,
    playId: null,
    blackCount: black,
    whiteCount: white,
    winner: null,
    validMoves: passedMoves,
    lastMove,
  };
}
