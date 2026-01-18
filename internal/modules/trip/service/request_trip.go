package service

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
)

// RequestTripService handles the business logic for a customer to request a trip.
type RequestTripService struct {
	TripRequestRepository repository.ITripRequestRepository
}

// Execute creates a new trip request.
func (uc *RequestTripService) Execute(customerID uuid.UUID, origin, destination string) (*domain.TripRequest, error) {
	tripRequest := &domain.TripRequest{
		CustomerID:  customerID,
		Origin:      origin,
		Destination: destination,
		Status:      domain.NO_DRIVER_FOUND, // default status
	}

	createdTripRequest, err := uc.TripRequestRepository.Create(tripRequest)
	if err != nil {
		return nil, err
	}

	return createdTripRequest, nil
}
