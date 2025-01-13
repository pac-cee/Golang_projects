package models

import "time"

// Book represents a book in the system
type Book struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	ISBN          string    `json:"isbn"`
	Description   string    `json:"description"`
	PublishedYear int       `json:"published_year"`
	UserID        int       `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BookInput represents the input for creating or updating a book
type BookInput struct {
	Title         string `json:"title" validate:"required"`
	Author        string `json:"author" validate:"required"`
	ISBN          string `json:"isbn" validate:"required,len=13"`
	Description   string `json:"description"`
	PublishedYear int    `json:"published_year" validate:"required,min=1000,max=9999"`
}
