package service

import (
	"github.com/google/uuid"
	authdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/domain"
	authmocks "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/repository/mocks"
	authservice "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/service"
	tripdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	tripmocks "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository/mocks"
)

const (
	testCustomerEmail         = "customer@example.com"
	testCustomerName          = "test customer"
	testDriverEmail           = "driver@example.com"
	testDriverName            = "test driver"
	testPassword              = "password123"
	testVehicleType           = "car"
	testInvalidVehicleType    = "plane"
	testVehicleRegistration   = "DHAKA-METRO-GA-12-1234"
	testTripOrigin            = "A"
	testTripDestination       = "B"
	testNotFoundErrorMessage  = "not found"
)

func fixtureAuthUser(id uuid.UUID, email string) *authdomain.User {
	return &authdomain.User{ID: id, Email: email}
}

func fixtureCustomer(authUserID uuid.UUID) *tripdomain.Customer {
	return &tripdomain.Customer{
		ID:         uuid.New(),
		Email:      testCustomerEmail,
		Name:       testCustomerName,
		AuthUserID: &authUserID,
	}
}

func fixtureDriver(authUserID uuid.UUID) *tripdomain.Driver {
	return &tripdomain.Driver{
		ID:                  uuid.New(),
		Email:               testDriverEmail,
		Name:                testDriverName,
		AuthUserID:          &authUserID,
		VehicleTypeEnumCode: int(tripdomain.VehicleEnumCar),
		VehicleRegistration: testVehicleRegistration,
	}
}

func fixtureTripRequest(customerID uuid.UUID) *tripdomain.TripRequest {
	return &tripdomain.TripRequest{
		ID:          uuid.New(),
		CustomerID:  customerID,
		Origin:      testTripOrigin,
		Destination: testTripDestination,
		Status:      tripdomain.NO_DRIVER_FOUND,
	}
}

func setupCustomerService() (*customerService, *tripmocks.CustomerRepository, *authmocks.UserRepository) {
	customerRepo := new(tripmocks.CustomerRepository)
	authUserRepo := new(authmocks.UserRepository)
	authSvc := authservice.NewUserService(authUserRepo, true)
	svc := NewCustomerService(customerRepo, authSvc).(*customerService)
	return svc, customerRepo, authUserRepo
}

func setupDriverService() (*driverService, *tripmocks.DriverRepository, *authmocks.UserRepository) {
	driverRepo := new(tripmocks.DriverRepository)
	authUserRepo := new(authmocks.UserRepository)
	authSvc := authservice.NewUserService(authUserRepo, true)
	svc := NewDriverService(driverRepo, authSvc).(*driverService)
	return svc, driverRepo, authUserRepo
}

func setupTripRequestService() (*tripRequestService, *tripmocks.TripRequestRepository) {
	tripRequestRepo := new(tripmocks.TripRequestRepository)
	svc := NewTripRequestService(tripRequestRepo).(*tripRequestService)
	return svc, tripRequestRepo
}

