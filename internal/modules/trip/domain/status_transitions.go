package domain

import "slices"

// TripRequestAllowedTransitions maps a trip request status to its allowed target statuses (TRIP.md §4).
var TripRequestAllowedTransitions = map[TripRequestStatus][]TripRequestStatus{
	NO_DRIVER_FOUND: {CUSTOMER_CANCELED, DRIVER_ACCEPTED, EXPIRED},
}

// TripAllowedTransitions maps a trip status to its allowed target statuses (TRIP.md §4).
var TripAllowedTransitions = map[TripStatus][]TripStatus{
	TRIP_ACCEPTED: {
		TRIP_IN_PROGRESS,
		TRIP_CANCELLED_BY_CUSTOMER,
		TRIP_CANCELLED_BY_DRIVER,
	},
	TRIP_IN_PROGRESS: {
		TRIP_COMPLETED,
		TRIP_CANCELLED_BY_CUSTOMER,
	},
}

func (from TripRequestStatus) CanTransitionTo(to TripRequestStatus) bool {
	return slices.Contains(TripRequestAllowedTransitions[from], to)
}

func (from TripStatus) CanTransitionTo(to TripStatus) bool {
	return slices.Contains(TripAllowedTransitions[from], to)
}
