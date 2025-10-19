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

type RedisConfig struct {
	URL string
}

type Config struct {
	Server ServerConfig
	Redis  RedisConfig
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
		Redis: RedisConfig{
			URL: getEnv("Redis__URL", "redis://0.0.0.0:6379"),
		},
	}
}
