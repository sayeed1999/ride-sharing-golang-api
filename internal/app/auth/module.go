package auth

import (
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/handler/http"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository/postgres"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"

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

	// Register HTTP routes
	http.RegisterRoutes(rg, registerUC, loginUC)
}
