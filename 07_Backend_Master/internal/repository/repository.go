package repository

import (
	"context"

	"backend-master/internal/model"
)

// Repository defines the interface for data access
type Repository interface {
	// User operations
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, id string, updates *model.UserUpdate) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, offset, limit int) ([]model.User, error)
	CountUsers(ctx context.Context) (int, error)

	// Transaction support
	BeginTx(ctx context.Context) (Repository, error)
	Commit() error
	Rollback() error
}

// Querier interface for database operations
type Querier interface {
	Query(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	Exec(ctx context.Context, query string, args ...interface{}) error
}

// Cache interface for caching operations
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl int) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}
