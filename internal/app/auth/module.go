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

// ExposeRoutes registers the auth module's HTTP routes to a given router group.
// It delegates the creation of HTTP handlers and their registration to the http handler package.
func ExposeRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	http.RegisterAllHTTPRoutes(rg, db, cfg)
}

// NewRegisterUsecase builds and returns a RegisterUsecase instance backed by the
// provided DB and configuration. Other modules can call this to obtain an
// instance without reaching into auth internals.
func NewRegisterUsecase(db *gorm.DB, cfg *config.Config) *usecase.RegisterUsecase {
	var userRepo repository.UserRepository = &postgres.UserRepo{DB: db}

	registerUC := &usecase.RegisterUsecase{
		UserRepo:                  userRepo,
		RequireRoleOnRegistration: cfg.FeatureFlags.RequireRoleOnRegistration,
	}

	return registerUC
}

// NewDeleteUserUsecase builds and returns a DeleteUserUsecase instance.
func NewDeleteUserUsecase(db *gorm.DB) *usecase.DeleteUserUsecase {
	var userRepo repository.UserRepository = &postgres.UserRepo{DB: db}

	deleteUC := &usecase.DeleteUserUsecase{
		UserRepo: userRepo,
	}

	return deleteUC
}
