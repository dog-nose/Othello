import type { CellState, Color } from '../types/game';
import './Cell.css';

interface CellProps {
  state: CellState;
  isValid: boolean;
  isLastMove: boolean;
  hintColor?: Color;
  onClick: () => void;
}

export function Cell({ state, isValid, isLastMove, hintColor, onClick }: CellProps) {
  return (
    <div
      className={`cell ${isValid ? 'cell--valid' : ''} ${isLastMove ? 'cell--last-move' : ''}`}
      onClick={isValid ? onClick : undefined}
    >
      {state && <div className={`stone stone--${state}`} />}
      {isValid && !state && <div className={`hint ${hintColor ? `hint--${hintColor}` : ''}`} />}
    </div>
  );
}
