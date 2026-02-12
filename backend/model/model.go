package model

import "time"

// Domain types

type Game struct {
	PlayID      string    `json:"play_id"`
	BlackCount  *int      `json:"black_count"`
	WhiteCount  *int      `json:"white_count"`
	Result      *string   `json:"result"`
	HostSecret  *string   `json:"host_secret,omitempty"`
	GuestSecret *string   `json:"guest_secret,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Move struct {
	ID        int64     `json:"id"`
	PlayID    string    `json:"play_id"`
	Color     string    `json:"color"`
	Col       int       `json:"col"`
	Row       int       `json:"row"`
	MoveOrder int       `json:"move_order"`
	CreatedAt time.Time `json:"created_at"`
}

// Request types

type PlaceStoneRequest struct {
	PlayID string `json:"play_id"`
	Color  string `json:"color"`
	Col    int    `json:"col"`
	Row    int    `json:"row"`
	Secret string `json:"secret,omitempty"`
}

type EndGameRequest struct {
	PlayID     string `json:"play_id"`
	BlackCount int    `json:"black_count"`
	WhiteCount int    `json:"white_count"`
}

// Response types

type StartGameResponse struct {
	PlayID     string `json:"play_id"`
	HostSecret string `json:"host_secret,omitempty"`
}

type JoinGameRequest struct {
	PlayID string `json:"play_id"`
}

type JoinGameResponse struct {
	GuestSecret string `json:"guest_secret"`
}

type PollMovesRequest struct {
	PlayID         string `json:"play_id"`
	AfterMoveOrder int    `json:"after_move_order"`
}

type PollMovesResponse struct {
	Moves []Move `json:"moves"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
