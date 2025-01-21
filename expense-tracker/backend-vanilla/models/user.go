package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name           string            `bson:"name" json:"name"`
	Email          string            `bson:"email" json:"email"`
	Password       string            `bson:"password" json:"-"`
	SpendingLimit  float64          `bson:"spending_limit" json:"spending_limit"`
	Provider       string            `bson:"provider,omitempty" json:"provider,omitempty"`
	ProviderID     string            `bson:"provider_id,omitempty" json:"provider_id,omitempty"`
	CreatedAt      time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time         `bson:"updated_at" json:"updated_at"`
}

type UserResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	SpendingLimit float64   `json:"spending_limit"`
	Provider      string    `json:"provider,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
