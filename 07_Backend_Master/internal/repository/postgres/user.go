package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"backend-master/internal/model"
	"backend-master/pkg/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

// UserRepository implements user-related database operations
type UserRepository struct {
	db *database.PostgreSQL
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.PostgreSQL) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (
			id, email, password_hash, first_name, last_name,
			role, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`

	err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Password, user.FirstName,
		user.LastName, user.Role, user.Status, user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		if database.IsUniqueViolation(err) {
			return fmt.Errorf("email already exists: %w", err)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name,
			   role, status, created_at, updated_at
		FROM users
		WHERE id = $1 AND status != 'deleted'`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName,
		&user.LastName, &user.Role, &user.Status, &user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name,
			   role, status, created_at, updated_at
		FROM users
		WHERE email = $1 AND status != 'deleted'`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName,
		&user.LastName, &user.Role, &user.Status, &user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(ctx context.Context, id string, updates *model.UserUpdate) error {
	query := `
		UPDATE users
		SET updated_at = NOW()`

	args := []interface{}{id}
	argCount := 1

	if updates.Email != nil {
		argCount++
		query += fmt.Sprintf(", email = $%d", argCount)
		args = append(args, *updates.Email)
	}

	if updates.Password != nil {
		argCount++
		query += fmt.Sprintf(", password_hash = $%d", argCount)
		args = append(args, *updates.Password)
	}

	if updates.FirstName != nil {
		argCount++
		query += fmt.Sprintf(", first_name = $%d", argCount)
		args = append(args, *updates.FirstName)
	}

	if updates.LastName != nil {
		argCount++
		query += fmt.Sprintf(", last_name = $%d", argCount)
		args = append(args, *updates.LastName)
	}

	if updates.Status != nil {
		argCount++
		query += fmt.Sprintf(", status = $%d", argCount)
		args = append(args, *updates.Status)
	}

	query += ` WHERE id = $1 AND status != 'deleted'`

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if database.IsUniqueViolation(err) {
			return fmt.Errorf("email already exists: %w", err)
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser soft deletes a user
func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET status = 'deleted', updated_at = NOW()
		WHERE id = $1 AND status != 'deleted'`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ListUsers retrieves a paginated list of users
func (r *UserRepository) ListUsers(ctx context.Context, offset, limit int) ([]model.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name,
			   role, status, created_at, updated_at
		FROM users
		WHERE status != 'deleted'
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Password, &user.FirstName,
			&user.LastName, &user.Role, &user.Status, &user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (r *UserRepository) CountUsers(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM users
		WHERE status != 'deleted'`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// BeginTx starts a new transaction
func (r *UserRepository) BeginTx(ctx context.Context) (*UserRepository, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &UserRepository{db: tx}, nil
}

// Commit commits the transaction
func (r *UserRepository) Commit() error {
	if tx, ok := r.db.(pgx.Tx); ok {
		return tx.Commit(context.Background())
	}
	return nil
}

// Rollback rolls back the transaction
func (r *UserRepository) Rollback() error {
	if tx, ok := r.db.(pgx.Tx); ok {
		return tx.Rollback(context.Background())
	}
	return nil
}
