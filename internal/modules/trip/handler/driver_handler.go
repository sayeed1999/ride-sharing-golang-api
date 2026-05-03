package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/dto"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/service"
)

type DriverHandler struct {
	DriverService      service.IDriverService
	TripRequestService service.ITripRequestService
	TripService        service.ITripService
}

func NewDriverHandler(driverService service.IDriverService, tripRequestService service.ITripRequestService, tripService service.ITripService) *DriverHandler {
	return &DriverHandler{DriverService: driverService, TripRequestService: tripRequestService, TripService: tripService}
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

func (h *DriverHandler) AcceptTripRequest(c *gin.Context) {
	driver, ok := c.MustGet("driver").(*domain.Driver)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "driver not found in context"})
		return
	}

	tripRequestID, err := uuid.Parse(c.Param("trip_request_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip_request_id"})
		return
	}

	trip, tr, err := h.TripService.AcceptTripRequest(c.Request.Context(), driver.ID, tripRequestID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTripRequestNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrTripRequestNotOpen):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"trip": trip, "trip_request": tr})
}

func (h *DriverHandler) StartTrip(c *gin.Context) {
	driver, ok := c.MustGet("driver").(*domain.Driver)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "driver not found in context"})
		return
	}

	tripID, err := uuid.Parse(c.Param("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip_id"})
		return
	}

	trip, err := h.TripService.StartTrip(c.Request.Context(), driver.ID, tripID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTripNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrTripWrongDriver):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrTripInvalidState):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrTripStartConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"trip": trip})
}
