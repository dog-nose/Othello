package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/dog-nose/othello-backend/model"
	"github.com/dog-nose/othello-backend/repository"
)

type Handler struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) StartGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	playID := uuid.New().String()
	if err := h.repo.CreateGame(playID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create game")
		return
	}

	respondJSON(w, http.StatusOK, model.StartGameResponse{PlayID: playID})
}

func (h *Handler) PlaceStone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.PlaceStoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PlayID == "" {
		respondError(w, http.StatusBadRequest, "play_id is required")
		return
	}
	if req.Color != "black" && req.Color != "white" {
		respondError(w, http.StatusBadRequest, "color must be 'black' or 'white'")
		return
	}
	if req.Col < 0 || req.Col > 7 || req.Row < 0 || req.Row > 7 {
		respondError(w, http.StatusBadRequest, "col and row must be between 0 and 7")
		return
	}

	moveCount, err := h.repo.GetMoveCount(req.PlayID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get move count")
		return
	}

	if err := h.repo.RecordMove(req.PlayID, req.Color, req.Col, req.Row, moveCount+1); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to record move")
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse{Success: true})
}

func (h *Handler) EndGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.EndGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PlayID == "" {
		respondError(w, http.StatusBadRequest, "play_id is required")
		return
	}

	var result string
	if req.BlackCount > req.WhiteCount {
		result = "black_win"
	} else if req.WhiteCount > req.BlackCount {
		result = "white_win"
	} else {
		result = "draw"
	}

	if err := h.repo.EndGame(req.PlayID, req.BlackCount, req.WhiteCount, result); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to end game")
		return
	}

	respondJSON(w, http.StatusOK, model.SuccessResponse{Success: true})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, model.SuccessResponse{Success: false, Message: message})
}
