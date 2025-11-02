package domain

// VehicleType maps allowed vehicle types to a simple lookup table.
// Fields: Name (human-friendly) and EnumCode (machine-friendly enum code).
type VehicleType struct {
    ID       string `gorm:"type:uuid;primary_key;" json:"id"`
    Name     string `gorm:"uniqueIndex;size:50;not null" json:"name"`
    EnumCode int    `gorm:"not null;uniqueIndex" json:"enum_code"`
}

func (VehicleType) TableName() string {
    return "trip.vehicle_types"
}
