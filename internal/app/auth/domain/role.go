package domain

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

// TableName sets the insert table name for this struct to the auth schema.
func (r *Role) TableName() string {
	return "auth.roles"
}
