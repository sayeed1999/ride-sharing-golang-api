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
	"gorm.io/gorm"
)

func TestAcceptTripRequest(t *testing.T) {
	driverID := uuid.New()
	customerID := uuid.New()

	t.Run("happy path: accepts open trip request and creates trip", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripService()

		tripRequestID := uuid.New()
		open := &tripdomain.TripRequest{
			ID:          tripRequestID,
			CustomerID:  customerID,
			Origin:      testTripOrigin,
			Destination: testTripDestination,
			Status:      tripdomain.NO_DRIVER_FOUND,
		}
		accepted := &tripdomain.TripRequest{
			ID:          tripRequestID,
			CustomerID:  customerID,
			Origin:      testTripOrigin,
			Destination: testTripDestination,
			Status:      tripdomain.DRIVER_ACCEPTED,
		}

		tripRequestRepo.On("FindByID", tripRequestID).Return(open, nil).Once()
		tripRequestRepo.On("UpdateTripRequestStatusIf", mock.Anything, tripRequestID,
			tripdomain.NO_DRIVER_FOUND, tripdomain.DRIVER_ACCEPTED).Return(true, nil).Once()
		tripRepo.On("Create", mock.Anything, mock.MatchedBy(func(tr *tripdomain.Trip) bool {
			return tr.TripRequestID == tripRequestID && tr.DriverID == driverID && tr.Status == tripdomain.TRIP_ACCEPTED
		})).Return(nil).Once()
		tripRequestRepo.On("FindByID", tripRequestID).Return(accepted, nil).Once()

		trip, trAfter, err := svc.AcceptTripRequest(context.Background(), driverID, tripRequestID)

		require.NoError(t, err)
		require.NotNil(t, trip)
		require.NotNil(t, trAfter)
		assert.Equal(t, tripRequestID, trip.TripRequestID)
		assert.Equal(t, driverID, trip.DriverID)
		assert.Equal(t, tripdomain.TRIP_ACCEPTED, trip.Status)
		assert.Equal(t, tripdomain.DRIVER_ACCEPTED, trAfter.Status)

		tripRequestRepo.AssertExpectations(t)
		tripRepo.AssertExpectations(t)
	})

	t.Run("trip request not found", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripService()

		tripRequestID := uuid.New()
		tripRequestRepo.On("FindByID", tripRequestID).Return(nil, gorm.ErrRecordNotFound).Once()

		trip, tr, err := svc.AcceptTripRequest(context.Background(), driverID, tripRequestID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripRequestNotFound)
		assert.Nil(t, trip)
		assert.Nil(t, tr)
		tripRequestRepo.AssertExpectations(t)
		tripRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("trip request not open", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripService()

		tripRequestID := uuid.New()
		tr := &tripdomain.TripRequest{
			ID:          tripRequestID,
			CustomerID:  customerID,
			Status:      tripdomain.DRIVER_ACCEPTED,
			Origin:      testTripOrigin,
			Destination: testTripDestination,
		}
		tripRequestRepo.On("FindByID", tripRequestID).Return(tr, nil).Once()

		trip, trAfter, err := svc.AcceptTripRequest(context.Background(), driverID, tripRequestID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripRequestNotOpen)
		assert.Nil(t, trip)
		assert.Nil(t, trAfter)
		tripRequestRepo.AssertExpectations(t)
		tripRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("concurrent accept: conditional update does not match", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripService()

		tripRequestID := uuid.New()
		open := &tripdomain.TripRequest{
			ID:          tripRequestID,
			CustomerID:  customerID,
			Origin:      testTripOrigin,
			Destination: testTripDestination,
			Status:      tripdomain.NO_DRIVER_FOUND,
		}
		tripRequestRepo.On("FindByID", tripRequestID).Return(open, nil).Once()
		tripRequestRepo.On("UpdateTripRequestStatusIf", mock.Anything, tripRequestID,
			tripdomain.NO_DRIVER_FOUND, tripdomain.DRIVER_ACCEPTED).Return(false, nil).Once()

		trip, trAfter, err := svc.AcceptTripRequest(context.Background(), driverID, tripRequestID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripRequestNotOpen)
		assert.Nil(t, trip)
		assert.Nil(t, trAfter)
		tripRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
		tripRequestRepo.AssertExpectations(t)
	})

	t.Run("conditional update fails with repository error", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripService()

		tripRequestID := uuid.New()
		open := fixtureTripRequest(customerID)
		open.ID = tripRequestID

		repoErr := errors.New("db update failed")
		tripRequestRepo.On("FindByID", tripRequestID).Return(open, nil).Once()
		tripRequestRepo.On("UpdateTripRequestStatusIf", mock.Anything, tripRequestID,
			tripdomain.NO_DRIVER_FOUND, tripdomain.DRIVER_ACCEPTED).Return(false, repoErr).Once()

		trip, trAfter, err := svc.AcceptTripRequest(context.Background(), driverID, tripRequestID)

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, trip)
		assert.Nil(t, trAfter)
		tripRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("create trip fails", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripService()

		tripRequestID := uuid.New()
		open := fixtureTripRequest(customerID)
		open.ID = tripRequestID

		createErr := errors.New("create failed")
		tripRequestRepo.On("FindByID", tripRequestID).Return(open, nil).Once()
		tripRequestRepo.On("UpdateTripRequestStatusIf", mock.Anything, tripRequestID,
			tripdomain.NO_DRIVER_FOUND, tripdomain.DRIVER_ACCEPTED).Return(true, nil).Once()
		tripRepo.On("Create", mock.Anything, mock.Anything).Return(createErr).Once()

		trip, trAfter, err := svc.AcceptTripRequest(context.Background(), driverID, tripRequestID)

		require.Error(t, err)
		assert.ErrorIs(t, err, createErr)
		assert.Nil(t, trip)
		assert.Nil(t, trAfter)
		tripRequestRepo.AssertExpectations(t)
		tripRepo.AssertExpectations(t)
	})
}

func TestStartTrip(t *testing.T) {
	driverID := uuid.New()
	otherDriverID := uuid.New()

	t.Run("happy path: trip becomes in progress and request is started", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripService()

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
		tripRequestRepo.On("UpdateTripRequestStatusIf", mock.Anything, tripRequestID,
			tripdomain.DRIVER_ACCEPTED, tripdomain.TRIP_STARTED).Return(true, nil).Once()
		tripRepo.On("FindByID", tripID).Return(startedTrip, nil).Once()

		got, err := svc.StartTrip(context.Background(), driverID, tripID)

		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, tripdomain.TRIP_IN_PROGRESS, got.Status)

		tripRepo.AssertExpectations(t)
		tripRequestRepo.AssertExpectations(t)
	})

	t.Run("trip not found", func(t *testing.T) {
		svc, _, tripRepo := setupTripService()

		tripID := uuid.New()
		tripRepo.On("FindByID", tripID).Return(nil, gorm.ErrRecordNotFound).Once()

		got, err := svc.StartTrip(context.Background(), driverID, tripID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripNotFound)
		assert.Nil(t, got)
	})

	t.Run("wrong driver", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripService()

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
		tripRequestRepo.AssertNotCalled(t, "UpdateTripRequestStatusIf", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("invalid state: not accepted", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripService()

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
		tripRequestRepo.AssertNotCalled(t, "UpdateTripRequestStatusIf", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("trip status update does not match", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripService()

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
		tripRequestRepo.AssertNotCalled(t, "UpdateTripRequestStatusIf", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("trip request status update does not match", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripService()

		tripID := uuid.New()
		tripRequestID := uuid.New()
		trip := &tripdomain.Trip{
			ID:            tripID,
			TripRequestID: tripRequestID,
			DriverID:      driverID,
			Status:        tripdomain.TRIP_ACCEPTED,
		}
		tripRepo.On("FindByID", tripID).Return(trip, nil).Once()
		tripRepo.On("UpdateTripStatus", mock.Anything, tripID, driverID,
			tripdomain.TRIP_ACCEPTED, tripdomain.TRIP_IN_PROGRESS).Return(true, nil).Once()
		tripRequestRepo.On("UpdateTripRequestStatusIf", mock.Anything, tripRequestID,
			tripdomain.DRIVER_ACCEPTED, tripdomain.TRIP_STARTED).Return(false, nil).Once()

		got, err := svc.StartTrip(context.Background(), driverID, tripID)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripStartConflict)
		assert.Nil(t, got)
	})
}
