package handler

import (
	"net/http"

	"backend-master/internal/model"

	"github.com/gin-gonic/gin"
)

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.UserCreate true "User registration data"
// @Success 201 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var input model.UserCreate

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, successResponse(user.SanitizeUser()))
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.UserLogin true "User login data"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 401 {object} response
// @Failure 500 {object} response
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var input model.UserLogin

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, successResponse(gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}))
}

// RefreshToken handles token refresh
// @Summary Refresh tokens
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body map[string]string true "Refresh token"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 401 {object} response
// @Failure 500 {object} response
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accessToken, refreshToken, err := h.service.RefreshToken(c.Request.Context(), input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, successResponse(gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}))
}
