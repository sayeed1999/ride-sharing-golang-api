package testhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

// DoJSONRequest is a small helper to marshal payload and perform an HTTP request
func DoJSONRequest(t testing.TB, handler http.Handler, method, path string, payload interface{}) *httptest.ResponseRecorder {
	t.Helper() // marks this function as a test helper

	var reqBody io.Reader
	if payload != nil {
		body, err := json.Marshal(payload)
		require.NoError(t, err)
		reqBody = bytes.NewReader(body)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

// DoJSONRequestWithAuth is a helper to make JSON requests with an Authorization header
func DoJSONRequestWithAuth(t *testing.T, router *gin.Engine, method, path string, payload interface{}, token string) *httptest.ResponseRecorder {
	t.Helper() // marks this function as a test helper

	var reqBody io.Reader
	if payload != nil {
		body, err := json.Marshal(payload)
		require.NoError(t, err)
		reqBody = bytes.NewReader(body)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func AssertAndLogErrors(t testing.TB, w *httptest.ResponseRecorder, expectedHttpStatus int) {
	t.Helper() // marks this function as a test helper

	if w.Code != expectedHttpStatus {
		t.Logf("Unexpected status code: %d\nBody: %s", w.Code, w.Body.String())
	}

	require.Equal(t, expectedHttpStatus, w.Code)
}

func AssertAndLogErrorsWithBody(t testing.TB, w *httptest.ResponseRecorder, expectedHttpStatus int, expectedBody string) {
	t.Helper() // marks this function as a test helper

	if w.Code != expectedHttpStatus {
		t.Logf("Unexpected status code: %d\nBody: %s", w.Code, w.Body.String())
	}

	require.Equal(t, expectedHttpStatus, w.Code)
	require.Contains(t, w.Body.String(), expectedBody)
}

// ExtractTokenFromResponse unmarshals the JSON response body and returns the "token" field.
func ExtractTokenFromResponse(t testing.TB, w *httptest.ResponseRecorder) string {
	t.Helper()
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	tokenStr, ok := resp["token"]
	require.True(t, ok, "token not found in response")
	return tokenStr
}

// AssertValidJWT is a small helper that parses and validates a JWT token string
// using the provided HMAC secret and asserts the subject claim equals expectedSub.
func AssertValidJWT(t testing.TB, tokenStr, secret, expectedSub string) {
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
