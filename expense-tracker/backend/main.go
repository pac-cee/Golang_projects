package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"expense-tracker/handlers"
)

var client *mongo.Client
var db *mongo.Database

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "expense_tracker"
	}
	db = client.Database(dbName)

	// Initialize handlers
	transactionHandler := handlers.NewTransactionHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	budgetHandler := handlers.NewBudgetHandler(db)
	reportHandler := handlers.NewReportHandler(db)

	// Initialize Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Initialize routes
	api := router.Group("/api")
	{
		// Transaction routes
		api.GET("/transactions", transactionHandler.GetTransactions)
		api.GET("/transactions/:id", transactionHandler.GetTransaction)
		api.POST("/transactions", transactionHandler.CreateTransaction)
		api.PUT("/transactions/:id", transactionHandler.UpdateTransaction)
		api.DELETE("/transactions/:id", transactionHandler.DeleteTransaction)

		// Category routes
		api.GET("/categories", categoryHandler.GetCategories)
		api.GET("/categories/:id", categoryHandler.GetCategory)
		api.POST("/categories", categoryHandler.CreateCategory)
		api.PUT("/categories/:id", categoryHandler.UpdateCategory)
		api.DELETE("/categories/:id", categoryHandler.DeleteCategory)

		// Budget routes
		api.GET("/budgets", budgetHandler.GetBudgets)
		api.GET("/budgets/:id", budgetHandler.GetBudget)
		api.POST("/budgets", budgetHandler.CreateBudget)
		api.PUT("/budgets/:id", budgetHandler.UpdateBudget)
		api.DELETE("/budgets/:id", budgetHandler.DeleteBudget)

		// Report routes
		api.GET("/reports/transactions", reportHandler.GetTransactionReport)
		api.GET("/reports/categories", reportHandler.GetCategoryReport)
		api.GET("/reports/budgets", reportHandler.GetBudgetReport)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
