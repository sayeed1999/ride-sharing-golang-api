package setup

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	tc "github.com/testcontainers/testcontainers-go"
	"gorm.io/gorm"
)

type TestApp struct {
	router         *gin.Engine
	pgContainer    tc.Container
	redisContainer tc.Container
	DB             *gorm.DB
	Redis          *redis.Client
}

func NewTestApp(ctx context.Context, t *testing.T, rrequireRoleOnRegistration bool) *TestApp {

	pgContainer, dbConfig := setupContainer(ctx, t)
	redisContainer, redisURL := setupRedisContainer(ctx, t)

	cfg := buildConfig(t, dbConfig, redisURL, rrequireRoleOnRegistration)
	db := setupTestDB(t, cfg)
	rdb := setupRedis(t, cfg)
	router := setupRouter(t, db, cfg, rdb)

	return &TestApp{
		router:         router,
		pgContainer:    pgContainer,
		redisContainer: redisContainer,
		DB:             db,
		Redis:          rdb,
	}

}

func (testApp *TestApp) Router() *gin.Engine {
	return testApp.router
}

func (testApp *TestApp) CleanUp(ctx context.Context, t *testing.T) {
	t.Helper()

	if testApp.Redis != nil {
		_ = testApp.Redis.Close()
	}
	if err := testApp.redisContainer.Terminate(ctx); err != nil {
		t.Fatalf("failed to terminate redis test container: %v", err)
	}
	if err := testApp.pgContainer.Terminate(ctx); err != nil {
		t.Fatalf("failed to terminate test container: %v", err)
	}
}
