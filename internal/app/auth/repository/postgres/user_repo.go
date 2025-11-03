package postgres

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/domain"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func (r *UserRepo) CreateUser(user *domain.User) (*domain.User, error) {
	if err := r.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) AssignRole(userID uuid.UUID, roleName string) (*domain.UserRole, error) {
	var role domain.Role
	if err := r.DB.Where("name = ?", roleName).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	// 🔍 Check if the user already has this role
	var existing domain.UserRole
	err := r.DB.Where("user_id = ? AND role_id = ?", userID, role.ID).First(&existing).Error
	if err == nil {
		// ✅ Already assigned — return existing record (no duplicate)
		return &existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err // unexpected DB error
	}

	// 🆕 Not assigned yet → create new relation
	userRole := domain.UserRole{
		UserID: userID,
		RoleID: role.ID,
	}
	if err := r.DB.Create(&userRole).Error; err != nil {
		return nil, err
	}

	return &userRole, nil
}

func (r *UserRepo) DeleteUser(userID uuid.UUID) error {
	// Also deletes associated user roles due to CASCADE constraint
	if err := r.DB.Delete(&domain.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}
