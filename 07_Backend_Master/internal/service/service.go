package service

import (
	"context"

	"backend-master/internal/model"
	"backend-master/internal/repository"
	"backend-master/pkg/auth"
	"backend-master/pkg/logger"
)

// Service defines the interface for business logic
type Service interface {
	// User operations
	CreateUser(ctx context.Context, input model.UserCreate) (*model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, id string, input model.UserUpdate) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, page, pageSize int) ([]model.User, int, error)

	// Auth operations
	Login(ctx context.Context, input model.UserLogin) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	ValidateToken(ctx context.Context, token string) (*auth.TokenClaims, error)
}

// Dependencies holds service dependencies
type Dependencies struct {
	Repo         repository.Repository
	TokenManager *auth.TokenManager
	Logger       *logger.Logger
}

// service implements the Service interface
type service struct {
	repo         repository.Repository
	tokenManager *auth.TokenManager
	logger       *logger.Logger
}

// NewService creates a new service instance
func NewService(deps Dependencies) Service {
	return &service{
		repo:         deps.Repo,
		tokenManager: deps.TokenManager,
		logger:       deps.Logger,
	}
}
