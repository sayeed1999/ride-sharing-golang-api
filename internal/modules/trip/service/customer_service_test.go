package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	authdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCustomerSignup(t *testing.T) {
	t.Run("customer signup [happy path]: registers auth user -> assign role -> creates customer", func(t *testing.T) {
		svc, customerRepo, authUserRepo := setupCustomerService()

		authUserID := uuid.New()
		authUserRepo.On("FindByEmail", testCustomerEmail).Return(nil, errors.New(testNotFoundErrorMessage)).Once()
		authUserRepo.On("CreateUser", mock.Anything).Return(fixtureAuthUser(authUserID, testCustomerEmail), nil).Once()
		authUserRepo.On("AssignRole", authUserID, "customer").Return(&authdomain.UserRole{ID: uuid.New(), UserID: authUserID}, nil).Once()
		customerRepo.On("CreateCustomer", mock.Anything).Return(fixtureCustomer(authUserID), nil).Once()

		customer, err := svc.Signup(testCustomerEmail, testCustomerName, testPassword)

		require.NoError(t, err)
		require.NotNil(t, customer)
		assert.Equal(t, testCustomerEmail, customer.Email)
		assert.Equal(t, testCustomerName, customer.Name)
		require.NotNil(t, customer.AuthUserID)
		assert.Equal(t, authUserID, *customer.AuthUserID)

		authUserRepo.AssertExpectations(t)
		customerRepo.AssertExpectations(t)
	})

	t.Run("auth register fails: customer is not created", func(t *testing.T) {
		svc, customerRepo, authUserRepo := setupCustomerService()

		authErr := errors.New("auth register failed")
		authUserRepo.On("FindByEmail", testCustomerEmail).Return(nil, errors.New(testNotFoundErrorMessage)).Once()
		authUserRepo.On("CreateUser", mock.Anything).Return(nil, authErr).Once()

		customer, err := svc.Signup(testCustomerEmail, testCustomerName, testPassword)

		require.Error(t, err)
		assert.ErrorIs(t, err, authErr)
		assert.Nil(t, customer)

		customerRepo.AssertNotCalled(t, "CreateCustomer", mock.Anything)
		authUserRepo.AssertExpectations(t)
		customerRepo.AssertExpectations(t)
	})
}
