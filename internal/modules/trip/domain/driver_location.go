package domain

import "time"

// DriverLocation is an ephemeral live position stored in Redis (not PostgreSQL).
type DriverLocation struct {
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Heading   *float64  `json:"heading,omitempty"`
	Speed     *float64  `json:"speed,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s TripStatus) IsTrackable() bool {
	return s == TRIP_ACCEPTED || s == TRIP_IN_PROGRESS
}
