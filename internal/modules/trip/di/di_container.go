package di

import (
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/handler"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
	tripPostgres "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository/postgres"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/usecase"
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

	// ======== Other module usecases =========
	authService := auth.NewUserService(db, cfg)

	// ======== Usecases or services =========
	customerSignupUsecase := &usecase.CustomerSignupUsecase{
		CustomerRepo: customerRepository,
		AuthService:  authService,
	}

	driverSignupUsecase := &usecase.DriverSignupUsecase{
		DriverRepo:  driverRepository,
		AuthService: authService,
	}

	requestTripUsecase := &usecase.RequestTripUsecase{
		TripRequestRepo: tripRequestRepository,
	}

	customerCancelTripUsecase := &usecase.CustomerCancelTrip{
		TripRequestRepo: tripRequestRepository,
	}

	// ======== Handlers =========
	customerHandler := handler.NewCustomerHandler(customerSignupUsecase)
	driverHandler := handler.NewDriverHandler(driverSignupUsecase)
	tripRequestHandler := handler.NewTripRequestHandler(
		requestTripUsecase,
		customerCancelTripUsecase,
		customerRepository)

	// ======== DI Container =========
	return &DIContainer{
		CustomerHandler:    customerHandler,
		DriverHandler:      driverHandler,
		TripRequestHandler: tripRequestHandler,
	}
}
