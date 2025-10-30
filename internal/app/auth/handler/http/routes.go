package http

import "github.com/gin-gonic/gin"

// RegisterRoutes registers routes for an existing AuthHandler instance.
// The handler should be constructed by the module wiring layer (module.go)
// after repositories and usecases have been created.
func RegisterRoutes(rg *gin.RouterGroup, h *AuthHandler) {
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)
}
