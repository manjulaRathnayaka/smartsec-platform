package mcp

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"

	"mcp-server/internal/query"
	"mcp-server/internal/schema"
)

var validate = validator.New()

// Handler handles MCP API requests
type Handler struct {
	queryEngine *query.QueryEngine
}

// NewHandler creates a new MCP handler
func NewHandler(queryEngine *query.QueryEngine) *Handler {
	return &Handler{
		queryEngine: queryEngine,
	}
}

// GetSchema returns the MCP schema
func (h *Handler) GetSchema(c *gin.Context) {
	schema := h.queryEngine.GetSchema()

	// Update timestamp
	schema.Timestamp = time.Now()

	c.JSON(http.StatusOK, gin.H{
		"schema": schema,
		"meta": gin.H{
			"version":   "1.0.0",
			"timestamp": time.Now().Format(time.RFC3339),
			"generated": "MCP Server for SmartSec Platform",
		},
	})
}

// PostQuery handles structured query requests
func (h *Handler) PostQuery(c *gin.Context) {
	var req schema.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to bind query request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate the query
	if err := h.queryEngine.ValidateQuery(&req); err != nil {
		log.Error().Err(err).Msg("Query validation failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Query validation failed",
			"details": err.Error(),
		})
		return
	}

	// Execute the query
	response, err := h.queryEngine.ExecuteQuery(&req)
	if err != nil {
		log.Error().Err(err).Str("query_id", req.ID).Msg("Query execution failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "Query execution failed",
			"query_id": req.ID,
			"details":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetQueryResult returns the result of a specific query
func (h *Handler) GetQueryResult(c *gin.Context) {
	queryID := c.Param("id")
	if queryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query ID is required",
		})
		return
	}

	response, err := h.queryEngine.GetQueryResult(queryID)
	if err != nil {
		log.Error().Err(err).Str("query_id", queryID).Msg("Failed to get query result")
		c.JSON(http.StatusNotFound, gin.H{
			"error":    "Query not found",
			"query_id": queryID,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetQueryExamples returns example queries for each entity
func (h *Handler) GetQueryExamples(c *gin.Context) {
	entity := c.Query("entity")

	schema := h.queryEngine.GetSchema()

	if entity != "" {
		// Return examples for specific entity
		examples := make([]interface{}, 0)

		for _, operation := range schema.Operations {
			for _, example := range operation.Examples {
				// For now, skip type checking of examples
				examples = append(examples, example)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"entity":   entity,
			"examples": examples,
		})
		return
	}

	// Return all examples organized by entity
	examplesByEntity := make(map[string][]interface{})

	for _, operation := range schema.Operations {
		for _, example := range operation.Examples {
			// For now, add all examples without type checking
			entityName := "unknown"
			examplesByEntity[entityName] = append(examplesByEntity[entityName], example)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"examples": examplesByEntity,
	})
}

// GetEntities returns information about available entities
func (h *Handler) GetEntities(c *gin.Context) {
	schema := h.queryEngine.GetSchema()

	entities := make(map[string]interface{})

	for name, entity := range schema.Entities {
		entities[name] = gin.H{
			"name":        entity.Name,
			"description": entity.Description,
			"fields":      entity.Fields,
			"indexes":     entity.Indexes,
			"relations":   entity.Relations,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"entities": entities,
		"total":    len(entities),
	})
}

// GetEntity returns information about a specific entity
func (h *Handler) GetEntity(c *gin.Context) {
	entityName := c.Param("name")
	if entityName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Entity name is required",
		})
		return
	}

	schema := h.queryEngine.GetSchema()

	entity, exists := schema.Entities[entityName]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "Entity not found",
			"entity": entityName,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"entity": gin.H{
			"name":        entity.Name,
			"description": entity.Description,
			"fields":      entity.Fields,
			"indexes":     entity.Indexes,
			"relations":   entity.Relations,
		},
	})
}

// GetHealth returns health status of the MCP server
func (h *Handler) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
		"service":   "MCP Server",
	})
}

// SetupRoutes configures the API routes
func SetupRoutes(router *gin.Engine, handler *Handler) {
	// MCP API routes
	mcp := router.Group("/mcp")
	{
		mcp.GET("/schema", handler.GetSchema)
		mcp.POST("/query", handler.PostQuery)
		mcp.GET("/query/:id/result", handler.GetQueryResult)
		mcp.GET("/examples", handler.GetQueryExamples)
		mcp.GET("/entities", handler.GetEntities)
		mcp.GET("/entities/:name", handler.GetEntity)
		mcp.GET("/health", handler.GetHealth)
	}

	// Root health check
	router.GET("/health", handler.GetHealth)

	// API documentation endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":     "SmartSec MCP Server",
			"version":     "1.0.0",
			"description": "Model Context Protocol interface for SmartSec telemetry data",
			"endpoints": gin.H{
				"GET /mcp/schema":           "Get the complete MCP schema",
				"POST /mcp/query":           "Execute a structured query",
				"GET /mcp/query/:id/result": "Get query result by ID",
				"GET /mcp/examples":         "Get query examples",
				"GET /mcp/entities":         "Get available entities",
				"GET /mcp/entities/:name":   "Get specific entity information",
				"GET /mcp/health":           "Health check",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})
}
