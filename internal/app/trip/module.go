package trip

import (
	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth"
	triphttp "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/handler/http"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/repository/postgres"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/usecase"

	"gorm.io/gorm"
)

// TODO: refactor and make the code clean.

// ExposeRoutes wires the trip module routes to a given router group
// It creates its own auth.RegisterUsecase (backed by the same DB) and uses
// it to call auth registration internally (modular-monolith style).
func ExposeRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// Repositories
	pg := &postgres.CustomerRepo{DB: db}
	var custRepo repository.CustomerRepository = pg

	// Obtain the auth RegisterUsecase via the auth package (exported constructor)
	registerUC := auth.NewRegisterUsecase(db, cfg)

	signupUC := &usecase.CustomerSignupUsecase{
		CustomerRepo: custRepo,
		AuthRegister: registerUC,
	}

	custHandler := triphttp.NewCustomerHandler(signupUC)

	customers := rg.Group("/customers")
	triphttp.RegisterCustomerRoutes(customers, custHandler)

	// --- Driver wiring ---
	driverRepo := &postgres.DriverRepo{DB: db}
	var drvRepo repository.DriverRepository = driverRepo
	driverSignupUC := &usecase.DriverSignupUsecase{
		DriverRepo:   drvRepo,
		AuthRegister: registerUC,
	}
	drvHandler := triphttp.NewDriverHandler(driverSignupUC)
	drivers := rg.Group("/drivers")
	triphttp.RegisterDriverRoutes(drivers, drvHandler)
}
