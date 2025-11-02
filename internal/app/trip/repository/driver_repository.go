package repository

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"
)

type DriverRepository interface {
	CreateDriver(d *domain.Driver) (*domain.Driver, error)
	FindByID(id uuid.UUID) (*domain.Driver, error)
	FindByEmail(email string) (*domain.Driver, error)
	DeleteDriver(id uuid.UUID) error
	UpdateAuthUserID(driverID uuid.UUID, authUserID uuid.UUID) error
}
