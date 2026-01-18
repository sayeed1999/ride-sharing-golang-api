package trip_request

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
	checker := TripRequestStatusTransitionCheckerMap()

	var request struct {
		FromStatus TripRequestStatus `json:"fromStatus"`
		ToStatus   TripRequestStatus `json:"toStatus"`
	}

	// Bind JSON request body into 'request' struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call transition checker
	valid, err := checker.IsTransitionValid(request.FromStatus, request.ToStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": valid})
}
