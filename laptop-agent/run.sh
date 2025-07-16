#!/bin/bash

# SmartSec Laptop Agent Launcher
# This script sets up environment variables and runs the laptop agent

echo "Starting SmartSec Laptop Agent..."

# Set default configuration
export API_ENDPOINT="${API_ENDPOINT:-https://api.smartsec.local/telemetry}"
export COLLECTION_INTERVAL="${COLLECTION_INTERVAL:-60}"
export LOG_ONLY="${LOG_ONLY:-false}"

# Show current configuration
if [ "$LOG_ONLY" = "true" ]; then
    echo "Running in LOG_ONLY mode - telemetry will be logged instead of sent to API"
else
    echo "Running in API mode - telemetry will be sent to: $API_ENDPOINT"
fi

# Create logs directory if it doesn't exist
mkdir -p logs

# Run the agent with logging
./laptop-agent 2>&1 | tee logs/agent-$(date +%Y%m%d-%H%M%S).log
