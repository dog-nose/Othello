package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dog-nose/othello-backend/model"
)

// Mock Repository

type mockRepository struct {
	createGameFn   func(playID string) error
	getGameFn      func(playID string) (*model.Game, error)
	recordMoveFn   func(playID, color string, col, row, moveOrder int) error
	getMoveCountFn func(playID string) (int, error)
	endGameFn      func(playID string, blackCount, whiteCount int, result string) error
}

func (m *mockRepository) CreateGame(playID string) error {
	if m.createGameFn != nil {
		return m.createGameFn(playID)
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

func TestStartGame(t *testing.T) {
	mock := &mockRepository{}
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
		createGameFn: func(playID string) error {
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
