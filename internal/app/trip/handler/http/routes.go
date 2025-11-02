package http

import "github.com/gin-gonic/gin"

// RegisterCustomerRoutes registers routes for a preconstructed CustomerHandler.
// The handler should be created by module wiring after repositories and
// usecases are available.
func RegisterCustomerRoutes(rg *gin.RouterGroup, h *CustomerHandler) {
	rg.POST("/signup", h.CustomerSignup)
}

// RegisterDriverRoutes registers routes for driver-related handlers.
func RegisterDriverRoutes(rg *gin.RouterGroup, h *DriverHandler) {
	rg.POST("/signup", h.DriverSignup)
}

// RegisterTripRequestRoutes registers routes for trip request-related handlers.
func RegisterTripRequestRoutes(rg *gin.RouterGroup, h *TripRequestHandler) {
	rg.POST("/request", h.RequestTrip)
}
