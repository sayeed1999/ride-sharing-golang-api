package repository

import "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"

type CustomerRepository interface {
	CreateCustomer(c *domain.Customer) (*domain.Customer, error)
	FindByID(id string) (*domain.Customer, error)
	DeleteCustomer(id string) error
	// UpdateAuthUserID updates the auth_user_id for a customer record
	UpdateAuthUserID(customerID string, authUserID uint) error
}
