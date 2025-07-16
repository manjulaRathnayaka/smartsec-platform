#!/bin/bash

# SmartSec Platform Integration Test Script
# This script tests the basic functionality of the web portal

echo "🚀 Starting SmartSec Platform Integration Test..."
echo "================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if service is running
check_service() {
    local service_name="$1"
    local url="$2"
    local expected_status="$3"

    echo -e "${YELLOW}Checking ${service_name}...${NC}"

    response=$(curl -s -o /dev/null -w "%{http_code}" "$url")

    if [ "$response" -eq "$expected_status" ]; then
        echo -e "${GREEN}✅ ${service_name} is running (HTTP ${response})${NC}"
        return 0
    else
        echo -e "${RED}❌ ${service_name} is not responding correctly (HTTP ${response})${NC}"
        return 1
    fi
}

# Function to check if process is running
check_process() {
    local process_name="$1"
    local process_pattern="$2"

    echo -e "${YELLOW}Checking ${process_name}...${NC}"

    if pgrep -f "$process_pattern" > /dev/null; then
        echo -e "${GREEN}✅ ${process_name} is running${NC}"
        return 0
    else
        echo -e "${RED}❌ ${process_name} is not running${NC}"
        return 1
    fi
}

# Test results
tests_passed=0
tests_failed=0

echo -e "\n📋 Testing Services..."
echo "----------------------"

# Test BFF Server
if check_service "BFF Server" "http://localhost:3001/health" 200; then
    ((tests_passed++))
else
    ((tests_failed++))
fi

# Test React Frontend
if check_service "React Frontend" "http://localhost:3000" 200; then
    ((tests_passed++))
else
    ((tests_failed++))
fi

echo -e "\n🔍 Testing Processes..."
echo "----------------------"

# Test BFF Process
if check_process "BFF Process" "node server.js"; then
    ((tests_passed++))
else
    ((tests_failed++))
fi

# Test React Process
if check_process "React Process" "react-scripts start"; then
    ((tests_passed++))
else
    ((tests_failed++))
fi

echo -e "\n📊 Testing API Endpoints..."
echo "----------------------------"

# Test BFF Health endpoint
echo -e "${YELLOW}Testing BFF health endpoint...${NC}"
health_response=$(curl -s "http://localhost:3001/health")
if echo "$health_response" | grep -q "healthy"; then
    echo -e "${GREEN}✅ BFF health endpoint working${NC}"
    ((tests_passed++))
else
    echo -e "${RED}❌ BFF health endpoint not working${NC}"
    ((tests_failed++))
fi

# Test if frontend serves static files
echo -e "${YELLOW}Testing frontend static files...${NC}"
if curl -s "http://localhost:3000/static/css/" > /dev/null || curl -s "http://localhost:3000/static/js/" > /dev/null || curl -s "http://localhost:3000/manifest.json" > /dev/null; then
    echo -e "${GREEN}✅ Frontend static files accessible${NC}"
    ((tests_passed++))
else
    echo -e "${RED}❌ Frontend static files not accessible${NC}"
    ((tests_failed++))
fi

echo -e "\n📁 Testing Project Structure..."
echo "-------------------------------"

# Check key directories
directories=(
    "web-portal/frontend/src"
    "web-portal/frontend/public"
    "web-portal/bff"
    "web-portal/bff/routes"
    "web-portal/bff/middleware"
)

for dir in "${directories[@]}"; do
    if [ -d "$dir" ]; then
        echo -e "${GREEN}✅ Directory exists: $dir${NC}"
        ((tests_passed++))
    else
        echo -e "${RED}❌ Directory missing: $dir${NC}"
        ((tests_failed++))
    fi
done

echo -e "\n📋 Summary"
echo "=========="
echo -e "Tests passed: ${GREEN}$tests_passed${NC}"
echo -e "Tests failed: ${RED}$tests_failed${NC}"
echo -e "Total tests: $((tests_passed + tests_failed))"

if [ $tests_failed -eq 0 ]; then
    echo -e "\n🎉 ${GREEN}All tests passed! The SmartSec Platform is running correctly.${NC}"
    exit 0
else
    echo -e "\n⚠️  ${RED}Some tests failed. Please check the output above.${NC}"
    exit 1
fi
