package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	authdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/domain"
	authmocks "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/repository/mocks"
	authservice "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/service"
	tripdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	tripmocks "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDriverSignup(t *testing.T) {
	driverRepo := new(tripmocks.DriverRepository)
	authUserRepo := new(authmocks.UserRepository)
	authSvc := authservice.NewUserService(authUserRepo, true)
	svc := NewDriverService(driverRepo, authSvc)

	authUserID := uuid.New()
	authUserRepo.On("FindByEmail", "driver@example.com").Return(nil, errors.New("not found")).Once()
	authUserRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(&authdomain.User{ID: authUserID, Email: "driver@example.com"}, nil).Once()
	authUserRepo.On("AssignRole", authUserID, "driver").Return(&authdomain.UserRole{ID: uuid.New(), UserID: authUserID}, nil).Once()
	driverRepo.On("CreateDriver", mock.MatchedBy(func(d *tripdomain.Driver) bool {
		return d.Email == "driver@example.com" &&
			d.Name == "test driver" &&
			d.AuthUserID != nil &&
			*d.AuthUserID == authUserID &&
			d.VehicleTypeEnumCode == int(tripdomain.VehicleEnumCar) &&
			d.VehicleRegistration == "DHAKA-METRO-GA-12-1234"
	})).Return(&tripdomain.Driver{
		ID:                  uuid.New(),
		Email:               "driver@example.com",
		Name:                "test driver",
		AuthUserID:          &authUserID,
		VehicleTypeEnumCode: int(tripdomain.VehicleEnumCar),
		VehicleRegistration: "DHAKA-METRO-GA-12-1234",
	}, nil).Once()

	driver, err := svc.Signup("driver@example.com", "test driver", "password123", "car", "DHAKA-METRO-GA-12-1234")

	require.NoError(t, err)
	require.NotNil(t, driver)
	assert.Equal(t, "driver@example.com", driver.Email)
	assert.Equal(t, "test driver", driver.Name)
	require.NotNil(t, driver.AuthUserID)
	assert.Equal(t, authUserID, *driver.AuthUserID)
	assert.Equal(t, int(tripdomain.VehicleEnumCar), driver.VehicleTypeEnumCode)

	authUserRepo.AssertExpectations(t)
	driverRepo.AssertExpectations(t)
}

