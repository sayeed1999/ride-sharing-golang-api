package http

import (
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/usecase"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, registerUC *usecase.RegisterUsecase, loginUC *usecase.LoginUsecase) {
	authHandler := NewAuthHandler(registerUC, loginUC)

	rg.POST("/register", authHandler.Register)
	rg.POST("/login", authHandler.Login)
}
