package postgres

import (
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"

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

func (r *CustomerRepo) FindByID(id uint) (*domain.Customer, error) {
	var c domain.Customer
	if err := r.DB.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepo) DeleteCustomer(id uint) error {
	return r.DB.Delete(&domain.Customer{}, id).Error
}

func (r *CustomerRepo) UpdateAuthUserID(customerID uint, authUserID uint) error {
	return r.DB.Model(&domain.Customer{}).Where("id = ?", customerID).Update("auth_user_id", authUserID).Error
}
