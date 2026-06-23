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
	TripHandler        *handler.TripHandler
}

func NewDIContainer(db *gorm.DB, cfg *config.Config) *DIContainer {

	// ======== Repositories =========
	var customerRepository repository.ICustomerRepository = &repository.CustomerRepository{DB: db}
	var driverRepository repository.IDriverRepository = &repository.DriverRepository{DB: db}
	var tripRequestRepository repository.ITripRequestRepository = &repository.TripRequestRepository{DB: db}
	var tripRepository repository.ITripRepository = &repository.TripRepository{DB: db}

	// ======== Other module services =========
	authService := auth.NewUserService(db, cfg)

	// ======== Services or services =========
	customerService := service.NewCustomerService(customerRepository, authService)
	driverService := service.NewDriverService(driverRepository, authService)
	tripRequestService := service.NewTripRequestService(tripRequestRepository)
	tripService := service.NewTripService(db, tripRequestRepository, tripRepository)

	// ======== Handlers =========
	customerHandler := handler.NewCustomerHandler(customerService)
	driverHandler := handler.NewDriverHandler(driverService, tripRequestService, tripService)
	tripRequestHandler := handler.NewTripRequestHandler(tripRequestService)
	tripHandler := handler.NewTripHandler(tripService, customerRepository, driverRepository)

	// ======== DI Container =========
	return &DIContainer{
		CustomerHandler:    customerHandler,
		DriverHandler:      driverHandler,
		TripRequestHandler: tripRequestHandler,
		TripHandler:        tripHandler,
	}
}
