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

func TestRequestTripByCustomer(t *testing.T) {
	t.Run("happy path: creates trip request with default status", func(t *testing.T) {
		svc, tripRequestRepo, _ := setupTripRequestService()

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
		svc, tripRequestRepo, _ := setupTripRequestService()

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
		svc, tripRequestRepo, _ := setupTripRequestService()

		tripReq := &tripdomain.TripRequest{
			ID:     uuid.New(),
			Status: tripdomain.NO_DRIVER_FOUND,
		}

		tripRequestRepo.On("UpdateTripRequestStatus", tripReq.ID, tripdomain.CUSTOMER_CANCELED).Return(nil).Once()

		err := svc.CancelTripRequest(context.Background(), tripReq)

		require.NoError(t, err)
		tripRequestRepo.AssertExpectations(t)
	})

	t.Run("cancel fails: cannot cancel trip request at this stage", func(t *testing.T) {
		svc, tripRequestRepo, _ := setupTripRequestService()

		tripReq := &tripdomain.TripRequest{
			ID:     uuid.New(),
			Status: tripdomain.DRIVER_ACCEPTED,
		}

		err := svc.CancelTripRequest(context.Background(), tripReq)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrTripRequestInvalidState)
		tripRequestRepo.AssertNotCalled(t, "UpdateTripRequestStatus", mock.Anything, mock.Anything)
	})
}

func TestListOpenTripRequests(t *testing.T) {
	t.Run("happy path: returns open trip requests from repository", func(t *testing.T) {
		svc, tripRequestRepo, _ := setupTripRequestService()

		customerID := uuid.New()
		open := []tripdomain.TripRequest{*fixtureTripRequest(customerID)}
		tripRequestRepo.On("ListOpenTripRequests", 20).Return(open, nil).Once()

		got, err := svc.ListOpenTripRequests(20)

		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, tripdomain.NO_DRIVER_FOUND, got[0].Status)
		tripRequestRepo.AssertExpectations(t)
	})

	t.Run("repository error: propagates", func(t *testing.T) {
		svc, tripRequestRepo, _ := setupTripRequestService()

		repoErr := errors.New("db error")
		tripRequestRepo.On("ListOpenTripRequests", 5).Return(nil, repoErr).Once()

		got, err := svc.ListOpenTripRequests(5)

		require.Error(t, err)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, got)
		tripRequestRepo.AssertExpectations(t)
	})
}

func TestAcceptTripRequest(t *testing.T) {
	driverID := uuid.New()
	customerID := uuid.New()

	t.Run("happy path: accepts open trip request and creates trip", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripRequestService()

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
			return tr.TripRequestID == tripRequestID && tr.CustomerID == customerID && tr.DriverID == driverID && tr.Status == tripdomain.TRIP_ACCEPTED
		})).Return(nil).Once()
		tripRequestRepo.On("FindByID", tripRequestID).Return(accepted, nil).Once()

		trip, trAfter, err := svc.AcceptTripRequest(context.Background(), driverID, tripRequestID)

		require.NoError(t, err)
		require.NotNil(t, trip)
		require.NotNil(t, trAfter)
		assert.Equal(t, tripRequestID, trip.TripRequestID)
		assert.Equal(t, customerID, trip.CustomerID)
		assert.Equal(t, driverID, trip.DriverID)
		assert.Equal(t, tripdomain.TRIP_ACCEPTED, trip.Status)
		assert.Equal(t, tripdomain.DRIVER_ACCEPTED, trAfter.Status)

		tripRequestRepo.AssertExpectations(t)
		tripRepo.AssertExpectations(t)
	})

	t.Run("trip request not found", func(t *testing.T) {
		svc, tripRequestRepo, tripRepo := setupTripRequestService()

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
		svc, tripRequestRepo, tripRepo := setupTripRequestService()

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
		svc, tripRequestRepo, tripRepo := setupTripRequestService()

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
		svc, tripRequestRepo, tripRepo := setupTripRequestService()

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
		svc, tripRequestRepo, tripRepo := setupTripRequestService()

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
