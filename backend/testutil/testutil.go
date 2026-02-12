package testutil

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dog-nose/othello-backend/config"
)

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	cfg := config.LoadTest()
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}
	return db
}

func CleanupTestDB(t *testing.T, db *sql.DB) {
	t.Helper()
	db.Exec("DELETE FROM moves")
	db.Exec("DELETE FROM games")
}
