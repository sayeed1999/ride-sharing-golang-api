package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	tripdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
	tripmocks "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupLocationService() (ILocationService, *tripmocks.TripRepository, *tripmocks.LocationRepository) {
	tripRepo := new(tripmocks.TripRepository)
	locRepo := new(tripmocks.LocationRepository)
	svc := NewLocationService(tripRepo, locRepo)
	return svc, tripRepo, locRepo
}

func TestUpdateDriverLocation(t *testing.T) {
	driverID := uuid.New()

	t.Run("happy path", func(t *testing.T) {
		svc, tripRepo, locRepo := setupLocationService()

		trip := &tripdomain.Trip{
			ID:       uuid.New(),
			DriverID: driverID,
			Status:   tripdomain.TRIP_ACCEPTED,
		}
		loc := &tripdomain.DriverLocation{Lat: 23.78, Lng: 90.41}

		tripRepo.On("FindActiveByDriverID", driverID).Return(trip, nil).Once()
		locRepo.On("SetDriverLocation", mock.Anything, driverID, mock.AnythingOfType("*domain.DriverLocation")).Return(nil).Once()

		err := svc.UpdateDriverLocation(context.Background(), driverID, loc)
		require.NoError(t, err)
		assert.False(t, loc.UpdatedAt.IsZero())

		tripRepo.AssertExpectations(t)
		locRepo.AssertExpectations(t)
	})

	t.Run("no active trip", func(t *testing.T) {
		svc, tripRepo, locRepo := setupLocationService()

		tripRepo.On("FindActiveByDriverID", driverID).Return(nil, gorm.ErrRecordNotFound).Once()

		err := svc.UpdateDriverLocation(context.Background(), driverID, &tripdomain.DriverLocation{Lat: 1, Lng: 1})
		require.ErrorIs(t, err, ErrNoActiveTrip)

		locRepo.AssertNotCalled(t, "SetDriverLocation", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("invalid coordinates", func(t *testing.T) {
		svc, tripRepo, locRepo := setupLocationService()

		err := svc.UpdateDriverLocation(context.Background(), driverID, &tripdomain.DriverLocation{Lat: 91, Lng: 0})
		require.ErrorIs(t, err, ErrInvalidCoordinates)

		tripRepo.AssertNotCalled(t, "FindActiveByDriverID", mock.Anything)
		locRepo.AssertNotCalled(t, "SetDriverLocation", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestGetDriverLocation(t *testing.T) {
	driverID := uuid.New()
	customerID := uuid.New()
	tripID := uuid.New()

	trip := &tripdomain.Trip{
		ID:         tripID,
		DriverID:   driverID,
		CustomerID: customerID,
		Status:     tripdomain.TRIP_IN_PROGRESS,
	}

	loc := &tripdomain.DriverLocation{
		Lat:       23.78,
		Lng:       90.41,
		UpdatedAt: time.Now().UTC(),
	}

	t.Run("customer can read", func(t *testing.T) {
		svc, _, locRepo := setupLocationService()
		locRepo.On("GetDriverLocation", mock.Anything, driverID).Return(loc, nil).Once()

		got, err := svc.GetDriverLocation(context.Background(), trip, &customerID, nil)
		require.NoError(t, err)
		assert.Equal(t, loc, got)
	})

	t.Run("driver can read", func(t *testing.T) {
		svc, _, locRepo := setupLocationService()
		locRepo.On("GetDriverLocation", mock.Anything, driverID).Return(loc, nil).Once()

		got, err := svc.GetDriverLocation(context.Background(), trip, nil, &driverID)
		require.NoError(t, err)
		assert.Equal(t, loc, got)
	})

	t.Run("forbidden actor", func(t *testing.T) {
		svc, _, locRepo := setupLocationService()
		other := uuid.New()

		got, err := svc.GetDriverLocation(context.Background(), trip, &other, nil)
		require.ErrorIs(t, err, ErrTripActorForbidden)
		assert.Nil(t, got)

		locRepo.AssertNotCalled(t, "GetDriverLocation", mock.Anything, mock.Anything)
	})

	t.Run("location missing", func(t *testing.T) {
		svc, _, locRepo := setupLocationService()
		locRepo.On("GetDriverLocation", mock.Anything, driverID).Return(nil, repository.ErrLocationNotFound).Once()

		got, err := svc.GetDriverLocation(context.Background(), trip, &customerID, nil)
		require.ErrorIs(t, err, ErrLocationNotAvailable)
		assert.Nil(t, got)
	})

	t.Run("trip not trackable", func(t *testing.T) {
		svc, _, locRepo := setupLocationService()
		completedTrip := *trip
		completedTrip.Status = tripdomain.TRIP_COMPLETED

		got, err := svc.GetDriverLocation(context.Background(), &completedTrip, &customerID, nil)
		require.ErrorIs(t, err, ErrTripNotTrackable)
		assert.Nil(t, got)

		locRepo.AssertNotCalled(t, "GetDriverLocation", mock.Anything, mock.Anything)
	})
}
