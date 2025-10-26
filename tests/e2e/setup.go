package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/database"
	auth "github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupContainer(ctx context.Context, t *testing.T) (tc.Container, *config.Config) {
	req := tc.ContainerRequest{
		Image:        "postgres:18-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "ridesharing_testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	pgC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{ContainerRequest: req, Started: true})
	require.NoError(t, err)

	host, err := pgC.Host(ctx)
	require.NoError(t, err)
	mappedPort, err := pgC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	cfg := &config.Config{
		Server: config.ServerConfig{Host: "0.0.0.0", Port: "7000"},
		Database: config.DatabaseConfig{
			User:     "testuser",
			Password: "testpass",
			Host:     host,
			Port:     mappedPort.Port(),
			DB:       "ridesharing_testdb",
		},
		FeatureFlags: config.FeatureFlags{RequireRoleOnRegistration: false},
	}

	return pgC, cfg
}

func setupRouterWithDB(t *testing.T, cfg *config.Config) *gin.Engine {
	db, err := database.InitDB(cfg)
	require.NoError(t, err)

	// run migrations
	require.NoError(t, database.AutoMigrate(db))

	gin.SetMode(gin.TestMode)
	r := gin.New()
	// expose routes at root so endpoints are /register and /login
	auth.ExposeRoutes(r.Group(""), db, cfg)
	return r
}
