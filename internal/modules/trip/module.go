package trip

import (
	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/config"

	"gorm.io/gorm"
)

// ExposeRoutes registers the trip module's HTTP routes to a given router group.
// It delegates the creation of HTTP handlers and their registration to the http handler package.
func ExposeRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	registerAllHTTPRoutes(rg, db, cfg)
}
