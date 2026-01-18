package di

import (
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/handler"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
	tripPostgres "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository/postgres"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/service"
	"gorm.io/gorm"
)

type DIContainer struct {
	CustomerHandler    *handler.CustomerHandler
	DriverHandler      *handler.DriverHandler
	TripRequestHandler *handler.TripRequestHandler
}

func NewDIContainer(db *gorm.DB, cfg *config.Config) *DIContainer {

	// ======== Repositories =========
	var customerRepository repository.CustomerRepository = &tripPostgres.CustomerRepo{DB: db}
	var driverRepository repository.DriverRepository = &tripPostgres.DriverRepo{DB: db}
	var tripRequestRepository repository.TripRequestRepository = &tripPostgres.TripRequestRepo{DB: db}

	// ======== Other module services =========
	authService := auth.NewUserService(db, cfg)

	// ======== Services or services =========
	customerSignupService := &service.CustomerSignupService{
		CustomerRepo: customerRepository,
		AuthService:  authService,
	}

	driverSignupService := &service.DriverSignupService{
		DriverRepo:  driverRepository,
		AuthService: authService,
	}

	requestTripService := &service.RequestTripService{
		TripRequestRepo: tripRequestRepository,
	}

	customerCancelTripService := &service.CustomerCancelTrip{
		TripRequestRepo: tripRequestRepository,
	}

	// ======== Handlers =========
	customerHandler := handler.NewCustomerHandler(customerSignupService)
	driverHandler := handler.NewDriverHandler(driverSignupService)
	tripRequestHandler := handler.NewTripRequestHandler(
		requestTripService,
		customerCancelTripService,
		customerRepository)

	// ======== DI Container =========
	return &DIContainer{
		CustomerHandler:    customerHandler,
		DriverHandler:      driverHandler,
		TripRequestHandler: tripRequestHandler,
	}
}
