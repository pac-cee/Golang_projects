package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Amount      float64           `json:"amount" bson:"amount"`
	Type        string            `json:"type" bson:"type"` // "income" or "expense"
	Category    string            `json:"category" bson:"category"`
	Subcategory string            `json:"subcategory,omitempty" bson:"subcategory,omitempty"`
	Account     string            `json:"account" bson:"account"`
	Description string            `json:"description" bson:"description"`
	Date        time.Time         `json:"date" bson:"date"`
	CreatedAt   time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt" bson:"updatedAt"`
}

type Category struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string            `json:"name" bson:"name"`
	Subcategories []string          `json:"subcategories" bson:"subcategories"`
	CreatedAt     time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt" bson:"updatedAt"`
}

type Budget struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Category  string            `json:"category" bson:"category"`
	Amount    float64           `json:"amount" bson:"amount"`
	Spent     float64           `json:"spent" bson:"spent"`
	Period    string            `json:"period" bson:"period"` // "weekly", "monthly", "yearly"
	StartDate time.Time         `json:"startDate" bson:"startDate"`
	EndDate   time.Time         `json:"endDate" bson:"endDate"`
	CreatedAt time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt" bson:"updatedAt"`
}

type Account struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string            `json:"name" bson:"name"`
	Type      string            `json:"type" bson:"type"` // "bank", "cash", "mobile_money"
	Balance   float64           `json:"balance" bson:"balance"`
	CreatedAt time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt" bson:"updatedAt"`
}

type TransactionReport struct {
	TotalIncome   float64                  `json:"totalIncome"`
	TotalExpenses float64                  `json:"totalExpenses"`
	NetAmount     float64                  `json:"netAmount"`
	Transactions  []Transaction            `json:"transactions"`
	Categories    []CategoryExpenseSummary `json:"categories"`
}

type CategoryExpenseSummary struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Count    int     `json:"count"`
}

type BudgetReport struct {
	Category     string    `json:"category"`
	BudgetAmount float64   `json:"budgetAmount"`
	SpentAmount  float64   `json:"spentAmount"`
	Percentage   float64   `json:"percentage"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
}
