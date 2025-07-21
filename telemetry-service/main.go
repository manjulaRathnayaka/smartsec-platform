package main

import (
	"database/sql"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"telemetry-service/internal/api"
	"telemetry-service/internal/config"
	"telemetry-service/internal/database"
	"telemetry-service/internal/repository"
	"telemetry-service/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Info().Msg("No .env file found")
	}

	// Setup logging
	setupLogging()

	// Load configuration
	log.Info().Msg("Starting SmartSec Telemetry Service")
	log.Info().Msg("Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Log configuration details (without sensitive data)
	log.Info().
		Str("server_port", cfg.Server.Port).
		Str("db_host", cfg.Database.Host).
		Str("db_port", cfg.Database.Port).
		Str("db_name", cfg.Database.Name).
		Str("db_user", cfg.Database.User).
		Str("db_ssl_mode", cfg.Database.SSLMode).
		Msg("Configuration loaded successfully")

	// Initialize database
	log.Info().Msg("Initializing database connection...")
	db, err := database.Initialize(cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).
			Str("database_url", maskDatabaseURL(cfg.Database.URL)).
			Msg("Failed to initialize database")
	}
	defer db.Close()

	log.Info().Msg("Database connection established successfully")

	// Run database migrations
	log.Info().Msg("Running database migrations...")
	if err := database.RunMigrations(cfg.Database.URL); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}
	log.Info().Msg("Database migrations completed successfully")

	// Initialize repositories
	deviceRepo := repository.NewDeviceRepository(db)
	processRepo := repository.NewProcessRepository(db)
	containerRepo := repository.NewContainerRepository(db)
	threatRepo := repository.NewThreatRepository(db)

	// Initialize services
	telemetryService := service.NewTelemetryService(deviceRepo, processRepo, containerRepo, threatRepo)

	// Initialize HTTP server
	router := setupRouter(db)
	api.SetupRoutes(router, telemetryService)

	// Start server
	log.Info().Str("port", cfg.Server.Port).Msg("Starting telemetry service")
	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func setupLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	// Set log level
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		// Default to debug level for better troubleshooting
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// maskDatabaseURL masks sensitive information in database URL for logging
func maskDatabaseURL(url string) string {
	if url == "" {
		return ""
	}

	// Find password in URL and mask it
	parts := strings.Split(url, "@")
	if len(parts) < 2 {
		return url // No password found
	}

	userInfo := parts[0]
	if strings.Contains(userInfo, ":") {
		userParts := strings.Split(userInfo, ":")
		if len(userParts) >= 3 {
			// postgres://user:password@host... -> postgres://user:***@host...
			userParts[2] = "***"
			userInfo = strings.Join(userParts, ":")
		}
	}

	return userInfo + "@" + parts[1]
}

func setupRouter(db *sql.DB) *gin.Engine {
	// Set gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		healthStatus := gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"service":   "telemetry-service",
		}

		// Test database connectivity
		if db != nil {
			if err := db.Ping(); err != nil {
				log.Error().Err(err).Msg("Database health check failed")
				healthStatus["database"] = gin.H{
					"status": "unhealthy",
					"error":  err.Error(),
				}
				healthStatus["status"] = "unhealthy"
				c.JSON(http.StatusServiceUnavailable, healthStatus)
				return
			} else {
				healthStatus["database"] = gin.H{
					"status": "healthy",
				}
			}
		}

		c.JSON(http.StatusOK, healthStatus)
	})

	return router
}
