import type { CellState } from '../types/game';
import './Cell.css';

interface CellProps {
  state: CellState;
  isValid: boolean;
  isLastMove: boolean;
  onClick: () => void;
}

export function Cell({ state, isValid, isLastMove, onClick }: CellProps) {
  return (
    <div
      className={`cell ${isValid ? 'cell--valid' : ''} ${isLastMove ? 'cell--last-move' : ''}`}
      onClick={isValid ? onClick : undefined}
    >
      {state && <div className={`stone stone--${state}`} />}
      {isValid && !state && <div className="hint" />}
    </div>
  );
}
