package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository"
	trippostgres "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository/postgres"
	tripusecase "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/usecase"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/pkg/middleware"
	public_middleware "github.com/sayeed1999/ride-sharing-golang-api/pkg/middleware"
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

	// Usecases
	authRegisterUC := auth.NewRegisterUsecase(db, cfg)
	authDeleteUC := auth.NewDeleteUserUsecase(db)

	signupUC := &tripusecase.CustomerSignupUsecase{
		CustomerRepo:   custRepo,
		AuthRegister:   authRegisterUC,
		AuthDeleteUser: authDeleteUC,
	}
	driverSignupUC := &tripusecase.DriverSignupUsecase{
		DriverRepo:     drvRepo,
		AuthRegister:   authRegisterUC,
		AuthDeleteUser: authDeleteUC,
	}
	requestTripUC := &tripusecase.RequestTripUsecase{TripRequestRepo: trRepo}
	customerCancelTripUC := &tripusecase.CustomerCancelTrip{TripRequestRepo: trRepo}

	// Handlers
	custHandler := NewCustomerHandler(signupUC)
	drvHandler := NewDriverHandler(driverSignupUC)
	tripRequestHandler := NewTripRequestHandler(requestTripUC, customerCancelTripUC, custRepo)

	return &Handlers{
		CustomerHandler:    custHandler,
		DriverHandler:      drvHandler,
		TripRequestHandler: tripRequestHandler,
	}
}

// RegisterAllHTTPRoutes registers all HTTP routes for the trip module.
// It performs dependency injection for the HTTP handlers internally.
func RegisterAllHTTPRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	var custRepo repository.CustomerRepository = &trippostgres.CustomerRepo{DB: db}
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
	tripRequests.Use(
		public_middleware.AuthMiddleware(cfg.Auth.JWTSecret),
		middleware.CustomerMiddleware(custRepo))
	{
		tripRequests.POST("/request", handlers.TripRequestHandler.RequestTrip)
		tripRequests.DELETE("/:tripID", handlers.TripRequestHandler.CancelTripRequest)
	}
}
