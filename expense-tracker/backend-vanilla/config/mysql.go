package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLDB holds the MySQL database instance
type MySQLDB struct {
	*sql.DB
}

// NewMySQLDB creates a new MySQL database connection
func NewMySQLDB() (*MySQLDB, error) {
	// Get database connection details from environment variables
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbname := os.Getenv("MYSQL_DB")

	// Create connection string
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, port, dbname)

	// Open database connection
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mysql database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging mysql database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	log.Println("Successfully connected to MySQL database")
	return &MySQLDB{db}, nil
}

// CreateTables creates all necessary tables for the expense tracker
func (db *MySQLDB) CreateTables() error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			name VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_email (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	if err != nil {
		return fmt.Errorf("error creating users table: %v", err)
	}

	// Create categories table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			budget DECIMAL(10,2),
			color VARCHAR(7),
			icon VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			INDEX idx_user_id (user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	if err != nil {
		return fmt.Errorf("error creating categories table: %v", err)
	}

	// Create expenses table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS expenses (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT,
			category_id INT,
			amount DECIMAL(10,2) NOT NULL,
			description TEXT,
			date DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
			INDEX idx_user_id (user_id),
			INDEX idx_category_id (category_id),
			INDEX idx_date (date)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	if err != nil {
		return fmt.Errorf("error creating expenses table: %v", err)
	}

	// Create budget_alerts table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS budget_alerts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT,
			category_id INT,
			threshold DECIMAL(10,2) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
			INDEX idx_user_id (user_id),
			INDEX idx_category_id (category_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	if err != nil {
		return fmt.Errorf("error creating budget_alerts table: %v", err)
	}

	// Create user_settings table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS user_settings (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT UNIQUE,
			currency VARCHAR(3) DEFAULT 'USD',
			theme VARCHAR(10) DEFAULT 'light',
			notification_enabled BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			INDEX idx_user_id (user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)
	if err != nil {
		return fmt.Errorf("error creating user_settings table: %v", err)
	}

	return nil
}

func (mdb *MySQLDB) Exec(s string) (any, any) {
	panic("unimplemented")
}
