package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Activity represents a student's planned activity
type Activity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string            `bson:"title" json:"title"`
	Description string            `bson:"description" json:"description"`
	Duration    int               `bson:"duration" json:"duration"` // in minutes
	Status      string            `bson:"status" json:"status"`    // planned, in-progress, completed
	StartTime   *time.Time        `bson:"startTime" json:"startTime,omitempty"`
	EndTime     *time.Time        `bson:"endTime" json:"endTime,omitempty"`
	CreatedAt   time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time         `bson:"updatedAt" json:"updatedAt"`
}
