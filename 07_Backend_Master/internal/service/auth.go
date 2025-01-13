package service

import (
	"context"
	"fmt"

	"backend-master/internal/model"
	"backend-master/pkg/auth"
)

// Login authenticates a user and returns access and refresh tokens
func (s *service) Login(ctx context.Context, input model.UserLogin) (string, string, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// Check password
	if !user.CheckPassword(input.Password) {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// Check user status
	if user.Status != "active" {
		return "", "", fmt.Errorf("user account is %s", user.Status)
	}

	// Generate tokens
	accessToken, refreshToken, err := s.tokenManager.GenerateTokenPair(
		user.ID,
		user.Email,
		user.Role,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	s.logger.Info("User logged in", map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	})

	return accessToken, refreshToken, nil
}

// RefreshToken validates refresh token and returns new token pair
func (s *service) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Validate and refresh tokens
	accessToken, newRefreshToken, err := s.tokenManager.RefreshTokens(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to refresh tokens: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

// ValidateToken validates a token and returns its claims
func (s *service) ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	claims, err := s.tokenManager.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
