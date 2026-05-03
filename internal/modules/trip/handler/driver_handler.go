package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/dto"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/service"
)

type DriverHandler struct {
	DriverService      service.IDriverService
	TripRequestService service.ITripRequestService
}

func NewDriverHandler(driverService service.IDriverService, tripRequestService service.ITripRequestService) *DriverHandler {
	return &DriverHandler{DriverService: driverService, TripRequestService: tripRequestService}
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

func (h *DriverHandler) ListOpenTripRequests(c *gin.Context) {
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
			if limit > 100 {
				limit = 100
			}
		}
	}

	list, err := h.TripRequestService.ListOpenTripRequests(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trip_requests": list})
}
