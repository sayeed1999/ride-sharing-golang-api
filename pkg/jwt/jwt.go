package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Service is a JWT token generator that holds signing secret and expiry.
type Service struct {
	secret string
	expiry time.Duration
}

// New creates a new JWT Service. secret should be provided by DI (config).
func New(secret string, expiry time.Duration) *Service {
	return &Service{secret: secret, expiry: expiry}
}

// GenerateToken signs a token for the provided subject (typically user id or email).
func (s *Service) GenerateToken(subject string) (string, error) {
	if s == nil || s.secret == "" {
		return "", errors.New("jwt secret not configured")
	}

	claims := jwt.MapClaims{
		"sub": subject,
		"exp": time.Now().Add(s.expiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}
