package usecase

import (
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/repository"
)

type DeleteUserUsecase struct {
	UserRepo repository.UserRepository
}

func (uc *DeleteUserUsecase) DeleteUser(userID uuid.UUID) error {
	return uc.UserRepo.DeleteUser(userID)
}
