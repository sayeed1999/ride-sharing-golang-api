package service

import (
	"errors"

	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/service"
	tripdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
)

// DriverSignupService handles driver signups including vehicle details.
type DriverSignupService struct {
	DriverRepo  repository.DriverRepository
	AuthService *service.UserService // For compensating actions
}

// Signup registers an auth user and then creates a corresponding driver record.
func (uc *DriverSignupService) Signup(email, name, password, vehicleType, vehicleRegistration string) (*tripdomain.Driver, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	// 1. Register auth user first
	authUser, err := uc.AuthService.Register(email, password, "driver")
	if err != nil {
		return nil, err
	}

	// 2. Validate vehicle type
	ve, _, ok := tripdomain.LookupVehicleEnum(vehicleType)
	if !ok {
		// Compensating action: delete the created auth user
		_ = uc.AuthService.DeleteUser(authUser.ID)
		return nil, errors.New("invalid vehicle type")
	}

	// 3. Create the driver record with the new AuthUserID
	driver := &tripdomain.Driver{
		Email:               email,
		Name:                name,
		AuthUserID:          &authUser.ID,
		VehicleTypeEnumCode: int(ve),
		VehicleRegistration: vehicleRegistration,
	}

	created, err := uc.DriverRepo.CreateDriver(driver)
	if err != nil {
		// Compensating action: delete the created auth user
		_ = uc.AuthService.DeleteUser(authUser.ID)
		return nil, err
	}

	return created, nil
}
