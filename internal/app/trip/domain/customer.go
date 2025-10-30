package domain

import "time"

// Customer represents a rider in the trip module. It stores a reference to the
// auth user via AuthUserID (nullable) but does not enforce a DB-level foreign key.
type Customer struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Email      string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Name       string    `gorm:"size:255" json:"name"`
	AuthUserID *uint     `gorm:"index" json:"auth_user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (Customer) TableName() string {
	return "trip.customers"
}
