package controllers

import (
	"context"
	"expense-tracker/backend-vanilla/config"
	"expense-tracker/backend-vanilla/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var input struct {
		Name     string  `json:"name" binding:"required"`
		Email    string  `json:"email" binding:"required,email"`
		Password string  `json:"password" binding:"required,min=6"`
		Limit    float64 `json:"spending_limit"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
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

	// Handle OAuth token exchange and user info fetching based on provider
	var userInfo struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	// TODO: Implement OAuth token exchange and user info fetching for each provider

	// Find or create user
	var user models.User
	err := config.DB.Collection("users").FindOne(context.Background(), bson.M{
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
