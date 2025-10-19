package usecase

import (
	"errors"

	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository"
)

type RegisterUsecase struct {
	UserRepo repository.UserRepository
}

func (uc *RegisterUsecase) Register(email, password, role string) error {
	// Check existing user
	if _, err := uc.UserRepo.FindByEmail(email); err == nil {
		return errors.New("user already exists")
	}

	// Generate salt + hash
	salt, err := GenerateSalt()
	if err != nil {
		return err
	}
	hash, err := HashPassword(password, salt)
	if err != nil {
		return err
	}

	// Create user
	user := &domain.User{
		Email:        email,
		PasswordHash: hash,
		PasswordSalt: salt,
	}

	newUser, err := uc.UserRepo.CreateUser(user)
	if err != nil {
		return err
	}

	if _, err := uc.UserRepo.AssignRole(newUser.ID, role); err != nil {
		return err
	}

	return nil
}
