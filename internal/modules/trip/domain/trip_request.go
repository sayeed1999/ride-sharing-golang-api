package domain

import (
	"time"

	"github.com/google/uuid"
)

type TripRequestStatus int

const (
	NO_DRIVER_FOUND TripRequestStatus = iota + 1
	CUSTOMER_CANCELED
	DRIVER_ACCEPTED
	CUSTOMER_REJECTED_DRIVER
	DRIVER_REJECTED_CUSTOMER
	TRIP_STARTED
	TRIP_REQUEST_REJECTED
)

// TripRequest represents a customer's request for a ride.
type TripRequest struct {
	ID          uuid.UUID         `gorm:"type:uuid;primary_key;" json:"id"`
	CustomerID  uuid.UUID         `gorm:"type:uuid;not null" json:"customer_id"`
	Origin      string            `gorm:"size:255;not null" json:"origin"`
	Destination string            `gorm:"size:255;not null" json:"destination"`
	Status      TripRequestStatus `gorm:"not null" json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (TripRequest) TableName() string {
	return "trip.trip_requests"
}
