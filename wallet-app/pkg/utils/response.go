package utils

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
    Message string `json:"message"`
}

func ErrorResponse(c *gin.Context, status int, message string) {
    c.JSON(status, ErrorResponse{Message: message})
}