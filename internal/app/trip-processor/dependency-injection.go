package tripprocessor

import (
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip-processor/trip"
	trip_request "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip-processor/trip-request"

	"github.com/gin-gonic/gin"
)

func InitEndpoints(r *gin.Engine) {
	rg := r.Group("/api/transition-checker")
	{
		rg.POST("/trip-request-status", trip_request.Handler)
		rg.POST("/trip-status", trip.Handler)
	}
}
