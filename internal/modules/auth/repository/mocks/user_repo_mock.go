package mocks

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/pkg/password"
)

// MockUserRepository is a simple mock implementation for testing
type MockUserRepository struct {
	users []domain.User
	roles []domain.Role
}

// NewMockUserRepository creates a new mock with 3 default users
func NewMockUserRepository() *MockUserRepository {
	// Generate real password hashes using existing functions
	hash1, salt1, _ := generateUserPassword("password123")
	hash2, salt2, _ := generateUserPassword("password456")
	hash3, salt3, _ := generateUserPassword("password789")

	return &MockUserRepository{
		users: []domain.User{
			{ID: uuid.New(), Email: "john@example.com", PasswordHash: hash1, PasswordSalt: salt1},
			{ID: uuid.New(), Email: "jane@example.com", PasswordHash: hash2, PasswordSalt: salt2},
			{ID: uuid.New(), Email: "admin@example.com", PasswordHash: hash3, PasswordSalt: salt3},
		},
		roles: []domain.Role{
			{ID: uuid.New(), Name: "customer"},
			{ID: uuid.New(), Name: "driver"},
			{ID: uuid.New(), Name: "admin"},
		},
	}
}

// generateUserPassword is a helper that uses existing password functions
func generateUserPassword(pass string) (string, string, error) {
	salt, err := password.GenerateSalt()
	if err != nil {
		return "", "", err
	}
	hash, err := password.HashPassword(pass, salt)
	if err != nil {
		return "", "", err
	}
	return hash, salt, nil
}

// CreateUser adds a new user to the mock
func (m *MockUserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	// Check if user already exists
	for _, existingUser := range m.users {
		if existingUser.Email == user.Email {
			return nil, errors.New("user already exists")
		}
	}

	// Assign ID and add user
	user.ID = uuid.New()
	m.users = append(m.users, *user)
	return user, nil
}

// FindByEmail finds a user by email
func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

// AssignRole assigns a role to a user (simplified - just returns success)
func (m *MockUserRepository) AssignRole(userID uuid.UUID, roleName string) (*domain.UserRole, error) {
	// Find the role
	var roleID uuid.UUID
	for _, role := range m.roles {
		if role.Name == roleName {
			roleID = role.ID
			break
		}
	}
	if roleID == uuid.Nil {
		return nil, errors.New(fmt.Sprintf("role %s not found", roleName))
	}

	// Return a mock UserRole
	return &domain.UserRole{
		ID:     uuid.New(),
		UserID: userID,
		RoleID: roleID,
	}, nil
}

// DeleteUser removes a user from the mock by ID
func (m *MockUserRepository) DeleteUser(userID uuid.UUID) error {
	for i, user := range m.users {
		if user.ID == userID {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found")
}

// GetUserCount returns the number of users
func (m *MockUserRepository) GetUserCount() int {
	return len(m.users)
}

// Ensure MockUserRepository implements IUserRepository interface
var _ repository.IUserRepository = (*MockUserRepository)(nil)
