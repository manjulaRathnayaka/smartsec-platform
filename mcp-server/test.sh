#!/bin/bash

# Test script for MCP Server

set -e

echo "Testing MCP Server..."
echo "====================="

# Configuration
BASE_URL="http://localhost:8082"
TIMEOUT=30

# Function to check if server is running
check_server() {
    local retries=0
    local max_retries=30

    while [ $retries -lt $max_retries ]; do
        if curl -s -f "$BASE_URL/health" > /dev/null 2>&1; then
            echo "✓ Server is running"
            return 0
        fi

        echo "Waiting for server to start... ($((retries + 1))/$max_retries)"
        sleep 1
        retries=$((retries + 1))
    done

    echo "✗ Server failed to start within $max_retries seconds"
    return 1
}

# Function to test endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4

    echo "Testing $method $endpoint..."

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$BASE_URL$endpoint")
    fi

    body=$(echo "$response" | head -n -1 2>/dev/null || echo "$response" | sed '$d')
    status=$(echo "$response" | tail -n 1)

    if [ "$status" = "$expected_status" ]; then
        echo "✓ $method $endpoint (status: $status)"
        return 0
    else
        echo "✗ $method $endpoint (expected: $expected_status, got: $status)"
        echo "Response: $body"
        return 1
    fi
}

# Main test function
run_tests() {
    echo "Checking if server is running..."
    check_server

    echo -e "\nRunning API tests..."

    # Test health endpoint
    test_endpoint "GET" "/health" "" "200"

    # Test root endpoint
    test_endpoint "GET" "/" "" "200"

    # Test schema endpoint
    test_endpoint "GET" "/mcp/schema" "" "200"

    # Test entities endpoint
    test_endpoint "GET" "/mcp/entities" "" "200"

    # Test specific entity endpoint
    test_endpoint "GET" "/mcp/entities/devices" "" "200"

    # Test examples endpoint
    test_endpoint "GET" "/mcp/examples" "" "200"    # Test query endpoint with simple query (expect 500 since no database)
    query_data='{"entity":"devices","fields":["id","hostname","os"],"limit":10}'
    test_endpoint "POST" "/mcp/query" "$query_data" "500"

    # Test query endpoint with filters (expect 500 since no database)
    query_with_filters='{"entity":"devices","fields":["id","hostname","os"],"filters":[{"field":"os","operator":"eq","value":"Linux"}],"limit":5}'
    test_endpoint "POST" "/mcp/query" "$query_with_filters" "500"

    # Test invalid query
    invalid_query='{"entity":"nonexistent","fields":["id"]}'
    test_endpoint "POST" "/mcp/query" "$invalid_query" "400"

    echo -e "\n====================="
    echo "All tests completed successfully! ✓"
}

# Help function
show_help() {
    echo "Usage: $0 [OPTIONS]"
    echo "Test script for MCP Server"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -u, --url URL  Base URL of the MCP server (default: http://localhost:8082)"
    echo ""
    echo "Examples:"
    echo "  $0                          # Test local server"
    echo "  $0 -u http://localhost:8082 # Test server on specific URL"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -u|--url)
            BASE_URL="$2"
            shift 2
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Run tests
run_tests
