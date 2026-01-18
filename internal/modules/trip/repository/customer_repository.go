package repository

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"

	"gorm.io/gorm"
)

type ICustomerRepository interface {
	CreateCustomer(c *domain.Customer) (*domain.Customer, error)
	FindByID(id uuid.UUID) (*domain.Customer, error)
	FindByEmail(email string) (*domain.Customer, error)
	DeleteCustomer(id uuid.UUID) error
	// UpdateAuthUserID updates the auth_user_id for a customer record
	UpdateAuthUserID(customerID uuid.UUID, authUserID uuid.UUID) error
}

type CustomerRepository struct {
	DB *gorm.DB
}

func (r *CustomerRepository) CreateCustomer(c *domain.Customer) (*domain.Customer, error) {
	if err := r.DB.Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CustomerRepository) FindByID(id uuid.UUID) (*domain.Customer, error) {
	var c domain.Customer
	if err := r.DB.First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepository) FindByEmail(email string) (*domain.Customer, error) {
	var c domain.Customer
	if err := r.DB.First(&c, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepository) DeleteCustomer(id uuid.UUID) error {
	return r.DB.Delete(&domain.Customer{}, "id = ?", id).Error
}

func (r *CustomerRepository) UpdateAuthUserID(customerID uuid.UUID, authUserID uuid.UUID) error {
	return r.DB.Model(&domain.Customer{}).Where("id = ?", customerID).Update("auth_user_id", authUserID).Error
}
