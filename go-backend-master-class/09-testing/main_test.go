package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockUserStore implements UserStore interface for testing
type MockUserStore struct {
	users map[int]User
}

func NewMockUserStore() *MockUserStore {
	return &MockUserStore{
		users: make(map[int]User),
	}
}

func (s *MockUserStore) Create(user User) error {
	s.users[user.ID] = user
	return nil
}

func (s *MockUserStore) Get(id int) (User, error) {
	user, exists := s.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (s *MockUserStore) Update(user User) error {
	s.users[user.ID] = user
	return nil
}

func (s *MockUserStore) Delete(id int) error {
	delete(s.users, id)
	return nil
}

// TestUserService_CreateUser tests the CreateUser method
func TestUserService_CreateUser(t *testing.T) {
	store := NewMockUserStore()
	service := NewUserService(store)

	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid user",
			user: User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			user: User{
				ID:    2,
				Email: "test@example.com",
			},
			wantErr: true,
		},
		{
			name: "missing email",
			user: User{
				ID:       3,
				Username: "testuser",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateUser(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestUserHandler_CreateUser tests the CreateUser handler
func TestUserHandler_CreateUser(t *testing.T) {
	store := NewMockUserStore()
	service := NewUserService(store)
	handler := NewUserHandler(service)

	tests := []struct {
		name           string
		user          User
		expectedCode  int
		expectedError bool
	}{
		{
			name: "valid user",
			user: User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectedCode:  http.StatusCreated,
			expectedError: false,
		},
		{
			name: "invalid user",
			user: User{
				ID: 2,
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userJSON, _ := json.Marshal(tt.user)
			req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedError {
				if w.Code == http.StatusCreated {
					t.Error("Expected error response, got success")
				}
			} else {
				var response User
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}

				if response.ID != tt.user.ID {
					t.Errorf("Expected user ID %d, got %d", tt.user.ID, response.ID)
				}
			}
		})
	}
}

// Benchmark example
func BenchmarkUserService_CreateUser(b *testing.B) {
	store := NewMockUserStore()
	service := NewUserService(store)
	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CreateUser(user)
	}
}

// Example of table-driven tests for GetUser
func TestUserService_GetUser(t *testing.T) {
	store := NewMockUserStore()
	service := NewUserService(store)

	// Setup test data
	testUser := User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	store.Create(testUser)

	tests := []struct {
		name    string
		id      int
		want    User
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      1,
			want:    testUser,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			id:      2,
			want:    User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetUser(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
