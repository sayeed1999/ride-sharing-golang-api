package setup

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	tc "github.com/testcontainers/testcontainers-go"
)

type TestApp struct {
	router        *gin.Engine
	testcontainer tc.Container
}

func NewTestApp(ctx context.Context, t *testing.T) *TestApp {

	pgContainer, dbConfig := setupContainer(ctx, t)

	cfg := buildConfig(t, dbConfig, false)
	db := setupTestDB(t, cfg)
	router := setupRouter(t, db, cfg)

	return &TestApp{
		router:        router,
		testcontainer: pgContainer,
	}

}

func (testApp *TestApp) Router() *gin.Engine {
	return testApp.router
}

func (testApp *TestApp) CleanUp(ctx context.Context, t *testing.T) {
	t.Helper()

	if err := testApp.testcontainer.Terminate(ctx); err != nil {
		t.Fatalf("failed to terminate test container: %v", err)
	}
}
