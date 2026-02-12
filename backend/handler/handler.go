package handler

import (
	"encoding/json"
	"errors"
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
	hostSecret := uuid.New().String()
	if err := h.repo.CreateGameWithSecret(playID, hostSecret); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create game")
		return
	}

	respondJSON(w, http.StatusOK, model.StartGameResponse{PlayID: playID, HostSecret: hostSecret})
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

	moveOrder := moveCount + 1

	// Secret validation for PvP games
	game, err := h.repo.GetGame(req.PlayID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get game")
		return
	}

	if game.HostSecret != nil {
		// PvP game: validate secret
		var expectedSecret string
		if moveOrder%2 == 1 {
			// Odd move_order = black (host)
			expectedSecret = *game.HostSecret
		} else {
			// Even move_order = white (guest)
			if game.GuestSecret != nil {
				expectedSecret = *game.GuestSecret
			}
		}
		if req.Secret != expectedSecret {
			respondError(w, http.StatusForbidden, "invalid secret")
			return
		}
	}

	if err := h.repo.RecordMove(req.PlayID, req.Color, req.Col, req.Row, moveOrder); err != nil {
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

func (h *Handler) JoinGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.JoinGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PlayID == "" {
		respondError(w, http.StatusBadRequest, "play_id is required")
		return
	}

	guestSecret := uuid.New().String()
	if err := h.repo.SetGuestSecret(req.PlayID, guestSecret); err != nil {
		if errors.Is(err, repository.ErrGuestAlreadyJoined) {
			respondError(w, http.StatusConflict, "guest already joined or game not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to join game")
		return
	}

	respondJSON(w, http.StatusOK, model.JoinGameResponse{GuestSecret: guestSecret})
}

func (h *Handler) PollMoves(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.PollMovesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PlayID == "" {
		respondError(w, http.StatusBadRequest, "play_id is required")
		return
	}

	moves, err := h.repo.GetMovesAfter(req.PlayID, req.AfterMoveOrder)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get moves")
		return
	}

	respondJSON(w, http.StatusOK, model.PollMovesResponse{Moves: moves})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, model.SuccessResponse{Success: false, Message: message})
}
