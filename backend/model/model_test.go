package model

import (
	"encoding/json"
	"testing"
)

func TestStartGameResponseJSON(t *testing.T) {
	resp := StartGameResponse{PlayID: "test-123"}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var decoded StartGameResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if decoded.PlayID != "test-123" {
		t.Fatalf("expected play_id test-123, got %s", decoded.PlayID)
	}
}

func TestSuccessResponseJSON(t *testing.T) {
	resp := SuccessResponse{Success: true, Message: "ok"}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var decoded SuccessResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if !decoded.Success {
		t.Fatal("expected success to be true")
	}
	if decoded.Message != "ok" {
		t.Fatalf("expected message 'ok', got %s", decoded.Message)
	}
}

func TestSuccessResponseOmitsEmptyMessage(t *testing.T) {
	resp := SuccessResponse{Success: true}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	str := string(data)
	if str != `{"success":true}` {
		t.Fatalf("expected message to be omitted, got %s", str)
	}
}

func TestPlaceStoneRequestJSON(t *testing.T) {
	jsonStr := `{"play_id":"abc","color":"black","col":2,"row":3}`
	var req PlaceStoneRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if req.PlayID != "abc" {
		t.Fatalf("expected play_id abc, got %s", req.PlayID)
	}
	if req.Color != "black" {
		t.Fatalf("expected color black, got %s", req.Color)
	}
	if req.Col != 2 || req.Row != 3 {
		t.Fatalf("expected col=2 row=3, got col=%d row=%d", req.Col, req.Row)
	}
}

func TestStartGameResponseWithHostSecretJSON(t *testing.T) {
	resp := StartGameResponse{PlayID: "test-123", HostSecret: "secret-abc"}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var decoded StartGameResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if decoded.HostSecret != "secret-abc" {
		t.Fatalf("expected host_secret secret-abc, got %s", decoded.HostSecret)
	}
}

func TestStartGameResponseOmitsEmptyHostSecret(t *testing.T) {
	resp := StartGameResponse{PlayID: "test-123"}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	str := string(data)
	if str != `{"play_id":"test-123"}` {
		t.Fatalf("expected host_secret to be omitted, got %s", str)
	}
}

func TestPlaceStoneRequestWithSecretJSON(t *testing.T) {
	jsonStr := `{"play_id":"abc","color":"black","col":2,"row":3,"secret":"my-secret"}`
	var req PlaceStoneRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if req.Secret != "my-secret" {
		t.Fatalf("expected secret my-secret, got %s", req.Secret)
	}
}

func TestJoinGameRequestJSON(t *testing.T) {
	jsonStr := `{"play_id":"game-123"}`
	var req JoinGameRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if req.PlayID != "game-123" {
		t.Fatalf("expected play_id game-123, got %s", req.PlayID)
	}
}

func TestJoinGameResponseJSON(t *testing.T) {
	resp := JoinGameResponse{GuestSecret: "guest-secret-xyz"}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var decoded JoinGameResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if decoded.GuestSecret != "guest-secret-xyz" {
		t.Fatalf("expected guest_secret guest-secret-xyz, got %s", decoded.GuestSecret)
	}
}

func TestPollMovesRequestJSON(t *testing.T) {
	jsonStr := `{"play_id":"game-123","after_move_order":5}`
	var req PollMovesRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if req.PlayID != "game-123" {
		t.Fatalf("expected play_id game-123, got %s", req.PlayID)
	}
	if req.AfterMoveOrder != 5 {
		t.Fatalf("expected after_move_order 5, got %d", req.AfterMoveOrder)
	}
}

func TestPollMovesResponseJSON(t *testing.T) {
	resp := PollMovesResponse{
		Moves: []Move{
			{PlayID: "g1", Color: "black", Col: 2, Row: 3, MoveOrder: 1},
		},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	var decoded PollMovesResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(decoded.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(decoded.Moves))
	}
	if decoded.Moves[0].Color != "black" {
		t.Fatalf("expected color black, got %s", decoded.Moves[0].Color)
	}
}

func TestEndGameRequestJSON(t *testing.T) {
	jsonStr := `{"play_id":"xyz","black_count":40,"white_count":24}`
	var req EndGameRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if req.PlayID != "xyz" {
		t.Fatalf("expected play_id xyz, got %s", req.PlayID)
	}
	if req.BlackCount != 40 || req.WhiteCount != 24 {
		t.Fatalf("expected 40/24, got %d/%d", req.BlackCount, req.WhiteCount)
	}
}
