package usecase_test

import (
	"testing"

	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository/mocks"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"
)

func TestRegisterUsecase_AddNewUser(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepository()
	registerUC := &usecase.RegisterUsecase{UserRepo: mockRepo}

	initialUserCount := mockRepo.GetUserCount()

	// Test adding a new user
	err := registerUC.Register("newuser@example.com", "password123", "customer")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err.Error())
	}

	// Check user count increased
	if mockRepo.GetUserCount() != 4 { // 3 default + 1 new
		t.Errorf("Expected %d users, but got %d", initialUserCount+1, mockRepo.GetUserCount())
	}

	// Verify user was created
	user, err := mockRepo.FindByEmail("newuser@example.com")
	if err != nil {
		t.Errorf("Expected to find new user, but got error: %s", err.Error())
	}
	if user.Email != "newuser@example.com" {
		t.Errorf("Expected email 'newuser@example.com', but got '%s'", user.Email)
	}
}

func TestRegisterUsecase_AddExistingUser(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepository()
	registerUC := &usecase.RegisterUsecase{UserRepo: mockRepo}

	// Test adding existing user (john@example.com is in default users)
	err := registerUC.Register("john@example.com", "password123", "customer")

	// Assert
	if err == nil {
		t.Error("Expected error for existing user, but got nil")
	}
	if err.Error() != "user already exists" {
		t.Errorf("Expected 'user already exists', but got: %s", err.Error())
	}

	// Check user count didn't increase
	if mockRepo.GetUserCount() != 3 { // Should still be 3 default users
		t.Errorf("Expected 3 users, but got %d", mockRepo.GetUserCount())
	}
}
