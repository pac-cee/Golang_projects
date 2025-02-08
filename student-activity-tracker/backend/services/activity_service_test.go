package services

import (
	"context"
	"testing"
	"time"

	"github.com/pac-cee/student-activity-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Collection, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	collection := client.Database("test_db").Collection("test_activities")

	return collection, func() {
		collection.Drop(ctx)
		client.Disconnect(ctx)
	}
}

func TestActivityService_CreateActivity(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	service := NewActivityService(collection)
	activity := &models.Activity{
		Title:       "Test Activity",
		Description: "Test Description",
		Duration:    30,
	}

	err := service.CreateActivity(context.Background(), activity)
	if err != nil {
		t.Errorf("Failed to create activity: %v", err)
	}

	if activity.ID.IsZero() {
		t.Error("Activity ID should not be zero")
	}
	if activity.Status != "planned" {
		t.Errorf("Expected status 'planned', got %s", activity.Status)
	}
}

func TestActivityService_GetActivities(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	service := NewActivityService(collection)
	
	// Create test activities
	activities := []models.Activity{
		{
			Title:       "Activity 1",
			Description: "Description 1",
			Duration:    30,
		},
		{
			Title:       "Activity 2",
			Description: "Description 2",
			Duration:    45,
		},
	}

	for _, a := range activities {
		err := service.CreateActivity(context.Background(), &a)
		if err != nil {
			t.Fatalf("Failed to create test activity: %v", err)
		}
	}

	// Test getting activities
	result, err := service.GetActivities(context.Background())
	if err != nil {
		t.Errorf("Failed to get activities: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 activities, got %d", len(result))
	}
}

func TestActivityService_StartActivity(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	service := NewActivityService(collection)
	activity := &models.Activity{
		Title:       "Test Activity",
		Description: "Test Description",
		Duration:    30,
	}

	err := service.CreateActivity(context.Background(), activity)
	if err != nil {
		t.Fatalf("Failed to create activity: %v", err)
	}

	err = service.StartActivity(context.Background(), activity.ID)
	if err != nil {
		t.Errorf("Failed to start activity: %v", err)
	}

	// Verify activity status
	var updatedActivity models.Activity
	err = collection.FindOne(context.Background(), bson.M{"_id": activity.ID}).Decode(&updatedActivity)
	if err != nil {
		t.Fatalf("Failed to get updated activity: %v", err)
	}

	if updatedActivity.Status != "in-progress" {
		t.Errorf("Expected status 'in-progress', got %s", updatedActivity.Status)
	}
	if updatedActivity.StartTime == nil {
		t.Error("Start time should not be nil")
	}
}

func TestActivityService_CompleteActivity(t *testing.T) {
	collection, cleanup := setupTestDB(t)
	defer cleanup()

	service := NewActivityService(collection)
	activity := &models.Activity{
		Title:       "Test Activity",
		Description: "Test Description",
		Duration:    30,
	}

	// Create and start activity
	err := service.CreateActivity(context.Background(), activity)
	if err != nil {
		t.Fatalf("Failed to create activity: %v", err)
	}

	err = service.StartActivity(context.Background(), activity.ID)
	if err != nil {
		t.Fatalf("Failed to start activity: %v", err)
	}

	// Complete activity
	err = service.CompleteActivity(context.Background(), activity.ID)
	if err != nil {
		t.Errorf("Failed to complete activity: %v", err)
	}

	// Verify activity status
	var updatedActivity models.Activity
	err = collection.FindOne(context.Background(), bson.M{"_id": activity.ID}).Decode(&updatedActivity)
	if err != nil {
		t.Fatalf("Failed to get updated activity: %v", err)
	}

	if updatedActivity.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", updatedActivity.Status)
	}
	if updatedActivity.EndTime == nil {
		t.Error("End time should not be nil")
	}
}
