# Telemetry Service

A Go-based telemetry service that accepts telemetry data from laptop agents and stores it in PostgreSQL.

## Features

- **Telemetry Ingestion**: Accepts POST requests with device, process, and container data
- **REST API**: Exposes endpoints for querying telemetry data
- **PostgreSQL Storage**: Normalized database schema for efficient storage
- **Data Cleanup**: Automatic cleanup of old telemetry data
- **Structured Logging**: Using zerolog for structured logging
- **Database Migrations**: Automated database schema management

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd smartsec-platform/telemetry-service
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up PostgreSQL database:
```bash
createdb telemetry
```

4. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your database credentials
```

5. Build the service:
```bash
go build -o telemetry-service
```

## Configuration

The service uses environment variables for configuration:

- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: PostgreSQL connection string
- `LOG_LEVEL`: Log level (debug, info, warn, error)

Example `.env` file:
```
PORT=8080
DATABASE_URL=postgres://username:password@localhost:5432/telemetry?sslmode=disable
LOG_LEVEL=info
```

## Running the Service

1. Start the service:
```bash
./telemetry-service
```

2. The service will:
   - Run database migrations automatically
   - Start HTTP server on the configured port
   - Begin accepting telemetry data

## API Endpoints

### Telemetry Ingestion

- **POST** `/api/telemetry` - Accept telemetry data from agents

Request body:
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "mac_address": "aa:bb:cc:dd:ee:ff",
  "host_metadata": {
    "hostname": "laptop-001",
    "os": "darwin",
    "platform": "darwin",
    "version": "14.0.0",
    "current_user": "john",
    "uptime": 3600
  },
  "processes": [
    {
      "pid": 1234,
      "name": "chrome",
      "cmdline": ["/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"],
      "username": "john",
      "exe_path": "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
      "start_time": 1642234800,
      "status": "running",
      "sha256": "abc123...",
      "version": "120.0.0",
      "file_size": 1024000
    }
  ],
  "containers": [
    {
      "id": "container123",
      "image": "nginx:latest",
      "names": ["/nginx-container"],
      "status": "running",
      "ports": ["80:8080"],
      "labels": {"app": "web"},
      "created": 1642234800
    }
  ]
}
```

### Query Endpoints

- **GET** `/api/telemetry?device_id=<uuid>&type=<type>&limit=<n>&offset=<n>` - Query telemetry data
- **GET** `/api/devices` - List all devices
- **GET** `/api/devices/:id` - Get device details
- **GET** `/api/devices/:id/processes` - Get processes for a device
- **GET** `/api/devices/:id/containers` - Get containers for a device
- **GET** `/api/devices/:id/threats` - Get threat findings for a device
- **GET** `/api/threats` - List threat findings
- **POST** `/api/threats` - Create threat finding

### Health Check

- **GET** `/health` - Health check endpoint

## Database Schema

The service uses the following tables:

- **devices**: Device information (hostname, OS, platform, etc.)
- **processes**: Process information collected from devices
- **containers**: Container information collected from devices
- **browser_sessions**: Browser session data (future feature)
- **threat_findings**: Security threat findings

## Usage with Laptop Agent

The laptop agent should be configured to send telemetry data to:
```
http://localhost:8080/api/telemetry
```

Set the `API_ENDPOINT` environment variable in the laptop agent:
```bash
export API_ENDPOINT="http://localhost:8080/api/telemetry"
```

## Development

### Running Tests

```bash
go test ./...
```

### Database Migrations

Migrations are automatically run on startup. To manually run migrations:

```bash
# The service handles migrations automatically
# Migration files are in the migrations/ directory
```

### Logging

The service uses structured logging with zerolog. Log levels can be configured via the `LOG_LEVEL` environment variable.

## Troubleshooting

1. **Database Connection Issues**: Ensure PostgreSQL is running and the `DATABASE_URL` is correct
2. **Migration Failures**: Check that the database user has sufficient permissions
3. **Port Conflicts**: Ensure the configured port is not already in use

## License

MIT License
