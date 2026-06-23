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

type TripRequestHandler struct {
	TripRequestService service.ITripRequestService
}

func NewTripRequestHandler(tripRequestService service.ITripRequestService) *TripRequestHandler {
	return &TripRequestHandler{TripRequestService: tripRequestService}
}

func (h *TripRequestHandler) RequestTrip(c *gin.Context) {
	customer, _ := c.MustGet("customer").(*domain.Customer) // assumed to be set by customer middleware

	var req dto.TripRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	customerUUID, err := uuid.Parse(customer.ID.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	tripRequest, err := h.TripRequestService.RequestTrip(customerUUID, req.Origin, req.Destination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"trip_request": tripRequest})
}

func (h *TripRequestHandler) CancelTripRequest(c *gin.Context) {
	tripRequest, ok := c.MustGet("trip_request").(*domain.TripRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trip request not found in context"})
		return
	}

	if err := h.TripRequestService.CancelTripRequest(c.Request.Context(), tripRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TripRequestHandler) GetDetails(c *gin.Context) {
	tripRequest, ok := c.MustGet("trip_request").(*domain.TripRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trip request not found in context"})
		return
	}

	// no need to call service layer since middleware has extracted the trip_request from db

	c.JSON(200, gin.H{"trip_request": tripRequest})
}

func (h *TripRequestHandler) ListOpenTripRequests(c *gin.Context) {
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

func (h *TripRequestHandler) AcceptTripRequest(c *gin.Context) {
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

	trip, tr, err := h.TripRequestService.AcceptTripRequest(c.Request.Context(), driver.ID, tripRequestID)
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
