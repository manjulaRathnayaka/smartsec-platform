package schema

import (
	"encoding/json"
	"time"
)

// MCPSchema represents the overall schema for the MCP interface
type MCPSchema struct {
	Version    string               `json:"version"`
	Timestamp  time.Time            `json:"timestamp"`
	Entities   map[string]Entity    `json:"entities"`
	Relations  map[string]Relation  `json:"relations"`
	Operations map[string]Operation `json:"operations"`
}

// Entity represents a data entity (table) in the telemetry system
type Entity struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Fields      map[string]Field `json:"fields"`
	Indexes     []Index          `json:"indexes"`
	Relations   []EntityRelation `json:"relations"`
}

// Field represents a field in an entity
type Field struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Description string      `json:"description"`
	Format      string      `json:"format,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Example     interface{} `json:"example,omitempty"`
}

// Index represents an index on an entity
type Index struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields"`
	Unique bool     `json:"unique"`
}

// EntityRelation represents a relationship between entities
type EntityRelation struct {
	Type         string `json:"type"` // "one-to-many", "many-to-one", "many-to-many"
	TargetEntity string `json:"target_entity"`
	ForeignKey   string `json:"foreign_key"`
	Description  string `json:"description"`
}

// Relation represents a global relation definition
type Relation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	FromEntity  string `json:"from_entity"`
	ToEntity    string `json:"to_entity"`
	Type        string `json:"type"`
}

// Operation represents an available operation
type Operation struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Method      string               `json:"method"`
	Path        string               `json:"path"`
	Parameters  map[string]Parameter `json:"parameters"`
	Response    ResponseSchema       `json:"response"`
	Examples    []OperationExample   `json:"examples"`
}

// Parameter represents an operation parameter
type Parameter struct {
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Description string      `json:"description"`
	Example     interface{} `json:"example,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
}

// ResponseSchema represents the response structure
type ResponseSchema struct {
	Type        string           `json:"type"`
	Description string           `json:"description"`
	Properties  map[string]Field `json:"properties"`
}

// OperationExample represents an example of an operation
type OperationExample struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Request     interface{} `json:"request"`
	Response    interface{} `json:"response"`
}

// QueryRequest represents a structured query request
type QueryRequest struct {
	ID         string                 `json:"id"`
	Entity     string                 `json:"entity"`
	Fields     []string               `json:"fields,omitempty"`
	Filters    []Filter               `json:"filters,omitempty"`
	Joins      []Join                 `json:"joins,omitempty"`
	OrderBy    []OrderBy              `json:"order_by,omitempty"`
	Limit      *int                   `json:"limit,omitempty"`
	Offset     *int                   `json:"offset,omitempty"`
	Aggregates []Aggregate            `json:"aggregates,omitempty"`
	GroupBy    []string               `json:"group_by,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt  time.Time              `json:"created_at"`
}

// Filter represents a filter condition
type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // "eq", "ne", "gt", "gte", "lt", "lte", "in", "like", "ilike", "is_null", "is_not_null"
	Value    interface{} `json:"value"`
	Logic    string      `json:"logic,omitempty"` // "and", "or"
}

// Join represents a join operation
type Join struct {
	Entity    string `json:"entity"`
	Type      string `json:"type"`      // "inner", "left", "right", "full"
	Condition string `json:"condition"` // e.g., "devices.id = processes.device_id"
}

// OrderBy represents ordering
type OrderBy struct {
	Field     string `json:"field"`
	Direction string `json:"direction"` // "asc", "desc"
}

// Aggregate represents an aggregation
type Aggregate struct {
	Function string `json:"function"` // "count", "sum", "avg", "min", "max"
	Field    string `json:"field"`
	Alias    string `json:"alias,omitempty"`
}

// QueryResponse represents the response to a query
type QueryResponse struct {
	ID          string                   `json:"id"`
	Status      string                   `json:"status"` // "pending", "running", "completed", "failed"
	Data        []map[string]interface{} `json:"data,omitempty"`
	Error       string                   `json:"error,omitempty"`
	Metadata    QueryMetadata            `json:"metadata"`
	CreatedAt   time.Time                `json:"created_at"`
	CompletedAt *time.Time               `json:"completed_at,omitempty"`
}

// QueryMetadata represents metadata about the query execution
type QueryMetadata struct {
	RowCount      int           `json:"row_count"`
	ExecutionTime time.Duration `json:"execution_time"`
	SQL           string        `json:"sql,omitempty"`
	QueryPlan     string        `json:"query_plan,omitempty"`
}

// BuildSchema creates the complete MCP schema
func BuildSchema() *MCPSchema {
	return &MCPSchema{
		Version:    "1.0.0",
		Timestamp:  time.Now(),
		Entities:   buildEntities(),
		Relations:  buildRelations(),
		Operations: buildOperations(),
	}
}

func buildEntities() map[string]Entity {
	entities := make(map[string]Entity)

	// Devices entity
	entities["devices"] = Entity{
		Name:        "devices",
		Description: "Devices in the system with their metadata",
		Fields: map[string]Field{
			"id":           {Name: "id", Type: "string", Required: true, Description: "Unique device identifier", Example: "dev-123"},
			"mac_address":  {Name: "mac_address", Type: "string", Required: true, Description: "MAC address of the device", Example: "00:11:22:33:44:55"},
			"hostname":     {Name: "hostname", Type: "string", Required: true, Description: "Device hostname", Example: "laptop-001"},
			"os":           {Name: "os", Type: "string", Required: true, Description: "Operating system", Example: "Linux"},
			"platform":     {Name: "platform", Type: "string", Required: true, Description: "Platform architecture", Example: "x86_64"},
			"version":      {Name: "version", Type: "string", Required: true, Description: "OS version", Example: "Ubuntu 22.04"},
			"current_user": {Name: "current_user", Type: "string", Required: true, Description: "Current logged-in user", Example: "john.doe"},
			"user_id":      {Name: "user_id", Type: "string", Required: false, Description: "User identifier", Example: "user-456"},
			"org_unit":     {Name: "org_unit", Type: "string", Required: false, Description: "Organizational unit", Example: "Engineering"},
			"created_at":   {Name: "created_at", Type: "datetime", Required: true, Description: "Creation timestamp", Example: "2025-01-15T10:30:00Z"},
			"updated_at":   {Name: "updated_at", Type: "datetime", Required: true, Description: "Last update timestamp", Example: "2025-01-15T11:30:00Z"},
			"last_seen_at": {Name: "last_seen_at", Type: "datetime", Required: true, Description: "Last seen timestamp", Example: "2025-01-15T12:30:00Z"},
		},
		Indexes: []Index{
			{Name: "idx_devices_mac_address", Fields: []string{"mac_address"}, Unique: true},
			{Name: "idx_devices_hostname", Fields: []string{"hostname"}, Unique: false},
			{Name: "idx_devices_last_seen", Fields: []string{"last_seen_at"}, Unique: false},
		},
		Relations: []EntityRelation{
			{Type: "one-to-many", TargetEntity: "processes", ForeignKey: "device_id", Description: "Processes running on this device"},
			{Type: "one-to-many", TargetEntity: "containers", ForeignKey: "device_id", Description: "Containers running on this device"},
			{Type: "one-to-many", TargetEntity: "threat_findings", ForeignKey: "device_id", Description: "Threat findings on this device"},
			{Type: "one-to-many", TargetEntity: "browser_sessions", ForeignKey: "device_id", Description: "Browser sessions on this device"},
		},
	}

	// Processes entity
	entities["processes"] = Entity{
		Name:        "processes",
		Description: "Process information collected from devices",
		Fields: map[string]Field{
			"id":           {Name: "id", Type: "string", Required: true, Description: "Unique process record identifier", Example: "proc-123"},
			"device_id":    {Name: "device_id", Type: "string", Required: true, Description: "Device identifier", Example: "dev-123"},
			"pid":          {Name: "pid", Type: "integer", Required: true, Description: "Process ID", Example: 1234},
			"name":         {Name: "name", Type: "string", Required: true, Description: "Process name", Example: "nginx"},
			"cmdline":      {Name: "cmdline", Type: "array", Required: false, Description: "Command line arguments", Example: []string{"nginx", "-g", "daemon off;"}},
			"username":     {Name: "username", Type: "string", Required: false, Description: "User running the process", Example: "www-data"},
			"exe_path":     {Name: "exe_path", Type: "string", Required: false, Description: "Executable path", Example: "/usr/sbin/nginx"},
			"start_time":   {Name: "start_time", Type: "integer", Required: false, Description: "Process start time (unix timestamp)", Example: 1704441600},
			"status":       {Name: "status", Type: "string", Required: false, Description: "Process status", Example: "running"},
			"sha256":       {Name: "sha256", Type: "string", Required: false, Description: "SHA256 hash of executable", Example: "abc123..."},
			"version":      {Name: "version", Type: "string", Required: false, Description: "Executable version", Example: "1.18.0"},
			"file_size":    {Name: "file_size", Type: "integer", Required: false, Description: "Executable file size in bytes", Example: 1048576},
			"collected_at": {Name: "collected_at", Type: "datetime", Required: true, Description: "Data collection timestamp", Example: "2025-01-15T10:30:00Z"},
			"created_at":   {Name: "created_at", Type: "datetime", Required: true, Description: "Record creation timestamp", Example: "2025-01-15T10:30:00Z"},
		},
		Indexes: []Index{
			{Name: "idx_processes_device_id", Fields: []string{"device_id"}, Unique: false},
			{Name: "idx_processes_pid_device", Fields: []string{"pid", "device_id"}, Unique: false},
			{Name: "idx_processes_collected_at", Fields: []string{"collected_at"}, Unique: false},
		},
		Relations: []EntityRelation{
			{Type: "many-to-one", TargetEntity: "devices", ForeignKey: "device_id", Description: "Device this process is running on"},
			{Type: "one-to-many", TargetEntity: "threat_findings", ForeignKey: "process_id", Description: "Threat findings related to this process"},
		},
	}

	// Containers entity
	entities["containers"] = Entity{
		Name:        "containers",
		Description: "Container information collected from devices",
		Fields: map[string]Field{
			"id":                {Name: "id", Type: "string", Required: true, Description: "Unique container record identifier", Example: "cont-123"},
			"device_id":         {Name: "device_id", Type: "string", Required: true, Description: "Device identifier", Example: "dev-123"},
			"container_id":      {Name: "container_id", Type: "string", Required: true, Description: "Docker container ID", Example: "abc123def456"},
			"image":             {Name: "image", Type: "string", Required: true, Description: "Container image", Example: "nginx:latest"},
			"names":             {Name: "names", Type: "array", Required: false, Description: "Container names", Example: []string{"web-server"}},
			"status":            {Name: "status", Type: "string", Required: false, Description: "Container status", Example: "running"},
			"ports":             {Name: "ports", Type: "array", Required: false, Description: "Exposed ports", Example: []string{"80/tcp", "443/tcp"}},
			"labels":            {Name: "labels", Type: "object", Required: false, Description: "Container labels", Example: map[string]string{"env": "prod"}},
			"container_created": {Name: "container_created", Type: "integer", Required: false, Description: "Container creation time (unix timestamp)", Example: 1704441600},
			"collected_at":      {Name: "collected_at", Type: "datetime", Required: true, Description: "Data collection timestamp", Example: "2025-01-15T10:30:00Z"},
			"created_at":        {Name: "created_at", Type: "datetime", Required: true, Description: "Record creation timestamp", Example: "2025-01-15T10:30:00Z"},
		},
		Indexes: []Index{
			{Name: "idx_containers_device_id", Fields: []string{"device_id"}, Unique: false},
			{Name: "idx_containers_container_id", Fields: []string{"container_id"}, Unique: false},
			{Name: "idx_containers_collected_at", Fields: []string{"collected_at"}, Unique: false},
		},
		Relations: []EntityRelation{
			{Type: "many-to-one", TargetEntity: "devices", ForeignKey: "device_id", Description: "Device this container is running on"},
			{Type: "one-to-many", TargetEntity: "threat_findings", ForeignKey: "container_id", Description: "Threat findings related to this container"},
		},
	}

	// Threat findings entity
	entities["threat_findings"] = Entity{
		Name:        "threat_findings",
		Description: "Security threat findings and alerts",
		Fields: map[string]Field{
			"id":           {Name: "id", Type: "string", Required: true, Description: "Unique threat finding identifier", Example: "threat-123"},
			"device_id":    {Name: "device_id", Type: "string", Required: true, Description: "Device identifier", Example: "dev-123"},
			"process_id":   {Name: "process_id", Type: "string", Required: false, Description: "Related process identifier", Example: "proc-456"},
			"container_id": {Name: "container_id", Type: "string", Required: false, Description: "Related container identifier", Example: "cont-789"},
			"description":  {Name: "description", Type: "string", Required: true, Description: "Threat description", Example: "Suspicious process detected"},
			"severity":     {Name: "severity", Type: "string", Required: true, Description: "Threat severity level", Enum: []string{"low", "medium", "high", "critical"}, Example: "high"},
			"rule_id":      {Name: "rule_id", Type: "string", Required: true, Description: "Detection rule identifier", Example: "rule-001"},
			"rule_name":    {Name: "rule_name", Type: "string", Required: true, Description: "Detection rule name", Example: "Malware Detection"},
			"timestamp":    {Name: "timestamp", Type: "datetime", Required: true, Description: "Threat detection timestamp", Example: "2025-01-15T10:30:00Z"},
			"created_at":   {Name: "created_at", Type: "datetime", Required: true, Description: "Record creation timestamp", Example: "2025-01-15T10:30:00Z"},
		},
		Indexes: []Index{
			{Name: "idx_threats_device_id", Fields: []string{"device_id"}, Unique: false},
			{Name: "idx_threats_severity", Fields: []string{"severity"}, Unique: false},
			{Name: "idx_threats_timestamp", Fields: []string{"timestamp"}, Unique: false},
		},
		Relations: []EntityRelation{
			{Type: "many-to-one", TargetEntity: "devices", ForeignKey: "device_id", Description: "Device where threat was detected"},
			{Type: "many-to-one", TargetEntity: "processes", ForeignKey: "process_id", Description: "Process related to the threat"},
			{Type: "many-to-one", TargetEntity: "containers", ForeignKey: "container_id", Description: "Container related to the threat"},
		},
	}

	// Browser sessions entity
	entities["browser_sessions"] = Entity{
		Name:        "browser_sessions",
		Description: "Browser session information collected from devices",
		Fields: map[string]Field{
			"id":                  {Name: "id", Type: "string", Required: true, Description: "Unique browser session identifier", Example: "session-123"},
			"device_id":           {Name: "device_id", Type: "string", Required: true, Description: "Device identifier", Example: "dev-123"},
			"browser_fingerprint": {Name: "browser_fingerprint", Type: "string", Required: true, Description: "Browser fingerprint", Example: "fp-abc123"},
			"user_agent":          {Name: "user_agent", Type: "string", Required: true, Description: "Browser user agent", Example: "Mozilla/5.0..."},
			"tabs":                {Name: "tabs", Type: "array", Required: false, Description: "Open tabs/URLs", Example: []string{"https://example.com"}},
			"user_id":             {Name: "user_id", Type: "string", Required: false, Description: "User identifier", Example: "user-456"},
			"collected_at":        {Name: "collected_at", Type: "datetime", Required: true, Description: "Data collection timestamp", Example: "2025-01-15T10:30:00Z"},
			"created_at":          {Name: "created_at", Type: "datetime", Required: true, Description: "Record creation timestamp", Example: "2025-01-15T10:30:00Z"},
		},
		Indexes: []Index{
			{Name: "idx_browser_sessions_device_id", Fields: []string{"device_id"}, Unique: false},
			{Name: "idx_browser_sessions_collected_at", Fields: []string{"collected_at"}, Unique: false},
		},
		Relations: []EntityRelation{
			{Type: "many-to-one", TargetEntity: "devices", ForeignKey: "device_id", Description: "Device where browser session was captured"},
		},
	}

	return entities
}

func buildRelations() map[string]Relation {
	return map[string]Relation{
		"device_processes": {
			Name:        "device_processes",
			Description: "Processes running on devices",
			FromEntity:  "devices",
			ToEntity:    "processes",
			Type:        "one-to-many",
		},
		"device_containers": {
			Name:        "device_containers",
			Description: "Containers running on devices",
			FromEntity:  "devices",
			ToEntity:    "containers",
			Type:        "one-to-many",
		},
		"device_threats": {
			Name:        "device_threats",
			Description: "Threat findings on devices",
			FromEntity:  "devices",
			ToEntity:    "threat_findings",
			Type:        "one-to-many",
		},
		"process_threats": {
			Name:        "process_threats",
			Description: "Threat findings related to processes",
			FromEntity:  "processes",
			ToEntity:    "threat_findings",
			Type:        "one-to-many",
		},
		"container_threats": {
			Name:        "container_threats",
			Description: "Threat findings related to containers",
			FromEntity:  "containers",
			ToEntity:    "threat_findings",
			Type:        "one-to-many",
		},
	}
}

func buildOperations() map[string]Operation {
	return map[string]Operation{
		"query_devices": {
			Name:        "query_devices",
			Description: "Query devices with filters, joins, and aggregations",
			Method:      "POST",
			Path:        "/mcp/query",
			Parameters: map[string]Parameter{
				"entity":   {Name: "entity", Type: "string", Required: true, Description: "Must be 'devices'", Example: "devices"},
				"fields":   {Name: "fields", Type: "array", Required: false, Description: "Fields to select", Example: []string{"id", "hostname", "os"}},
				"filters":  {Name: "filters", Type: "array", Required: false, Description: "Filter conditions"},
				"joins":    {Name: "joins", Type: "array", Required: false, Description: "Join operations"},
				"order_by": {Name: "order_by", Type: "array", Required: false, Description: "Ordering specification"},
				"limit":    {Name: "limit", Type: "integer", Required: false, Description: "Maximum number of results", Example: 100},
				"offset":   {Name: "offset", Type: "integer", Required: false, Description: "Number of results to skip", Example: 0},
			},
			Response: ResponseSchema{
				Type:        "object",
				Description: "Query response with data and metadata",
				Properties: map[string]Field{
					"id":     {Name: "id", Type: "string", Description: "Query identifier"},
					"status": {Name: "status", Type: "string", Description: "Query status"},
					"data":   {Name: "data", Type: "array", Description: "Query results"},
				},
			},
			Examples: []OperationExample{
				{
					Name:        "Get all Linux devices",
					Description: "Query devices running Linux OS",
					Request: QueryRequest{
						Entity: "devices",
						Fields: []string{"id", "hostname", "os", "last_seen_at"},
						Filters: []Filter{
							{Field: "os", Operator: "eq", Value: "Linux"},
						},
						OrderBy: []OrderBy{
							{Field: "last_seen_at", Direction: "desc"},
						},
						Limit: &[]int{10}[0],
					},
				},
			},
		},
		"query_processes": {
			Name:        "query_processes",
			Description: "Query processes with filters, joins, and aggregations",
			Method:      "POST",
			Path:        "/mcp/query",
			Parameters: map[string]Parameter{
				"entity": {Name: "entity", Type: "string", Required: true, Description: "Must be 'processes'", Example: "processes"},
			},
			Response: ResponseSchema{
				Type:        "object",
				Description: "Query response with data and metadata",
			},
			Examples: []OperationExample{
				{
					Name:        "Get processes with high CPU usage",
					Description: "Query processes that might be consuming high resources",
					Request: QueryRequest{
						Entity: "processes",
						Fields: []string{"id", "name", "pid", "exe_path", "device_id"},
						Joins: []Join{
							{Entity: "devices", Type: "inner", Condition: "processes.device_id = devices.id"},
						},
						OrderBy: []OrderBy{
							{Field: "collected_at", Direction: "desc"},
						},
						Limit: &[]int{50}[0],
					},
				},
			},
		},
		"query_threats": {
			Name:        "query_threats",
			Description: "Query threat findings with filters and aggregations",
			Method:      "POST",
			Path:        "/mcp/query",
			Parameters: map[string]Parameter{
				"entity": {Name: "entity", Type: "string", Required: true, Description: "Must be 'threat_findings'", Example: "threat_findings"},
			},
			Response: ResponseSchema{
				Type:        "object",
				Description: "Query response with data and metadata",
			},
			Examples: []OperationExample{
				{
					Name:        "Get critical threats",
					Description: "Query critical severity threat findings",
					Request: QueryRequest{
						Entity: "threat_findings",
						Filters: []Filter{
							{Field: "severity", Operator: "eq", Value: "critical"},
						},
						OrderBy: []OrderBy{
							{Field: "timestamp", Direction: "desc"},
						},
						Limit: &[]int{20}[0],
					},
				},
			},
		},
	}
}

// ToJSON converts the schema to JSON string
func (s *MCPSchema) ToJSON() (string, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
