package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()

	r := gin.Default()

	// Public routes
	r.POST("/register", registerHandler)
	r.POST("/login", loginHandler)

	// Authenticated routes
	auth := r.Group("/api")
	auth.Use(AuthMiddleware())
	{
		auth.GET("/habits", getHabitsHandler)
		auth.POST("/habits", createHabitHandler)
		auth.PUT("/habits/:id", updateHabitHandler)
		auth.DELETE("/habits/:id", deleteHabitHandler)
		auth.POST("/habits/:id/mark", markHabitHandler)
	}

	r.Run(":8080")
}
