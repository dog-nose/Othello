package repository

import (
	"database/sql"

	"github.com/dog-nose/othello-backend/model"
)

type Repository interface {
	CreateGame(playID string) error
	GetGame(playID string) (*model.Game, error)
	RecordMove(playID, color string, col, row, moveOrder int) error
	GetMoveCount(playID string) (int, error)
	EndGame(playID string, blackCount, whiteCount int, result string) error
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

func (r *MySQLRepository) GetGame(playID string) (*model.Game, error) {
	game := &model.Game{}
	err := r.db.QueryRow(
		"SELECT play_id, black_count, white_count, result, created_at, updated_at FROM games WHERE play_id = ?",
		playID,
	).Scan(&game.PlayID, &game.BlackCount, &game.WhiteCount, &game.Result, &game.CreatedAt, &game.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return game, nil
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
