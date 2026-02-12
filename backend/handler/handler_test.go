package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dog-nose/othello-backend/model"
	"github.com/dog-nose/othello-backend/repository"
)

// Mock Repository

type mockRepository struct {
	createGameFn           func(playID string) error
	createGameWithSecretFn func(playID, hostSecret string) error
	getGameFn              func(playID string) (*model.Game, error)
	recordMoveFn           func(playID, color string, col, row, moveOrder int) error
	getMoveCountFn         func(playID string) (int, error)
	endGameFn              func(playID string, blackCount, whiteCount int, result string) error
	setGuestSecretFn       func(playID, guestSecret string) error
	getMovesAfterFn        func(playID string, afterMoveOrder int) ([]model.Move, error)
}

func (m *mockRepository) CreateGame(playID string) error {
	if m.createGameFn != nil {
		return m.createGameFn(playID)
	}
	return nil
}

func (m *mockRepository) CreateGameWithSecret(playID, hostSecret string) error {
	if m.createGameWithSecretFn != nil {
		return m.createGameWithSecretFn(playID, hostSecret)
	}
	return nil
}

func (m *mockRepository) GetGame(playID string) (*model.Game, error) {
	if m.getGameFn != nil {
		return m.getGameFn(playID)
	}
	return &model.Game{PlayID: playID}, nil
}

func (m *mockRepository) RecordMove(playID, color string, col, row, moveOrder int) error {
	if m.recordMoveFn != nil {
		return m.recordMoveFn(playID, color, col, row, moveOrder)
	}
	return nil
}

func (m *mockRepository) GetMoveCount(playID string) (int, error) {
	if m.getMoveCountFn != nil {
		return m.getMoveCountFn(playID)
	}
	return 0, nil
}

func (m *mockRepository) EndGame(playID string, blackCount, whiteCount int, result string) error {
	if m.endGameFn != nil {
		return m.endGameFn(playID, blackCount, whiteCount, result)
	}
	return nil
}

func (m *mockRepository) SetGuestSecret(playID, guestSecret string) error {
	if m.setGuestSecretFn != nil {
		return m.setGuestSecretFn(playID, guestSecret)
	}
	return nil
}

func (m *mockRepository) GetMovesAfter(playID string, afterMoveOrder int) ([]model.Move, error) {
	if m.getMovesAfterFn != nil {
		return m.getMovesAfterFn(playID, afterMoveOrder)
	}
	return []model.Move{}, nil
}

func TestStartGame(t *testing.T) {
	var calledPlayID, calledSecret string
	mock := &mockRepository{
		createGameWithSecretFn: func(playID, hostSecret string) error {
			calledPlayID = playID
			calledSecret = hostSecret
			return nil
		},
	}
	h := New(mock)

	req := httptest.NewRequest(http.MethodPost, "/start-game", nil)
	rec := httptest.NewRecorder()

	h.StartGame(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp model.StartGameResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.PlayID == "" {
		t.Fatal("expected play_id to be non-empty")
	}
	if resp.HostSecret == "" {
		t.Fatal("expected host_secret to be non-empty")
	}
	if calledPlayID == "" || calledSecret == "" {
		t.Fatal("expected CreateGameWithSecret to be called")
	}
}

func TestStartGame_MethodNotAllowed(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	req := httptest.NewRequest(http.MethodGet, "/start-game", nil)
	rec := httptest.NewRecorder()

	h.StartGame(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestPlaceStone(t *testing.T) {
	var recordedColor string
	var recordedCol, recordedRow, recordedOrder int

	mock := &mockRepository{
		getMoveCountFn: func(playID string) (int, error) {
			return 3, nil
		},
		recordMoveFn: func(playID, color string, col, row, moveOrder int) error {
			recordedColor = color
			recordedCol = col
			recordedRow = row
			recordedOrder = moveOrder
			return nil
		},
	}
	h := New(mock)

	body := model.PlaceStoneRequest{
		PlayID: "test-id",
		Color:  "black",
		Col:    2,
		Row:    3,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/place-stone", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if recordedColor != "black" {
		t.Fatalf("expected color black, got %s", recordedColor)
	}
	if recordedCol != 2 || recordedRow != 3 {
		t.Fatalf("expected col=2, row=3, got col=%d, row=%d", recordedCol, recordedRow)
	}
	if recordedOrder != 4 {
		t.Fatalf("expected move_order 4, got %d", recordedOrder)
	}
}

func TestPlaceStone_ValidSecret(t *testing.T) {
	hostSecret := "host-secret-123"
	mock := &mockRepository{
		getMoveCountFn: func(playID string) (int, error) {
			return 0, nil // move_order will be 1 (black/host)
		},
		getGameFn: func(playID string) (*model.Game, error) {
			return &model.Game{PlayID: playID, HostSecret: &hostSecret}, nil
		},
	}
	h := New(mock)

	body := model.PlaceStoneRequest{
		PlayID: "test-id",
		Color:  "black",
		Col:    2,
		Row:    3,
		Secret: "host-secret-123",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/place-stone", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestPlaceStone_WrongSecret(t *testing.T) {
	hostSecret := "host-secret-123"
	mock := &mockRepository{
		getMoveCountFn: func(playID string) (int, error) {
			return 0, nil // move_order will be 1 (black/host)
		},
		getGameFn: func(playID string) (*model.Game, error) {
			return &model.Game{PlayID: playID, HostSecret: &hostSecret}, nil
		},
	}
	h := New(mock)

	body := model.PlaceStoneRequest{
		PlayID: "test-id",
		Color:  "black",
		Col:    2,
		Row:    3,
		Secret: "wrong-secret",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/place-stone", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
}

func TestPlaceStone_GuestSecret(t *testing.T) {
	hostSecret := "host-secret-123"
	guestSecret := "guest-secret-456"
	mock := &mockRepository{
		getMoveCountFn: func(playID string) (int, error) {
			return 1, nil // move_order will be 2 (white/guest)
		},
		getGameFn: func(playID string) (*model.Game, error) {
			return &model.Game{PlayID: playID, HostSecret: &hostSecret, GuestSecret: &guestSecret}, nil
		},
	}
	h := New(mock)

	body := model.PlaceStoneRequest{
		PlayID: "test-id",
		Color:  "white",
		Col:    2,
		Row:    3,
		Secret: "guest-secret-456",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/place-stone", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestPlaceStone_NoSecretLegacyGame(t *testing.T) {
	// Legacy games (no host_secret) should skip secret validation
	mock := &mockRepository{
		getMoveCountFn: func(playID string) (int, error) {
			return 0, nil
		},
		getGameFn: func(playID string) (*model.Game, error) {
			return &model.Game{PlayID: playID}, nil
		},
	}
	h := New(mock)

	body := model.PlaceStoneRequest{
		PlayID: "test-id",
		Color:  "black",
		Col:    2,
		Row:    3,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/place-stone", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestPlaceStone_InvalidColor(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	body := model.PlaceStoneRequest{
		PlayID: "test-id",
		Color:  "red",
		Col:    0,
		Row:    0,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/place-stone", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestPlaceStone_InvalidPosition(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	body := model.PlaceStoneRequest{
		PlayID: "test-id",
		Color:  "black",
		Col:    8,
		Row:    0,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/place-stone", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestEndGame_BlackWin(t *testing.T) {
	var recordedResult string
	mock := &mockRepository{
		endGameFn: func(playID string, blackCount, whiteCount int, result string) error {
			recordedResult = result
			return nil
		},
	}
	h := New(mock)

	body := model.EndGameRequest{
		PlayID:     "test-id",
		BlackCount: 40,
		WhiteCount: 24,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/end-game", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.EndGame(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if recordedResult != "black_win" {
		t.Fatalf("expected result black_win, got %s", recordedResult)
	}
}

func TestEndGame_WhiteWin(t *testing.T) {
	var recordedResult string
	mock := &mockRepository{
		endGameFn: func(playID string, blackCount, whiteCount int, result string) error {
			recordedResult = result
			return nil
		},
	}
	h := New(mock)

	body := model.EndGameRequest{
		PlayID:     "test-id",
		BlackCount: 20,
		WhiteCount: 44,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/end-game", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.EndGame(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if recordedResult != "white_win" {
		t.Fatalf("expected result white_win, got %s", recordedResult)
	}
}

func TestEndGame_Draw(t *testing.T) {
	var recordedResult string
	mock := &mockRepository{
		endGameFn: func(playID string, blackCount, whiteCount int, result string) error {
			recordedResult = result
			return nil
		},
	}
	h := New(mock)

	body := model.EndGameRequest{
		PlayID:     "test-id",
		BlackCount: 32,
		WhiteCount: 32,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/end-game", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.EndGame(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if recordedResult != "draw" {
		t.Fatalf("expected result draw, got %s", recordedResult)
	}
}

func TestEndGame_MissingPlayID(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	body := model.EndGameRequest{
		BlackCount: 32,
		WhiteCount: 32,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/end-game", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.EndGame(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestStartGame_CreateGameError(t *testing.T) {
	mock := &mockRepository{
		createGameWithSecretFn: func(playID, hostSecret string) error {
			return fmt.Errorf("db error")
		},
	}
	h := New(mock)

	req := httptest.NewRequest(http.MethodPost, "/start-game", nil)
	rec := httptest.NewRecorder()

	h.StartGame(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func TestPlaceStone_MethodNotAllowed(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	req := httptest.NewRequest(http.MethodGet, "/place-stone", nil)
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestPlaceStone_InvalidBody(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	req := httptest.NewRequest(http.MethodPost, "/place-stone", strings.NewReader("invalid json"))
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestPlaceStone_MissingPlayID(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	body := model.PlaceStoneRequest{Color: "black", Col: 0, Row: 0}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/place-stone", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestPlaceStone_GetMoveCountError(t *testing.T) {
	mock := &mockRepository{
		getMoveCountFn: func(playID string) (int, error) {
			return 0, fmt.Errorf("db error")
		},
	}
	h := New(mock)

	body := model.PlaceStoneRequest{PlayID: "test-id", Color: "black", Col: 0, Row: 0}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/place-stone", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func TestPlaceStone_RecordMoveError(t *testing.T) {
	mock := &mockRepository{
		recordMoveFn: func(playID, color string, col, row, moveOrder int) error {
			return fmt.Errorf("db error")
		},
	}
	h := New(mock)

	body := model.PlaceStoneRequest{PlayID: "test-id", Color: "black", Col: 0, Row: 0}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/place-stone", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PlaceStone(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func TestEndGame_MethodNotAllowed(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	req := httptest.NewRequest(http.MethodGet, "/end-game", nil)
	rec := httptest.NewRecorder()

	h.EndGame(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

func TestEndGame_InvalidBody(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	req := httptest.NewRequest(http.MethodPost, "/end-game", strings.NewReader("bad json"))
	rec := httptest.NewRecorder()

	h.EndGame(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestEndGame_EndGameError(t *testing.T) {
	mock := &mockRepository{
		endGameFn: func(playID string, blackCount, whiteCount int, result string) error {
			return fmt.Errorf("db error")
		},
	}
	h := New(mock)

	body := model.EndGameRequest{PlayID: "test-id", BlackCount: 32, WhiteCount: 32}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/end-game", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.EndGame(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

// JoinGame tests

func TestJoinGame(t *testing.T) {
	var calledPlayID, calledGuestSecret string
	mock := &mockRepository{
		setGuestSecretFn: func(playID, guestSecret string) error {
			calledPlayID = playID
			calledGuestSecret = guestSecret
			return nil
		},
	}
	h := New(mock)

	body := model.JoinGameRequest{PlayID: "game-123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/join-game", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.JoinGame(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp model.JoinGameResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.GuestSecret == "" {
		t.Fatal("expected guest_secret to be non-empty")
	}
	if calledPlayID != "game-123" {
		t.Fatalf("expected play_id game-123, got %s", calledPlayID)
	}
	if calledGuestSecret == "" {
		t.Fatal("expected SetGuestSecret to be called with non-empty secret")
	}
}

func TestJoinGame_AlreadyJoined(t *testing.T) {
	mock := &mockRepository{
		setGuestSecretFn: func(playID, guestSecret string) error {
			return repository.ErrGuestAlreadyJoined
		},
	}
	h := New(mock)

	body := model.JoinGameRequest{PlayID: "game-123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/join-game", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.JoinGame(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", rec.Code)
	}
}

func TestJoinGame_MissingPlayID(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	body := model.JoinGameRequest{}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/join-game", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.JoinGame(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestJoinGame_MethodNotAllowed(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	req := httptest.NewRequest(http.MethodGet, "/join-game", nil)
	rec := httptest.NewRecorder()

	h.JoinGame(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}

// PollMoves tests

func TestPollMoves(t *testing.T) {
	mock := &mockRepository{
		getMovesAfterFn: func(playID string, afterMoveOrder int) ([]model.Move, error) {
			return []model.Move{
				{PlayID: playID, Color: "black", Col: 2, Row: 3, MoveOrder: 2},
				{PlayID: playID, Color: "white", Col: 4, Row: 5, MoveOrder: 3},
			}, nil
		},
	}
	h := New(mock)

	body := model.PollMovesRequest{PlayID: "game-123", AfterMoveOrder: 1}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/poll-moves", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PollMoves(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp model.PollMovesResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Moves) != 2 {
		t.Fatalf("expected 2 moves, got %d", len(resp.Moves))
	}
}

func TestPollMoves_NoNewMoves(t *testing.T) {
	mock := &mockRepository{
		getMovesAfterFn: func(playID string, afterMoveOrder int) ([]model.Move, error) {
			return []model.Move{}, nil
		},
	}
	h := New(mock)

	body := model.PollMovesRequest{PlayID: "game-123", AfterMoveOrder: 5}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/poll-moves", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PollMoves(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp model.PollMovesResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp.Moves) != 0 {
		t.Fatalf("expected 0 moves, got %d", len(resp.Moves))
	}
}

func TestPollMoves_MissingPlayID(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	body := model.PollMovesRequest{}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/poll-moves", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.PollMoves(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestPollMoves_MethodNotAllowed(t *testing.T) {
	mock := &mockRepository{}
	h := New(mock)

	req := httptest.NewRequest(http.MethodGet, "/poll-moves", nil)
	rec := httptest.NewRecorder()

	h.PollMoves(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}
