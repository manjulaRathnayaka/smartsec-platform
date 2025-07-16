package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            string
	DatabaseURL     string
	TelemetryAPIURL string
	LogLevel        string
}

func Load() *Config {
	config := &Config{
		Port:            getEnv("PORT", "8082"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://telemetry_user:telemetry_password@localhost:5432/telemetry?sslmode=disable"),
		TelemetryAPIURL: getEnv("TELEMETRY_API_URL", "http://localhost:8080"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
