package mocks

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type TripRepository struct {
	mock.Mock
}

func (m *TripRepository) Create(db *gorm.DB, t *domain.Trip) error {
	args := m.Called(db, t)
	return args.Error(0)
}

func (m *TripRepository) FindByID(id uuid.UUID) (*domain.Trip, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Trip), args.Error(1)
}

func (m *TripRepository) FindByTripRequestID(tripRequestID uuid.UUID) (*domain.Trip, error) {
	args := m.Called(tripRequestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Trip), args.Error(1)
}

func (m *TripRepository) UpdateTripStatus(db *gorm.DB, tripID, driverID uuid.UUID, from, to domain.TripStatus) (bool, error) {
	args := m.Called(db, tripID, driverID, from, to)
	return args.Bool(0), args.Error(1)
}

func (m *TripRepository) UpdateTripStatusIf(db *gorm.DB, tripID uuid.UUID, from, to domain.TripStatus) (bool, error) {
	args := m.Called(db, tripID, from, to)
	return args.Bool(0), args.Error(1)
}

var _ repository.ITripRepository = (*TripRepository)(nil)
