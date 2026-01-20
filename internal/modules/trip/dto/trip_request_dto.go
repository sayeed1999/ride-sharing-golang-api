package dto

type TripRequestDTO struct {
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

func NewTripRequestDTO(origin, destination string) *TripRequestDTO {
	return &TripRequestDTO{
		Origin:      origin,
		Destination: destination,
	}
}
