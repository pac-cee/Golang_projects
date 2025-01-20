package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/expense-tracker/models"
)

type ReportHandler struct {
	db *mongo.Database
}

func NewReportHandler(db *mongo.Database) *ReportHandler {
	return &ReportHandler{
		db: db,
	}
}

func (h *ReportHandler) GetTransactionReport(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse date range
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		start = time.Now().AddDate(0, -1, 0) // Default to last month
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		end = time.Now() // Default to now
	}

	// Build date filter
	dateFilter := bson.M{
		"date": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	// Get transactions
	transactionsColl := h.db.Collection("transactions")
	cursor, err := transactionsColl.Find(ctx, dateFilter, options.Find().SetSort(bson.D{{Key: "date", Value: -1}}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate totals
	var totalIncome, totalExpenses float64
	categoryMap := make(map[string]*models.CategoryExpenseSummary)

	for _, t := range transactions {
		if t.Amount >= 0 {
			totalIncome += t.Amount
		} else {
			totalExpenses += -t.Amount
			// Track category expenses
			if summary, exists := categoryMap[t.Category]; exists {
				summary.Amount += -t.Amount
				summary.Count++
			} else {
				categoryMap[t.Category] = &models.CategoryExpenseSummary{
					Category: t.Category,
					Amount:   -t.Amount,
					Count:    1,
				}
			}
		}
	}

	// Convert category map to slice
	categories := make([]models.CategoryExpenseSummary, 0, len(categoryMap))
	for _, summary := range categoryMap {
		categories = append(categories, *summary)
	}

	report := models.TransactionReport{
		TotalIncome:   totalIncome,
		TotalExpenses: totalExpenses,
		NetAmount:     totalIncome - totalExpenses,
		Transactions:  transactions,
		Categories:    categories,
	}

	c.JSON(http.StatusOK, report)
}

func (h *ReportHandler) GetCategoryReport(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		start = time.Now().AddDate(0, -1, 0)
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		end = time.Now()
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"date": bson.M{
					"$gte": start,
					"$lte": end,
				},
				"amount": bson.M{"$lt": 0}, // Only expenses
			},
		},
		{
			"$group": bson.M{
				"_id": "$category",
				"amount": bson.M{
					"$sum": bson.M{
						"$abs": "$amount",
					},
				},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$project": bson.M{
				"category": "$_id",
				"amount":   1,
				"count":    1,
				"_id":      0,
			},
		},
	}

	cursor, err := h.db.Collection("transactions").Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var categories []models.CategoryExpenseSummary
	if err = cursor.All(ctx, &categories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *ReportHandler) GetBudgetReport(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all budgets
	budgetsColl := h.db.Collection("budgets")
	cursor, err := budgetsColl.Find(ctx, bson.M{})
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

	// Calculate current spending for each budget
	var budgetReports []models.BudgetReport
	for _, budget := range budgets {
		// Calculate date range based on budget period
		var startDate time.Time
		endDate := time.Now()

		switch budget.Period {
		case "weekly":
			startDate = time.Now().AddDate(0, 0, -7)
		case "monthly":
			startDate = time.Now().AddDate(0, -1, 0)
		case "yearly":
			startDate = time.Now().AddDate(-1, 0, 0)
		}

		// Get total expenses for this category in the period
		pipeline := []bson.M{
			{
				"$match": bson.M{
					"category": budget.Category,
					"amount":   bson.M{"$lt": 0},
					"date": bson.M{
						"$gte": startDate,
						"$lte": endDate,
					},
				},
			},
			{
				"$group": bson.M{
					"_id": nil,
					"total": bson.M{
						"$sum": bson.M{
							"$abs": "$amount",
						},
					},
				},
			},
		}

		cursor, err := h.db.Collection("transactions").Aggregate(ctx, pipeline)
		if err != nil {
			continue
		}

		var result []bson.M
		if err = cursor.All(ctx, &result); err != nil {
			continue
		}

		var spentAmount float64
		if len(result) > 0 {
			spentAmount = result[0]["total"].(float64)
		}

		percentage := (spentAmount / budget.Amount) * 100

		budgetReports = append(budgetReports, models.BudgetReport{
			Category:     budget.Category,
			BudgetAmount: budget.Amount,
			SpentAmount:  spentAmount,
			Percentage:   percentage,
			StartDate:    startDate,
			EndDate:      endDate,
		})
	}

	c.JSON(http.StatusOK, budgetReports)
}
