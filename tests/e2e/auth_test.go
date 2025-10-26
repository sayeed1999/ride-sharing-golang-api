package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterAndLogin_E2E(t *testing.T) {
	ctx := context.Background()
	pgC, dbConfig := setupContainer(ctx, t)
	defer pgC.Terminate(ctx)
	cfg := buildConfig(t, dbConfig, false)
	db := setupTestDB(t, cfg)
	router := setupRouter(t, db, cfg)

	// Register user
	regPayload := map[string]string{"email": "e2e-user@example.com", "password": "pass123", "role": ""}
	w := doJSONRequest(t, router, http.MethodPost, "/register", regPayload)
	assertAndLogErrors(t, w, http.StatusCreated)

	// Login with wrong password (should fail)
	badLogin := map[string]string{"email": "e2e-user@example.com", "password": "wrong"}
	w = doJSONRequest(t, router, http.MethodPost, "/login", badLogin)
	assertAndLogErrors(t, w, http.StatusUnauthorized)

	// Login with correct password (last)
	loginPayload := map[string]string{"email": "e2e-user@example.com", "password": "pass123"}
	w = doJSONRequest(t, router, http.MethodPost, "/login", loginPayload)
	assertAndLogErrors(t, w, http.StatusOK)
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
