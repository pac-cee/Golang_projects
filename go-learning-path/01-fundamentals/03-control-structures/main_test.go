package main

import (
	"strings"
	"testing"
	"time"
)

// TestUserRoles tests the UserRole constants
func TestUserRoles(t *testing.T) {
	// Test role ordering
	if Guest >= User || User >= Moderator || Moderator >= Admin {
		t.Error("User roles are not properly ordered")
	}

	// Test role values
	roles := []UserRole{Guest, User, Moderator, Admin}
	for i, role := range roles {
		if int(role) != i {
			t.Errorf("Expected role value %d, got %d", i, role)
		}
	}
}

// TestCheckUserAccess tests the checkUserAccess function
func TestCheckUserAccess(t *testing.T) {
	tests := []struct {
		name           string
		user          User
		resource      string
		expectedParts []string
	}{
		{
			name: "admin access",
			user: User{
				Name:     "Alice",
				Role:     Admin,
				JoinDate: time.Now(),
			},
			resource: "system",
			expectedParts: []string{
				"New user detected",
				"Admin",
				"full access",
			},
		},
		{
			name: "old user access",
			user: User{
				Name:     "Bob",
				Role:     User,
				JoinDate: time.Now().Add(-48 * time.Hour),
			},
			resource: "system",
			expectedParts: []string{
				"User",
				"read access",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output
			var output strings.Builder
			// Temporarily replace stdout
			oldStdout := fmt.Stdout
			fmt.Stdout = &output
			defer func() { fmt.Stdout = oldStdout }()

			checkUserAccess(tt.user, tt.resource)

			// Check if output contains expected parts
			for _, part := range tt.expectedParts {
				if !strings.Contains(output.String(), part) {
					t.Errorf("Expected output to contain %q", part)
				}
			}
		})
	}
}

// TestSafeOperation tests the safeOperation function
func TestSafeOperation(t *testing.T) {
	tests := []struct {
		name          string
		operation     string
		expectError   bool
		errorContains string
	}{
		{
			name:          "dangerous operation",
			operation:     "dangerous",
			expectError:   true,
			errorContains: "dangerous operation failed",
		},
		{
			name:        "safe operation",
			operation:   "safe",
			expectError: false,
		},
		{
			name:        "risky operation",
			operation:   "risky",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := safeOperation(tt.operation)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, got %v", 
						tt.errorContains, err)
				}
			} else if err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// TestDeferOrder tests the order of defer statements
func TestDeferOrder(t *testing.T) {
	var output strings.Builder
	// Temporarily replace stdout
	oldStdout := fmt.Stdout
	fmt.Stdout = &output
	defer func() { fmt.Stdout = oldStdout }()

	demonstrateDefer()

	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	
	// Check if we have the expected number of lines
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines of output, got %d", len(lines))
	}

	// Check the order of defer statements
	expectedOrder := []string{
		"Demonstrating defer:",
		"Main function body",
		"Third defer",
		"Second defer",
		"First defer",
	}

	for i, expected := range expectedOrder {
		if !strings.Contains(lines[i], expected) {
			t.Errorf("Line %d: expected %q, got %q", i, expected, lines[i])
		}
	}
}

// BenchmarkSafeOperation benchmarks the safeOperation function
func BenchmarkSafeOperation(b *testing.B) {
	operations := []string{"safe", "risky", "dangerous"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := operations[i%len(operations)]
		_ = safeOperation(op)
	}
}
