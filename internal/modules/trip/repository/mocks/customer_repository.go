package mocks

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
	"github.com/stretchr/testify/mock"
)

type CustomerRepository struct {
	mock.Mock
}

func (m *CustomerRepository) CreateCustomer(c *domain.Customer) (*domain.Customer, error) {
	args := m.Called(c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *CustomerRepository) FindByID(id uuid.UUID) (*domain.Customer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *CustomerRepository) FindByEmail(email string) (*domain.Customer, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *CustomerRepository) DeleteCustomer(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *CustomerRepository) UpdateAuthUserID(customerID uuid.UUID, authUserID uuid.UUID) error {
	args := m.Called(customerID, authUserID)
	return args.Error(0)
}

// proof of mock implementation
var _ repository.ICustomerRepository = (*CustomerRepository)(nil)
