const API_BASE = '/api';

export async function startGame(): Promise<{ play_id: string; host_secret: string }> {
  const res = await fetch(`${API_BASE}/start-game`, { method: 'POST' });
  return res.json();
}

export async function placeStone(playId: string, color: string, col: number, row: number, secret?: string): Promise<{ success: boolean; message?: string }> {
  const body: Record<string, unknown> = { play_id: playId, color, col, row };
  if (secret) {
    body.secret = secret;
  }
  const res = await fetch(`${API_BASE}/place-stone`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  return res.json();
}

export async function endGame(playId: string, blackCount: number, whiteCount: number): Promise<{ success: boolean; message?: string }> {
  const res = await fetch(`${API_BASE}/end-game`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ play_id: playId, black_count: blackCount, white_count: whiteCount }),
  });
  return res.json();
}

export async function joinGame(playId: string): Promise<{ guest_secret: string }> {
  const res = await fetch(`${API_BASE}/join-game`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ play_id: playId }),
  });
  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.message || 'Failed to join game');
  }
  return res.json();
}

export interface PollMovesMove {
  play_id: string;
  color: string;
  col: number;
  row: number;
  move_order: number;
}

export async function pollMoves(playId: string, afterMoveOrder: number): Promise<{ moves: PollMovesMove[] }> {
  const res = await fetch(`${API_BASE}/poll-moves`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ play_id: playId, after_move_order: afterMoveOrder }),
  });
  return res.json();
}
