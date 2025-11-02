package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository/postgres"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"
	jwtpkg "github.com/sayeed1999/ride-sharing-golang-api/pkg/jwt"

	"gorm.io/gorm"
)

// RegisterAllHTTPRoutes registers all HTTP routes for the auth module.
// It performs dependency injection for the HTTP handlers internally.
func RegisterAllHTTPRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
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

	authHandler := NewAuthHandler(registerUC, loginUC, jwtService)

	rg.POST("/register", authHandler.Register)
	rg.POST("/login", authHandler.Login)
}
