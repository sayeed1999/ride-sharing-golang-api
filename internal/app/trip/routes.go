package trip

import (
	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/di"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository"
	trippostgres "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository/postgres"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/pkg/middleware"
	public_middleware "github.com/sayeed1999/ride-sharing-golang-api/pkg/middleware"
	"gorm.io/gorm"
)

// RegisterAllHTTPRoutes registers all HTTP routes for the trip module.
// It performs dependency injection for the HTTP diContainer internally.
func registerAllHTTPRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	var custRepo repository.CustomerRepository = &trippostgres.CustomerRepo{DB: db}
	var tripRequestRepo repository.TripRequestRepository = &trippostgres.TripRequestRepo{DB: db}

	diContainer := di.NewDIContainer(db, cfg)

	customers := rg.Group("/customers")
	{
		customers.POST("/signup", diContainer.CustomerHandler.CustomerSignup)
	}

	drivers := rg.Group("/drivers")
	{
		drivers.POST("/signup", diContainer.DriverHandler.DriverSignup)
	}

	tripRequests := rg.Group("/trip-requests")
	tripRequests.Use(
		public_middleware.AuthMiddleware(cfg.Auth.JWTSecret),
		middleware.CustomerMiddleware(custRepo))
	{
		tripRequests.POST("", diContainer.TripRequestHandler.RequestTrip)
		tripRequests.DELETE("/:tripID",
			middleware.TripRequestMiddleware(tripRequestRepo),
			diContainer.TripRequestHandler.CancelTripRequest)
	}
}
