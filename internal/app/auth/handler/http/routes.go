package http

import (
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"
	jwtpkg "github.com/sayeed1999/ride-sharing-golang-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, registerUC *usecase.RegisterUsecase, loginUC *usecase.LoginUsecase, jwtService *jwtpkg.Service) {
	authHandler := NewAuthHandler(registerUC, loginUC, jwtService)

	rg.POST("/register", authHandler.Register)
	rg.POST("/login", authHandler.Login)
}
