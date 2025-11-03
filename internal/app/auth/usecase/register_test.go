package usecase_test

import (
	"testing"

	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository/mocks"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"
)

func TestRegisterUsecase_AddNewUser(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepository()
	registerUC := &usecase.RegisterUsecase{
		UserRepo:                  mockRepo,
		RequireRoleOnRegistration: true,
	}

	initialUserCount := mockRepo.GetUserCount()

	// Test adding a new user
	_, err := registerUC.Register("newuser@example.com", "password123", "customer")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %s", err.Error())
	}

	// Check user count increased
	if mockRepo.GetUserCount() != initialUserCount+1 {
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
	registerUC := &usecase.RegisterUsecase{
		UserRepo:                  mockRepo,
		RequireRoleOnRegistration: true,
	}

	// Test adding existing user (john@example.com is in default users)
	_, err := registerUC.Register("john@example.com", "password123", "customer")

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

func TestRegisterUsecase_RoleRequiredWhenFeatureEnabled(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepository()
	registerUC := &usecase.RegisterUsecase{
		UserRepo:                  mockRepo,
		RequireRoleOnRegistration: true,
	}

	// Test registration without role
	_, err := registerUC.Register("newuser@example.com", "password123", "")

	// Assert
	if err == nil {
		t.Error("Expected error for missing role, but got nil")
	}
	if err.Error() != "role is required for registration" {
		t.Errorf("Expected 'role is required for registration', but got: %s", err.Error())
	}

	// Test registration with empty/whitespace role
	_, err = registerUC.Register("newuser2@example.com", "password123", "   ")

	// Assert
	if err == nil {
		t.Error("Expected error for whitespace-only role, but got nil")
	}
	if err.Error() != "role is required for registration" {
		t.Errorf("Expected 'role is required for registration', but got: %s", err.Error())
	}
}

func TestRegisterUsecase_RoleNotRequiredWhenFeatureDisabled(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepository()
	registerUC := &usecase.RegisterUsecase{
		UserRepo:                  mockRepo,
		RequireRoleOnRegistration: false,
	}

	initialUserCount := mockRepo.GetUserCount()

	// Test registration without role (should succeed)
	_, err := registerUC.Register("newuser@example.com", "password123", "")

	// Assert
	if err != nil {
		t.Errorf("Expected no error when role not required, but got: %s", err.Error())
	}

	// Check user count increased
	if mockRepo.GetUserCount() != initialUserCount+1 {
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

func TestRegisterUsecase_RoleAssignmentWhenFeatureEnabled(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepository()
	registerUC := &usecase.RegisterUsecase{
		UserRepo:                  mockRepo,
		RequireRoleOnRegistration: true,
	}

	// Test registration with valid role
	_, err := registerUC.Register("newuser@example.com", "password123", "customer")

	// Assert
	if err != nil {
		t.Errorf("Expected no error with valid role, but got: %s", err.Error())
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

func TestRegisterUsecase_RoleAssignmentWhenFeatureDisabled(t *testing.T) {
	// Setup
	mockRepo := mocks.NewMockUserRepository()
	registerUC := &usecase.RegisterUsecase{
		UserRepo:                  mockRepo,
		RequireRoleOnRegistration: false,
	}

	// Test registration with role (should succeed and assign role)
	_, err := registerUC.Register("newuser@example.com", "password123", "customer")

	// Assert
	if err != nil {
		t.Errorf("Expected no error when feature disabled, but got: %s", err.Error())
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
