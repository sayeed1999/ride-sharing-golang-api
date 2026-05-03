package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip/repository"
)

func DriverMiddleware(driverRepo repository.IDriverRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail, ok := c.MustGet("x-user-email").(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user email not found in context"})
			c.Abort()
			return
		}

		driver, err := driverRepo.FindByEmail(userEmail)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "driver not found"})
			c.Abort()
			return
		}

		c.Set("driver", driver)
		c.Next()
	}
}
