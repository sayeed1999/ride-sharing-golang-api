package usecase

import (
	"errors"

	authusecase "github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository"
)

// CustomerSignupUsecase handles the business logic for customer signups.
// It first registers a user in the auth module and then creates a customer
// record in the trip module. If customer creation fails, it deletes the
// previously created auth user as a compensating action.
type CustomerSignupUsecase struct {
	CustomerRepo   repository.CustomerRepository
	AuthRegister   *authusecase.RegisterUsecase
	AuthDeleteUser *authusecase.DeleteUserUsecase // For compensating actions
}

// Signup registers an auth user and then creates a corresponding customer record.
func (uc *CustomerSignupUsecase) Signup(email, name, password string) (*domain.Customer, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	// 1. Register auth user first
	authUser, err := uc.AuthRegister.Register(email, password, "customer")
	if err != nil {
		return nil, err
	}

	// 2. Create the customer record with the new AuthUserID
	customer := &domain.Customer{
		Email:      email,
		Name:       name,
		AuthUserID: &authUser.ID,
	}

	newCustomer, err := uc.CustomerRepo.CreateCustomer(customer)
	if err != nil {
		// Compensating action: delete the created auth user
		_ = uc.AuthDeleteUser.DeleteUser(authUser.ID)
		return nil, err
	}

	return newCustomer, nil
}
