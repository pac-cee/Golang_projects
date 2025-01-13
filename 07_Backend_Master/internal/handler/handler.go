package handler

import (
	"backend-master/internal/service"
	"backend-master/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Handler holds all HTTP handlers
type Handler struct {
	service service.Service
	logger  *logger.Logger
}

// NewHandler creates a new handler instance
func NewHandler(service service.Service, logger *logger.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// response represents a standard API response
type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// errorResponse creates an error response
func errorResponse(err error) *response {
	return &response{
		Success: false,
		Error:   err.Error(),
	}
}

// successResponse creates a success response
func successResponse(data interface{}) *response {
	return &response{
		Success: true,
		Data:    data,
	}
}

// HealthCheck handles health check requests
// @Summary Health check endpoint
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(200, successResponse(map[string]string{
		"status": "ok",
	}))
}
