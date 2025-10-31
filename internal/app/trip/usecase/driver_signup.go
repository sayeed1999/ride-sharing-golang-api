package usecase

import (
	"errors"

	authusecase "github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"
	tripdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository"
)

// DriverSignupUsecase handles driver signups including vehicle details.
type DriverSignupUsecase struct {
	DriverRepo   repository.DriverRepository
	AuthRegister *authusecase.RegisterUsecase
}

// Signup registers a driver record and creates a corresponding auth user
// with role "driver". VehicleType should be one of the seeded values.
func (uc *DriverSignupUsecase) Signup(email, name, password, vehicleType, vehicleRegistration string) (*tripdomain.Driver, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	// validate & map vehicle type to enum code and canonical name
	ve, _, ok := tripdomain.LookupVehicleEnum(vehicleType)
	if !ok {
		return nil, errors.New("invalid vehicle type")
	}

	driver := &tripdomain.Driver{
		Email:               email,
		Name:                name,
		VehicleTypeEnumCode: int(ve),
		VehicleRegistration: vehicleRegistration,
	}

	created, err := uc.DriverRepo.CreateDriver(driver)
	if err != nil {
		return nil, err
	}

	// register auth user with hardcoded role 'driver'
	if err := uc.AuthRegister.Register(email, password, "driver"); err != nil {
		// compensating delete
		_ = uc.DriverRepo.DeleteDriver(created.ID)
		return nil, err
	}

	return created, nil
}
