import type { GameState, PvPState } from '../types/game';

interface GameInfoProps {
  gameState: GameState;
  pvpState?: PvPState | null;
}

export function GameInfo({ gameState, pvpState }: GameInfoProps) {
  const { currentPlayer, blackCount, whiteCount, isGameOver, winner, playId } = gameState;

  return (
    <div className="game-info">
      {playId && <p className="play-id">Play ID: {playId}</p>}
      {pvpState && (
        <p className="pvp-role">
          You are: {pvpState.myColor === 'black' ? 'Black' : 'White'}
        </p>
      )}
      <div className="score">
        <span className={`score-item ${currentPlayer === 'black' && !isGameOver ? 'score-item--active' : ''}`}>
          <span className="stone-icon stone-icon--black" /> {blackCount}
        </span>
        <span className="score-separator">-</span>
        <span className={`score-item ${currentPlayer === 'white' && !isGameOver ? 'score-item--active' : ''}`}>
          {whiteCount} <span className="stone-icon stone-icon--white" />
        </span>
      </div>
      {pvpState?.isWaitingForOpponent ? (
        <p className="turn-indicator">Waiting for opponent to join...</p>
      ) : isGameOver ? (
        <p className="game-result">
          {winner === 'draw' ? 'Draw!' : winner === 'black' ? 'Black wins!' : 'White wins!'}
        </p>
      ) : pvpState ? (
        <p className="turn-indicator">
          {pvpState.isMyTurn ? 'Your turn' : "Opponent's turn..."}
        </p>
      ) : (
        <p className="turn-indicator">
          {currentPlayer === 'black' ? 'Black' : 'White'}&apos;s turn
        </p>
      )}
    </div>
  );
}
