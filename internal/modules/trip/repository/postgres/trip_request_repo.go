package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"gorm.io/gorm"
)

type TripRequestRepo struct {
	DB *gorm.DB
}

func (r *TripRequestRepo) Create(tr *domain.TripRequest) (*domain.TripRequest, error) {
	if err := r.DB.Create(tr).Error; err != nil {
		return nil, err
	}
	return tr, nil
}

func (r *TripRequestRepo) FindByID(id uuid.UUID) (*domain.TripRequest, error) {
	var tr domain.TripRequest
	if err := r.DB.First(&tr, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tr, nil
}

func (r *TripRequestRepo) Update(tr *domain.TripRequest) (*domain.TripRequest, error) {
	if err := r.DB.Save(tr).Error; err != nil {
		return nil, err
	}
	return tr, nil
}

func (r *TripRequestRepo) UpdateTripRequestStatus(tripID uuid.UUID, status domain.TripRequestStatus) error {
	return r.DB.Model(&domain.TripRequest{}).Where("id = ?", tripID).Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}
