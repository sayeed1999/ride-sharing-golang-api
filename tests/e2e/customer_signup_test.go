package e2e

import (
	"context"
	"net/http"
	"testing"

	testhelper "github.com/sayeed1999/ride-sharing-golang-api/pkg/test_helper"
	"github.com/sayeed1999/ride-sharing-golang-api/tests/e2e/setup"
	"github.com/stretchr/testify/require"
)

func TestCustomerSignupAndLogin_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, true)
	defer testApp.CleanUp(ctx, t)

	email := "e2e-customer@example.com"

	// Signup as customer
	signupPayload := map[string]string{"email": email, "name": "E2E Customer", "password": "pass123"}
	w := testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/customers/signup", signupPayload)
	testhelper.AssertAndLogErrors(t, w, http.StatusCreated)

	// Verify trip.customers has the record
	var tripRec struct{ Email string }
	err := testApp.DB.Raw("SELECT email FROM trip.customers WHERE email = ?", email).Scan(&tripRec).Error
	require.NoError(t, err)
	if tripRec.Email != email {
		t.Errorf("expected customer in trip schema")
	}

	// Verify auth.users has the record
	var authRec struct{ Email string }
	err = testApp.DB.Raw("SELECT email FROM auth.users WHERE email = ?", email).Scan(&authRec).Error
	require.NoError(t, err)
	if tripRec.Email != email {
		t.Errorf("expected customer in auth schema")
	}
	// Try login
	loginPayload := map[string]string{"email": email, "password": "pass123"}
	w = testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/login", loginPayload)
	testhelper.AssertAndLogErrors(t, w, http.StatusOK)

	// validate JWT
	tokenStr := extractTokenFromResponse(t, w)
	assertValidJWT(t, tokenStr, "test_jwt_secret_change_me", email)
}

func TestCustomerSignup_Validation_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, true)
	defer testApp.CleanUp(ctx, t)

	// Signup as customer to test duplicate email
	signupPayload := map[string]string{"email": "duplicate-customer@example.com", "name": "E2E Customer", "password": "pass123"}
	w := testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/customers/signup", signupPayload)
	testhelper.AssertAndLogErrors(t, w, http.StatusCreated)

	cases := []struct {
		name    string
		payload map[string]string
	}{
		{
			name:    "duplicate email",
			payload: map[string]string{"email": "duplicate-customer@example.com", "name": "E2E Customer", "password": "pass123"},
		},
		{
			name:    "invalid email",
			payload: map[string]string{"email": "e2e-customer", "name": "E2E Customer", "password": "pass123"},
		},
		{
			name:    "missing name",
			payload: map[string]string{"email": "e2e-customer-no-name@example.com", "password": "pass123"},
		},
		{
			name:    "missing password",
			payload: map[string]string{"email": "e2e-customer-no-password@example.com", "name": "E2E Customer"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/customers/signup", tc.payload)
			testhelper.AssertAndLogErrors(t, w, http.StatusBadRequest)
		})
	}
}
