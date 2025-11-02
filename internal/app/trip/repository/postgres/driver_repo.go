package postgres

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"

	"gorm.io/gorm"
)

type DriverRepo struct {
	DB *gorm.DB
}

func (r *DriverRepo) CreateDriver(d *domain.Driver) (*domain.Driver, error) {
	if err := r.DB.Create(d).Error; err != nil {
		return nil, err
	}
	return d, nil
}

func (r *DriverRepo) FindByID(id uuid.UUID) (*domain.Driver, error) {
	var d domain.Driver
	if err := r.DB.First(&d, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DriverRepo) FindByEmail(email string) (*domain.Driver, error) {
	var d domain.Driver
	if err := r.DB.First(&d, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DriverRepo) DeleteDriver(id uuid.UUID) error {
	return r.DB.Delete(&domain.Driver{}, "id = ?", id).Error
}

func (r *DriverRepo) UpdateAuthUserID(driverID uuid.UUID, authUserID uuid.UUID) error {
	return r.DB.Model(&domain.Driver{}).Where("id = ?", driverID).Update("auth_user_id", authUserID).Error
}
