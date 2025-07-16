#!/bin/bash

# Test script to verify the agent works
echo "Testing SmartSec Laptop Agent..."

# Set test environment variables for log-only mode
export LOG_ONLY="true"
export COLLECTION_INTERVAL="5"

echo "Building agent..."
go build -o laptop-agent .

echo "Running agent in log-only mode for 15 seconds..."
./laptop-agent &
AGENT_PID=$!
sleep 15
kill $AGENT_PID 2>/dev/null
echo "Agent test completed"

echo "Test finished. Check the output above for telemetry data."
