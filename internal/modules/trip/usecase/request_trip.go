package usecase

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
)

// RequestTripUsecase handles the business logic for a customer to request a trip.
type RequestTripUsecase struct {
	TripRequestRepo repository.TripRequestRepository
}

// Execute creates a new trip request.
func (uc *RequestTripUsecase) Execute(customerID uuid.UUID, origin, destination string) (*domain.TripRequest, error) {
	tripRequest := &domain.TripRequest{
		CustomerID:  customerID,
		Origin:      origin,
		Destination: destination,
		Status:      domain.NO_DRIVER_FOUND, // default status
	}

	createdTripRequest, err := uc.TripRequestRepo.Create(tripRequest)
	if err != nil {
		return nil, err
	}

	return createdTripRequest, nil
}
