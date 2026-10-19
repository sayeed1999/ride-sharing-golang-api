package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/stretchr/testify/mock"
)

type LocationRepository struct {
	mock.Mock
}

func (m *LocationRepository) SetDriverLocation(ctx context.Context, driverID uuid.UUID, loc *domain.DriverLocation) error {
	args := m.Called(ctx, driverID, loc)
	return args.Error(0)
}

func (m *LocationRepository) GetDriverLocation(ctx context.Context, driverID uuid.UUID) (*domain.DriverLocation, error) {
	args := m.Called(ctx, driverID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DriverLocation), args.Error(1)
}
