package model

import "time"

// Domain types

type Game struct {
	PlayID     string    `json:"play_id"`
	BlackCount *int      `json:"black_count"`
	WhiteCount *int      `json:"white_count"`
	Result     *string   `json:"result"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
}

type EndGameRequest struct {
	PlayID     string `json:"play_id"`
	BlackCount int    `json:"black_count"`
	WhiteCount int    `json:"white_count"`
}

// Response types

type StartGameResponse struct {
	PlayID string `json:"play_id"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
