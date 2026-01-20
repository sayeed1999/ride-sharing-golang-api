package repository

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"

	"gorm.io/gorm"
	"time"
)

type ITripRequestRepository interface {
	Create(tr *domain.TripRequest) (*domain.TripRequest, error)
	FindByID(id uuid.UUID) (*domain.TripRequest, error)
	Update(tr *domain.TripRequest) (*domain.TripRequest, error)
	UpdateTripRequestStatus(tripRequestID uuid.UUID, status domain.TripRequestStatus) error
}

type TripRequestRepository struct {
	DB *gorm.DB
}

func (r *TripRequestRepository) Create(tr *domain.TripRequest) (*domain.TripRequest, error) {
	if err := r.DB.Create(tr).Error; err != nil {
		return nil, err
	}
	return tr, nil
}

func (r *TripRequestRepository) FindByID(id uuid.UUID) (*domain.TripRequest, error) {
	var tr domain.TripRequest
	if err := r.DB.First(&tr, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tr, nil
}

func (r *TripRequestRepository) Update(tr *domain.TripRequest) (*domain.TripRequest, error) {
	if err := r.DB.Save(tr).Error; err != nil {
		return nil, err
	}
	return tr, nil
}

func (r *TripRequestRepository) UpdateTripRequestStatus(tripRequestID uuid.UUID, status domain.TripRequestStatus) error {
	return r.DB.Model(&domain.TripRequest{}).Where("id = ?", tripRequestID).Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}
