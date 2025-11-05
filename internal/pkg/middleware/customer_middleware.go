package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository"
)

func CustomerMiddleware(customerRepo repository.CustomerRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail, ok := c.MustGet("x-user-email").(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user email not found in context"})
			c.Abort()
			return
		}

		customer, err := customerRepo.FindByEmail(userEmail)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "customer not found"})
			c.Abort()
			return
		}

		c.Set("customer", customer)
		c.Next()
	}
}
