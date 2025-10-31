package e2e

import (
	"context"
	"net/http"
	"testing"

	"github.com/sayeed1999/ride-sharing-golang-api/tests/e2e/setup"
	"github.com/stretchr/testify/require"
)

func TestDriverSignup_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, true)
	defer testApp.CleanUp(ctx, t)

	email := "e2e-driver@example.com"

	// Missing vehicle type
	signupPayload := map[string]string{
		"email":                "e2e-driver-no-vehicle-type@example.com",
		"name":                 "E2E Driver",
		"password":             "pass123",
		"vehicle_registration": "ABC-123",
	}
	w := doJSONRequest(t, testApp.Router(), http.MethodPost, "/drivers/signup", signupPayload)
	assertAndLogErrors(t, w, http.StatusBadRequest)

	// Missing vehicle registration
	signupPayload = map[string]string{
		"email":        "e2e-driver-no-vehicle-reg@example.com",
		"name":         "E2E Driver",
		"password":     "pass123",
		"vehicle_type": "sedan",
	}
	w = doJSONRequest(t, testApp.Router(), http.MethodPost, "/drivers/signup", signupPayload)
	assertAndLogErrors(t, w, http.StatusBadRequest)

	// Invalid vehicle type
	signupPayload = map[string]string{
		"email":                "e2e-driver-invalid-vehicle-type@example.com",
		"name":                 "E2E Driver",
		"password":             "pass123",
		"vehicle_type":         "rickshaw",
		"vehicle_registration": "ABC-123",
	}
	w = doJSONRequest(t, testApp.Router(), http.MethodPost, "/drivers/signup", signupPayload)
	assertAndLogErrorsWithBody(t, w, http.StatusBadRequest, "invalid vehicle type")

	// Successful signup
	signupPayload = map[string]string{
		"email":                email,
		"name":                 "E2E Driver",
		"password":             "pass123",
		"vehicle_type":         "car",
		"vehicle_registration": "ABC-123",
	}
	w = doJSONRequest(t, testApp.Router(), http.MethodPost, "/drivers/signup", signupPayload)
	assertAndLogErrors(t, w, http.StatusCreated)

	// Verify trip.drivers has the record
	var tripRec struct{ Email string }
	err := testApp.DB.Raw("SELECT email FROM trip.drivers WHERE email = ?", email).Scan(&tripRec).Error
	require.NoError(t, err)
	if tripRec.Email != email {
		t.Errorf("expected driver in trip schema")
	}

	// Verify auth.users has the record
	var authRec struct{ Email string }
	err = testApp.DB.Raw("SELECT email FROM auth.users WHERE email = ?", email).Scan(&authRec).Error
	require.NoError(t, err)
	if authRec.Email != email {
		t.Errorf("expected driver in auth schema")
	}
}
