package domain

import (
	"time"

	"github.com/google/uuid"
)

type TripStatus int

const (
	TRIP_ACCEPTED TripStatus = iota + 1
	TRIP_IN_PROGRESS
	TRIP_COMPLETED
	TRIP_CANCELLED_BY_CUSTOMER
	TRIP_CANCELLED_BY_DRIVER
)

// Trip is an active assignment of a driver to a trip request after acceptance.
type Trip struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TripRequestID uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"trip_request_id"`
	DriverID      uuid.UUID  `gorm:"type:uuid;not null" json:"driver_id"`
	Status        TripStatus `gorm:"not null" json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Trip) TableName() string {
	return "trip.trips"
}
