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
