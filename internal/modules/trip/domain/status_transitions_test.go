package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTripRequestStatusCanTransitionTo(t *testing.T) {
	assert.True(t, NO_DRIVER_FOUND.CanTransitionTo(CUSTOMER_CANCELED))
	assert.True(t, NO_DRIVER_FOUND.CanTransitionTo(DRIVER_ACCEPTED))
	assert.True(t, NO_DRIVER_FOUND.CanTransitionTo(EXPIRED))

	assert.False(t, NO_DRIVER_FOUND.CanTransitionTo(NO_DRIVER_FOUND))
	assert.False(t, DRIVER_ACCEPTED.CanTransitionTo(CUSTOMER_CANCELED))
	assert.False(t, CUSTOMER_CANCELED.CanTransitionTo(DRIVER_ACCEPTED))
}

func TestTripStatusCanTransitionTo(t *testing.T) {
	assert.True(t, TRIP_ACCEPTED.CanTransitionTo(TRIP_IN_PROGRESS))
	assert.True(t, TRIP_ACCEPTED.CanTransitionTo(TRIP_CANCELLED_BY_CUSTOMER))
	assert.True(t, TRIP_ACCEPTED.CanTransitionTo(TRIP_CANCELLED_BY_DRIVER))

	assert.True(t, TRIP_IN_PROGRESS.CanTransitionTo(TRIP_COMPLETED))
	assert.True(t, TRIP_IN_PROGRESS.CanTransitionTo(TRIP_CANCELLED_BY_CUSTOMER))

	assert.False(t, TRIP_IN_PROGRESS.CanTransitionTo(TRIP_CANCELLED_BY_DRIVER))
	assert.False(t, TRIP_COMPLETED.CanTransitionTo(TRIP_IN_PROGRESS))
	assert.False(t, TRIP_CANCELLED_BY_CUSTOMER.CanTransitionTo(TRIP_ACCEPTED))
}
