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
	ErrTripRequestNotOpen  = errors.New("trip request is not open for acceptance")
	ErrTripRequestNotFound = errors.New("trip request not found")
	ErrTripNotFound        = errors.New("trip not found")
	ErrTripWrongDriver     = errors.New("trip does not belong to this driver")
	ErrTripInvalidState    = errors.New("trip is not in a valid state for this action")
	ErrTripStartConflict   = errors.New("could not mark trip as in progress")
)

type ITripService interface {
	AcceptTripRequest(ctx context.Context, driverID, tripRequestID uuid.UUID) (*domain.Trip, *domain.TripRequest, error)
	StartTrip(ctx context.Context, driverID, tripID uuid.UUID) (*domain.Trip, error)
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
		ok2, err := s.tripReqRepo.UpdateTripRequestStatusIf(tx, trip.TripRequestID, domain.DRIVER_ACCEPTED, domain.TRIP_STARTED)
		if err != nil {
			return err
		}
		if !ok2 {
			return ErrTripStartConflict
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.tripRepo.FindByID(tripID)
}
