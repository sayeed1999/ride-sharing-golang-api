package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
	"gorm.io/gorm"
)

var (
	ErrNoActiveTrip         = errors.New("driver has no active trip")
	ErrInvalidCoordinates   = errors.New("invalid coordinates")
	ErrLocationNotAvailable = errors.New("driver location not available")
	ErrTripNotTrackable     = errors.New("trip is not trackable")
	ErrTripActorForbidden   = errors.New("actor not authorized for this trip")
)

type ILocationService interface {
	UpdateDriverLocation(ctx context.Context, driverID uuid.UUID, loc *domain.DriverLocation) error
	GetDriverLocation(ctx context.Context, trip *domain.Trip, actorCustomerID, actorDriverID *uuid.UUID) (*domain.DriverLocation, error)
}

type locationService struct {
	tripRepo repository.ITripRepository
	locRepo  repository.ILocationRepository
}

func NewLocationService(tripRepo repository.ITripRepository, locRepo repository.ILocationRepository) ILocationService {
	return &locationService{
		tripRepo: tripRepo,
		locRepo:  locRepo,
	}
}

func (s *locationService) UpdateDriverLocation(ctx context.Context, driverID uuid.UUID, loc *domain.DriverLocation) error {
	if err := validateCoordinates(loc.Lat, loc.Lng); err != nil {
		return err
	}
	if err := validateOptionalHeadingSpeed(loc.Heading, loc.Speed); err != nil {
		return err
	}

	trip, err := s.tripRepo.FindActiveByDriverID(driverID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNoActiveTrip
		}
		return err
	}
	if !trip.Status.IsTrackable() {
		return ErrNoActiveTrip
	}

	loc.UpdatedAt = time.Now().UTC()
	return s.locRepo.SetDriverLocation(ctx, driverID, loc)
}

func (s *locationService) GetDriverLocation(
	ctx context.Context,
	trip *domain.Trip,
	actorCustomerID, actorDriverID *uuid.UUID,
) (*domain.DriverLocation, error) {
	if !trip.Status.IsTrackable() {
		return nil, ErrTripNotTrackable
	}

	authorized := actorDriverID != nil && *actorDriverID == trip.DriverID
	if actorCustomerID != nil && *actorCustomerID == trip.CustomerID {
		authorized = true
	}
	if !authorized {
		return nil, ErrTripActorForbidden
	}

	loc, err := s.locRepo.GetDriverLocation(ctx, trip.DriverID)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			return nil, ErrLocationNotAvailable
		}
		return nil, err
	}
	return loc, nil
}

func validateCoordinates(lat, lng float64) error {
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return ErrInvalidCoordinates
	}
	return nil
}

func validateOptionalHeadingSpeed(heading, speed *float64) error {
	if heading != nil && (*heading < 0 || *heading > 360) {
		return ErrInvalidCoordinates
	}
	if speed != nil && *speed < 0 {
		return ErrInvalidCoordinates
	}
	return nil
}
