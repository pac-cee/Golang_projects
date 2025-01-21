package controllers

import (
	"context"
	"expense-tracker/backend-vanilla/config"
	"expense-tracker/backend-vanilla/models"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Password validation rules
var passwordRules = struct {
	MinLength      int
	MaxLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
	AllowedSpecial string
}{
	MinLength:      8,
	MaxLength:      32,
	RequireUpper:   true,
	RequireLower:   true,
	RequireNumber:  true,
	RequireSpecial: true,
	AllowedSpecial: "!@#$%^&*(),.?\":{}|<>",
}

func validatePassword(password string) error {
	if len(password) < passwordRules.MinLength {
		return fmt.Errorf("password must be at least %d characters long", passwordRules.MinLength)
	}

	if len(password) > passwordRules.MaxLength {
		return fmt.Errorf("password must be less than %d characters", passwordRules.MaxLength)
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case strings.ContainsRune(passwordRules.AllowedSpecial, char):
			hasSpecial = true
		}
	}

	if passwordRules.RequireUpper && !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if passwordRules.RequireLower && !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if passwordRules.RequireNumber && !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}

	if passwordRules.RequireSpecial && !hasSpecial {
		return fmt.Errorf("password must contain at least one special character (%s)", passwordRules.AllowedSpecial)
	}

	return nil
}

func Register(c *gin.Context) {
	var input struct {
		Name     string  `json:"name" binding:"required"`
		Email    string  `json:"email" binding:"required,email"`
		Password string  `json:"password" binding:"required"`
		Limit    float64 `json:"spending_limit"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password
	if err := validatePassword(input.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var existingUser models.User
	err := config.DB.Collection("users").FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		Name:          input.Name,
		Email:         input.Email,
		Password:      string(hashedPassword),
		SpendingLimit: input.Limit,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	result, err := config.DB.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, models.UserResponse{
		ID:            user.ID.Hex(),
		Name:          user.Name,
		Email:         user.Email,
		SpendingLimit: user.SpendingLimit,
		CreatedAt:     user.CreatedAt,
	})
}

func Login(c *gin.Context) {
	var input struct {
		Email      string `json:"email" binding:"required,email"`
		Password   string `json:"password" binding:"required"`
		RememberMe bool   `json:"remember_me"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := config.DB.Collection("users").FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": models.UserResponse{
			ID:            user.ID.Hex(),
			Name:          user.Name,
			Email:         user.Email,
			SpendingLimit: user.SpendingLimit,
			CreatedAt:     user.CreatedAt,
		},
	})
}

func SocialCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code required"})
		return
	}

	// Handle OAuth token exchange and user info fetching based on provider
	var userInfo struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// Exchange the authorization code for access token
	var accessToken string
	var err error

	switch provider {
	case "google":
		accessToken, err = exchangeGoogleCode(code)
	case "github":
		accessToken, err = exchangeGithubCode(code)
	case "twitter":
		accessToken, err = exchangeTwitterCode(code)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange authorization code"})
		return
	}

	// Fetch user info using the access token
	switch provider {
	case "google":
		userInfo, err = fetchGoogleUserInfo(accessToken)
	case "github":
		userInfo, err = fetchGithubUserInfo(accessToken)
	case "twitter":
		userInfo, err = fetchTwitterUserInfo(accessToken)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	// Find or create user
	var user models.User
	err = config.DB.Collection("users").FindOne(context.Background(), bson.M{
		"provider":    provider,
		"provider_id": userInfo.ID,
	}).Decode(&user)

	if err != nil {
		// Create new user
		user = models.User{
			Name:       userInfo.Name,
			Email:      userInfo.Email,
			Provider:   provider,
			ProviderID: userInfo.ID,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		result, err := config.DB.Collection("users").InsertOne(context.Background(), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		user.ID = result.InsertedID.(primitive.ObjectID)
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": models.UserResponse{
			ID:            user.ID.Hex(),
			Name:          user.Name,
			Email:         user.Email,
			SpendingLimit: user.SpendingLimit,
			Provider:      user.Provider,
			CreatedAt:     user.CreatedAt,
		},
	})
}

// Helper functions for OAuth token exchange
func exchangeGoogleCode(code string) (string, error) {
	// TODO: Implement Google OAuth token exchange
	return "", nil
}

func exchangeGithubCode(code string) (string, error) {
	// TODO: Implement GitHub OAuth token exchange
	return "", nil
}

func exchangeTwitterCode(code string) (string, error) {
	// TODO: Implement Twitter OAuth token exchange
	return "", nil
}

// Helper functions for fetching user info
func fetchGoogleUserInfo(accessToken string) (struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}, error) {
	// TODO: Implement Google user info fetching
	return struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{}, nil
}

func fetchGithubUserInfo(accessToken string) (struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}, error) {
	// TODO: Implement GitHub user info fetching
	return struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{}, nil
}

func fetchTwitterUserInfo(accessToken string) (struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}, error) {
	// TODO: Implement Twitter user info fetching
	return struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{}, nil
}
