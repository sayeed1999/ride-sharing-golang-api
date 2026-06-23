package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	tripdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestStartTrip(t *testing.T) {
	driverID := uuid.New()
	otherDriverID := uuid.New()

	t.Run("happy path: trip becomes in progress", func(t *testing.T) {
		svc, tripRepo := setupTripService()

		tripID := uuid.New()
		tripRequestID := uuid.New()
		acceptedTrip := &tripdomain.Trip{
			ID:            tripID,
			TripRequestID: tripRequestID,
			DriverID:      driverID,
			Status:        tripdomain.TRIP_ACCEPTED,
		}
		startedTrip := &tripdomain.Trip{
			ID:            tripID,
			TripRequestID: tripRequestID,
			DriverID:      driverID,
			Status:        tripdomain.TRIP_IN_PROGRESS,
		}

		tripRepo.On("FindByID", tripID).Return(acceptedTrip, nil).Once()
		tripRepo.On("UpdateTripStatus", mock.Anything, tripID, driverID,
			tripdomain.TRIP_ACCEPTED, tripdomain.TRIP_IN_PROGRESS).Return(true, nil).Once()
		tripRepo.On("FindByID", tripID).Return(startedTrip, nil).Once()

		got, err := svc.StartTrip(context.Background(), driverID, tripID)

		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, tripdomain.TRIP_IN_PROGRESS, got.Status)

		tripRepo.AssertExpectations(t)
	})

	t.Run("trip not found", func(t *testing.T) {
		svc, tripRepo := setupTripService()

		tripID := uuid.New()
		tripRepo.On("FindByID", tripID).Return(nil, gorm.ErrRecordNotFound).Once()

		got, err := svc.StartTrip(context.Background(), driverID, tripID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripNotFound)
		assert.Nil(t, got)
	})

	t.Run("wrong driver", func(t *testing.T) {
		svc, tripRepo := setupTripService()

		tripID := uuid.New()
		trip := &tripdomain.Trip{
			ID:            tripID,
			TripRequestID: uuid.New(),
			DriverID:      otherDriverID,
			Status:        tripdomain.TRIP_ACCEPTED,
		}
		tripRepo.On("FindByID", tripID).Return(trip, nil).Once()

		got, err := svc.StartTrip(context.Background(), driverID, tripID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripWrongDriver)
		assert.Nil(t, got)
		tripRepo.AssertExpectations(t)
	})

	t.Run("invalid state: not accepted", func(t *testing.T) {
		svc, tripRepo := setupTripService()

		tripID := uuid.New()
		trip := &tripdomain.Trip{
			ID:            tripID,
			TripRequestID: uuid.New(),
			DriverID:      driverID,
			Status:        tripdomain.TRIP_IN_PROGRESS,
		}
		tripRepo.On("FindByID", tripID).Return(trip, nil).Once()

		got, err := svc.StartTrip(context.Background(), driverID, tripID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripInvalidState)
		assert.Nil(t, got)
	})

	t.Run("trip status update does not match", func(t *testing.T) {
		svc, tripRepo := setupTripService()

		tripID := uuid.New()
		trip := &tripdomain.Trip{
			ID:            tripID,
			TripRequestID: uuid.New(),
			DriverID:      driverID,
			Status:        tripdomain.TRIP_ACCEPTED,
		}
		tripRepo.On("FindByID", tripID).Return(trip, nil).Once()
		tripRepo.On("UpdateTripStatus", mock.Anything, tripID, driverID,
			tripdomain.TRIP_ACCEPTED, tripdomain.TRIP_IN_PROGRESS).Return(false, nil).Once()

		got, err := svc.StartTrip(context.Background(), driverID, tripID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripStartConflict)
		assert.Nil(t, got)
	})
}

func TestCompleteTrip(t *testing.T) {
	driverID := uuid.New()
	tripID := uuid.New()
	tripRequestID := uuid.New()

	t.Run("happy path", func(t *testing.T) {
		svc, tripRepo := setupTripService()

		inProgress := &tripdomain.Trip{
			ID:            tripID,
			TripRequestID: tripRequestID,
			DriverID:      driverID,
			Status:        tripdomain.TRIP_IN_PROGRESS,
		}
		completed := &tripdomain.Trip{
			ID:            tripID,
			TripRequestID: tripRequestID,
			DriverID:      driverID,
			Status:        tripdomain.TRIP_COMPLETED,
		}

		tripRepo.On("FindByID", tripID).Return(inProgress, nil).Once()
		tripRepo.On("UpdateTripStatus", mock.Anything, tripID, driverID,
			tripdomain.TRIP_IN_PROGRESS, tripdomain.TRIP_COMPLETED).Return(true, nil).Once()
		tripRepo.On("FindByID", tripID).Return(completed, nil).Once()

		got, err := svc.CompleteTrip(context.Background(), driverID, tripID)

		require.NoError(t, err)
		assert.Equal(t, tripdomain.TRIP_COMPLETED, got.Status)
	})

	t.Run("invalid state", func(t *testing.T) {
		svc, tripRepo := setupTripService()

		trip := &tripdomain.Trip{
			ID:            tripID,
			TripRequestID: tripRequestID,
			DriverID:      driverID,
			Status:        tripdomain.TRIP_ACCEPTED,
		}
		tripRepo.On("FindByID", tripID).Return(trip, nil).Once()

		got, err := svc.CompleteTrip(context.Background(), driverID, tripID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripInvalidState)
		assert.Nil(t, got)
	})
}

func TestCancelTripByCustomer(t *testing.T) {
	customerID := uuid.New()
	otherCustomerID := uuid.New()
	tripID := uuid.New()
	tripRequestID := uuid.New()

	t.Run("happy path from accepted", func(t *testing.T) {
		svc, tripRepo := setupTripService()

		trip := &tripdomain.Trip{
			ID:            tripID,
			TripRequestID: tripRequestID,
			CustomerID:    customerID,
			DriverID:      uuid.New(),
			Status:        tripdomain.TRIP_ACCEPTED,
		}
		cancelled := &tripdomain.Trip{ID: tripID, TripRequestID: tripRequestID, CustomerID: customerID, Status: tripdomain.TRIP_CANCELLED_BY_CUSTOMER}

		tripRepo.On("FindByID", tripID).Return(trip, nil).Once()
		tripRepo.On("UpdateTripStatusIf", mock.Anything, tripID,
			tripdomain.TRIP_ACCEPTED, tripdomain.TRIP_CANCELLED_BY_CUSTOMER).Return(true, nil).Once()
		tripRepo.On("FindByID", tripID).Return(cancelled, nil).Once()

		got, err := svc.CancelTripByCustomer(context.Background(), customerID, tripID)

		require.NoError(t, err)
		assert.Equal(t, tripdomain.TRIP_CANCELLED_BY_CUSTOMER, got.Status)
	})

	t.Run("not owned by customer", func(t *testing.T) {
		svc, tripRepo := setupTripService()

		trip := &tripdomain.Trip{ID: tripID, TripRequestID: tripRequestID, CustomerID: otherCustomerID, Status: tripdomain.TRIP_ACCEPTED}

		tripRepo.On("FindByID", tripID).Return(trip, nil).Once()

		got, err := svc.CancelTripByCustomer(context.Background(), customerID, tripID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripNotOwnedByCustomer)
		assert.Nil(t, got)
	})
}

func TestCancelTripByDriver(t *testing.T) {
	driverID := uuid.New()
	tripID := uuid.New()

	t.Run("happy path", func(t *testing.T) {
		svc, tripRepo := setupTripService()

		trip := &tripdomain.Trip{ID: tripID, TripRequestID: uuid.New(), DriverID: driverID, Status: tripdomain.TRIP_ACCEPTED}
		cancelled := &tripdomain.Trip{ID: tripID, DriverID: driverID, Status: tripdomain.TRIP_CANCELLED_BY_DRIVER}

		tripRepo.On("FindByID", tripID).Return(trip, nil).Once()
		tripRepo.On("UpdateTripStatus", mock.Anything, tripID, driverID,
			tripdomain.TRIP_ACCEPTED, tripdomain.TRIP_CANCELLED_BY_DRIVER).Return(true, nil).Once()
		tripRepo.On("FindByID", tripID).Return(cancelled, nil).Once()

		got, err := svc.CancelTripByDriver(context.Background(), driverID, tripID)

		require.NoError(t, err)
		assert.Equal(t, tripdomain.TRIP_CANCELLED_BY_DRIVER, got.Status)
	})

	t.Run("cannot cancel after start", func(t *testing.T) {
		svc, tripRepo := setupTripService()

		trip := &tripdomain.Trip{ID: tripID, DriverID: driverID, Status: tripdomain.TRIP_IN_PROGRESS}
		tripRepo.On("FindByID", tripID).Return(trip, nil).Once()

		got, err := svc.CancelTripByDriver(context.Background(), driverID, tripID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripInvalidState)
		assert.Nil(t, got)
	})
}
