package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/usecase"
)

type TripRequestRequest struct {
	CustomerID  string `json:"customer_id" binding:"required"`
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

type TripRequestHandler struct {
	RequestTripUC        *usecase.RequestTripUsecase
	CustomerCancelTripUC *usecase.CustomerCancelTrip
	CustomerRepo         repository.CustomerRepository
}

func NewTripRequestHandler(requestTripUC *usecase.RequestTripUsecase, customerCancelTripUC *usecase.CustomerCancelTrip, customerRepo repository.CustomerRepository) *TripRequestHandler {
	return &TripRequestHandler{RequestTripUC: requestTripUC, CustomerCancelTripUC: customerCancelTripUC, CustomerRepo: customerRepo}
}

func (h *TripRequestHandler) RequestTrip(c *gin.Context) {
	var req TripRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

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
	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip ID"})
		return
	}

	userEmail, ok := c.MustGet("x-user-email").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user email not found in context"})
		return
	}

	customer, err := h.CustomerRepo.FindByEmail(userEmail)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "customer not found"})
		return
	}

	if err := h.CustomerCancelTripUC.Execute(c.Request.Context(), tripID, customer.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
