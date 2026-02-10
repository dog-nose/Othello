const API_BASE = '/api';

export async function startGame(): Promise<{ play_id: string }> {
  const res = await fetch(`${API_BASE}/start-game`, { method: 'POST' });
  return res.json();
}

export async function placeStone(playId: string, color: string, col: number, row: number): Promise<{ success: boolean; message?: string }> {
  const res = await fetch(`${API_BASE}/place-stone`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ play_id: playId, color, col, row }),
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
