package trip

import (
	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/di"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/pkg/middleware"
	public_middleware "github.com/sayeed1999/ride-sharing-golang-api/pkg/middleware"
	"gorm.io/gorm"
)

// registerAllHTTPRoutes registers all HTTP routes for the trip module.
func registerAllHTTPRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	c := di.NewDIContainer(db, cfg)

	authMW := public_middleware.AuthMiddleware(cfg.Auth.JWTSecret)
	customerMW := middleware.CustomerMiddleware(c.CustomerRepository)
	driverMW := middleware.DriverMiddleware(c.DriverRepository)

	registerCustomersRoutes(rg, c)
	registerDriversRoutes(rg, c)
	registerCustomerTripRequestRoutes(rg, c, authMW, customerMW)
	registerDriverTripRequestRoutes(rg, c, authMW, driverMW)
	registerTripRoutes(rg, c, authMW, driverMW)
}

func registerCustomersRoutes(rg *gin.RouterGroup, c *di.DIContainer) {
	rg.Group("/customers").
		POST("/signup", c.CustomerHandler.CustomerSignup)
}

func registerDriversRoutes(rg *gin.RouterGroup, c *di.DIContainer) {
	rg.Group("/drivers").
		POST("/signup", c.DriverHandler.DriverSignup)
}

func registerCustomerTripRequestRoutes(rg *gin.RouterGroup, c *di.DIContainer, authMW, customerMW gin.HandlerFunc) {
	g := rg.Group("/trip-requests")
	g.Use(authMW, customerMW) // All routes must be authenticated and the user must be a customer

	tripRequestMW := middleware.TripRequestMiddleware(c.TripRequestRepository)
	g.Group("/:trip_request_id").
		Use(tripRequestMW). // trip request middleware to check if the trip request belongs to the customer
		GET("", c.TripRequestHandler.GetDetails).
		DELETE("", c.TripRequestHandler.CancelTripRequest)

	// Added at the end to avoid matching routes with the ones above
	g.POST("", c.TripRequestHandler.RequestTrip)
}

func registerDriverTripRequestRoutes(rg *gin.RouterGroup, c *di.DIContainer, auth, driverMW gin.HandlerFunc) {
	_ = rg.Group("/trip-requests").
		Use(auth, driverMW). // All routes must be authenticated and the user must be a driver
		GET("/open", c.TripRequestHandler.ListOpenTripRequests).
		POST("/:trip_request_id/accept", c.TripRequestHandler.AcceptTripRequest)
}

func registerTripRoutes(rg *gin.RouterGroup, c *di.DIContainer, auth, driverMW gin.HandlerFunc) {
	tripMW := middleware.TripMiddleware(c.TripRepository)

	g := rg.Group("/trips/:trip_id").Use(auth)
	g.POST("/start", driverMW, tripMW, c.TripHandler.StartTrip)
	g.POST("/complete", driverMW, tripMW, c.TripHandler.CompleteTrip)
	g.POST("/cancel", tripMW, c.TripHandler.CancelTrip)
}
