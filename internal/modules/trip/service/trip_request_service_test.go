package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	tripdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRequestTripByCustomer(t *testing.T) {
	t.Run("happy path: creates trip request with default status", func(t *testing.T) {
		svc, tripRequestRepo := setupTripRequestService()

		customerID := uuid.New()
		expected := fixtureTripRequest(customerID)

		tripRequestRepo.On("Create", mock.Anything).Return(expected, nil).Once()

		created, err := svc.RequestTrip(customerID, testTripOrigin, testTripDestination)

		require.NoError(t, err)
		require.NotNil(t, created)
		assert.Equal(t, expected.ID, created.ID)
		assert.Equal(t, tripdomain.NO_DRIVER_FOUND, created.Status)

		tripRequestRepo.AssertExpectations(t)
	})

	t.Run("create fails: returns repository error", func(t *testing.T) {
		svc, tripRequestRepo := setupTripRequestService()

		customerID := uuid.New()
		repoErr := errors.New("db create failed")
		tripRequestRepo.On("Create", mock.Anything).Return(nil, repoErr).Once()

		created, err := svc.RequestTrip(customerID, testTripOrigin, testTripDestination)

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, created)
		tripRequestRepo.AssertExpectations(t)
	})
}

func TestCancelTripRequestByCustomer(t *testing.T) {
	t.Run("happy path: cancels trip in no-driver-found stage", func(t *testing.T) {
		svc, tripRequestRepo := setupTripRequestService()

		tripReq := &tripdomain.TripRequest{
			ID:     uuid.New(),
			Status: tripdomain.NO_DRIVER_FOUND,
		}

		tripRequestRepo.On("UpdateTripRequestStatus", tripReq.ID, tripdomain.CUSTOMER_CANCELED).Return(nil).Once()

		err := svc.CancelTripRequest(context.Background(), tripReq)

		require.NoError(t, err)
		tripRequestRepo.AssertExpectations(t)
	})

	t.Run("non cancellable status: returns stage error", func(t *testing.T) {
		svc, tripRequestRepo := setupTripRequestService()

		tripReq := &tripdomain.TripRequest{
			ID:     uuid.New(),
			Status: tripdomain.DRIVER_ACCEPTED,
		}

		err := svc.CancelTripRequest(context.Background(), tripReq)

		require.Error(t, err)
		assert.EqualError(t, err, "trip cannot be cancelled at this stage")
		tripRequestRepo.AssertNotCalled(t, "UpdateTripRequestStatus", mock.Anything, mock.Anything)
	})
}

