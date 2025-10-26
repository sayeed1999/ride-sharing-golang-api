package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sayeed1999/ride-sharing-golang-api/tests/e2e/setup"
	"github.com/stretchr/testify/require"
)

func TestRegisterAndLogin_E2E(t *testing.T) {
	ctx := context.Background()
	testApp := setup.NewTestApp(ctx, t)
	defer testApp.CleanUp(ctx, t)

	// Register user
	regPayload := map[string]string{"email": "e2e-user@example.com", "password": "pass123", "role": ""}
	w := doJSONRequest(t, testApp.Router(), http.MethodPost, "/register", regPayload)
	assertAndLogErrors(t, w, http.StatusCreated)

	// Login with wrong password (should fail)
	badLogin := map[string]string{"email": "e2e-user@example.com", "password": "wrong"}
	w = doJSONRequest(t, testApp.Router(), http.MethodPost, "/login", badLogin)
	assertAndLogErrors(t, w, http.StatusUnauthorized)

	// Login with correct password (last)
	loginPayload := map[string]string{"email": "e2e-user@example.com", "password": "pass123"}
	w = doJSONRequest(t, testApp.Router(), http.MethodPost, "/login", loginPayload)
	assertAndLogErrors(t, w, http.StatusOK)

	// verify returned token is a valid JWT signed with the test secret
	tokenStr := extractTokenFromResponse(t, w)

	// one-liner helper validates token and subject
	assertValidJWT(t, tokenStr, "test_jwt_secret_change_me", "e2e-user@example.com")
}

func assertAndLogErrors(t testing.TB, w *httptest.ResponseRecorder, expectedHttpStatus int) {
	t.Helper() // marks this function as a test helper

	if w.Code != expectedHttpStatus {
		t.Logf("Unexpected status code: %d\nBody: %s", w.Code, w.Body.String())
	}

	require.Equal(t, expectedHttpStatus, w.Code)
}

// doJSONRequest is a small helper to marshal payload and perform an HTTP request
func doJSONRequest(t testing.TB, handler http.Handler, method, path string, payload interface{}) *httptest.ResponseRecorder {
	t.Helper() // marks this function as a test helper

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
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
