package usecase

import (
	"errors"
	"strings"

	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/pkg/password"
)

type RegisterUsecase struct {
	UserRepo                  repository.UserRepository
	RequireRoleOnRegistration bool
}

func (uc *RegisterUsecase) Register(email, pass, role string) error {
	role = strings.TrimSpace(role)

	// Validate role requirement if feature flag is enabled
	if uc.RequireRoleOnRegistration {
		if role == "" {
			return errors.New("role is required for registration")
		}
	}

	// Check existing user
	if _, err := uc.UserRepo.FindByEmail(email); err == nil {
		return errors.New("user already exists")
	}

	// Generate salt + hash
	salt, err := password.GenerateSalt()
	if err != nil {
		return err
	}
	hash, err := password.HashPassword(pass, salt)
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

	// Assign role if provided (regardless of feature flag)
	if role != "" {
		if _, err := uc.UserRepo.AssignRole(newUser.ID, role); err != nil {
			return err
		}
	}

	return nil
}
