package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	auth "github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth" // Import auth module
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository"
	trippostgres "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository/postgres"
	tripusecase "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/usecase"
	"gorm.io/gorm"
)

// Handlers struct holds all the HTTP handlers for the trip module.
type Handlers struct {
	CustomerHandler    *CustomerHandler
	DriverHandler      *DriverHandler
	TripRequestHandler *TripRequestHandler
}

// newHTTPHandlers creates and injects all dependencies for the HTTP handlers.
func newHTTPHandlers(db *gorm.DB, cfg *config.Config) *Handlers {
	// Repositories
	var custRepo repository.CustomerRepository = &trippostgres.CustomerRepo{DB: db}
	var drvRepo repository.DriverRepository = &trippostgres.DriverRepo{DB: db}
	var trRepo repository.TripRequestRepository = &trippostgres.TripRequestRepo{DB: db}

	// Auth Usecase (shared)
	registerUC := auth.NewRegisterUsecase(db, cfg)

	// Usecases
	signupUC := &tripusecase.CustomerSignupUsecase{
		CustomerRepo: custRepo,
		AuthRegister: registerUC,
	}
	driverSignupUC := &tripusecase.DriverSignupUsecase{
		DriverRepo:   drvRepo,
		AuthRegister: registerUC,
	}
	requestTripUC := &tripusecase.RequestTripUsecase{TripRequestRepo: trRepo}

	// Handlers
	custHandler := NewCustomerHandler(signupUC)
	drvHandler := NewDriverHandler(driverSignupUC)
	tripRequestHandler := NewTripRequestHandler(requestTripUC)

	return &Handlers{
		CustomerHandler:    custHandler,
		DriverHandler:      drvHandler,
		TripRequestHandler: tripRequestHandler,
	}
}

// RegisterAllHTTPRoutes registers all HTTP routes for the trip module.
// It performs dependency injection for the HTTP handlers internally.
func RegisterAllHTTPRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	handlers := newHTTPHandlers(db, cfg)

	customers := rg.Group("/customers")
	{
		customers.POST("/signup", handlers.CustomerHandler.CustomerSignup)
	}

	drivers := rg.Group("/drivers")
	{
		drivers.POST("/signup", handlers.DriverHandler.DriverSignup)
	}

	tripRequests := rg.Group("/trip-requests")
	{
		tripRequests.POST("/request", handlers.TripRequestHandler.RequestTrip)
	}
}
