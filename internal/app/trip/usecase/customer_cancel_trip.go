package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository"
)

type CustomerCancelTrip struct {
	TripRequestRepo repository.TripRequestRepository
}

func (uc *CustomerCancelTrip) Execute(ctx context.Context, tripID uuid.UUID, customerID uuid.UUID) error {
	tripRequest, err := uc.TripRequestRepo.FindByID(tripID)
	if err != nil {
		return err
	}

	// trip request middleware validates that the trip belongs to the customer

	// Only allow cancellation if the trip is in NO_DRIVER_FOUND state
	if tripRequest.Status != domain.NO_DRIVER_FOUND {
		return errors.New("trip cannot be cancelled at this stage")
	}

	return uc.TripRequestRepo.UpdateTripRequestStatus(tripID, domain.CUSTOMER_CANCELED)
}
