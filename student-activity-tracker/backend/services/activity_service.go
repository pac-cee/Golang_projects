package services

import (
	"context"
	"time"

	"github.com/pac-cee/student-activity-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ActivityService handles business logic for activities
type ActivityService struct {
	collection *mongo.Collection
}

// NewActivityService creates a new ActivityService
func NewActivityService(collection *mongo.Collection) *ActivityService {
	return &ActivityService{
		collection: collection,
	}
}

// CreateActivity creates a new activity
func (s *ActivityService) CreateActivity(ctx context.Context, activity *models.Activity) error {
	activity.ID = primitive.NewObjectID()
	activity.Status = "planned"
	activity.CreatedAt = time.Now()
	activity.UpdatedAt = time.Now()

	_, err := s.collection.InsertOne(ctx, activity)
	return err
}

// GetActivities retrieves all activities
func (s *ActivityService) GetActivities(ctx context.Context) ([]models.Activity, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var activities []models.Activity
	if err := cursor.All(ctx, &activities); err != nil {
		return nil, err
	}

	return activities, nil
}

// GetActivity retrieves a single activity by ID
func (s *ActivityService) GetActivity(ctx context.Context, id primitive.ObjectID) (*models.Activity, error) {
	var activity models.Activity
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&activity)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

// UpdateActivity updates an existing activity
func (s *ActivityService) UpdateActivity(ctx context.Context, id primitive.ObjectID, activity *models.Activity) error {
	activity.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"title":       activity.Title,
			"description": activity.Description,
			"duration":    activity.Duration,
			"updatedAt":   activity.UpdatedAt,
		},
	}

	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// DeleteActivity deletes an activity
func (s *ActivityService) DeleteActivity(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// StartActivity starts an activity
func (s *ActivityService) StartActivity(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":    "in-progress",
			"startTime": now,
			"updatedAt": now,
		},
	}

	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// CompleteActivity completes an activity
func (s *ActivityService) CompleteActivity(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":    "completed",
			"endTime":   now,
			"updatedAt": now,
		},
	}

	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}
