package handler

import (
	"net/http"

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

	var req dto.TripRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	req.CustomerID = customer.ID.String()

	customerUUID, err := uuid.Parse(req.CustomerID)
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
	tripIDStr := c.Param("tripID")
	tripID, _ := uuid.Parse(tripIDStr) // assumed to be validated by trip request middleware!!

	customer, _ := c.MustGet("customer").(*domain.Customer) // assumed to be set by middleware

	if err := h.TripRequestService.CancelTripRequest(c.Request.Context(), tripID, customer.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
