package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware is a placeholder for JWT authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement JWT authentication
		// On failure:
		// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		// return
		c.Next()
	}
}
