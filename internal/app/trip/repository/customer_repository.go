package repository

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"
)

type CustomerRepository interface {
	CreateCustomer(c *domain.Customer) (*domain.Customer, error)
	FindByID(id uuid.UUID) (*domain.Customer, error)
	FindByEmail(email string) (*domain.Customer, error)
	DeleteCustomer(id uuid.UUID) error
	// UpdateAuthUserID updates the auth_user_id for a customer record
	UpdateAuthUserID(customerID uuid.UUID, authUserID uuid.UUID) error
}
