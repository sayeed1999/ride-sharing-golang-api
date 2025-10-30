package auth

import (
	"time"

	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/handler/http"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository/postgres"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"
	jwtpkg "github.com/sayeed1999/ride-sharing-golang-api/pkg/jwt"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

// ExposeRoutes wires the auth module routes to a given router group
func ExposeRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// Repositories
	postgresRepo := &postgres.UserRepo{DB: db}
	var userRepo repository.UserRepository = postgresRepo

	// Usecases
	registerUC := &usecase.RegisterUsecase{
		UserRepo:                  userRepo,
		RequireRoleOnRegistration: cfg.FeatureFlags.RequireRoleOnRegistration,
	}
	loginUC := &usecase.LoginUsecase{UserRepo: userRepo}

	// JWT service (injected)
	jwtService := jwtpkg.New(cfg.Auth.JWTSecret, 24*time.Hour)

	// Register HTTP routes (pass jwt service)
	http.RegisterRoutes(rg, registerUC, loginUC, jwtService)
}

// NewRegisterUsecase builds and returns a RegisterUsecase instance backed by the
// provided DB and configuration. Other modules can call this to obtain an
// instance without reaching into auth internals.
func NewRegisterUsecase(db *gorm.DB, cfg *config.Config) *usecase.RegisterUsecase {
	postgresRepo := &postgres.UserRepo{DB: db}
	var userRepo repository.UserRepository = postgresRepo

	registerUC := &usecase.RegisterUsecase{
		UserRepo:                  userRepo,
		RequireRoleOnRegistration: cfg.FeatureFlags.RequireRoleOnRegistration,
	}

	return registerUC
}
