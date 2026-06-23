package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/service"
)

type TripHandler struct {
	TripService    service.ITripService
	CustomerRepo   repository.ICustomerRepository
	DriverRepo     repository.IDriverRepository
}

func NewTripHandler(
	tripService service.ITripService,
	customerRepo repository.ICustomerRepository,
	driverRepo repository.IDriverRepository,
) *TripHandler {
	return &TripHandler{
		TripService:  tripService,
		CustomerRepo: customerRepo,
		DriverRepo:   driverRepo,
	}
}

func (h *TripHandler) StartTrip(c *gin.Context) {
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
		writeTripError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"trip": trip})
}

func (h *TripHandler) CompleteTrip(c *gin.Context) {
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

	trip, err := h.TripService.CompleteTrip(c.Request.Context(), driver.ID, tripID)
	if err != nil {
		writeTripError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"trip": trip})
}

func (h *TripHandler) CancelTrip(c *gin.Context) {
	tripID, err := uuid.Parse(c.Param("trip_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip_id"})
		return
	}

	userEmail, ok := c.MustGet("x-user-email").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user email not found in context"})
		return
	}

	if driver, err := h.DriverRepo.FindByEmail(userEmail); err == nil {
		trip, err := h.TripService.CancelTripByDriver(c.Request.Context(), driver.ID, tripID)
		if err != nil {
			writeTripError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"trip": trip})
		return
	}

	customer, err := h.CustomerRepo.FindByEmail(userEmail)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "actor not found"})
		return
	}

	trip, err := h.TripService.CancelTripByCustomer(c.Request.Context(), customer.ID, tripID)
	if err != nil {
		writeTripError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"trip": trip})
}

func writeTripError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrTripNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrTripWrongDriver), errors.Is(err, service.ErrTripNotOwnedByCustomer):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrTripInvalidState):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrTripStartConflict),
		errors.Is(err, service.ErrTripCompleteConflict),
		errors.Is(err, service.ErrTripCancelConflict):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
