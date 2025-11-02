package domain

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name string    `gorm:"unique;not null"`
}

// TableName sets the insert table name for this struct to the auth schema.
func (r *Role) TableName() string {
	return "auth.roles"
}
