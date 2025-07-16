package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL      string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnvOrDefault("PORT", "8080"),
		},
		Database: DatabaseConfig{
			// Use Choreo-defined environment variables
			Host:     getEnvOrDefault("CHOREO_TELEMETRYDB_HOSTNAME", "localhost"),
			Port:     getEnvOrDefault("CHOREO_TELEMETRYDB_PORT", "5432"),
			Name:     getEnvOrDefault("CHOREO_TELEMETRYDB_DATABASENAME", "smartsec_db"),
			User:     getEnvOrDefault("CHOREO_TELEMETRYDB_USERNAME", "smartsec_user"),
			Password: getEnvOrDefault("CHOREO_TELEMETRYDB_PASSWORD", ""),
			SSLMode:  getEnvOrDefault("DB_SSL_MODE", "require"),
		},
	}

	// Build DATABASE_URL from Choreo components or use direct URL
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		config.Database.URL = databaseURL
	} else {
		config.Database.URL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
			config.Database.SSLMode,
		)
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
