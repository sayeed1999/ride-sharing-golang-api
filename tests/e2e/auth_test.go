package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	testhelper "github.com/sayeed1999/ride-sharing-golang-api/pkg/test_helper"
	"github.com/sayeed1999/ride-sharing-golang-api/tests/e2e/setup"
	"github.com/stretchr/testify/require"
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
	tokenStr := extractTokenFromResponse(t, w)

	// one-liner helper validates token and subject
	assertValidJWT(t, tokenStr, "test_jwt_secret_change_me", "e2e-user@example.com")
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

// extractTokenFromResponse unmarshals the JSON response body and returns the "token" field.
func extractTokenFromResponse(t testing.TB, w *httptest.ResponseRecorder) string {
	t.Helper()
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	tokenStr, ok := resp["token"]
	require.True(t, ok, "token not found in response")
	return tokenStr
}

// assertValidJWT is a small helper that parses and validates a JWT token string
// using the provided HMAC secret and asserts the subject claim equals expectedSub.
func assertValidJWT(t testing.TB, tokenStr, secret, expectedSub string) {
	t.Helper()
	parsed, err := jwt.Parse(tokenStr, func(tkn *jwt.Token) (any, error) {
		if tkn.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", tkn.Method.Alg())
		}
		return []byte(secret), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid, "token is not valid")

	claims, ok := parsed.Claims.(jwt.MapClaims)
	require.True(t, ok, "unexpected claims type")
	require.Equal(t, expectedSub, claims["sub"])
}
