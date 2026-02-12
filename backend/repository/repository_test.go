package repository

import (
	"errors"
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

func TestCreateGameWithSecret(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()
	defer testutil.CleanupTestDB(t, db)

	repo := NewMySQLRepository(db)

	err := repo.CreateGameWithSecret("test-secret-1", "host-secret-abc")
	if err != nil {
		t.Fatalf("failed to create game with secret: %v", err)
	}

	game, err := repo.GetGame("test-secret-1")
	if err != nil {
		t.Fatalf("failed to get game: %v", err)
	}
	if game.PlayID != "test-secret-1" {
		t.Fatalf("expected play_id test-secret-1, got %s", game.PlayID)
	}
	if game.HostSecret == nil || *game.HostSecret != "host-secret-abc" {
		t.Fatalf("expected host_secret host-secret-abc, got %v", game.HostSecret)
	}
	if game.GuestSecret != nil {
		t.Fatalf("expected guest_secret to be nil, got %v", *game.GuestSecret)
	}
}

func TestSetGuestSecret(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()
	defer testutil.CleanupTestDB(t, db)

	repo := NewMySQLRepository(db)

	err := repo.CreateGameWithSecret("test-guest-1", "host-secret-123")
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	err = repo.SetGuestSecret("test-guest-1", "guest-secret-456")
	if err != nil {
		t.Fatalf("failed to set guest secret: %v", err)
	}

	game, err := repo.GetGame("test-guest-1")
	if err != nil {
		t.Fatalf("failed to get game: %v", err)
	}
	if game.GuestSecret == nil || *game.GuestSecret != "guest-secret-456" {
		t.Fatalf("expected guest_secret guest-secret-456, got %v", game.GuestSecret)
	}
}

func TestSetGuestSecret_AlreadySet(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()
	defer testutil.CleanupTestDB(t, db)

	repo := NewMySQLRepository(db)

	err := repo.CreateGameWithSecret("test-guest-2", "host-secret-aaa")
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	err = repo.SetGuestSecret("test-guest-2", "guest-secret-first")
	if err != nil {
		t.Fatalf("failed to set guest secret: %v", err)
	}

	err = repo.SetGuestSecret("test-guest-2", "guest-secret-second")
	if !errors.Is(err, ErrGuestAlreadyJoined) {
		t.Fatalf("expected ErrGuestAlreadyJoined, got %v", err)
	}
}

func TestGetMovesAfter(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()
	defer testutil.CleanupTestDB(t, db)

	repo := NewMySQLRepository(db)

	err := repo.CreateGame("test-moves-1")
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	repo.RecordMove("test-moves-1", "black", 2, 3, 1)
	repo.RecordMove("test-moves-1", "white", 4, 5, 2)
	repo.RecordMove("test-moves-1", "black", 1, 2, 3)

	moves, err := repo.GetMovesAfter("test-moves-1", 1)
	if err != nil {
		t.Fatalf("failed to get moves after: %v", err)
	}
	if len(moves) != 2 {
		t.Fatalf("expected 2 moves, got %d", len(moves))
	}
	if moves[0].MoveOrder != 2 {
		t.Fatalf("expected first move order 2, got %d", moves[0].MoveOrder)
	}
	if moves[1].MoveOrder != 3 {
		t.Fatalf("expected second move order 3, got %d", moves[1].MoveOrder)
	}
}

func TestGetMovesAfter_NoMoves(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()
	defer testutil.CleanupTestDB(t, db)

	repo := NewMySQLRepository(db)

	err := repo.CreateGame("test-moves-2")
	if err != nil {
		t.Fatalf("failed to create game: %v", err)
	}

	moves, err := repo.GetMovesAfter("test-moves-2", 0)
	if err != nil {
		t.Fatalf("failed to get moves after: %v", err)
	}
	if len(moves) != 0 {
		t.Fatalf("expected 0 moves, got %d", len(moves))
	}
}
