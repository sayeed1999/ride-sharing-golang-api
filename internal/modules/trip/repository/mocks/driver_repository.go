package mocks

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
	"github.com/stretchr/testify/mock"
)

type DriverRepository struct {
	mock.Mock
}

func (m *DriverRepository) CreateDriver(d *domain.Driver) (*domain.Driver, error) {
	args := m.Called(d)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *DriverRepository) FindByID(id uuid.UUID) (*domain.Driver, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *DriverRepository) FindByEmail(email string) (*domain.Driver, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Driver), args.Error(1)
}

func (m *DriverRepository) DeleteDriver(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *DriverRepository) UpdateAuthUserID(driverID uuid.UUID, authUserID uuid.UUID) error {
	args := m.Called(driverID, authUserID)
	return args.Error(0)
}

// proof of mock implementation
var _ repository.IDriverRepository = (*DriverRepository)(nil)
