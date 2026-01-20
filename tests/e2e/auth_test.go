package e2e

import (
	"context"
	"net/http"
	"testing"

	testhelper "github.com/sayeed1999/ride-sharing-golang-api/pkg/test_helper"
	"github.com/sayeed1999/ride-sharing-golang-api/tests/e2e/setup"
)

func TestRegisterAndLogin_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, false)
	defer testApp.CleanUp(ctx, t)

	// Register user
	regPayload := map[string]string{"email": "e2e-user@example.com", "password": "pass123", "role": ""}
	w := testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/register", regPayload)
	testhelper.AssertAndLogErrors(t, w, http.StatusCreated)

	// Login with wrong password (should fail)
	badLogin := map[string]string{"email": "e2e-user@example.com", "password": "wrong"}
	w = testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/login", badLogin)
	testhelper.AssertAndLogErrors(t, w, http.StatusUnauthorized)

	// Login with non-existent email (should fail)
	badLogin = map[string]string{"email": "non-existent-user@example.com", "password": "pass123"}
	w = testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/login", badLogin)
	testhelper.AssertAndLogErrors(t, w, http.StatusUnauthorized)

	// Login with correct password (last)
	loginPayload := map[string]string{"email": "e2e-user@example.com", "password": "pass123"}
	w = testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/login", loginPayload)
	testhelper.AssertAndLogErrors(t, w, http.StatusOK)

	// verify returned token is a valid JWT signed with the test secret
	tokenStr := testhelper.ExtractTokenFromResponse(t, w)

	// one-liner helper validates token and subject
	testhelper.AssertValidJWT(t, tokenStr, "test_jwt_secret_change_me", "e2e-user@example.com")
}

func TestRegisterWithDuplicateEmail_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, false)
	defer testApp.CleanUp(ctx, t)

	// Register user
	regPayload := map[string]string{"email": "duplicate-user@example.com", "password": "pass123", "role": ""}
	w := testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/register", regPayload)
	testhelper.AssertAndLogErrors(t, w, http.StatusCreated)

	// Register with the same email again
	w = testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/register", regPayload)
	testhelper.AssertAndLogErrors(t, w, http.StatusBadRequest)
}

func TestRegister_Validation_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t, false)
	defer testApp.CleanUp(ctx, t)

	cases := []struct {
		name    string
		payload map[string]string
	}{
		{
			name:    "invalid email",
			payload: map[string]string{"email": "e2e-user", "password": "pass123"},
		},
		{
			name:    "weak password",
			payload: map[string]string{"email": "e2e-user@example.com", "password": "p"},
		},
		{
			name:    "empty email",
			payload: map[string]string{"email": "", "password": "pass123"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := testhelper.DoJSONRequest(t, testApp.Router(), http.MethodPost, "/register", tc.payload)
			testhelper.AssertAndLogErrors(t, w, http.StatusBadRequest)
		})
	}
}
