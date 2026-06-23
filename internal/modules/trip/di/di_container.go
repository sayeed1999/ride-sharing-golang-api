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
	// repositories
	CustomerRepository    repository.ICustomerRepository
	DriverRepository      repository.IDriverRepository
	TripRequestRepository repository.ITripRequestRepository
	TripRepository        repository.ITripRepository

	// services
	CustomerService    service.ICustomerService
	DriverService      service.IDriverService
	TripRequestService service.ITripRequestService
	TripService        service.ITripService

	// handlers
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
	tripRequestService := service.NewTripRequestService(db, tripRequestRepository, tripRepository)
	tripService := service.NewTripService(db, tripRepository)

	// ======== Handlers =========
	customerHandler := handler.NewCustomerHandler(customerService)
	driverHandler := handler.NewDriverHandler(driverService)
	tripRequestHandler := handler.NewTripRequestHandler(tripRequestService)
	tripHandler := handler.NewTripHandler(tripService, customerRepository, driverRepository)

	// ======== DI Container =========
	return &DIContainer{
		CustomerRepository:    customerRepository,
		DriverRepository:      driverRepository,
		TripRequestRepository: tripRequestRepository,
		TripRepository:        tripRepository,

		CustomerService:    customerService,
		DriverService:      driverService,
		TripRequestService: tripRequestService,
		TripService:        tripService,

		CustomerHandler:    customerHandler,
		DriverHandler:      driverHandler,
		TripRequestHandler: tripRequestHandler,
		TripHandler:        tripHandler,
	}
}
