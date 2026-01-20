package dto

type TripRequestRequest struct {
	CustomerID  string `json:"customer_id" binding:"required"`
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

func NewTripRequestRequest(customerID, origin, destination string) *TripRequestRequest {
	return &TripRequestRequest{
		CustomerID:  customerID,
		Origin:      origin,
		Destination: destination,
	}
}
