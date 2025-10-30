package http

import "github.com/gin-gonic/gin"

// RegisterRoutes registers routes for a preconstructed CustomerHandler.
// The handler should be created by module wiring after repositories and
// usecases are available.
func RegisterRoutes(rg *gin.RouterGroup, h *CustomerHandler) {
	rg.POST("/signup", h.CustomerSignup)
}
