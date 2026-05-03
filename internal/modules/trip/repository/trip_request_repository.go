package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"gorm.io/gorm"
)

type ITripRequestRepository interface {
	Create(tr *domain.TripRequest) (*domain.TripRequest, error)
	FindByID(id uuid.UUID) (*domain.TripRequest, error)
	Update(tr *domain.TripRequest) (*domain.TripRequest, error)
	UpdateTripRequestStatus(tripRequestID uuid.UUID, status domain.TripRequestStatus) error
	UpdateTripRequestStatusIf(db *gorm.DB, tripRequestID uuid.UUID, currentStatus, newStatus domain.TripRequestStatus) (bool, error)
	ListOpenTripRequests(limit int) ([]domain.TripRequest, error)
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

func (r *TripRequestRepository) conn(db *gorm.DB) *gorm.DB {
	if db != nil {
		return db
	}
	return r.DB
}

func (r *TripRequestRepository) UpdateTripRequestStatusIf(db *gorm.DB, tripRequestID uuid.UUID, currentStatus, newStatus domain.TripRequestStatus) (bool, error) {
	res := r.conn(db).Model(&domain.TripRequest{}).
		Where("id = ? AND status = ?", tripRequestID, currentStatus).
		Updates(map[string]interface{}{"status": newStatus, "updated_at": time.Now()})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

func (r *TripRequestRepository) ListOpenTripRequests(limit int) ([]domain.TripRequest, error) {
	var rows []domain.TripRequest
	err := r.DB.Where("status = ?", domain.NO_DRIVER_FOUND).
		Order("created_at DESC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
