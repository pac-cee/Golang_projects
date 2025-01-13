package handler

import (
	"context"
	"time"

	"user-service/db"
	pb "user-service/proto"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	jwtSecret string
}

func NewUserServer(jwtSecret string) *UserServer {
	return &UserServer{
		jwtSecret: jwtSecret,
	}
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	// Check if user already exists
	if _, err := db.GetUserByEmail(ctx, req.Email); err == nil {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	// Create user
	user := &db.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := db.CreateUser(ctx, user, req.Password); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	createdAt, err := ptypes.TimestampProto(user.CreatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	updatedAt, err := ptypes.TimestampProto(user.UpdatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	user, err := db.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	createdAt, err := ptypes.TimestampProto(user.CreatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	updatedAt, err := ptypes.TimestampProto(user.UpdatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	// Get existing user
	user, err := db.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	// Prepare updates
	updates := make(map[string]interface{})
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}

	// Update user
	if err := db.UpdateUser(ctx, req.UserId, updates); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	// Get updated user
	user, err = db.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get updated user: %v", err)
	}

	createdAt, err := ptypes.TimestampProto(user.CreatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	updatedAt, err := ptypes.TimestampProto(user.UpdatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (s *UserServer) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	// Get user by email
	user, err := db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invalid email or password")
	}

	// Validate password
	if !db.ValidatePassword(user.PasswordHash, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	// Generate JWT token (implementation not shown)
	token := "dummy-jwt-token" // Replace with actual JWT token generation

	return &pb.AuthenticateUserResponse{
		Token: token,
		User: &pb.User{
			Id:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}, nil
}
