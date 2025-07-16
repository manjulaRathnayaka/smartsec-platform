package main

import (
	"net/http"
	"os"
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
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database
	db, err := database.Initialize(cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer db.Close()

	// Run database migrations
	if err := database.RunMigrations(cfg.Database.URL); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Initialize repositories
	deviceRepo := repository.NewDeviceRepository(db)
	processRepo := repository.NewProcessRepository(db)
	containerRepo := repository.NewContainerRepository(db)
	threatRepo := repository.NewThreatRepository(db)

	// Initialize services
	telemetryService := service.NewTelemetryService(deviceRepo, processRepo, containerRepo, threatRepo)

	// Initialize HTTP server
	router := setupRouter()
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
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func setupRouter() *gin.Engine {
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
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
		})
	})

	return router
}
