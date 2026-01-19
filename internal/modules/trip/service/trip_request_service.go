package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
)

type ITripRequestService interface {
	RequestTrip(customerID uuid.UUID, origin, destination string) (*domain.TripRequest, error)
	CancelTripRequest(ctx context.Context, tripID uuid.UUID, customerID uuid.UUID) error
}

type tripRequestService struct {
	TripRequestRepository repository.ITripRequestRepository
}

func NewTripRequestService(tripRequestRepo repository.ITripRequestRepository) ITripRequestService {
	return &tripRequestService{
		TripRequestRepository: tripRequestRepo,
	}
}

// RequestTrip creates a new trip request.
func (s *tripRequestService) RequestTrip(customerID uuid.UUID, origin, destination string) (*domain.TripRequest, error) {
	tripRequest := &domain.TripRequest{
		CustomerID:  customerID,
		Origin:      origin,
		Destination: destination,
		Status:      domain.NO_DRIVER_FOUND, // default status
	}

	createdTripRequest, err := s.TripRequestRepository.Create(tripRequest)
	if err != nil {
		return nil, err
	}

	return createdTripRequest, nil
}

// CustomerCancelTrip handles the business logic for a customer to cancel a trip request.
func (s *tripRequestService) CancelTripRequest(ctx context.Context, tripID uuid.UUID, customerID uuid.UUID) error {
	tripRequest, err := s.TripRequestRepository.FindByID(tripID)
	if err != nil {
		return err
	}

	// trip request middleware validates that the trip belongs to the customer

	// Only allow cancellation if the trip is in NO_DRIVER_FOUND state
	if tripRequest.Status != domain.NO_DRIVER_FOUND {
		return errors.New("trip cannot be cancelled at this stage")
	}

	return s.TripRequestRepository.UpdateTripRequestStatus(tripID, domain.CUSTOMER_CANCELED)
}
