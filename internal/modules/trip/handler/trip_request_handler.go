package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/service"
)

type TripRequestRequest struct {
	CustomerID  string `json:"customer_id" binding:"required"`
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

type TripRequestHandler struct {
	RequestTripUC        *service.RequestTripService
	CustomerCancelTripUC *service.CustomerCancelTrip
	CustomerRepo         repository.CustomerRepository
}

func NewTripRequestHandler(requestTripUC *service.RequestTripService, customerCancelTripUC *service.CustomerCancelTrip, customerRepo repository.CustomerRepository) *TripRequestHandler {
	return &TripRequestHandler{RequestTripUC: requestTripUC, CustomerCancelTripUC: customerCancelTripUC, CustomerRepo: customerRepo}
}

func (h *TripRequestHandler) RequestTrip(c *gin.Context) {
	customer, _ := c.MustGet("customer").(*domain.Customer) // assumed to be set by customer middleware

	var req TripRequestRequest
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

	tripRequest, err := h.RequestTripUC.Execute(customerUUID, req.Origin, req.Destination)
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

	if err := h.CustomerCancelTripUC.Execute(c.Request.Context(), tripID, customer.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
