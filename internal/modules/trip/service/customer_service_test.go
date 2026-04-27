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

func TestCustomerSignup(t *testing.T) {
	customerRepo := new(tripmocks.CustomerRepository)
	authUserRepo := new(authmocks.UserRepository)
	authSvc := authservice.NewUserService(authUserRepo, true)
	svc := NewCustomerService(customerRepo, authSvc)

	authUserID := uuid.New()
	authUserRepo.On("FindByEmail", "customer@example.com").Return(nil, errors.New("not found")).Once()
	authUserRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(&authdomain.User{ID: authUserID, Email: "customer@example.com"}, nil).Once()
	authUserRepo.On("AssignRole", authUserID, "customer").Return(&authdomain.UserRole{ID: uuid.New(), UserID: authUserID}, nil).Once()
	customerRepo.On("CreateCustomer", mock.MatchedBy(func(c *tripdomain.Customer) bool {
		return c.Email == "customer@example.com" && c.Name == "test customer" && c.AuthUserID != nil && *c.AuthUserID == authUserID
	})).Return(&tripdomain.Customer{
		ID:         uuid.New(),
		Email:      "customer@example.com",
		Name:       "test customer",
		AuthUserID: &authUserID,
	}, nil).Once()

	customer, err := svc.Signup("customer@example.com", "test customer", "password123")

	require.NoError(t, err)
	require.NotNil(t, customer)
	assert.Equal(t, "customer@example.com", customer.Email)
	assert.Equal(t, "test customer", customer.Name)
	require.NotNil(t, customer.AuthUserID)
	assert.Equal(t, authUserID, *customer.AuthUserID)

	authUserRepo.AssertExpectations(t)
	customerRepo.AssertExpectations(t)
}

