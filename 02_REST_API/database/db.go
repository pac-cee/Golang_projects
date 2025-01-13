package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// InitDB initializes the database connection
func InitDB() (*sql.DB, error) {
	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// Create books table if it doesn't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("error creating tables: %v", err)
	}

	return db, nil
}

// createTables creates necessary database tables
func createTables(db *sql.DB) error {
	// Create users table
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(userTable); err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	// Create books table
	bookTable := `
	CREATE TABLE IF NOT EXISTS books (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		author VARCHAR(100) NOT NULL,
		isbn VARCHAR(13) UNIQUE NOT NULL,
		description TEXT,
		published_year INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		user_id INTEGER REFERENCES users(id)
	);`

	if _, err := db.Exec(bookTable); err != nil {
		return fmt.Errorf("error creating books table: %v", err)
	}

	return nil
}
