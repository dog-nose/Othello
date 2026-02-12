import { useState, useCallback, useEffect, useRef } from 'react';
import type { GameState, PvPState, Color } from '../types/game';
import { createInitialGameState, handleMove, applyOpponentMove } from '../logic/game';
import * as api from '../api/client';

export function usePvPGame() {
  const [gameState, setGameState] = useState<GameState>(createInitialGameState);
  const [pvpState, setPvPState] = useState<PvPState | null>(null);
  const [isStarted, setIsStarted] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const gameStateRef = useRef(gameState);
  const pvpStateRef = useRef(pvpState);

  gameStateRef.current = gameState;
  pvpStateRef.current = pvpState;

  const createGame = useCallback(async () => {
    try {
      setError(null);
      const { play_id, host_secret } = await api.startGame();
      const state = createInitialGameState();
      state.playId = play_id;
      setGameState(state);
      setPvPState({
        role: 'host',
        myColor: 'black',
        secret: host_secret,
        lastKnownMoveOrder: 0,
        isWaitingForOpponent: true,
        isMyTurn: true,
      });
      setIsStarted(true);
    } catch (err) {
      setError('Failed to create game');
      console.error('Failed to create game:', err);
    }
  }, []);

  const joinGame = useCallback(async (playId: string) => {
    try {
      setError(null);
      const { guest_secret } = await api.joinGame(playId);
      const state = createInitialGameState();
      state.playId = playId;
      setGameState(state);
      setPvPState({
        role: 'guest',
        myColor: 'white',
        secret: guest_secret,
        lastKnownMoveOrder: 0,
        isWaitingForOpponent: false,
        isMyTurn: false,
      });
      setIsStarted(true);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to join game';
      setError(message);
      console.error('Failed to join game:', err);
    }
  }, []);

  const makeMove = useCallback(async (row: number, col: number) => {
    const currentPvP = pvpStateRef.current;
    if (!currentPvP || !currentPvP.isMyTurn) return;

    setGameState(prev => {
      const next = handleMove(prev, row, col);
      if (next === prev) return prev;

      const newMoveOrder = currentPvP.lastKnownMoveOrder + 1;

      if (prev.playId) {
        api.placeStone(prev.playId, prev.currentPlayer, col, row, currentPvP.secret).catch(err => {
          console.error('Failed to record move:', err);
        });
      }

      if (next.isGameOver && prev.playId) {
        api.endGame(prev.playId, next.blackCount, next.whiteCount).catch(err => {
          console.error('Failed to end game:', err);
        });
      }

      // Determine if it's still my turn (pass case)
      const stillMyTurn = !next.isGameOver && next.currentPlayer === currentPvP.myColor;

      setPvPState(prevPvP => prevPvP ? {
        ...prevPvP,
        lastKnownMoveOrder: newMoveOrder,
        isMyTurn: stillMyTurn,
      } : prevPvP);

      return { ...next, playId: prev.playId };
    });
  }, []);

  // Polling for opponent moves
  useEffect(() => {
    if (!isStarted || !pvpState) return;

    const interval = setInterval(async () => {
      const currentPvP = pvpStateRef.current;
      const currentGame = gameStateRef.current;

      if (!currentPvP || !currentGame.playId) return;
      if (currentPvP.isMyTurn && !currentPvP.isWaitingForOpponent) return;
      if (currentGame.isGameOver) return;

      try {
        const { moves } = await api.pollMoves(currentGame.playId, currentPvP.lastKnownMoveOrder);

        if (moves.length === 0) {
          // Check if opponent has joined (host waiting)
          if (currentPvP.isWaitingForOpponent && currentPvP.role === 'host') {
            // We can't directly check if guest joined, but once they make a move we'll know
            // For now, keep polling
          }
          return;
        }

        // If we're waiting for opponent and we got moves, they've joined
        let updatedGame = currentGame;
        let latestMoveOrder = currentPvP.lastKnownMoveOrder;

        for (const move of moves) {
          const color = move.color as Color;
          updatedGame = applyOpponentMove(updatedGame, move.row, move.col, color);
          updatedGame = { ...updatedGame, playId: currentGame.playId };
          latestMoveOrder = move.move_order;
        }

        setGameState(updatedGame);

        const isNowMyTurn = !updatedGame.isGameOver && updatedGame.currentPlayer === currentPvP.myColor;

        setPvPState(prev => prev ? {
          ...prev,
          lastKnownMoveOrder: latestMoveOrder,
          isMyTurn: isNowMyTurn,
          isWaitingForOpponent: false,
        } : prev);

        if (updatedGame.isGameOver && currentGame.playId) {
          api.endGame(currentGame.playId, updatedGame.blackCount, updatedGame.whiteCount).catch(err => {
            console.error('Failed to end game:', err);
          });
        }
      } catch (err) {
        console.error('Poll error:', err);
      }
    }, 1000);

    return () => clearInterval(interval);
  }, [isStarted, pvpState]);

  const restart = useCallback(() => {
    setIsStarted(false);
    setGameState(createInitialGameState());
    setPvPState(null);
    setError(null);
  }, []);

  return { gameState, pvpState, isStarted, error, createGame, joinGame, makeMove, restart };
}
