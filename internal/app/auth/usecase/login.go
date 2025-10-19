package usecase

import (
	"errors"

	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository"
)

type LoginUsecase struct {
	UserRepo repository.UserRepository
}

func (uc *LoginUsecase) Login(email, password string) error {
	user, err := uc.UserRepo.FindByEmail(email)
	if err != nil {
		return errors.New("invalid credentials")
	}

	if !verifyPassword(password, user.PasswordSalt, user.PasswordHash) {
		return errors.New("invalid credentials")
	}

	return nil
}
