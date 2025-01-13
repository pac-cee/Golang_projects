package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"book-api/models"

	"github.com/gorilla/mux"
	"github.com/go-playground/validator/v10"
)

type BookHandler struct {
	db        *sql.DB
	validator *validator.Validate
}

func NewBookHandler(db *sql.DB) *BookHandler {
	return &BookHandler{
		db:        db,
		validator: validator.New(),
	}
}

// GetBooks returns all books
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("userID").(int)

	// Query books
	rows, err := h.db.Query("SELECT id, title, author, isbn, description, published_year, created_at, updated_at FROM books WHERE user_id = $1", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Description, &book.PublishedYear, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

// GetBook returns a specific book
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	var book models.Book
	err = h.db.QueryRow("SELECT id, title, author, isbn, description, published_year, created_at, updated_at FROM books WHERE id = $1 AND user_id = $2",
		id, userID).Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Description, &book.PublishedYear, &book.CreatedAt, &book.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// CreateBook creates a new book
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var input models.BookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	var book models.Book
	err := h.db.QueryRow(
		"INSERT INTO books (title, author, isbn, description, published_year, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at",
		input.Title, input.Author, input.ISBN, input.Description, input.PublishedYear, userID,
	).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	book.Title = input.Title
	book.Author = input.Author
	book.ISBN = input.ISBN
	book.Description = input.Description
	book.PublishedYear = input.PublishedYear
	book.UserID = userID

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// UpdateBook updates an existing book
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var input models.BookInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	result, err := h.db.Exec(
		"UPDATE books SET title = $1, author = $2, isbn = $3, description = $4, published_year = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6 AND user_id = $7",
		input.Title, input.Author, input.ISBN, input.Description, input.PublishedYear, id, userID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteBook deletes a book
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int)

	result, err := h.db.Exec("DELETE FROM books WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
