package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/service"
)

type DriverSignupRequest struct {
	Email               string `json:"email" binding:"required,email"`
	Name                string `json:"name" binding:"required"`
	Password            string `json:"password" binding:"required,min=6"`
	VehicleType         string `json:"vehicle_type" binding:"required"`
	VehicleRegistration string `json:"vehicle_registration" binding:"required"`
}

type DriverHandler struct {
	DriverService *service.DriverService
}

func NewDriverHandler(driverService *service.DriverService) *DriverHandler {
	return &DriverHandler{DriverService: driverService}
}

func (h *DriverHandler) DriverSignup(c *gin.Context) {
	var req DriverSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	drv, err := h.DriverService.Signup(req.Email, req.Name, req.Password, req.VehicleType, req.VehicleRegistration)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"driver": drv})
}
