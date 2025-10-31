package trip

import "strings"

// VehicleEnum is the numeric enum stored in the DB.
type VehicleEnum int

const (
	VehicleEnumBike VehicleEnum = 1
	VehicleEnumCNG  VehicleEnum = 2
	VehicleEnumCar  VehicleEnum = 3
)

// LookupVehicleEnum maps a user-provided vehicle type string (case-insensitive)
// to its numeric enum and canonical lower-case name. It's intentionally small
// and implemented with a short switch for clarity and low overhead.
func LookupVehicleEnum(s string) (VehicleEnum, string, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "bike":
		return VehicleEnumBike, "bike", true
	case "cng":
		return VehicleEnumCNG, "cng", true
	case "car":
		return VehicleEnumCar, "car", true
	default:
		return 0, "", false
	}
}
