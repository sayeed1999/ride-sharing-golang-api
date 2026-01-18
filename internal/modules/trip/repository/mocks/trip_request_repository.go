package mocks

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/stretchr/testify/mock"
)

type TripRequestRepository struct {
	mock.Mock
}

func (m *TripRequestRepository) Create(tr *domain.TripRequest) (*domain.TripRequest, error) {
	args := m.Called(tr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TripRequest), args.Error(1)
}

func (m *TripRequestRepository) FindByID(id uuid.UUID) (*domain.TripRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TripRequest), args.Error(1)
}

func (m *TripRequestRepository) Update(tr *domain.TripRequest) (*domain.TripRequest, error) {
	args := m.Called(tr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TripRequest), args.Error(1)
}

func (m *TripRequestRepository) UpdateTripRequestStatus(tripID uuid.UUID, status domain.TripRequestStatus) error {
	args := m.Called(tripID, status)
	return args.Error(0)
}
