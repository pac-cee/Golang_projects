package models

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // "-" means this field won't be included in JSON
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserInput represents the input for user registration
type UserInput struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
}

// LoginInput represents the input for user login
type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
