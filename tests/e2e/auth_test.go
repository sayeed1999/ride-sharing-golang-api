package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRegisterAndLogin_E2E(t *testing.T) {
	ctx := context.Background()
	pgC, cfg := setupContainer(ctx, t)
	defer func() {
		_ = pgC.Terminate(ctx)
	}()

	// allow a short buffer for DB to be fully ready
	time.Sleep(1 * time.Second)

	router := setupRouterWithDB(t, cfg)

	// Register user
	regPayload := map[string]string{"email": "e2e-user@example.com", "password": "pass123", "role": ""}
	w := doJSONRequest(t, router, http.MethodPost, "/register", regPayload)
	require.Equal(t, http.StatusCreated, w.Code)

	// Login with wrong password (should fail)
	badLogin := map[string]string{"email": "e2e-user@example.com", "password": "wrong"}
	w = doJSONRequest(t, router, http.MethodPost, "/login", badLogin)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// Login with correct password (last)
	loginPayload := map[string]string{"email": "e2e-user@example.com", "password": "pass123"}
	w = doJSONRequest(t, router, http.MethodPost, "/login", loginPayload)
	require.Equal(t, http.StatusOK, w.Code)
}

// doJSONRequest is a small helper to marshal payload and perform an HTTP request
func doJSONRequest(t *testing.T, handler http.Handler, method, path string, payload interface{}) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}
