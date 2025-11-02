package domain

type UserRole struct {
	ID     string `gorm:"type:uuid;primary_key;"`
	UserID string `gorm:"type:uuid;uniqueIndex:idx_user_role"`
	RoleID string `gorm:"type:uuid;uniqueIndex:idx_user_role"`

	User User `gorm:"foreignKey:UserID"`
	Role Role `gorm:"foreignKey:RoleID"`
}

// TableName sets the insert table name for this struct to the auth schema.
func (ur *UserRole) TableName() string {
	return "auth.user_roles"
}
