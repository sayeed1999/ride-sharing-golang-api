package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/dto"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/service"
)

type LocationHandler struct {
	LocationService service.ILocationService
	CustomerRepo    repository.ICustomerRepository
	DriverRepo      repository.IDriverRepository
}

func NewLocationHandler(
	locationService service.ILocationService,
	customerRepo repository.ICustomerRepository,
	driverRepo repository.IDriverRepository,
) *LocationHandler {
	return &LocationHandler{
		LocationService: locationService,
		CustomerRepo:    customerRepo,
		DriverRepo:      driverRepo,
	}
}

func (h *LocationHandler) UpdateDriverLocation(c *gin.Context) {
	driver, ok := c.MustGet("driver").(*domain.Driver)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "driver not found in context"})
		return
	}

	var req dto.UpdateDriverLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	loc := &domain.DriverLocation{
		Lat:     req.Lat,
		Lng:     req.Lng,
		Heading: req.Heading,
		Speed:   req.Speed,
	}

	if err := h.LocationService.UpdateDriverLocation(c.Request.Context(), driver.ID, loc); err != nil {
		writeLocationError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *LocationHandler) GetDriverLocation(c *gin.Context) {
	trip, ok := c.MustGet("trip").(*domain.Trip)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trip not found in context"})
		return
	}

	userEmail, ok := c.MustGet("x-user-email").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user email not found in context"})
		return
	}

	var actorCustomerID, actorDriverID *uuid.UUID

	if driver, err := h.DriverRepo.FindByEmail(userEmail); err == nil {
		actorDriverID = &driver.ID
	} else if customer, err := h.CustomerRepo.FindByEmail(userEmail); err == nil {
		actorCustomerID = &customer.ID
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "actor not found"})
		return
	}

	loc, err := h.LocationService.GetDriverLocation(c.Request.Context(), trip, actorCustomerID, actorDriverID)
	if err != nil {
		writeLocationError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"driver_location": loc})
}

func writeLocationError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidCoordinates), errors.Is(err, service.ErrNoActiveTrip):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrTripActorForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrLocationNotAvailable), errors.Is(err, service.ErrTripNotTrackable):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
