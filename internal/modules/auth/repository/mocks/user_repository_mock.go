package mocks

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/repository"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepository) AssignRole(userID uuid.UUID, roleName string) (*domain.UserRole, error) {
	args := m.Called(userID, roleName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserRole), args.Error(1)
}

func (m *UserRepository) DeleteUser(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

var _ repository.IUserRepository = (*UserRepository)(nil)

