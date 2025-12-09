package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository"
)

func TripRequestMiddleware(tripRequestRepo repository.TripRequestRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Step 1: Extract trip ID from URL params & validate it
		tripIdStr := c.Param("tripID")

		tripId, err := uuid.Parse(tripIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip ID"})
			c.Abort()
			return
		}

		// Step 2: Check a valid trip request exists with given ID
		tripRequest, err := tripRequestRepo.FindByID(tripId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "trip request not found"})
			c.Abort()
			return
		}

		// Step 3: Get authenticated customer from context (set by customer middleware)
		// Problem here: it makes this middleware dependent on customer middleware !! (TODO) R&D needed to improve this!
		customer, ok := c.MustGet("customer").(*domain.Customer)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "customer not found in context"})
			c.Abort()
			return
		}

		// Step 4: Ensure the trip request belongs to the authenticated customer
		if tripRequest.CustomerID != customer.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "trip request does not belong to customer"})
			c.Abort()
			return
		}

		c.Next()
	}
}
