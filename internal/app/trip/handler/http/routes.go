package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/usecase"
)

func RegisterRoutes(rg *gin.RouterGroup, signupUC *usecase.SignupUsecase) {
	h := NewCustomerHandler(signupUC)
	rg.POST("/signup", h.CustomerSignup)
}
