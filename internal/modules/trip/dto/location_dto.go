package dto

type UpdateDriverLocationRequest struct {
	Lat     float64  `json:"lat" binding:"required"`
	Lng     float64  `json:"lng" binding:"required"`
	Heading *float64 `json:"heading"`
	Speed   *float64 `json:"speed"`
}

func NewUpdateDriverLocationRequest(lat, lng float64, heading, speed *float64) *UpdateDriverLocationRequest {
	return &UpdateDriverLocationRequest{
		Lat:     lat,
		Lng:     lng,
		Heading: heading,
		Speed:   speed,
	}
}
