package main

import "time"

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Hashed password
}

// Habit represents a habit tracked by a user
type Habit struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// HabitMark represents a day a habit was completed
type HabitMark struct {
	ID      int       `json:"id"`
	HabitID int       `json:"habit_id"`
	Date    time.Time `json:"date"`
}
