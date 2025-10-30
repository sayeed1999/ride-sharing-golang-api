package e2e

import (
	"context"
	"net/http"
	"testing"

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
	w := doJSONRequest(t, testApp.Router(), http.MethodPost, "/customers/signup", signupPayload)
	assertAndLogErrors(t, w, http.StatusCreated)

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
	w = doJSONRequest(t, testApp.Router(), http.MethodPost, "/login", loginPayload)
	assertAndLogErrors(t, w, http.StatusOK)

	// validate JWT
	tokenStr := extractTokenFromResponse(t, w)
	assertValidJWT(t, tokenStr, "test_jwt_secret_change_me", email)
}
