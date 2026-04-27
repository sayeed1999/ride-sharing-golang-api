package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	authdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/domain"
	tripdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestDriverSignup(t *testing.T) {
	t.Run("happy path: registers auth user and creates driver", func(t *testing.T) {
		svc, driverRepo, authUserRepo := setupDriverService()

		authUserID := uuid.New()
		authUserRepo.On("FindByEmail", testDriverEmail).Return(nil, errors.New(testNotFoundErrorMessage)).Once()
		authUserRepo.On("CreateUser", mock.Anything).Return(fixtureAuthUser(authUserID, testDriverEmail), nil).Once()
		authUserRepo.On("AssignRole", authUserID, "driver").Return(&authdomain.UserRole{ID: uuid.New(), UserID: authUserID}, nil).Once()
		driverRepo.On("CreateDriver", mock.Anything).Return(fixtureDriver(authUserID), nil).Once()

		driver, err := svc.Signup(testDriverEmail, testDriverName, testPassword, testVehicleType, testVehicleRegistration)

		require.NoError(t, err)
		require.NotNil(t, driver)
		assert.Equal(t, testDriverEmail, driver.Email)
		assert.Equal(t, testDriverName, driver.Name)
		require.NotNil(t, driver.AuthUserID)
		assert.Equal(t, authUserID, *driver.AuthUserID)
		assert.Equal(t, int(tripdomain.VehicleEnumCar), driver.VehicleTypeEnumCode)

		authUserRepo.AssertExpectations(t)
		driverRepo.AssertExpectations(t)
	})

	t.Run("invalid vehicle type: deletes created auth user", func(t *testing.T) {
		svc, driverRepo, authUserRepo := setupDriverService()

		authUserID := uuid.New()
		authUserRepo.On("FindByEmail", testDriverEmail).Return(nil, errors.New(testNotFoundErrorMessage)).Once()
		authUserRepo.On("CreateUser", mock.Anything).Return(fixtureAuthUser(authUserID, testDriverEmail), nil).Once()
		authUserRepo.On("AssignRole", authUserID, "driver").Return(&authdomain.UserRole{ID: uuid.New(), UserID: authUserID}, nil).Once()
		authUserRepo.On("DeleteUser", authUserID).Return(nil).Once()

		driver, err := svc.Signup(testDriverEmail, testDriverName, testPassword, testInvalidVehicleType, testVehicleRegistration)

		require.Error(t, err)
		assert.EqualError(t, err, "invalid vehicle type")
		assert.Nil(t, driver)
		driverRepo.AssertNotCalled(t, "CreateDriver", mock.Anything)

		authUserRepo.AssertExpectations(t)
		driverRepo.AssertExpectations(t)
	})
}

