package service

import (
	"context"
	"fmt"

	"backend-master/internal/model"
)

// CreateUser creates a new user
func (s *service) CreateUser(ctx context.Context, input model.UserCreate) (*model.User, error) {
	// Create user instance
	user := &model.User{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      input.Role,
		Status:    "active",
	}

	// Hash password
	if err := user.HashPassword(input.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Start transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create user
	if err := tx.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("User created", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
	})

	return user, nil
}

// GetUser retrieves a user by ID
func (s *service) GetUser(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// UpdateUser updates user information
func (s *service) UpdateUser(ctx context.Context, id string, input model.UserUpdate) error {
	// Start transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current user
	user, err := tx.GetUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Hash password if provided
	if input.Password != nil {
		if err := user.HashPassword(*input.Password); err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		input.Password = &user.Password
	}

	// Update user
	if err := tx.UpdateUser(ctx, id, &input); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("User updated", map[string]interface{}{
		"user_id": id,
	})

	return nil
}

// DeleteUser deletes a user
func (s *service) DeleteUser(ctx context.Context, id string) error {
	// Start transaction
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete user
	if err := tx.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("User deleted", map[string]interface{}{
		"user_id": id,
	})

	return nil
}

// ListUsers retrieves a paginated list of users
func (s *service) ListUsers(ctx context.Context, page, pageSize int) ([]model.User, int, error) {
	offset := (page - 1) * pageSize

	// Get users
	users, err := s.repo.ListUsers(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	// Get total count
	total, err := s.repo.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, total, nil
}
