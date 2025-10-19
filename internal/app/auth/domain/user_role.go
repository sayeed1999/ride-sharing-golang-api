package domain

type UserRole struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"uniqueIndex:idx_user_role"`
	RoleID uint `gorm:"uniqueIndex:idx_user_role"`

	User User `gorm:"foreignKey:UserID"`
	Role Role `gorm:"foreignKey:RoleID"`
}
