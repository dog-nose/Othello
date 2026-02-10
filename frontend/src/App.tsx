import { Board } from './components/Board';
import { GameInfo } from './components/GameInfo';
import { GameControls } from './components/GameControls';
import { useGame } from './hooks/useGame';
import './App.css';

function App() {
  const { gameState, isStarted, startGame, makeMove, restart } = useGame();

  return (
    <div className="app">
      <h1 className="app-title">Othello</h1>
      <GameInfo gameState={gameState} />
      <Board
        board={gameState.board}
        validMoves={isStarted && !gameState.isGameOver ? gameState.validMoves : []}
        lastMove={gameState.lastMove}
        onCellClick={makeMove}
      />
      <GameControls
        isGameStarted={isStarted}
        isGameOver={gameState.isGameOver}
        onStart={startGame}
        onRestart={restart}
      />
    </div>
  );
}

export default App;
