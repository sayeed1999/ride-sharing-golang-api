package di

import (
	"github.com/redis/go-redis/v9"
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
	LocationRepository    repository.ILocationRepository

	// services
	CustomerService    service.ICustomerService
	DriverService      service.IDriverService
	TripRequestService service.ITripRequestService
	TripService        service.ITripService
	LocationService    service.ILocationService

	// handlers
	CustomerHandler    *handler.CustomerHandler
	DriverHandler      *handler.DriverHandler
	TripRequestHandler *handler.TripRequestHandler
	TripHandler        *handler.TripHandler
	LocationHandler    *handler.LocationHandler
}

func NewDIContainer(db *gorm.DB, cfg *config.Config, rdb *redis.Client) *DIContainer {

	// ======== Repositories =========
	var customerRepository repository.ICustomerRepository = &repository.CustomerRepository{DB: db}
	var driverRepository repository.IDriverRepository = &repository.DriverRepository{DB: db}
	var tripRequestRepository repository.ITripRequestRepository = &repository.TripRequestRepository{DB: db}
	var tripRepository repository.ITripRepository = &repository.TripRepository{DB: db}
	var locationRepository repository.ILocationRepository = &repository.LocationRepository{Client: rdb}

	// ======== Other module services =========
	authService := auth.NewUserService(db, cfg)

	// ======== Services or services =========
	customerService := service.NewCustomerService(customerRepository, authService)
	driverService := service.NewDriverService(driverRepository, authService)
	tripRequestService := service.NewTripRequestService(db, tripRequestRepository, tripRepository)
	tripService := service.NewTripService(db, tripRepository)
	locationService := service.NewLocationService(tripRepository, locationRepository)

	// ======== Handlers =========
	customerHandler := handler.NewCustomerHandler(customerService)
	driverHandler := handler.NewDriverHandler(driverService)
	tripRequestHandler := handler.NewTripRequestHandler(tripRequestService)
	tripHandler := handler.NewTripHandler(tripService, customerRepository, driverRepository)
	locationHandler := handler.NewLocationHandler(locationService, customerRepository, driverRepository)

	// ======== DI Container =========
	return &DIContainer{
		CustomerRepository:    customerRepository,
		DriverRepository:      driverRepository,
		TripRequestRepository: tripRequestRepository,
		TripRepository:        tripRepository,
		LocationRepository:    locationRepository,

		CustomerService:    customerService,
		DriverService:      driverService,
		TripRequestService: tripRequestService,
		TripService:        tripService,
		LocationService:    locationService,

		CustomerHandler:    customerHandler,
		DriverHandler:      driverHandler,
		TripRequestHandler: tripRequestHandler,
		TripHandler:        tripHandler,
		LocationHandler:    locationHandler,
	}
}
