package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"expense-tracker/models"
)

type BudgetHandler struct {
	collection *mongo.Collection
}

func NewBudgetHandler(db *mongo.Database) *BudgetHandler {
	return &BudgetHandler{
		collection: db.Collection("budgets"),
	}
}

func (h *BudgetHandler) GetBudgets(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := h.collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var budgets []models.Budget
	if err = cursor.All(ctx, &budgets); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budgets)
}

func (h *BudgetHandler) GetBudget(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var budget models.Budget
	err = h.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&budget)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budget)
}

func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var budget models.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set initial values
	budget.CreatedAt = time.Now()
	budget.UpdatedAt = time.Now()
	budget.Spent = 0 // Initialize spent amount to 0

	// Validate budget period
	switch budget.Period {
	case "weekly", "monthly", "yearly":
		// Valid period
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget period. Must be 'weekly', 'monthly', or 'yearly'"})
		return
	}

	// Set start and end dates based on period if not provided
	if budget.StartDate.IsZero() {
		budget.StartDate = time.Now()
	}
	if budget.EndDate.IsZero() {
		switch budget.Period {
		case "weekly":
			budget.EndDate = budget.StartDate.AddDate(0, 0, 7)
		case "monthly":
			budget.EndDate = budget.StartDate.AddDate(0, 1, 0)
		case "yearly":
			budget.EndDate = budget.StartDate.AddDate(1, 0, 0)
		}
	}

	result, err := h.collection.InsertOne(ctx, budget)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	budget.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, budget)
}

func (h *BudgetHandler) UpdateBudget(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var budget models.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	budget.UpdatedAt = time.Now()

	// Validate budget period if provided
	if budget.Period != "" {
		switch budget.Period {
		case "weekly", "monthly", "yearly":
			// Valid period
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget period. Must be 'weekly', 'monthly', or 'yearly'"})
			return
		}
	}

	update := bson.M{
		"$set": budget,
	}

	result := h.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Err().Error()})
		return
	}

	if err := result.Decode(&budget); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, budget)
}

func (h *BudgetHandler) DeleteBudget(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := h.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted successfully"})
}
