package postgres

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"

	"gorm.io/gorm"
)

type CustomerRepo struct {
	DB *gorm.DB
}

func (r *CustomerRepo) CreateCustomer(c *domain.Customer) (*domain.Customer, error) {
	if err := r.DB.Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CustomerRepo) FindByID(id uuid.UUID) (*domain.Customer, error) {
	var c domain.Customer
	if err := r.DB.First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepo) FindByEmail(email string) (*domain.Customer, error) {
	var c domain.Customer
	if err := r.DB.First(&c, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepo) DeleteCustomer(id uuid.UUID) error {
	return r.DB.Delete(&domain.Customer{}, "id = ?", id).Error
}

func (r *CustomerRepo) UpdateAuthUserID(customerID uuid.UUID, authUserID uuid.UUID) error {
	return r.DB.Model(&domain.Customer{}).Where("id = ?", customerID).Update("auth_user_id", authUserID).Error
}
