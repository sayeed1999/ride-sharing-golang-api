package domain

type UserRole struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint
	RoleID uint

	User User `gorm:"foreignKey:UserID"`
	Role Role `gorm:"foreignKey:RoleID"`
}
