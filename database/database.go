package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sayeed1999/ride-sharing-golang-api/config"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes and returns a database connection
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Build DSN from config fields
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DB)
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
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

// RunMigrations runs the database migrations
func RunMigrations(db *sql.DB, cfg *config.Config) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{
		DatabaseName: cfg.Database.DB,
	})
	if err != nil {
		return err
	}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	migrationsPath := filepath.Join(basepath, "migrations")

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// RunMigrationsWithErrorHandling runs migrations with error handling
func RunMigrationsWithErrorHandling(db *sql.DB, cfg *config.Config) {
	if err := RunMigrations(db, cfg); err != nil {
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
