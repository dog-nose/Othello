import { useState } from 'react';

interface PvPLobbyProps {
  onCreateGame: () => void;
  onJoinGame: (playId: string) => void;
  error: string | null;
}

export function PvPLobby({ onCreateGame, onJoinGame, error }: PvPLobbyProps) {
  const [playId, setPlayId] = useState('');

  const handleJoin = () => {
    if (playId.trim()) {
      onJoinGame(playId.trim());
    }
  };

  return (
    <div className="pvp-lobby">
      <div className="lobby-section">
        <button className="btn btn--start" onClick={onCreateGame}>
          Create Game
        </button>
      </div>
      <div className="lobby-divider">or</div>
      <div className="lobby-section">
        <input
          className="lobby-input"
          type="text"
          placeholder="Enter Play ID"
          value={playId}
          onChange={e => setPlayId(e.target.value)}
          onKeyDown={e => e.key === 'Enter' && handleJoin()}
        />
        <button className="btn btn--join" onClick={handleJoin} disabled={!playId.trim()}>
          Join Game
        </button>
      </div>
      {error && <p className="lobby-error">{error}</p>}
    </div>
  );
}
