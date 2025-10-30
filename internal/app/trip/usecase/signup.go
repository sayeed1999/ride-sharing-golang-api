package usecase

import (
	"errors"

	authusecase "github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository"
)

// SignupUsecase handles business logic for customer signups in trip module.
// It creates a customer record and then registers an auth user by calling
// the auth module's Register usecase directly. If registering the auth user
// fails, the created customer is deleted (compensating action).
type SignupUsecase struct {
	CustomerRepo repository.CustomerRepository
	AuthRegister *authusecase.RegisterUsecase
}

// Signup creates a customer and registers an auth user with the hardcoded
// role "customer". The created customer's AuthUserID is set when auth
// registration succeeds (currently auth.Register does not return the ID).
func (uc *SignupUsecase) Signup(email, name, password string) (*domain.Customer, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	// Create customer record first
	customer := &domain.Customer{
		Email: email,
		Name:  name,
	}
	newCustomer, err := uc.CustomerRepo.CreateCustomer(customer)
	if err != nil {
		return nil, err
	}

	// Register auth user via auth module (internal call) with hardcoded role
	if err := uc.AuthRegister.Register(email, password, "customer"); err != nil {
		// Compensating delete: remove the created customer
		_ = uc.CustomerRepo.DeleteCustomer(newCustomer.ID)
		return nil, err
	}

	// Note: auth.Register does not return created user ID. If you need the
	// numeric auth user ID stored in customer.AuthUserID we would need the
	// Register usecase to return the new user's ID. For now we leave
	// AuthUserID nil to avoid changing the auth module contract.

	return newCustomer, nil
}
