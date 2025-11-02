
package repository

import "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"

type TripRequestRepository interface {
	Create(tr *domain.TripRequest) (*domain.TripRequest, error)
	FindByID(id string) (*domain.TripRequest, error)
	Update(tr *domain.TripRequest) (*domain.TripRequest, error)
}
