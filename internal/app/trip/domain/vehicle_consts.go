package domain

import "strings"

// VehicleEnum is the numeric enum stored in the DB.
type VehicleEnum int

const (
	VehicleEnumBike VehicleEnum = 1
	VehicleEnumCNG  VehicleEnum = 2
	VehicleEnumCar  VehicleEnum = 3
)

var vehicleNameMap = map[string]struct {
	enum VehicleEnum
	name string
}{
	"bike": {VehicleEnumBike, "bike"},
	"cng":  {VehicleEnumCNG, "cng"},
	"car":  {VehicleEnumCar, "car"},
}

// LookupVehicleEnum looks up the enum and canonical name by a user-provided
// vehicle type string (case-insensitive). Returns (enum, name, found).
func LookupVehicleEnum(s string) (VehicleEnum, string, bool) {
	key := strings.ToLower(strings.TrimSpace(s))
	v, ok := vehicleNameMap[key]
	if !ok {
		return 0, "", false
	}
	return v.enum, v.name, true
}
