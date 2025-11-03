package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/tests/e2e/setup"
	"github.com/stretchr/testify/require"
)

func TestRequestTrip_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, true)
	defer testApp.CleanUp(ctx, t)

	email := "e2e-customer-trip@example.com"
	password := "pass123"
	name := "E2E Customer Trip"

	// 1. Signup as customer
	signupPayload := map[string]string{"email": email, "name": name, "password": password}
	w := doJSONRequest(t, testApp.Router(), http.MethodPost, "/customers/signup", signupPayload)
	assertAndLogErrors(t, w, http.StatusCreated)

	// Extract customer ID
	var customer struct {
		ID string `json:"id"`
	}
	err := testApp.DB.Raw("SELECT id FROM trip.customers WHERE email = ?", email).Scan(&customer.ID).Error
	require.NoError(t, err)
	require.NotEmpty(t, customer.ID, "Customer ID should not be empty")

	// 2. Request a trip
	tripRequestPayload := map[string]interface{}{
		"customer_id": customer.ID,
		"origin":      "123 Main St",
		"destination": "456 Oak Ave",
	}
	w = doJSONRequest(t, testApp.Router(), http.MethodPost, "/trip-requests/request", tripRequestPayload)
	assertAndLogErrors(t, w, http.StatusCreated)

	// 3. Verify trip.trip_requests has the record
	var tripRequestRec domain.TripRequest
	err = testApp.DB.Raw("SELECT customer_id, origin, destination, status FROM trip.trip_requests WHERE customer_id = ?", customer.ID).Scan(&tripRequestRec).Error
	require.NoError(t, err)
	require.Equal(t, customer.ID, tripRequestRec.CustomerID.String())
	require.Equal(t, "123 Main St", tripRequestRec.Origin)
	require.Equal(t, "456 Oak Ave", tripRequestRec.Destination)
	require.Equal(t, domain.NO_DRIVER_FOUND, tripRequestRec.Status)
}

func TestRequestTrip_Validation_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, true)
	defer testApp.CleanUp(ctx, t)

	cases := []struct {
		name    string
		payload map[string]interface{}
	}{
		{
			name: "invalid customer id",
			payload: map[string]interface{}{
				"customer_id": "invalid-uuid",
				"origin":      "123 Main St",
				"destination": "456 Oak Ave",
			},
		},
		{
			name: "missing origin",
			payload: map[string]interface{}{
				"customer_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
				"destination": "456 Oak Ave",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := doJSONRequest(t, testApp.Router(), http.MethodPost, "/trip-requests/request", tc.payload)
			assertAndLogErrors(t, w, http.StatusBadRequest)
		})
	}
}

// doJSONRequestWithAuth is a helper to make JSON requests with an Authorization header
func doJSONRequestWithAuth(t *testing.T, router *gin.Engine, method, path string, payload interface{}, token string) *httptest.ResponseRecorder {
	var reqBody io.Reader
	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonPayload)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
