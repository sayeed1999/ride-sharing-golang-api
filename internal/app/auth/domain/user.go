package domain

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	PasswordSalt string `gorm:"not null"`
	CreatedAt    time.Time
}
