package dto

type DriverSignupRequest struct {
	Email               string `json:"email" binding:"required,email"`
	Name                string `json:"name" binding:"required"`
	Password            string `json:"password" binding:"required,min=6"`
	VehicleType         string `json:"vehicle_type" binding:"required"`
	VehicleRegistration string `json:"vehicle_registration" binding:"required"`
}

func NewDriverSignupRequest(email, name, password, vehicleType, vehicleRegistration string) *DriverSignupRequest {
	return &DriverSignupRequest{
		Email:               email,
		Name:                name,
		Password:            password,
		VehicleType:         vehicleType,
		VehicleRegistration: vehicleRegistration,
	}
}
