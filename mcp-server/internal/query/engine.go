package query

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"mcp-server/internal/schema"
)

// QueryEngine handles structured query execution
type QueryEngine struct {
	db     *sqlx.DB
	schema *schema.MCPSchema
	cache  map[string]*schema.QueryResponse
}

// NewQueryEngine creates a new query engine
func NewQueryEngine(db *sqlx.DB) *QueryEngine {
	return &QueryEngine{
		db:     db,
		schema: schema.BuildSchema(),
		cache:  make(map[string]*schema.QueryResponse),
	}
}

// ExecuteQuery executes a structured query
func (qe *QueryEngine) ExecuteQuery(req *schema.QueryRequest) (*schema.QueryResponse, error) {
	// Generate unique ID if not provided
	if req.ID == "" {
		req.ID = uuid.New().String()
	}

	// Set creation time
	req.CreatedAt = time.Now()

	// Initialize response
	response := &schema.QueryResponse{
		ID:        req.ID,
		Status:    "running",
		CreatedAt: req.CreatedAt,
		Metadata: schema.QueryMetadata{
			RowCount: 0,
		},
	}

	// Store in cache
	qe.cache[req.ID] = response

	// Check if database is available
	if qe.db == nil {
		response.Status = "failed"
		response.Error = "Database connection not available - running in schema-only mode"
		completedAt := time.Now()
		response.CompletedAt = &completedAt
		return response, fmt.Errorf("database not available")
	}

	// Validate entity exists
	entity, exists := qe.schema.Entities[req.Entity]
	if !exists {
		response.Status = "failed"
		response.Error = fmt.Sprintf("Entity '%s' not found", req.Entity)
		completedAt := time.Now()
		response.CompletedAt = &completedAt
		return response, fmt.Errorf("entity not found: %s", req.Entity)
	}

	// Build SQL query
	sqlQuery, err := qe.buildSQL(req, entity)
	if err != nil {
		response.Status = "failed"
		response.Error = err.Error()
		completedAt := time.Now()
		response.CompletedAt = &completedAt
		return response, err
	}

	response.Metadata.SQL = sqlQuery

	// Execute query
	start := time.Now()
	rows, err := qe.db.Query(sqlQuery)
	if err != nil {
		response.Status = "failed"
		response.Error = err.Error()
		completedAt := time.Now()
		response.CompletedAt = &completedAt
		response.Metadata.ExecutionTime = time.Since(start)
		return response, err
	}
	defer rows.Close()

	// Process results
	data, err := qe.processRows(rows)
	if err != nil {
		response.Status = "failed"
		response.Error = err.Error()
		completedAt := time.Now()
		response.CompletedAt = &completedAt
		response.Metadata.ExecutionTime = time.Since(start)
		return response, err
	}

	// Update response
	response.Status = "completed"
	response.Data = data
	response.Metadata.RowCount = len(data)
	response.Metadata.ExecutionTime = time.Since(start)
	completedAt := time.Now()
	response.CompletedAt = &completedAt

	return response, nil
}

// GetQueryResult retrieves a cached query result
func (qe *QueryEngine) GetQueryResult(queryID string) (*schema.QueryResponse, error) {
	response, exists := qe.cache[queryID]
	if !exists {
		return nil, fmt.Errorf("query not found: %s", queryID)
	}
	return response, nil
}

// buildSQL constructs SQL query from structured request
func (qe *QueryEngine) buildSQL(req *schema.QueryRequest, entity schema.Entity) (string, error) {
	var query strings.Builder

	// SELECT clause
	query.WriteString("SELECT ")
	if len(req.Fields) > 0 {
		fields := make([]string, 0, len(req.Fields))
		for _, field := range req.Fields {
			if _, exists := entity.Fields[field]; exists {
				fields = append(fields, fmt.Sprintf("%s.%s", req.Entity, field))
			} else {
				return "", fmt.Errorf("field '%s' not found in entity '%s'", field, req.Entity)
			}
		}
		query.WriteString(strings.Join(fields, ", "))
	} else {
		// Select all fields
		fields := make([]string, 0, len(entity.Fields))
		for fieldName := range entity.Fields {
			fields = append(fields, fmt.Sprintf("%s.%s", req.Entity, fieldName))
		}
		query.WriteString(strings.Join(fields, ", "))
	}

	// Handle aggregates
	if len(req.Aggregates) > 0 {
		for _, agg := range req.Aggregates {
			alias := agg.Alias
			if alias == "" {
				alias = fmt.Sprintf("%s_%s", agg.Function, agg.Field)
			}
			query.WriteString(fmt.Sprintf(", %s(%s.%s) AS %s",
				strings.ToUpper(agg.Function), req.Entity, agg.Field, alias))
		}
	}

	// FROM clause
	query.WriteString(fmt.Sprintf(" FROM %s", req.Entity))

	// JOIN clauses
	for _, join := range req.Joins {
		joinType := strings.ToUpper(join.Type)
		query.WriteString(fmt.Sprintf(" %s JOIN %s ON %s", joinType, join.Entity, join.Condition))
	}

	// WHERE clause
	if len(req.Filters) > 0 {
		query.WriteString(" WHERE ")
		whereConditions := make([]string, 0, len(req.Filters))
		for i, filter := range req.Filters {
			condition, err := qe.buildFilterCondition(filter, req.Entity)
			if err != nil {
				return "", err
			}

			if i > 0 {
				logic := filter.Logic
				if logic == "" {
					logic = "AND"
				}
				whereConditions = append(whereConditions, fmt.Sprintf(" %s %s", strings.ToUpper(logic), condition))
			} else {
				whereConditions = append(whereConditions, condition)
			}
		}
		query.WriteString(strings.Join(whereConditions, ""))
	}

	// GROUP BY clause
	if len(req.GroupBy) > 0 {
		query.WriteString(fmt.Sprintf(" GROUP BY %s", strings.Join(req.GroupBy, ", ")))
	}

	// ORDER BY clause
	if len(req.OrderBy) > 0 {
		query.WriteString(" ORDER BY ")
		orderItems := make([]string, 0, len(req.OrderBy))
		for _, order := range req.OrderBy {
			direction := strings.ToUpper(order.Direction)
			if direction != "ASC" && direction != "DESC" {
				direction = "ASC"
			}
			orderItems = append(orderItems, fmt.Sprintf("%s.%s %s", req.Entity, order.Field, direction))
		}
		query.WriteString(strings.Join(orderItems, ", "))
	}

	// LIMIT clause
	if req.Limit != nil {
		query.WriteString(fmt.Sprintf(" LIMIT %d", *req.Limit))
	}

	// OFFSET clause
	if req.Offset != nil {
		query.WriteString(fmt.Sprintf(" OFFSET %d", *req.Offset))
	}

	return query.String(), nil
}

// buildFilterCondition constructs WHERE condition from filter
func (qe *QueryEngine) buildFilterCondition(filter schema.Filter, entity string) (string, error) {
	field := fmt.Sprintf("%s.%s", entity, filter.Field)

	switch filter.Operator {
	case "eq":
		return fmt.Sprintf("%s = %s", field, qe.formatValue(filter.Value)), nil
	case "ne":
		return fmt.Sprintf("%s != %s", field, qe.formatValue(filter.Value)), nil
	case "gt":
		return fmt.Sprintf("%s > %s", field, qe.formatValue(filter.Value)), nil
	case "gte":
		return fmt.Sprintf("%s >= %s", field, qe.formatValue(filter.Value)), nil
	case "lt":
		return fmt.Sprintf("%s < %s", field, qe.formatValue(filter.Value)), nil
	case "lte":
		return fmt.Sprintf("%s <= %s", field, qe.formatValue(filter.Value)), nil
	case "like":
		return fmt.Sprintf("%s LIKE %s", field, qe.formatValue(filter.Value)), nil
	case "ilike":
		return fmt.Sprintf("%s ILIKE %s", field, qe.formatValue(filter.Value)), nil
	case "in":
		if values, ok := filter.Value.([]interface{}); ok {
			formattedValues := make([]string, len(values))
			for i, v := range values {
				formattedValues[i] = qe.formatValue(v)
			}
			return fmt.Sprintf("%s IN (%s)", field, strings.Join(formattedValues, ", ")), nil
		}
		return "", fmt.Errorf("invalid value for IN operator")
	case "is_null":
		return fmt.Sprintf("%s IS NULL", field), nil
	case "is_not_null":
		return fmt.Sprintf("%s IS NOT NULL", field), nil
	default:
		return "", fmt.Errorf("unsupported operator: %s", filter.Operator)
	}
}

// formatValue formats a value for SQL query
func (qe *QueryEngine) formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format(time.RFC3339))
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// processRows processes SQL result rows into JSON-compatible data
func (qe *QueryEngine) processRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}

	for rows.Next() {
		// Create slice to hold column values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan row into values
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Convert to map
		rowMap := make(map[string]interface{})
		for i, column := range columns {
			val := values[i]

			// Handle different types
			switch v := val.(type) {
			case []byte:
				// Try to parse as JSON first
				var jsonVal interface{}
				if err := json.Unmarshal(v, &jsonVal); err == nil {
					rowMap[column] = jsonVal
				} else {
					rowMap[column] = string(v)
				}
			case time.Time:
				rowMap[column] = v.Format(time.RFC3339)
			case nil:
				rowMap[column] = nil
			default:
				rowMap[column] = v
			}
		}

		result = append(result, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// ValidateQuery validates a query request
func (qe *QueryEngine) ValidateQuery(req *schema.QueryRequest) error {
	// Check if entity exists
	entity, exists := qe.schema.Entities[req.Entity]
	if !exists {
		return fmt.Errorf("entity '%s' not found", req.Entity)
	}

	// Validate fields
	for _, field := range req.Fields {
		if _, exists := entity.Fields[field]; !exists {
			return fmt.Errorf("field '%s' not found in entity '%s'", field, req.Entity)
		}
	}

	// Validate filters
	for _, filter := range req.Filters {
		if _, exists := entity.Fields[filter.Field]; !exists {
			return fmt.Errorf("filter field '%s' not found in entity '%s'", filter.Field, req.Entity)
		}

		// Validate operator
		validOperators := []string{"eq", "ne", "gt", "gte", "lt", "lte", "in", "like", "ilike", "is_null", "is_not_null"}
		validOperator := false
		for _, op := range validOperators {
			if filter.Operator == op {
				validOperator = true
				break
			}
		}
		if !validOperator {
			return fmt.Errorf("invalid filter operator: %s", filter.Operator)
		}
	}

	// Validate joins
	for _, join := range req.Joins {
		if _, exists := qe.schema.Entities[join.Entity]; !exists {
			return fmt.Errorf("join entity '%s' not found", join.Entity)
		}

		// Validate join type
		validJoinTypes := []string{"inner", "left", "right", "full"}
		validJoinType := false
		for _, jt := range validJoinTypes {
			if strings.ToLower(join.Type) == jt {
				validJoinType = true
				break
			}
		}
		if !validJoinType {
			return fmt.Errorf("invalid join type: %s", join.Type)
		}
	}

	// Validate order by
	for _, order := range req.OrderBy {
		if _, exists := entity.Fields[order.Field]; !exists {
			return fmt.Errorf("order by field '%s' not found in entity '%s'", order.Field, req.Entity)
		}

		if order.Direction != "asc" && order.Direction != "desc" && order.Direction != "" {
			return fmt.Errorf("invalid order direction: %s", order.Direction)
		}
	}

	// Validate aggregates
	for _, agg := range req.Aggregates {
		validFunctions := []string{"count", "sum", "avg", "min", "max"}
		validFunction := false
		for _, fn := range validFunctions {
			if strings.ToLower(agg.Function) == fn {
				validFunction = true
				break
			}
		}
		if !validFunction {
			return fmt.Errorf("invalid aggregate function: %s", agg.Function)
		}

		if agg.Field != "*" {
			if _, exists := entity.Fields[agg.Field]; !exists {
				return fmt.Errorf("aggregate field '%s' not found in entity '%s'", agg.Field, req.Entity)
			}
		}
	}

	return nil
}

// GetSchema returns the MCP schema
func (qe *QueryEngine) GetSchema() *schema.MCPSchema {
	return qe.schema
}

// ClearCache removes old query results from cache
func (qe *QueryEngine) ClearCache() {
	// Remove queries older than 1 hour
	cutoff := time.Now().Add(-1 * time.Hour)
	for id, response := range qe.cache {
		if response.CreatedAt.Before(cutoff) {
			delete(qe.cache, id)
		}
	}
}
