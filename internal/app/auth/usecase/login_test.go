package usecase_test

import (
	"testing"

	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository/mocks"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"
)

func TestLoginUsecase_RightLogin(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepository()
	loginUC := &usecase.LoginUsecase{UserRepo: mockRepo}

	// Test correct login with existing user and correct password
	err := loginUC.Login("john@example.com", "password123")

	// Assert
	if err != nil {
		t.Errorf("Expected no error for correct login, but got: %s", err.Error())
	}
}

func TestLoginUsecase_WrongLogin(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepository()
	loginUC := &usecase.LoginUsecase{UserRepo: mockRepo}

	// Test wrong password
	err := loginUC.Login("john@example.com", "wrongpassword")

	// Assert
	if err == nil {
		t.Error("Expected error for wrong password, but got nil")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials', but got: %s", err.Error())
	}
}

func TestLoginUsecase_NonExistentUser(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepository()
	loginUC := &usecase.LoginUsecase{UserRepo: mockRepo}

	// Test login with non-existent user
	err := loginUC.Login("nonexistent@example.com", "password123")

	// Assert
	if err == nil {
		t.Error("Expected error for non-existent user, but got nil")
	}
	if err.Error() != "invalid credentials" {
		t.Errorf("Expected 'invalid credentials', but got: %s", err.Error())
	}
}

func TestLoginUsecase_MultipleUsersCorrectPasswords(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepository()
	loginUC := &usecase.LoginUsecase{UserRepo: mockRepo}

	// Test all users with their correct passwords
	testCases := []struct {
		email    string
		password string
	}{
		{"john@example.com", "password123"},
		{"jane@example.com", "password456"},
		{"admin@example.com", "password789"},
	}

	for _, tc := range testCases {
		err := loginUC.Login(tc.email, tc.password)
		if err != nil {
			t.Errorf("Expected no error for user %s with correct password, but got: %s", tc.email, err.Error())
		}
	}
}
