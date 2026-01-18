package domain

import "strings"

// VehicleEnum is the numeric enum stored in the DB.
type VehicleEnum int

const (
	VehicleEnumBike VehicleEnum = 1
	VehicleEnumCNG  VehicleEnum = 2
	VehicleEnumCar  VehicleEnum = 3
)

var vehicleNameMap = map[string]VehicleEnum{
	"bike": VehicleEnumBike,
	"cng":  VehicleEnumCNG,
	"car":  VehicleEnumCar,
}

// LookupVehicleEnum looks up the enum and canonical name by a user-provided
// vehicle type string (case-insensitive). Returns (enum, name, found).
func LookupVehicleEnum(s string) (VehicleEnum, string, bool) {
	key := strings.ToLower(strings.TrimSpace(s))
	enum, ok := vehicleNameMap[key]
	if !ok {
		return 0, "", false
	}
	return enum, key, true
}
