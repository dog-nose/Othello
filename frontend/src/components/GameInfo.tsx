import type { GameState } from '../types/game';

interface GameInfoProps {
  gameState: GameState;
}

export function GameInfo({ gameState }: GameInfoProps) {
  const { currentPlayer, blackCount, whiteCount, isGameOver, winner, playId } = gameState;

  return (
    <div className="game-info">
      {playId && <p className="play-id">Play ID: {playId}</p>}
      <div className="score">
        <span className={`score-item ${currentPlayer === 'black' && !isGameOver ? 'score-item--active' : ''}`}>
          <span className="stone-icon stone-icon--black" /> {blackCount}
        </span>
        <span className="score-separator">-</span>
        <span className={`score-item ${currentPlayer === 'white' && !isGameOver ? 'score-item--active' : ''}`}>
          {whiteCount} <span className="stone-icon stone-icon--white" />
        </span>
      </div>
      {isGameOver ? (
        <p className="game-result">
          {winner === 'draw' ? 'Draw!' : winner === 'black' ? 'Black wins!' : 'White wins!'}
        </p>
      ) : (
        <p className="turn-indicator">
          {currentPlayer === 'black' ? 'Black' : 'White'}&apos;s turn
        </p>
      )}
    </div>
  );
}
