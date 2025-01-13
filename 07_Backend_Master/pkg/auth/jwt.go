package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// TokenManager handles JWT token operations
type TokenManager struct {
	secretKey      string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// TokenClaims represents JWT claims
type TokenClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.StandardClaims
}

// NewTokenManager creates a new JWT token manager
func NewTokenManager(secretKey string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{
		secretKey:   secretKey,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (m *TokenManager) GenerateTokenPair(userID, email, role string) (accessToken, refreshToken string, err error) {
	// Generate access token
	accessToken, err = m.GenerateToken(userID, email, role, "access", m.accessTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err = m.GenerateToken(userID, email, role, "refresh", m.refreshTTL)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// GenerateToken generates a new JWT token
func (m *TokenManager) GenerateToken(userID, email, role, tokenType string, ttl time.Duration) (string, error) {
	claims := TokenClaims{
		UserID:    userID,
		Email:     email,
		Role:      role,
		TokenType: tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken validates and parses a JWT token
func (m *TokenManager) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrExpiredToken
			}
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshTokens validates a refresh token and generates new token pair
func (m *TokenManager) RefreshTokens(refreshToken string) (string, string, error) {
	// Validate refresh token
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.TokenType != "refresh" {
		return "", "", errors.New("invalid token type")
	}

	// Generate new token pair
	accessToken, newRefreshToken, err := m.GenerateTokenPair(
		claims.UserID,
		claims.Email,
		claims.Role,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new tokens: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

// ExtractUserID extracts user ID from token
func (m *TokenManager) ExtractUserID(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// ExtractRole extracts user role from token
func (m *TokenManager) ExtractRole(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}
