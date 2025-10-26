package setup

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
	"gorm.io/gorm"
)

func setupContainer(ctx context.Context, t testing.TB) (tc.Container, config.DatabaseConfig) {
	t.Helper() // marks this function as a test helper

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

	dbConfig := config.DatabaseConfig{
		Host:     host,
		Port:     mappedPort.Port(),
		User:     "testuser",
		Password: "testpass",
		DB:       "ridesharing_testdb",
	}

	// allow a short buffer for DB to be fully ready
	time.Sleep(1 * time.Second)
	return pgC, dbConfig
}

func buildConfig(t testing.TB, dbConfig config.DatabaseConfig, requireRoleOnRegistration bool) *config.Config {
	t.Helper() // marks this function as a test helper

	return &config.Config{
		Server: config.ServerConfig{
			Host: "0.0.0.0",
			Port: "7000",
		},
		Database: dbConfig,
		FeatureFlags: config.FeatureFlags{
			RequireRoleOnRegistration: requireRoleOnRegistration,
		},
		Auth: config.AuthConfig{
			JWTSecret: "test_jwt_secret_change_me",
		},
	}
}

func setupTestDB(t testing.TB, cfg *config.Config) *gorm.DB {
	t.Helper() // marks this function as a test helper

	db, err := database.InitDB(cfg)
	require.NoError(t, err)

	// run migrations
	require.NoError(t, database.AutoMigrate(db))

	return db
}

func setupRouter(t testing.TB, db *gorm.DB, cfg *config.Config) *gin.Engine {
	t.Helper() // marks this function as a test helper

	gin.SetMode(gin.TestMode)
	r := gin.New()
	// expose routes at root so endpoints are /register and /login
	auth.ExposeRoutes(r.Group(""), db, cfg)
	return r
}
