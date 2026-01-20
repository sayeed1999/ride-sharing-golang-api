package e2e

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/dto"
	"github.com/sayeed1999/ride-sharing-golang-api/tests/e2e/setup"
	testhelper "github.com/sayeed1999/ride-sharing-golang-api/pkg/test_helper"
)

func Test_ValidWorkflow_CustomerCancelRideBeforeDriverFound_E2E(t *testing.T) {
	// Setup test app environment
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, true)
	defer testApp.CleanUp(ctx, t)

	// Register a valid customer
	customerSignupPayload := dto.NewCustomerSignupRequest("customer1@example.com", "First Customer", "password123")
	w := testhelper.DoJSONRequest(t, testApp.Router(), "POST", "/customers/signup", customerSignupPayload)
	testhelper.AssertAndLogErrors(t, w, 201)

	// Login as a customer
	loginPayload := dto.NewLoginRequest("customer1@example.com", "password123")
	w = testhelper.DoJSONRequest(t, testApp.Router(), "POST", "/login", loginPayload)
	testhelper.AssertAndLogErrors(t, w, 200)
	// Extract customer JWT token
	customerJwtToken := testhelper.ExtractTokenFromResponse(t, w)
	// Valid token check
	testhelper.AssertValidJWT(t, customerJwtToken, "test_jwt_secret_change_me", "customer1@example.com")

	// Register a valid driver
	driverSignupPayload := dto.NewDriverSignupRequest("driver1@example.com", "First Driver", "password12345", "bike", "ABC1234")
	w = testhelper.DoJSONRequest(t, testApp.Router(), "POST", "/drivers/signup", driverSignupPayload)
	testhelper.AssertAndLogErrors(t, w, 201)

	// Login as a driver
	loginPayloadDriver := dto.NewLoginRequest("driver1@example.com", "password12345")
	w = testhelper.DoJSONRequest(t, testApp.Router(), "POST", "/login", loginPayloadDriver)
	testhelper.AssertAndLogErrors(t, w, 200)
	// Extract driver JWT token
	driverJwtToken := testhelper.ExtractTokenFromResponse(t, w)
	// Valid token check
	testhelper.AssertValidJWT(t, driverJwtToken, "test_jwt_secret_change_me", "driver1@example.com")

	// Request a ride
	tripRequestPayload := dto.NewTripRequestDTO("Point A", "Point B")
	w = testhelper.DoJSONRequestWithAuth(t, testApp.Router(), "POST", "/trip-requests", tripRequestPayload, customerJwtToken)
	testhelper.AssertAndLogErrors(t, w, 201)

	// Extract trip request from response
	tripRequest := extractTripRequestFromResponse(t, w)

	// Cancel the ride before a driver is found
	w = testhelper.DoJSONRequestWithAuth(t, testApp.Router(), "DELETE", "/trip-requests/"+tripRequest.ID.String(), nil, customerJwtToken)
	testhelper.AssertAndLogErrors(t, w, 204)

	// Verify the ride status is "cancelled"
	w = testhelper.DoJSONRequestWithAuth(t, testApp.Router(), "GET", "/trip-requests/"+tripRequest.ID.String(), nil, customerJwtToken)
	testhelper.AssertAndLogErrors(t, w, 200)
	tripRequest = extractTripRequestFromResponse(t, w)
	if tripRequest.Status != domain.CUSTOMER_CANCELED {
		t.Error("expected status: CUSTOMER_CANCELED, got something else")
	}
}

func extractTripRequestFromResponse(t *testing.T, w *httptest.ResponseRecorder) domain.TripRequest {
	t.Helper()

	var response struct {
		TripRequest domain.TripRequest `json:"trip_request"`
	}

	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode trip request response: %v", err)
	}

	return response.TripRequest
}

// func Test_ValidWorkflow_CustomerAndDriverCompleteRide_E2E(t *testing.T) {
// 	ctx := context.Background()
// 	testApp := setup.NewTestApp(ctx, t, true)
// 	defer testApp.CleanUp(ctx, t)

// 	// Register a valid customers

// 	// Register a valid driver

// 	// Login as a customer

// 	// Request a ride

// 	// Login as a driver

// 	// Accept the ride

// 	// Complete the ride as customer and driver

// 	// Verify the ride status is "completed"
// }
