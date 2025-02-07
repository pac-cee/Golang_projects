package main

import (
	"testing"
)

// TestCalculateRank tests the rank calculation function
func TestCalculateRank(t *testing.T) {
	tests := []struct {
		name     string
		score    int
		expected float64
	}{
		{
			name:     "perfect score",
			score:    100,
			expected: 1.0,
		},
		{
			name:     "zero score",
			score:    0,
			expected: 0.0,
		},
		{
			name:     "partial score",
			score:    75,
			expected: 0.75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateRank(tt.score)
			if result != tt.expected {
				t.Errorf("calculateRank(%d) = %f, want %f", 
					tt.score, result, tt.expected)
			}
		})
	}
}

// TestUserScore tests the UserScore struct
func TestUserScore(t *testing.T) {
	// Test struct initialization
	user := UserScore{
		Username: "TestUser",
		Score:    85,
		Active:   true,
	}

	// Test Username
	if user.Username != "TestUser" {
		t.Errorf("Expected username 'TestUser', got %s", user.Username)
	}

	// Test Score
	if user.Score != 85 {
		t.Errorf("Expected score 85, got %d", user.Score)
	}

	// Test Active status
	if !user.Active {
		t.Error("Expected user to be active")
	}

	// Test Rank calculation
	user.Rank = calculateRank(user.Score)
	expectedRank := 0.85
	if user.Rank != expectedRank {
		t.Errorf("Expected rank %f, got %f", expectedRank, user.Rank)
	}
}

// TestGameLevels tests the game level constants
func TestGameLevels(t *testing.T) {
	// Test level values
	if BeginnerLevel != 0 {
		t.Errorf("Expected BeginnerLevel to be 0, got %d", BeginnerLevel)
	}

	if IntermediateLevel != 1 {
		t.Errorf("Expected IntermediateLevel to be 1, got %d", IntermediateLevel)
	}

	if AdvancedLevel != 2 {
		t.Errorf("Expected AdvancedLevel to be 2, got %d", AdvancedLevel)
	}

	if ExpertLevel != 3 {
		t.Errorf("Expected ExpertLevel to be 3, got %d", ExpertLevel)
	}

	// Test level ordering
	if BeginnerLevel >= IntermediateLevel {
		t.Error("BeginnerLevel should be less than IntermediateLevel")
	}

	if IntermediateLevel >= AdvancedLevel {
		t.Error("IntermediateLevel should be less than AdvancedLevel")
	}

	if AdvancedLevel >= ExpertLevel {
		t.Error("AdvancedLevel should be less than ExpertLevel")
	}
}

// BenchmarkCalculateRank benchmarks the rank calculation
func BenchmarkCalculateRank(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateRank(85)
	}
}
