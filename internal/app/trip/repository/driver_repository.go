package repository

import "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"

type DriverRepository interface {
    CreateDriver(d *domain.Driver) (*domain.Driver, error)
    FindByID(id string) (*domain.Driver, error)
    DeleteDriver(id string) error
}
