package domain

import "github.com/google/uuid"

type UserRole struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_role"`
	RoleID uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_user_role"`

	User User `gorm:"foreignKey:UserID"`
	Role Role `gorm:"foreignKey:RoleID"`
}

// TableName sets the insert table name for this struct to the auth schema.
func (ur *UserRole) TableName() string {
	return "auth.user_roles"
}
