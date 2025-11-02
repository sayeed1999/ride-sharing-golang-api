package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCustomerCancelTrip_Execute(t *testing.T) {
	mockRepo := new(mocks.TripRequestRepository)
	uc := &CustomerCancelTrip{TripRequestRepo: mockRepo}

	ctx := context.Background()

	t.Run("successfully cancels trip", func(t *testing.T) {
		tripID := uuid.New()
		customerID := uuid.New()
		tripRequest := &domain.TripRequest{
			ID:         tripID,
			CustomerID: customerID,
			Status:     domain.NO_DRIVER_FOUND,
		}

		mockRepo.On("FindByID", tripID).Return(tripRequest, nil).Once()
		mockRepo.On("UpdateTripRequestStatus", tripID, domain.CUSTOMER_CANCELED).Return(nil).Once()

		err := uc.Execute(ctx, tripID, customerID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error if trip not found", func(t *testing.T) {
		tripID := uuid.New()
		customerID := uuid.New()

		mockRepo.On("FindByID", tripID).Return(nil, errors.New("not found")).Once()

		err := uc.Execute(ctx, tripID, customerID)
		assert.Error(t, err)
		assert.EqualError(t, err, "not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error if customer ID does not match", func(t *testing.T) {
		tripID := uuid.New()
		customerID := uuid.New()
		wrongCustomerID := uuid.New()
		tripRequest := &domain.TripRequest{
			ID:         tripID,
			CustomerID: customerID,
			Status:     domain.NO_DRIVER_FOUND,
		}

		mockRepo.On("FindByID", tripID).Return(tripRequest, nil).Once()

		err := uc.Execute(ctx, tripID, wrongCustomerID)
		assert.Error(t, err)
		assert.EqualError(t, err, "trip request does not belong to customer")
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error if trip status is not NO_DRIVER_FOUND", func(t *testing.T) {
		tripID := uuid.New()
		customerID := uuid.New()
		tripRequest := &domain.TripRequest{
			ID:         tripID,
			CustomerID: customerID,
			Status:     domain.DRIVER_ACCEPTED,
		}

		mockRepo.On("FindByID", tripID).Return(tripRequest, nil).Once()

		err := uc.Execute(ctx, tripID, customerID)
		assert.Error(t, err)
		assert.EqualError(t, err, "trip cannot be cancelled at this stage")
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error if repository update fails", func(t *testing.T) {
		tripID := uuid.New()
		customerID := uuid.New()
		tripRequest := &domain.TripRequest{
			ID:         tripID,
			CustomerID: customerID,
			Status:     domain.NO_DRIVER_FOUND,
		}

		mockRepo.On("FindByID", tripID).Return(tripRequest, nil).Once()
		mockRepo.On("UpdateTripRequestStatus", tripID, domain.CUSTOMER_CANCELED).Return(errors.New("db error")).Once()

		err := uc.Execute(ctx, tripID, customerID)
		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}
