package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/pkg/password"
)

type UserService struct {
	UserRepository                  repository.IUserRepository
	RequireRoleOnRegistration bool
}

func NewUserService(userRepo repository.IUserRepository, requireRoleOnRegistration bool) *UserService {
	return &UserService{
		UserRepository:                  userRepo,
		RequireRoleOnRegistration: requireRoleOnRegistration,
	}
}

func (s *UserService) Register(email, pass, role string) (*domain.User, error) {
	role = strings.TrimSpace(role)

	if s.RequireRoleOnRegistration {
		if role == "" {
			return nil, errors.New("role is required for registration")
		}
	}

	if _, err := s.UserRepository.FindByEmail(email); err == nil {
		return nil, errors.New("user already exists")
	}

	salt, err := password.GenerateSalt()
	if err != nil {
		return nil, err
	}
	hash, err := password.HashPassword(pass, salt)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: hash,
		PasswordSalt: salt,
	}

	newUser, err := s.UserRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	if role != "" {
		if _, err := s.UserRepository.AssignRole(newUser.ID, role); err != nil {
			return nil, err
		}
	}

	return newUser, nil
}

func (s *UserService) Login(email, pass string) error {
	user, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return errors.New("invalid credentials")
	}

	if !password.VerifyPassword(pass, user.PasswordSalt, user.PasswordHash) {
		return errors.New("invalid credentials")
	}

	return nil
}

func (s *UserService) DeleteUser(userID uuid.UUID) error {
	return s.UserRepository.DeleteUser(userID)
}
