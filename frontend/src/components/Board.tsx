import type { Board as BoardType, Position } from '../types/game';
import { Cell } from './Cell';
import './Board.css';

interface BoardProps {
  board: BoardType;
  validMoves: Position[];
  lastMove: Position | null;
  onCellClick: (row: number, col: number) => void;
}

export function Board({ board, validMoves, lastMove, onCellClick }: BoardProps) {
  const validSet = new Set(validMoves.map(m => `${m.row},${m.col}`));

  return (
    <div className="board">
      {board.map((row, r) => (
        <div key={r} className="board-row">
          {row.map((cell, c) => (
            <Cell
              key={c}
              state={cell}
              isValid={validSet.has(`${r},${c}`)}
              isLastMove={lastMove !== null && lastMove.row === r && lastMove.col === c}
              onClick={() => onCellClick(r, c)}
            />
          ))}
        </div>
      ))}
    </div>
  );
}
