package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"telemetry-service/internal/models"
	"telemetry-service/internal/service"
)

var validate = validator.New()

type TelemetryHandler struct {
	service *service.TelemetryService
}

func NewTelemetryHandler(service *service.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{service: service}
}

func SetupRoutes(router *gin.Engine, service *service.TelemetryService) {
	handler := NewTelemetryHandler(service)

	api := router.Group("/api")
	{
		telemetry := api.Group("/telemetry")
		{
			telemetry.POST("", handler.PostTelemetry)
			telemetry.GET("", handler.GetTelemetry)
		}

		devices := api.Group("/devices")
		{
			devices.GET("", handler.GetDevices)
			devices.GET("/:id", handler.GetDevice)
			devices.GET("/:id/processes", handler.GetProcesses)
			devices.GET("/:id/containers", handler.GetContainers)
			devices.GET("/:id/threats", handler.GetThreatFindings)
		}

		threats := api.Group("/threats")
		{
			threats.GET("", handler.GetThreatsBySeverity)
			threats.POST("", handler.CreateThreatFinding)
		}
	}
}

func (h *TelemetryHandler) PostTelemetry(c *gin.Context) {
	var req models.TelemetryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind telemetry request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := validate.Struct(&req); err != nil {
		log.Error().Err(err).Msg("Failed to validate telemetry request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	err := h.service.ProcessTelemetry(&req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process telemetry")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process telemetry"})
		return
	}

	log.Info().
		Str("mac_address", req.MacAddress).
		Str("hostname", req.HostMetadata.Hostname).
		Int("processes", len(req.Processes)).
		Int("containers", len(req.Containers)).
		Msg("Telemetry processed successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Telemetry processed successfully",
		"stats": gin.H{
			"processes":  len(req.Processes),
			"containers": len(req.Containers),
		},
	})
}

func (h *TelemetryHandler) GetTelemetry(c *gin.Context) {
	deviceID := c.Query("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id parameter is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	dataType := c.DefaultQuery("type", "all")

	response := gin.H{}

	switch dataType {
	case "processes":
		processes, err := h.service.GetProcesses(deviceID, limit, offset)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get processes")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get processes"})
			return
		}
		response["processes"] = processes

	case "containers":
		containers, err := h.service.GetContainers(deviceID, limit, offset)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get containers")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get containers"})
			return
		}
		response["containers"] = containers

	case "threats":
		threats, err := h.service.GetThreatFindings(deviceID, limit, offset)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get threat findings")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get threat findings"})
			return
		}
		response["threats"] = threats

	default: // "all"
		processes, err := h.service.GetProcesses(deviceID, limit, offset)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get processes")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get processes"})
			return
		}
		response["processes"] = processes

		containers, err := h.service.GetContainers(deviceID, limit, offset)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get containers")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get containers"})
			return
		}
		response["containers"] = containers

		threats, err := h.service.GetThreatFindings(deviceID, limit, offset)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get threat findings")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get threat findings"})
			return
		}
		response["threats"] = threats
	}

	c.JSON(http.StatusOK, response)
}

func (h *TelemetryHandler) GetDevices(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	devices, err := h.service.GetDevices(limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get devices")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get devices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
		"count":   len(devices),
	})
}

func (h *TelemetryHandler) GetDevice(c *gin.Context) {
	deviceID := c.Param("id")

	device, err := h.service.GetDevice(deviceID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get device")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get device"})
		return
	}

	if device == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}

func (h *TelemetryHandler) GetProcesses(c *gin.Context) {
	deviceID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	processes, err := h.service.GetProcesses(deviceID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get processes")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get processes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"processes": processes,
		"count":     len(processes),
	})
}

func (h *TelemetryHandler) GetContainers(c *gin.Context) {
	deviceID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	containers, err := h.service.GetContainers(deviceID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get containers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get containers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"containers": containers,
		"count":      len(containers),
	})
}

func (h *TelemetryHandler) GetThreatFindings(c *gin.Context) {
	deviceID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	threats, err := h.service.GetThreatFindings(deviceID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get threat findings")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get threat findings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"threats": threats,
		"count":   len(threats),
	})
}

func (h *TelemetryHandler) GetThreatsBySeverity(c *gin.Context) {
	severity := c.Query("severity")
	if severity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "severity parameter is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	threats, err := h.service.GetThreatFindingsBySeverity(severity, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get threats by severity")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get threats by severity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"threats": threats,
		"count":   len(threats),
	})
}

func (h *TelemetryHandler) CreateThreatFinding(c *gin.Context) {
	var threat models.ThreatFinding

	if err := c.ShouldBindJSON(&threat); err != nil {
		log.Error().Err(err).Msg("Failed to bind threat finding request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := validate.Struct(&threat); err != nil {
		log.Error().Err(err).Msg("Failed to validate threat finding request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	err := h.service.CreateThreatFinding(&threat)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create threat finding")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create threat finding"})
		return
	}

	log.Info().
		Str("device_id", threat.DeviceID).
		Str("severity", threat.Severity).
		Str("rule_id", threat.RuleID).
		Msg("Threat finding created successfully")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Threat finding created successfully",
		"id":      threat.ID,
	})
}
