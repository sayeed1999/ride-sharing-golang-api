package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/dto"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/service"
)

type DriverHandler struct {
	DriverService service.IDriverService
}

func NewDriverHandler(driverService service.IDriverService) *DriverHandler {
	return &DriverHandler{DriverService: driverService}
}

func (h *DriverHandler) DriverSignup(c *gin.Context) {
	var req dto.DriverSignupRequest
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
