package database

import (
	"fmt"
	"log"

	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/app/auth/domain"

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

	err := db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.UserRole{},
	)

	if err != nil {
		return err
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
