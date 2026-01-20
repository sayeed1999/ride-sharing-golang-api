package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/dto"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/service"
)

type CustomerHandler struct {
	CustomerService service.ICustomerService
}

func NewCustomerHandler(customerService service.ICustomerService) *CustomerHandler {
	return &CustomerHandler{CustomerService: customerService}
}

func (h *CustomerHandler) CustomerSignup(c *gin.Context) {
	var req dto.CustomerSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	cust, err := h.CustomerService.Signup(req.Email, req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"customer": cust})
}
