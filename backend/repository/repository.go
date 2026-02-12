package repository

import (
	"database/sql"
	"errors"

	"github.com/dog-nose/othello-backend/model"
)

var ErrGuestAlreadyJoined = errors.New("guest already joined or game not found")

type Repository interface {
	CreateGame(playID string) error
	CreateGameWithSecret(playID, hostSecret string) error
	GetGame(playID string) (*model.Game, error)
	RecordMove(playID, color string, col, row, moveOrder int) error
	GetMoveCount(playID string) (int, error)
	EndGame(playID string, blackCount, whiteCount int, result string) error
	SetGuestSecret(playID, guestSecret string) error
	GetMovesAfter(playID string, afterMoveOrder int) ([]model.Move, error)
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) CreateGame(playID string) error {
	_, err := r.db.Exec("INSERT INTO games (play_id) VALUES (?)", playID)
	return err
}

func (r *MySQLRepository) CreateGameWithSecret(playID, hostSecret string) error {
	_, err := r.db.Exec("INSERT INTO games (play_id, host_secret) VALUES (?, ?)", playID, hostSecret)
	return err
}

func (r *MySQLRepository) GetGame(playID string) (*model.Game, error) {
	game := &model.Game{}
	err := r.db.QueryRow(
		"SELECT play_id, black_count, white_count, result, host_secret, guest_secret, created_at, updated_at FROM games WHERE play_id = ?",
		playID,
	).Scan(&game.PlayID, &game.BlackCount, &game.WhiteCount, &game.Result, &game.HostSecret, &game.GuestSecret, &game.CreatedAt, &game.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (r *MySQLRepository) SetGuestSecret(playID, guestSecret string) error {
	result, err := r.db.Exec(
		"UPDATE games SET guest_secret = ? WHERE play_id = ? AND guest_secret IS NULL",
		guestSecret, playID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrGuestAlreadyJoined
	}
	return nil
}

func (r *MySQLRepository) GetMovesAfter(playID string, afterMoveOrder int) ([]model.Move, error) {
	rows, err := r.db.Query(
		"SELECT id, play_id, color, col, `row`, move_order, created_at FROM moves WHERE play_id = ? AND move_order > ? ORDER BY move_order ASC",
		playID, afterMoveOrder,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var moves []model.Move
	for rows.Next() {
		var m model.Move
		if err := rows.Scan(&m.ID, &m.PlayID, &m.Color, &m.Col, &m.Row, &m.MoveOrder, &m.CreatedAt); err != nil {
			return nil, err
		}
		moves = append(moves, m)
	}
	if moves == nil {
		moves = []model.Move{}
	}
	return moves, rows.Err()
}

func (r *MySQLRepository) RecordMove(playID, color string, col, row, moveOrder int) error {
	_, err := r.db.Exec(
		"INSERT INTO moves (play_id, color, col, `row`, move_order) VALUES (?, ?, ?, ?, ?)",
		playID, color, col, row, moveOrder,
	)
	return err
}

func (r *MySQLRepository) GetMoveCount(playID string) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM moves WHERE play_id = ?", playID).Scan(&count)
	return count, err
}

func (r *MySQLRepository) EndGame(playID string, blackCount, whiteCount int, result string) error {
	_, err := r.db.Exec(
		"UPDATE games SET black_count = ?, white_count = ?, result = ? WHERE play_id = ?",
		blackCount, whiteCount, result, playID,
	)
	return err
}
