package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// PostgresDB holds the PostgreSQL database instance
type PostgresDB struct {
	*sql.DB
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB() (*PostgresDB, error) {
	// Get database connection details from environment variables
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging postgres database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Successfully connected to PostgreSQL database")
	return &PostgresDB{db}, nil
}

// CreateTables creates all necessary tables for the expense tracker
func (db *PostgresDB) CreateTables() error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	// Create categories table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			budget DECIMAL(10,2),
			color VARCHAR(7),
			icon VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating categories table: %v", err)
	}

	// Create expenses table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			category_id INTEGER REFERENCES categories(id),
			amount DECIMAL(10,2) NOT NULL,
			description TEXT,
			date DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating expenses table: %v", err)
	}

	// Create budget_alerts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS budget_alerts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			category_id INTEGER REFERENCES categories(id),
			threshold DECIMAL(10,2) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating budget_alerts table: %v", err)
	}

	// Create user_settings table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_settings (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) UNIQUE,
			currency VARCHAR(3) DEFAULT 'USD',
			theme VARCHAR(10) DEFAULT 'light',
			notification_enabled BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating user_settings table: %v", err)
	}

	return nil
}
