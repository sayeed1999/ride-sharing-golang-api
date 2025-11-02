package repository

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"
)

type TripRequestRepository interface {
	Create(tr *domain.TripRequest) (*domain.TripRequest, error)
	FindByID(id uuid.UUID) (*domain.TripRequest, error)
	Update(tr *domain.TripRequest) (*domain.TripRequest, error)
}
