package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
	"gorm.io/gorm"
)

var (
	ErrTripRequestNotOpen      = errors.New("trip request is not open for acceptance")
	ErrTripRequestNotFound     = errors.New("trip request not found")
	ErrTripNotFound            = errors.New("trip not found")
	ErrTripWrongDriver         = errors.New("trip does not belong to this driver")
	ErrTripNotOwnedByCustomer  = errors.New("trip does not belong to this customer")
	ErrTripInvalidState        = errors.New("trip is not in a valid state for this action")
	ErrTripStartConflict       = errors.New("could not mark trip as in progress")
	ErrTripCompleteConflict    = errors.New("could not mark trip as completed")
	ErrTripCancelConflict      = errors.New("could not cancel trip")
)

type ITripService interface {
	AcceptTripRequest(ctx context.Context, driverID, tripRequestID uuid.UUID) (*domain.Trip, *domain.TripRequest, error)
	StartTrip(ctx context.Context, driverID, tripID uuid.UUID) (*domain.Trip, error)
	CompleteTrip(ctx context.Context, driverID, tripID uuid.UUID) (*domain.Trip, error)
	CancelTripByCustomer(ctx context.Context, customerID, tripID uuid.UUID) (*domain.Trip, error)
	CancelTripByDriver(ctx context.Context, driverID, tripID uuid.UUID) (*domain.Trip, error)
}

// transactionRunner abstracts GORM transactions so unit tests can run the callback with mocks (noop runner).
type transactionRunner interface {
	Transaction(func(tx *gorm.DB) error) error
}

type gormTxRunner struct {
	db *gorm.DB
}

func (g *gormTxRunner) Transaction(fn func(tx *gorm.DB) error) error {
	return g.db.Transaction(fn)
}

type tripService struct {
	txRunner    transactionRunner
	tripReqRepo repository.ITripRequestRepository
	tripRepo    repository.ITripRepository
}

func NewTripService(db *gorm.DB, tripReqRepo repository.ITripRequestRepository, tripRepo repository.ITripRepository) ITripService {
	return newTripService(&gormTxRunner{db: db}, tripReqRepo, tripRepo)
}

func newTripService(txRunner transactionRunner, tripReqRepo repository.ITripRequestRepository, tripRepo repository.ITripRepository) ITripService {
	return &tripService{
		txRunner:    txRunner,
		tripReqRepo: tripReqRepo,
		tripRepo:    tripRepo,
	}
}

func (s *tripService) AcceptTripRequest(ctx context.Context, driverID, tripRequestID uuid.UUID) (*domain.Trip, *domain.TripRequest, error) {
	_ = ctx

	tr, err := s.tripReqRepo.FindByID(tripRequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrTripRequestNotFound
		}
		return nil, nil, err
	}
	if tr.Status != domain.NO_DRIVER_FOUND {
		return nil, nil, ErrTripRequestNotOpen
	}

	var tripOut *domain.Trip
	err = s.txRunner.Transaction(func(tx *gorm.DB) error {
		ok, err := s.tripReqRepo.UpdateTripRequestStatusIf(tx, tripRequestID, domain.NO_DRIVER_FOUND, domain.DRIVER_ACCEPTED)
		if err != nil {
			return err
		}
		if !ok {
			return ErrTripRequestNotOpen
		}
		trip := &domain.Trip{
			TripRequestID: tripRequestID,
			DriverID:      driverID,
			Status:        domain.TRIP_ACCEPTED,
		}
		if err := s.tripRepo.Create(tx, trip); err != nil {
			return err
		}
		tripOut = trip
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	trAfter, err := s.tripReqRepo.FindByID(tripRequestID)
	if err != nil {
		return tripOut, nil, err
	}
	return tripOut, trAfter, nil
}

func (s *tripService) StartTrip(ctx context.Context, driverID, tripID uuid.UUID) (*domain.Trip, error) {
	_ = ctx

	trip, err := s.tripRepo.FindByID(tripID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}
	if trip.DriverID != driverID {
		return nil, ErrTripWrongDriver
	}
	if trip.Status != domain.TRIP_ACCEPTED {
		return nil, ErrTripInvalidState
	}

	err = s.txRunner.Transaction(func(tx *gorm.DB) error {
		ok, err := s.tripRepo.UpdateTripStatus(tx, tripID, driverID, domain.TRIP_ACCEPTED, domain.TRIP_IN_PROGRESS)
		if err != nil {
			return err
		}
		if !ok {
			return ErrTripStartConflict
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.tripRepo.FindByID(tripID)
}

func (s *tripService) CompleteTrip(ctx context.Context, driverID, tripID uuid.UUID) (*domain.Trip, error) {
	_ = ctx

	trip, err := s.tripRepo.FindByID(tripID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}
	if trip.DriverID != driverID {
		return nil, ErrTripWrongDriver
	}
	if trip.Status != domain.TRIP_IN_PROGRESS {
		return nil, ErrTripInvalidState
	}

	err = s.txRunner.Transaction(func(tx *gorm.DB) error {
		ok, err := s.tripRepo.UpdateTripStatus(tx, tripID, driverID, domain.TRIP_IN_PROGRESS, domain.TRIP_COMPLETED)
		if err != nil {
			return err
		}
		if !ok {
			return ErrTripCompleteConflict
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.tripRepo.FindByID(tripID)
}

func (s *tripService) CancelTripByCustomer(ctx context.Context, customerID, tripID uuid.UUID) (*domain.Trip, error) {
	_ = ctx

	trip, err := s.tripRepo.FindByID(tripID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}

	tr, err := s.tripReqRepo.FindByID(trip.TripRequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}
	if tr.CustomerID != customerID {
		return nil, ErrTripNotOwnedByCustomer
	}

	var toStatus domain.TripStatus
	switch trip.Status {
	case domain.TRIP_ACCEPTED, domain.TRIP_IN_PROGRESS:
		toStatus = domain.TRIP_CANCELLED_BY_CUSTOMER
	default:
		return nil, ErrTripInvalidState
	}

	err = s.txRunner.Transaction(func(tx *gorm.DB) error {
		ok, err := s.tripRepo.UpdateTripStatusIf(tx, tripID, trip.Status, toStatus)
		if err != nil {
			return err
		}
		if !ok {
			return ErrTripCancelConflict
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.tripRepo.FindByID(tripID)
}

func (s *tripService) CancelTripByDriver(ctx context.Context, driverID, tripID uuid.UUID) (*domain.Trip, error) {
	_ = ctx

	trip, err := s.tripRepo.FindByID(tripID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTripNotFound
		}
		return nil, err
	}
	if trip.DriverID != driverID {
		return nil, ErrTripWrongDriver
	}
	if trip.Status != domain.TRIP_ACCEPTED {
		return nil, ErrTripInvalidState
	}

	err = s.txRunner.Transaction(func(tx *gorm.DB) error {
		ok, err := s.tripRepo.UpdateTripStatus(tx, tripID, driverID, domain.TRIP_ACCEPTED, domain.TRIP_CANCELLED_BY_DRIVER)
		if err != nil {
			return err
		}
		if !ok {
			return ErrTripCancelConflict
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.tripRepo.FindByID(tripID)
}
