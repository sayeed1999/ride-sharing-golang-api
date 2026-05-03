package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"

	"gorm.io/gorm"
)

type ITripRepository interface {
	Create(db *gorm.DB, t *domain.Trip) error
	FindByID(id uuid.UUID) (*domain.Trip, error)
	FindByTripRequestID(tripRequestID uuid.UUID) (*domain.Trip, error)
	UpdateTripStatus(db *gorm.DB, tripID, driverID uuid.UUID, from, to domain.TripStatus) (bool, error)
}

type TripRepository struct {
	DB *gorm.DB
}

func (r *TripRepository) conn(db *gorm.DB) *gorm.DB {
	if db != nil {
		return db
	}
	return r.DB
}

func (r *TripRepository) Create(db *gorm.DB, t *domain.Trip) error {
	return r.conn(db).Create(t).Error
}

func (r *TripRepository) FindByID(id uuid.UUID) (*domain.Trip, error) {
	var row domain.Trip
	if err := r.DB.First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *TripRepository) FindByTripRequestID(tripRequestID uuid.UUID) (*domain.Trip, error) {
	var row domain.Trip
	if err := r.DB.First(&row, "trip_request_id = ?", tripRequestID).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *TripRepository) UpdateTripStatus(db *gorm.DB, tripID, driverID uuid.UUID, from, to domain.TripStatus) (bool, error) {
	res := r.conn(db).Model(&domain.Trip{}).
		Where("id = ? AND driver_id = ? AND status = ?", tripID, driverID, from).
		Updates(map[string]interface{}{"status": to, "updated_at": time.Now()})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}
