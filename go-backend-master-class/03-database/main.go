package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// User represents a user in our system
type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
}

// Database configuration
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "gobackend"
)

// UserRepository handles all database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Initialize database connection
func initDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Create tables
func (r *UserRepository) createTables() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	_, err := r.db.Exec(query)
	return err
}

// CreateUser inserts a new user
func (r *UserRepository) CreateUser(username, email string) (*User, error) {
	query := `
		INSERT INTO users (username, email)
		VALUES ($1, $2)
		RETURNING id, username, email, created_at`

	user := &User{}
	err := r.db.QueryRow(query, username, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (r *UserRepository) GetUser(id int) (*User, error) {
	user := &User{}
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a user's email
func (r *UserRepository) UpdateUser(id int, email string) error {
	query := `
		UPDATE users
		SET email = $2
		WHERE id = $1`

	result, err := r.db.Exec(query, id, email)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

// DeleteUser removes a user
func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

func main() {
	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Create repository
	repo := NewUserRepository(db)

	// Create tables
	err = repo.createTables()
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}

	// Example usage
	user, err := repo.CreateUser("johndoe", "john@example.com")
	if err != nil {
		log.Printf("Error creating user: %v", err)
	} else {
		log.Printf("Created user: %+v", user)
	}

	// Get user
	retrievedUser, err := repo.GetUser(user.ID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		log.Printf("Retrieved user: %+v", retrievedUser)
	}

	// Update user
	err = repo.UpdateUser(user.ID, "john.doe@example.com")
	if err != nil {
		log.Printf("Error updating user: %v", err)
	}

	// Delete user
	err = repo.DeleteUser(user.ID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
	}
}
