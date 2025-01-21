package controllers

import (
	"context"
	"expense-tracker/backend-vanilla/config"
	"expense-tracker/backend-vanilla/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateExpense(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)

	var input struct {
		Amount      float64   `json:"amount" binding:"required"`
		Category    string    `json:"category" binding:"required"`
		Description string    `json:"description"`
		Date        time.Time `json:"date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense := models.Expense{
		UserID:      userID,
		Amount:      input.Amount,
		Category:    input.Category,
		Description: input.Description,
		Date:        input.Date,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := config.DB.Collection("expenses").InsertOne(context.Background(), expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	expense.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, models.ExpenseResponse{
		ID:          expense.ID.Hex(),
		Amount:      expense.Amount,
		Category:    expense.Category,
		Description: expense.Description,
		Date:        expense.Date,
		CreatedAt:   expense.CreatedAt,
	})
}

func GetExpenses(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)

	var query struct {
		StartDate string `form:"start_date"`
		EndDate   string `form:"end_date"`
		Category  string `form:"category"`
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"user_id": userID}

	if query.StartDate != "" && query.EndDate != "" {
		startDate, _ := time.Parse(time.RFC3339, query.StartDate)
		endDate, _ := time.Parse(time.RFC3339, query.EndDate)
		filter["date"] = bson.M{
			"$gte": startDate,
			"$lte": endDate,
		}
	}

	if query.Category != "" {
		filter["category"] = query.Category
	}

	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})

	cursor, err := config.DB.Collection("expenses").Find(context.Background(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}
	defer cursor.Close(context.Background())

	var expenses []models.Expense
	if err := cursor.All(context.Background(), &expenses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode expenses"})
		return
	}

	response := make([]models.ExpenseResponse, len(expenses))
	for i, expense := range expenses {
		response[i] = models.ExpenseResponse{
			ID:          expense.ID.Hex(),
			Amount:      expense.Amount,
			Category:    expense.Category,
			Description: expense.Description,
			Date:        expense.Date,
			CreatedAt:   expense.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

func UpdateExpense(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	expenseID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	var input struct {
		Amount      float64   `json:"amount"`
		Category    string    `json:"category"`
		Description string    `json:"description"`
		Date        time.Time `json:"date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"amount":      input.Amount,
			"category":    input.Category,
			"description": input.Description,
			"date":        input.Date,
			"updated_at":  time.Now(),
		},
	}

	result := config.DB.Collection("expenses").FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": expenseID, "user_id": userID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var expense models.Expense
	if err := result.Decode(&expense); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, models.ExpenseResponse{
		ID:          expense.ID.Hex(),
		Amount:      expense.Amount,
		Category:    expense.Category,
		Description: expense.Description,
		Date:        expense.Date,
		CreatedAt:   expense.CreatedAt,
	})
}

func DeleteExpense(c *gin.Context) {
	userID := c.MustGet("user_id").(primitive.ObjectID)
	expenseID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	result, err := config.DB.Collection("expenses").DeleteOne(
		context.Background(),
		bson.M{"_id": expenseID, "user_id": userID},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
