package repository

import (
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/domain"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUser(user *domain.User) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	AssignRole(userID uint, roleName string) (*domain.UserRole, error)
}
