#!/bin/bash

# Test script for telemetry service
echo "Testing Telemetry Service..."

# Check if service is running
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ Telemetry service is not running at http://localhost:8080"
    echo "Please start the service first: ./telemetry-service"
    exit 1
fi

echo "✅ Telemetry service is running"

# Test health endpoint
echo "Testing health endpoint..."
curl -s http://localhost:8080/health | jq '.'

# Test telemetry submission
echo -e "\nTesting telemetry submission..."
curl -X POST http://localhost:8080/api/telemetry \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2024-01-15T10:30:00Z",
    "mac_address": "aa:bb:cc:dd:ee:ff",
    "host_metadata": {
      "hostname": "test-laptop",
      "os": "darwin",
      "platform": "darwin",
      "version": "14.0.0",
      "current_user": "testuser",
      "uptime": 3600
    },
    "processes": [
      {
        "pid": 1234,
        "name": "chrome",
        "cmdline": ["/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"],
        "username": "testuser",
        "exe_path": "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
        "start_time": 1642234800,
        "status": "running",
        "sha256": "abc123def456ghi789jkl012mno345pqr678stu901vwx234yz567890abcdef12",
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
  }' | jq '.'

echo -e "\n✅ Test completed!"
echo "You can now run the laptop agent to send real telemetry data."
