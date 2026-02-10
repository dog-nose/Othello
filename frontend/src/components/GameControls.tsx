interface GameControlsProps {
  isGameStarted: boolean;
  isGameOver: boolean;
  onStart: () => void;
  onRestart: () => void;
}

export function GameControls({ isGameStarted, isGameOver, onStart, onRestart }: GameControlsProps) {
  return (
    <div className="game-controls">
      {!isGameStarted ? (
        <button className="btn btn--start" onClick={onStart}>
          Start Game
        </button>
      ) : isGameOver ? (
        <button className="btn btn--restart" onClick={onRestart}>
          Play Again
        </button>
      ) : null}
    </div>
  );
}
