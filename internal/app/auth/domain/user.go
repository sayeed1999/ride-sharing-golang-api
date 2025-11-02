package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	PasswordSalt string    `gorm:"not null"`
	CreatedAt    time.Time
}

// TableName sets the insert table name for this struct to the auth schema.
func (u *User) TableName() string {
	return "auth.users"
}
