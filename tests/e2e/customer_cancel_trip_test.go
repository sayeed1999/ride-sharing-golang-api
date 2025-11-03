package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"
	testhelper "github.com/sayeed1999/ride-sharing-golang-api/pkg/test_helper"
	"github.com/sayeed1999/ride-sharing-golang-api/tests/e2e/setup"
	"github.com/stretchr/testify/require"
)

func TestCancelTrip_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, true)
	defer testApp.CleanUp(ctx, t)

	email := "e2e-customer-cancel@example.com"
	password := "pass123"
	name := "E2E Customer Cancel"

	// 1. Signup as customer
	signupPayload := map[string]string{"email": email, "name": name, "password": password}
	w := testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/customers/signup", signupPayload)
	testhelper.AssertAndLogErrors(t, w, http.StatusCreated)

	// 2. Login as customer to get JWT token
	loginPayload := map[string]string{"email": email, "password": password}
	w = testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/login", loginPayload)
	testhelper.AssertAndLogErrors(t, w, http.StatusOK)

	var loginResponse struct {
		Token string `json:"token"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	require.NotEmpty(t, loginResponse.Token, "JWT token should not be empty")
	jwtToken := loginResponse.Token

	// Extract customer ID
	var customer struct {
		ID string `json:"id"`
	}
	err = testApp.DB.Raw("SELECT id FROM trip.customers WHERE email = ?", email).Scan(&customer.ID).Error
	require.NoError(t, err)
	require.NotEmpty(t, customer.ID, "Customer ID should not be empty")

	// 3. Request a trip
	tripRequestPayload := map[string]interface{}{
		"customer_id": customer.ID,
		"origin":      "789 Pine St",
		"destination": "101 Elm Ave",
	}
	w = testhelper.DoJSONRequestWithAuth(t, testApp.Router(), http.MethodPost, "/trip-requests/request", tripRequestPayload, jwtToken)
	testhelper.AssertAndLogErrors(t, w, http.StatusCreated)

	var tripRequestResponse struct {
		TripRequest struct {
			ID string `json:"id"`
		} `json:"trip_request"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &tripRequestResponse)
	require.NoError(t, err)
	require.NotEmpty(t, tripRequestResponse.TripRequest.ID, "Trip Request ID should not be empty")
	tripID := tripRequestResponse.TripRequest.ID

	// 4. Cancel the trip
	w = testhelper.DoJSONRequestWithAuth(t, testApp.Router(), http.MethodDelete, fmt.Sprintf("/trip-requests/%s", tripID), nil, jwtToken)
	testhelper.AssertAndLogErrors(t, w, http.StatusNoContent)

	// 5. Verify trip status in DB
	var tripRequestRec domain.TripRequest
	err = testApp.DB.Raw("SELECT status FROM trip.trip_requests WHERE id = ?", tripID).Scan(&tripRequestRec).Error
	require.NoError(t, err)
	require.Equal(t, domain.CUSTOMER_CANCELED, tripRequestRec.Status)
}

func TestCancelTrip_Validation_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, true)
	defer testApp.CleanUp(ctx, t)

	// 1. Signup and login to get a valid token
	email := "e2e-customer-cancel-validation@example.com"
	password := "pass123"
	name := "E2E Customer Cancel Validation"

	signupPayload := map[string]string{"email": email, "name": name, "password": password}
	w := testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/customers/signup", signupPayload)
	testhelper.AssertAndLogErrors(t, w, http.StatusCreated)

	loginPayload := map[string]string{"email": email, "password": password}
	w = testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/login", loginPayload)
	testhelper.AssertAndLogErrors(t, w, http.StatusOK)

	var loginResponse struct {
		Token string `json:"token"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	require.NotEmpty(t, loginResponse.Token, "JWT token should not be empty")
	jwtToken := loginResponse.Token

	// 2. Test with invalid trip ID
	w = testhelper.DoJSONRequestWithAuth(t, testApp.Router(), http.MethodDelete, "/trip-requests/invalid-trip-id", nil, jwtToken)
	testhelper.AssertAndLogErrors(t, w, http.StatusBadRequest)

	// 3. Create a trip to get a valid trip ID
	var customer struct {
		ID string `json:"id"`
	}
	err = testApp.DB.Raw("SELECT id FROM trip.customers WHERE email = ?", email).Scan(&customer.ID).Error
	require.NoError(t, err)
	require.NotEmpty(t, customer.ID, "Customer ID should not be empty")

	tripRequestPayload := map[string]interface{}{
		"customer_id": customer.ID,
		"origin":      "789 Pine St",
		"destination": "101 Elm Ave",
	}
	w = testhelper.DoJSONRequestWithAuth(t, testApp.Router(), http.MethodPost, "/trip-requests/request", tripRequestPayload, jwtToken)
	testhelper.AssertAndLogErrors(t, w, http.StatusCreated)

	var tripRequestResponse struct {
		TripRequest struct {
			ID string `json:"id"`
		} `json:"trip_request"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &tripRequestResponse)
	require.NoError(t, err)
	require.NotEmpty(t, tripRequestResponse.TripRequest.ID, "Trip Request ID should not be empty")
	tripID := tripRequestResponse.TripRequest.ID

	// 4. Test with invalid JWT token
	w = testhelper.DoJSONRequestWithAuth(t, testApp.Router(), http.MethodDelete, fmt.Sprintf("/trip-requests/%s", tripID), nil, "invalid-token")
	testhelper.AssertAndLogErrors(t, w, http.StatusUnauthorized)
}

func TestCancelTrip_Unauthorized(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, true)
	defer testApp.CleanUp(ctx, t)

	// 1. Create user A and their trip
	userAEmail := "userA@example.com"
	userAPassword := "pass123"
	signupPayloadA := map[string]string{"email": userAEmail, "name": "User A", "password": userAPassword}
	w := testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/customers/signup", signupPayloadA)
	testhelper.AssertAndLogErrors(t, w, http.StatusCreated)

	var customerA struct {
		ID string `json:"id"`
	}
	err := testApp.DB.Raw("SELECT id FROM trip.customers WHERE email = ?", userAEmail).Scan(&customerA.ID).Error
	require.NoError(t, err)

	loginPayloadA := map[string]string{"email": userAEmail, "password": userAPassword}
	w = testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/login", loginPayloadA)
	testhelper.AssertAndLogErrors(t, w, http.StatusOK)

	var loginResponseA struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &loginResponseA)
	require.NoError(t, err)
	jwtTokenA := loginResponseA.Token

	tripRequestPayload := map[string]interface{}{
		"customer_id": customerA.ID,
		"origin":      "123 Main St",
		"destination": "456 Oak Ave",
	}
	w = testhelper.DoJSONRequestWithAuth(t, testApp.Router(), http.MethodPost, "/trip-requests/request", tripRequestPayload, jwtTokenA)
	testhelper.AssertAndLogErrors(t, w, http.StatusCreated)

	var tripRequestResponse struct {
		TripRequest struct {
			ID string `json:"id"`
		} `json:"trip_request"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &tripRequestResponse)
	require.NoError(t, err)
	tripID := tripRequestResponse.TripRequest.ID

	// 2. Create user B and get their token
	userBEmail := "userB@example.com"
	userBPassword := "pass123"
	signupPayloadB := map[string]string{"email": userBEmail, "name": "User B", "password": userBPassword}
	w = testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/customers/signup", signupPayloadB)
	testhelper.AssertAndLogErrors(t, w, http.StatusCreated)

	loginPayloadB := map[string]string{"email": userBEmail, "password": userBPassword}
	w = testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/login", loginPayloadB)
	testhelper.AssertAndLogErrors(t, w, http.StatusOK)

	var loginResponseB struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &loginResponseB)
	require.NoError(t, err)
	jwtTokenB := loginResponseB.Token

	// 3. User B attempts to cancel User A's trip
	w = testhelper.DoJSONRequestWithAuth(t, testApp.Router(), http.MethodDelete, fmt.Sprintf("/trip-requests/%s", tripID), nil, jwtTokenB)
	testhelper.AssertAndLogErrors(t, w, http.StatusBadRequest)
}

func TestCancelTrip_Unauthenticated_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, true)
	defer testApp.CleanUp(ctx, t)

	// Attempt to cancel a trip without a token
	w := testhelper.DoJSONRequest(t, testApp.Router(), http.MethodDelete, "/trip-requests/some-trip-id", nil)
	testhelper.AssertAndLogErrors(t, w, http.StatusUnauthorized)
}