package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type FeatureFlags struct {
	RequireRoleOnRegistration bool
}

type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	Redis        RedisConfig
	FeatureFlags FeatureFlags
}

var AppConfig *Config

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return &Config{
		Server: ServerConfig{
			Host: getEnv("Server__Host", "0.0.0.0"),
			Port: getEnv("Server__Port", "8080"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/ride_sharing_db?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("Redis__URL", "redis://0.0.0.0:6379"),
		},
		FeatureFlags: FeatureFlags{
			RequireRoleOnRegistration: getEnv("REQUIRE_ROLE_ON_REGISTRATION", "true") == "true",
		},
	}
}
