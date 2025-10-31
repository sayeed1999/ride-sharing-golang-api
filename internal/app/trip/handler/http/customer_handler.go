package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/usecase"
)

type CustomerSignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type CustomerHandler struct {
	SignupUC *usecase.SignupUsecase
}

func NewCustomerHandler(signupUC *usecase.SignupUsecase) *CustomerHandler {
	return &CustomerHandler{SignupUC: signupUC}
}

func (h *CustomerHandler) CustomerSignup(c *gin.Context) {
	var req CustomerSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	cust, err := h.SignupUC.Signup(req.Email, req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"customer": cust})
}
