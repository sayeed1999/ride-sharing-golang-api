package auth

import (
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/handler/http"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/repository/postgres"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/service"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

// ExposeRoutes registers the auth module's HTTP routes to a given router group.
// It delegates the creation of HTTP handlers and their registration to the http handler package.
func ExposeRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	http.RegisterAllHTTPRoutes(rg, db, cfg)
}

// NewUserService builds and returns a UserService instance backed by the
// provided DB and configuration. Other modules can call this to obtain an
// instance without reaching into auth internals.
func NewUserService(db *gorm.DB, cfg *config.Config) *service.UserService {
	var userRepo repository.UserRepository = &postgres.UserRepo{DB: db}

	userService := service.NewUserService(userRepo, cfg.FeatureFlags.RequireRoleOnRegistration)

	return userService
}
