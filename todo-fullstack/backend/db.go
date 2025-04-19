package main

import (
	"fmt"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		// Use PostgreSQL if DATABASE_URL is set
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Sprintf("failed to connect to Postgres: %v", err))
		}
		return db
	}
	// Default to SQLite
	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to SQLite")
	}
	return db
}
