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
	User     string
	Password string
	Host     string
	Port     string
	DB       string
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
			User:     getEnv("POSTGRES_USER", getEnv("DB_USER", "user")),
			Password: getEnv("POSTGRES_PASSWORD", getEnv("DB_PASSWORD", "password")),
			Host:     getEnv("POSTGRES_HOST", getEnv("DB_HOST", "localhost")),
			Port:     getEnv("POSTGRES_PORT", getEnv("DB_PORT", "5432")),
			DB:       getEnv("POSTGRES_DB", getEnv("DB_NAME", "ride_sharing_db")),
		},
		Redis: RedisConfig{
			URL: getEnv("Redis__URL", "redis://0.0.0.0:6379"),
		},
		FeatureFlags: FeatureFlags{
			RequireRoleOnRegistration: getEnv("REQUIRE_ROLE_ON_REGISTRATION", "true") == "true",
		},
	}
}
