package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/service"
	jwtpkg "github.com/sayeed1999/ride-sharing-golang-api/pkg/jwt"

	"gorm.io/gorm"
)

// RegisterAllHTTPRoutes registers all HTTP routes for the auth module.
// It performs dependency injection for the HTTP handlers internally.
func RegisterAllHTTPRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// Repositories
	postgresRepo := &repository.UserRepository{DB: db}
	var userRepo repository.IUserRepository = postgresRepo

	// Services
	userService := service.NewUserService(userRepo, cfg.FeatureFlags.RequireRoleOnRegistration)

	// JWT service (injected)
	jwtService := jwtpkg.New(cfg.Auth.JWTSecret, 24*time.Hour)

	authHandler := NewAuthHandler(userService, jwtService)

	rg.POST("/register", authHandler.Register)
	rg.POST("/login", authHandler.Login)
}
