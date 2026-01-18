package di

import (
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/handler"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
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
	var customerRepository repository.ICustomerRepository = &repository.CustomerRepository{DB: db}
	var driverRepository repository.IDriverRepository = &repository.DriverRepository{DB: db}
	var tripRequestRepository repository.ITripRequestRepository = &repository.TripRequestRepository{DB: db}

	// ======== Other module services =========
	authService := auth.NewUserService(db, cfg)

	// ======== Services or services =========
	customerSignupService := &service.CustomerSignupService{
		CustomerRepository: customerRepository,
		AuthService:        authService,
	}

	driverSignupService := &service.DriverSignupService{
		DriverRepository: driverRepository,
		AuthService:      authService,
	}

	requestTripService := &service.RequestTripService{
		TripRequestRepository: tripRequestRepository,
	}

	customerCancelTripService := &service.CustomerCancelTrip{
		TripRequestRepository: tripRequestRepository,
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
