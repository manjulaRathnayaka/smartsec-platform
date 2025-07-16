package config

import (
	"fmt"
	"os"
	"strconv"
	
	"github.com/rs/zerolog/log"
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
	log.Debug().Msg("Loading configuration from environment variables")
	
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

	// Log which environment variables were found
	log.Debug().
		Bool("choreo_hostname_set", os.Getenv("CHOREO_TELEMETRYDB_HOSTNAME") != "").
		Bool("choreo_port_set", os.Getenv("CHOREO_TELEMETRYDB_PORT") != "").
		Bool("choreo_db_set", os.Getenv("CHOREO_TELEMETRYDB_DATABASENAME") != "").
		Bool("choreo_user_set", os.Getenv("CHOREO_TELEMETRYDB_USERNAME") != "").
		Bool("choreo_password_set", os.Getenv("CHOREO_TELEMETRYDB_PASSWORD") != "").
		Bool("database_url_set", os.Getenv("DATABASE_URL") != "").
		Msg("Environment variables status")

	// Build DATABASE_URL from Choreo components or use direct URL
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Info().Msg("Using direct DATABASE_URL from environment")
		config.Database.URL = databaseURL
	} else {
		log.Info().Msg("Building DATABASE_URL from Choreo environment variables")
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
