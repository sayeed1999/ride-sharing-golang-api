package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	tripdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	tripmocks "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRequestTripByCustomer(t *testing.T) {
	tripRequestRepo := new(tripmocks.TripRequestRepository)
	svc := NewTripRequestService(tripRequestRepo)

	customerID := uuid.New()
	expected := &tripdomain.TripRequest{
		ID:          uuid.New(),
		CustomerID:  customerID,
		Origin:      "A",
		Destination: "B",
		Status:      tripdomain.NO_DRIVER_FOUND,
	}

	tripRequestRepo.On("Create", mock.MatchedBy(func(tr *tripdomain.TripRequest) bool {
		return tr.CustomerID == customerID &&
			tr.Origin == "A" &&
			tr.Destination == "B" &&
			tr.Status == tripdomain.NO_DRIVER_FOUND
	})).Return(expected, nil).Once()

	created, err := svc.RequestTrip(customerID, "A", "B")

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, expected.ID, created.ID)
	assert.Equal(t, tripdomain.NO_DRIVER_FOUND, created.Status)

	tripRequestRepo.AssertExpectations(t)
}

func TestCancelTripRequestByCustomer(t *testing.T) {
	tripRequestRepo := new(tripmocks.TripRequestRepository)
	svc := NewTripRequestService(tripRequestRepo)

	tripReq := &tripdomain.TripRequest{
		ID:     uuid.New(),
		Status: tripdomain.NO_DRIVER_FOUND,
	}

	tripRequestRepo.On("UpdateTripRequestStatus", tripReq.ID, tripdomain.CUSTOMER_CANCELED).Return(nil).Once()

	err := svc.CancelTripRequest(context.Background(), tripReq)

	require.NoError(t, err)
	tripRequestRepo.AssertExpectations(t)
}

