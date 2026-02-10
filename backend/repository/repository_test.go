package repository

import (
	"testing"

	"github.com/dog-nose/othello-backend/testutil"
)

func TestCreateGame(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()
	defer testutil.CleanupTestDB(t, db)

	repo := NewMySQLRepository(db)

	err := repo.CreateGame("test-play-id-1")
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	game, err := repo.GetGame("test-play-id-1")
	if err != nil {
		t.Fatalf("failed to get game: %v", err)
	}
	if game.PlayID != "test-play-id-1" {
		t.Fatalf("expected play_id test-play-id-1, got %s", game.PlayID)
	}
	if game.Result != nil {
		t.Fatalf("expected result to be nil, got %v", *game.Result)
	}
}

func TestRecordMove(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()
	defer testutil.CleanupTestDB(t, db)

	repo := NewMySQLRepository(db)

	err := repo.CreateGame("test-play-id-2")
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	err = repo.RecordMove("test-play-id-2", "black", 2, 3, 1)
	if err != nil {
		t.Fatalf("failed to record move: %v", err)
	}

	count, err := repo.GetMoveCount("test-play-id-2")
	if err != nil {
		t.Fatalf("failed to get move count: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected move count 1, got %d", count)
	}

	err = repo.RecordMove("test-play-id-2", "white", 4, 5, 2)
	if err != nil {
		t.Fatalf("failed to record second move: %v", err)
	}

	count, err = repo.GetMoveCount("test-play-id-2")
	if err != nil {
		t.Fatalf("failed to get move count: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected move count 2, got %d", count)
	}
}

func TestEndGame(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()
	defer testutil.CleanupTestDB(t, db)

	repo := NewMySQLRepository(db)

	err := repo.CreateGame("test-play-id-3")
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	err = repo.EndGame("test-play-id-3", 40, 24, "black_win")
	if err != nil {
		t.Fatalf("failed to end game: %v", err)
	}

	game, err := repo.GetGame("test-play-id-3")
	if err != nil {
		t.Fatalf("failed to get game: %v", err)
	}
	if game.Result == nil || *game.Result != "black_win" {
		t.Fatalf("expected result black_win, got %v", game.Result)
	}
	if game.BlackCount == nil || *game.BlackCount != 40 {
		t.Fatalf("expected black_count 40, got %v", game.BlackCount)
	}
	if game.WhiteCount == nil || *game.WhiteCount != 24 {
		t.Fatalf("expected white_count 24, got %v", game.WhiteCount)
	}
}
