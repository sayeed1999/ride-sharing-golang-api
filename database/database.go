package database

import (
	"fmt"
	"log"

	"github.com/sayeed1999/ride-sharing-golang-api/config"
	authdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/domain"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip"
	tripdomain "github.com/sayeed1999/ride-sharing-golang-api/internal/app/trip/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes and returns a database connection
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Build DSN from config fields
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DB)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// InitDBWithErrorHandling initializes database connection with error handling
func InitDBWithErrorHandling(cfg *config.Config) *gorm.DB {
	db, err := InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db
}

// AutoMigrate runs database migrations for all domain models
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Ensure auth and trip schemas exist before migrating tables.
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS auth").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS trip").Error; err != nil {
		return err
	}

	err := db.AutoMigrate(
		// auth module
		&authdomain.User{},
		&authdomain.Role{},
		&authdomain.UserRole{},
		// trip module
		&tripdomain.Customer{},
		&tripdomain.Driver{},
		&tripdomain.VehicleType{},
	)

	if err != nil {
		return err
	}

	// Seed default roles if they don't exist. Keep idempotent so repeated
	// AutoMigrate calls are safe.
	if err := db.Exec("INSERT INTO auth.roles (name) SELECT $1 WHERE NOT EXISTS (SELECT 1 FROM auth.roles WHERE name = $1)", "customer").Error; err != nil {
		return err
	}
	if err := db.Exec("INSERT INTO auth.roles (name) SELECT $1 WHERE NOT EXISTS (SELECT 1 FROM auth.roles WHERE name = $1)", "driver").Error; err != nil {
		return err
	}

	// Seed vehicle types in trip schema
	vehicleTypes := []struct {
		Name     string
		EnumCode int
	}{
		{"bike", int(trip.VehicleEnumBike)},
		{"cng", int(trip.VehicleEnumCNG)},
		{"car", int(trip.VehicleEnumCar)},
	}
	for _, vt := range vehicleTypes {
		if err := db.Exec("INSERT INTO trip.vehicle_types (name, enum_code) SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM trip.vehicle_types WHERE name = $1)", vt.Name, vt.EnumCode).Error; err != nil {
			return err
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// AutoMigrateWithErrorHandling runs migrations with error handling
func AutoMigrateWithErrorHandling(db *gorm.DB) {
	if err := AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

// CloseDB closes the database connection
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// CloseDBWithErrorHandling closes database connection with error handling
func CloseDBWithErrorHandling(db *gorm.DB) {
	if err := CloseDB(db); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}
