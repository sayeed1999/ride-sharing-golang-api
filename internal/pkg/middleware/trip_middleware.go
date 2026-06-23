package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
)

func TripMiddleware(tripRepo repository.ITripRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Get the trip ID from the URL
		tripIDStr := c.Param("trip_id")

		tripID, err := uuid.Parse(tripIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip_id"})
			c.Abort()
			return
		}

		// Step 2: Check if the trip exists
		trip, err := tripRepo.FindByID(tripID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "trip not found"})
			c.Abort()
			return
		}

		// Step 3: Check if the user is a driver,
		// if yes, check if the trip belongs to the driver!
		if driverVal, ok := c.Get("driver"); ok {
			driver, ok := driverVal.(*domain.Driver)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "driver not found in context"})
				c.Abort()
				return
			}
			if trip.DriverID != driver.ID {
				c.JSON(http.StatusForbidden, gin.H{"error": "trip does not belong to driver"})
				c.Abort()
				return
			}
		}

		// Step 4: Check if the user is a customer,
		// if yes, check if the trip belongs to the customer!
		if customerVal, ok := c.Get("customer"); ok {
			customer, ok := customerVal.(*domain.Customer)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "customer not found in context"})
				c.Abort()
				return
			}
			if trip.CustomerID != customer.ID {
				c.JSON(http.StatusForbidden, gin.H{"error": "trip does not belong to customer"})
				c.Abort()
				return
			}
		}

		c.Set("trip", trip)
		c.Next()
	}
}
