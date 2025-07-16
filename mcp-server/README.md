# SmartSec MCP Server

A Model Context Protocol (MCP) server that provides a structured interface for querying SmartSec platform telemetry data. This server enables LLMs to compose valid structured queries over telemetry data by exposing a comprehensive schema and query API.

## Features

- **Schema Discovery**: Exposes complete data model schema including entities, fields, relations, and operations
- **Structured Queries**: Supports complex queries with filters, joins, aggregations, and sorting
- **Query Management**: Async query execution with result caching and retrieval
- **Entity Information**: Detailed information about available data entities
- **Example Queries**: Pre-built query examples for common use cases

## API Endpoints

### Schema and Discovery

- `GET /mcp/schema` - Get the complete MCP schema
- `GET /mcp/entities` - List all available entities
- `GET /mcp/entities/:name` - Get specific entity details
- `GET /mcp/examples` - Get query examples

### Query Execution

- `POST /mcp/query` - Execute a structured query
- `GET /mcp/query/:id/result` - Get query result by ID

### Health and Status

- `GET /health` - Health check endpoint
- `GET /` - API documentation

## Data Entities

The MCP server exposes the following entities:

### Devices
- Device metadata, hardware information, and status
- Relations: processes, containers, threat_findings, browser_sessions

### Processes
- Process information from endpoint devices
- Relations: devices (parent), threat_findings

### Containers
- Container runtime information
- Relations: devices (parent), threat_findings

### Threat Findings
- Security alerts and threat detections
- Relations: devices, processes, containers

### Browser Sessions
- Browser session and activity data
- Relations: devices

## Query Structure

```json
{
  "id": "query-123",
  "entity": "devices",
  "fields": ["id", "hostname", "os", "last_seen_at"],
  "filters": [
    {
      "field": "os",
      "operator": "eq",
      "value": "Linux"
    }
  ],
  "joins": [
    {
      "entity": "processes",
      "type": "inner",
      "condition": "devices.id = processes.device_id"
    }
  ],
  "order_by": [
    {
      "field": "last_seen_at",
      "direction": "desc"
    }
  ],
  "limit": 100
}
```

## Supported Operations

### Filters
- `eq`, `ne` - Equality/inequality
- `gt`, `gte`, `lt`, `lte` - Comparison
- `like`, `ilike` - Pattern matching
- `in` - List membership
- `is_null`, `is_not_null` - Null checks

### Joins
- `inner`, `left`, `right`, `full` - Join types

### Aggregations
- `count`, `sum`, `avg`, `min`, `max` - Aggregate functions

## Configuration

Environment variables:

- `PORT` - Server port (default: 8082)
- `DATABASE_URL` - PostgreSQL connection string
- `TELEMETRY_API_URL` - Telemetry service URL (default: http://localhost:8081)
- `LOG_LEVEL` - Logging level (default: info)

## Running the Server

### Prerequisites
- Go 1.22+
- PostgreSQL database with telemetry data
- Access to telemetry service

### Build and Run

```bash
# Build the server
go build -o mcp-server

# Run with environment variables
export DATABASE_URL="postgres://user:password@localhost:5432/smartsec?sslmode=disable"
export PORT=8082
./mcp-server
```

### Docker

```bash
# Build Docker image
docker build -t smartsec-mcp-server .

# Run container
docker run -p 8082:8082 \
  -e DATABASE_URL="postgres://user:password@host:5432/smartsec?sslmode=disable" \
  smartsec-mcp-server
```

## Example Usage

### Get Schema
```bash
curl -X GET http://localhost:8082/mcp/schema
```

### Query Devices
```bash
curl -X POST http://localhost:8082/mcp/query \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "devices",
    "fields": ["id", "hostname", "os"],
    "filters": [
      {"field": "os", "operator": "eq", "value": "Linux"}
    ],
    "limit": 10
  }'
```

### Get Threat Findings
```bash
curl -X POST http://localhost:8082/mcp/query \
  -H "Content-Type: application/json" \
  -d '{
    "entity": "threat_findings",
    "filters": [
      {"field": "severity", "operator": "eq", "value": "critical"}
    ],
    "order_by": [
      {"field": "timestamp", "direction": "desc"}
    ],
    "limit": 20
  }'
```

## Development

### Project Structure
```
mcp-server/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── mcp/               # MCP API handlers
│   ├── query/             # Query engine
│   └── schema/            # Schema definitions
├── go.mod                 # Go module definition
└── README.md             # This file
```

### Testing
```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Model Context Protocol

This server implements the Model Context Protocol (MCP) specification, providing:

1. **Schema Discovery**: Complete data model with entities, fields, and relationships
2. **Structured Queries**: Type-safe query construction with validation
3. **Result Management**: Async query execution with caching
4. **Documentation**: Self-documenting API with examples

The MCP interface enables LLMs to:
- Understand available data structures
- Compose valid queries
- Execute complex analytical queries
- Retrieve structured results

## Security Considerations

- Database connections use connection pooling
- Query validation prevents SQL injection
- CORS headers configured for web access
- Request logging for audit trails
- Input validation on all endpoints

## License

This project is part of the SmartSec platform and follows the same licensing terms.
