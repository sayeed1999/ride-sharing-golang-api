package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/usecase"
)

type TripRequestRequest struct {
	CustomerID  string `json:"customer_id" binding:"required"`
	Origin      string `json:"origin" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

type TripRequestHandler struct {
	RequestTripUC *usecase.RequestTripUsecase
}

func NewTripRequestHandler(requestTripUC *usecase.RequestTripUsecase) *TripRequestHandler {
	return &TripRequestHandler{RequestTripUC: requestTripUC}
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
