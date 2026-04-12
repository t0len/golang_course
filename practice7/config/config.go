package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort string
	PG       PGConfig
}

type PGConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	return &Config{
		HTTPPort: getEnv("HTTP_PORT", "8090"),
		PG: PGConfig{
			Host:     getEnv("PG_HOST", "localhost"),
			Port:     getEnv("PG_PORT", "5432"),
			User:     getEnv("PG_USER", "postgres"),
			Password: getEnv("PG_PASSWORD", "postgres"),
			DBName:   getEnv("PG_DBNAME", "auth_service"),
			SSLMode:  getEnv("PG_SSLMODE", "disable"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
