import { useState, useCallback } from 'react';
import type { GameState } from '../types/game';
import { createInitialGameState, handleMove } from '../logic/game';
import * as api from '../api/client';

export function useGame() {
  const [gameState, setGameState] = useState<GameState>(createInitialGameState);
  const [isStarted, setIsStarted] = useState(false);

  const startGame = useCallback(async () => {
    try {
      const { play_id } = await api.startGame();
      const state = createInitialGameState();
      state.playId = play_id;
      setGameState(state);
      setIsStarted(true);
    } catch (err) {
      console.error('Failed to start game:', err);
    }
  }, []);

  const makeMove = useCallback(async (row: number, col: number) => {
    setGameState(prev => {
      const next = handleMove(prev, row, col);
      if (next === prev) return prev;

      if (prev.playId) {
        api.placeStone(prev.playId, prev.currentPlayer, col, row).catch(err => {
          console.error('Failed to record move:', err);
        });
      }

      if (next.isGameOver && prev.playId) {
        api.endGame(prev.playId, next.blackCount, next.whiteCount).catch(err => {
          console.error('Failed to end game:', err);
        });
      }

      return { ...next, playId: prev.playId };
    });
  }, []);

  const restart = useCallback(async () => {
    setIsStarted(false);
    setGameState(createInitialGameState());
  }, []);

  return { gameState, isStarted, startGame, makeMove, restart };
}
