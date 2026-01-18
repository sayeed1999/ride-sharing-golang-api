package repository

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"

	"gorm.io/gorm"
)

type IDriverRepository interface {
	CreateDriver(d *domain.Driver) (*domain.Driver, error)
	FindByID(id uuid.UUID) (*domain.Driver, error)
	FindByEmail(email string) (*domain.Driver, error)
	DeleteDriver(id uuid.UUID) error
	UpdateAuthUserID(driverID uuid.UUID, authUserID uuid.UUID) error
}

type DriverRepository struct {
	DB *gorm.DB
}

func (r *DriverRepository) CreateDriver(d *domain.Driver) (*domain.Driver, error) {
	if err := r.DB.Create(d).Error; err != nil {
		return nil, err
	}
	return d, nil
}

func (r *DriverRepository) FindByID(id uuid.UUID) (*domain.Driver, error) {
	var d domain.Driver
	if err := r.DB.First(&d, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DriverRepository) FindByEmail(email string) (*domain.Driver, error) {
	var d domain.Driver
	if err := r.DB.First(&d, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DriverRepository) DeleteDriver(id uuid.UUID) error {
	return r.DB.Delete(&domain.Driver{}, "id = ?", id).Error
}

func (r *DriverRepository) UpdateAuthUserID(driverID uuid.UUID, authUserID uuid.UUID) error {
	return r.DB.Model(&domain.Driver{}).Where("id = ?", driverID).Update("auth_user_id", authUserID).Error
}
