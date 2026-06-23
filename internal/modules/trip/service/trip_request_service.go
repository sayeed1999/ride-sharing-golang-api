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
)

type ITripRequestService interface {
	RequestTrip(customerID uuid.UUID, origin, destination string) (*domain.TripRequest, error)
	CancelTripRequest(ctx context.Context, tripRequest *domain.TripRequest) error
	ListOpenTripRequests(limit int) ([]domain.TripRequest, error)
	AcceptTripRequest(ctx context.Context, driverID, tripRequestID uuid.UUID) (*domain.Trip, *domain.TripRequest, error)
}

type tripRequestService struct {
	txRunner              transactionRunner
	tripRequestRepository repository.ITripRequestRepository
	tripRepository        repository.ITripRepository
}

func NewTripRequestService(db *gorm.DB, tripRequestRepo repository.ITripRequestRepository, tripRepo repository.ITripRepository) ITripRequestService {
	return newTripRequestService(&gormTxRunner{db: db}, tripRequestRepo, tripRepo)
}

func newTripRequestService(txRunner transactionRunner, tripRequestRepo repository.ITripRequestRepository, tripRepo repository.ITripRepository) ITripRequestService {
	return &tripRequestService{
		txRunner:              txRunner,
		tripRequestRepository: tripRequestRepo,
		tripRepository:        tripRepo,
	}
}

// RequestTrip creates a new trip request.
func (s *tripRequestService) RequestTrip(customerID uuid.UUID, origin, destination string) (*domain.TripRequest, error) {
	tripRequest := &domain.TripRequest{
		CustomerID:  customerID,
		Origin:      origin,
		Destination: destination,
		Status:      domain.NO_DRIVER_FOUND, // default status
	}

	createdTripRequest, err := s.tripRequestRepository.Create(tripRequest)
	if err != nil {
		return nil, err
	}

	return createdTripRequest, nil
}

// CustomerCancelTrip handles the business logic for a customer to cancel a trip request.
func (s *tripRequestService) CancelTripRequest(ctx context.Context, tripRequest *domain.TripRequest) error {
	// trip request middleware validates that the trip belongs to the customer

	// Only allow cancellation if the trip is in NO_DRIVER_FOUND state
	if tripRequest.Status != domain.NO_DRIVER_FOUND {
		return errors.New("trip request cannot be cancelled at this stage")
	}

	return s.tripRequestRepository.UpdateTripRequestStatus(tripRequest.ID, domain.CUSTOMER_CANCELED)
}

func (s *tripRequestService) ListOpenTripRequests(limit int) ([]domain.TripRequest, error) {
	return s.tripRequestRepository.ListOpenTripRequests(limit)
}

func (s *tripRequestService) AcceptTripRequest(ctx context.Context, driverID, tripRequestID uuid.UUID) (*domain.Trip, *domain.TripRequest, error) {
	_ = ctx

	tr, err := s.tripRequestRepository.FindByID(tripRequestID)
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
		ok, err := s.tripRequestRepository.UpdateTripRequestStatusIf(tx, tripRequestID, domain.NO_DRIVER_FOUND, domain.DRIVER_ACCEPTED)
		if err != nil {
			return err
		}
		if !ok {
			return ErrTripRequestNotOpen
		}
		trip := &domain.Trip{
			TripRequestID: tripRequestID,
			CustomerID:    tr.CustomerID,
			DriverID:      driverID,
			Status:        domain.TRIP_ACCEPTED,
		}
		if err := s.tripRepository.Create(tx, trip); err != nil {
			return err
		}
		tripOut = trip
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	trAfter, err := s.tripRequestRepository.FindByID(tripRequestID)
	if err != nil {
		return tripOut, nil, err
	}
	return tripOut, trAfter, nil
}
