import { useState } from 'react';
import { Board } from './components/Board';
import { GameInfo } from './components/GameInfo';
import { GameControls } from './components/GameControls';
import { PvPLobby } from './components/PvPLobby';
import { useGame } from './hooks/useGame';
import { usePvPGame } from './hooks/usePvPGame';
import './App.css';

type GameMode = 'select' | 'local' | 'pvp';

function App() {
  const [mode, setMode] = useState<GameMode>('select');
  const localGame = useGame();
  const pvpGame = usePvPGame();

  const handleBackToMenu = () => {
    setMode('select');
    localGame.restart();
    pvpGame.restart();
  };

  if (mode === 'select') {
    return (
      <div className="app">
        <h1 className="app-title">Othello</h1>
        <div className="mode-select">
          <button className="btn btn--start" onClick={() => setMode('local')}>
            Local Game
          </button>
          <button className="btn btn--pvp" onClick={() => setMode('pvp')}>
            Online PvP
          </button>
        </div>
      </div>
    );
  }

  if (mode === 'local') {
    return (
      <div className="app">
        <h1 className="app-title">Othello</h1>
        <GameInfo gameState={localGame.gameState} />
        <Board
          board={localGame.gameState.board}
          validMoves={localGame.isStarted && !localGame.gameState.isGameOver ? localGame.gameState.validMoves : []}
          lastMove={localGame.gameState.lastMove}
          onCellClick={localGame.makeMove}
        />
        <GameControls
          isGameStarted={localGame.isStarted}
          isGameOver={localGame.gameState.isGameOver}
          onStart={localGame.startGame}
          onRestart={localGame.restart}
        />
        <button className="btn btn--back" onClick={handleBackToMenu}>
          Back to Menu
        </button>
      </div>
    );
  }

  // PvP mode
  const { gameState, pvpState, isStarted, error, createGame, joinGame, makeMove, restart } = pvpGame;
  const isGameOver = gameState.isGameOver;
  const showValidMoves = pvpState?.isMyTurn && !isGameOver && !pvpState.isWaitingForOpponent;

  return (
    <div className="app">
      <h1 className="app-title">Othello</h1>
      {!isStarted ? (
        <PvPLobby onCreateGame={createGame} onJoinGame={joinGame} error={error} />
      ) : (
        <>
          <GameInfo gameState={gameState} pvpState={pvpState} />
          <Board
            board={gameState.board}
            validMoves={showValidMoves ? gameState.validMoves : []}
            lastMove={gameState.lastMove}
            hintColor={pvpState?.myColor}
            onCellClick={makeMove}
          />
          {isGameOver && (
            <button className="btn btn--restart" onClick={restart}>
              Play Again
            </button>
          )}
        </>
      )}
      <button className="btn btn--back" onClick={handleBackToMenu}>
        Back to Menu
      </button>
    </div>
  );
}

export default App;
