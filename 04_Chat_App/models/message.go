package models

import (
	"time"
)

type Message struct {
	ID        int64     `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // "message", "join", "leave"
	Timestamp time.Time `json:"timestamp"`
}

type Room struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UserCount int       `json:"user_count"`
}
