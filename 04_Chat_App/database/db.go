package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Message struct {
	ID        int64     `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func InitDB() error {
	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Open database connection
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	// Create messages table if it doesn't exist
	if err := createTables(); err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	log.Println("Successfully connected to database")
	return nil
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		room_id VARCHAR(255) NOT NULL,
		user_id VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		type VARCHAR(50) NOT NULL,
		timestamp TIMESTAMP NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_messages_room_id ON messages(room_id);
	CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);
	`

	_, err := db.Exec(query)
	return err
}

func SaveMessage(msg *Message) error {
	query := `
	INSERT INTO messages (room_id, user_id, username, content, type, timestamp)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`

	err := db.QueryRow(
		query,
		msg.RoomID,
		msg.UserID,
		msg.Username,
		msg.Content,
		msg.Type,
		msg.Timestamp,
	).Scan(&msg.ID)

	return err
}

func GetMessages(roomID string, limit int) ([]Message, error) {
	query := `
	SELECT id, room_id, user_id, username, content, type, timestamp
	FROM messages
	WHERE room_id = $1
	ORDER BY timestamp DESC
	LIMIT $2`

	rows, err := db.Query(query, roomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.ID,
			&msg.RoomID,
			&msg.UserID,
			&msg.Username,
			&msg.Content,
			&msg.Type,
			&msg.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func Close() {
	if db != nil {
		db.Close()
	}
}
