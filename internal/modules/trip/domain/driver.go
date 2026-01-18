package domain

import (
	"time"

	"github.com/google/uuid"
)

// Driver represents a driver in the trip module. It stores a reference to the
// auth user via AuthUserID (nullable) and vehicle details.
type Driver struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email               string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Name                string     `gorm:"size:255" json:"name"`
	AuthUserID          *uuid.UUID `gorm:"type:uuid;index" json:"auth_user_id"`
	VehicleTypeEnumCode int        `gorm:"index" json:"vehicle_type_enum_code"`
	VehicleRegistration string     `gorm:"size:100" json:"vehicle_registration"`
	CreatedAt           time.Time  `json:"created_at"`
}

func (Driver) TableName() string {
	return "trip.drivers"
}
