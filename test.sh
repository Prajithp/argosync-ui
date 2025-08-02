#!/bin/bash

# Test script for Heirloom Deployment API
# This script tests the API endpoints using curl

# Set the API base URL
API_URL="http://localhost:8080"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Check if the server is already running
echo -e "${YELLOW}Checking if the server is already running...${NC}"
if curl -s "$API_URL" > /dev/null; then
  echo -e "${GREEN}Server is already running${NC}"
else
  echo -e "${YELLOW}Starting the server with SQLite database...${NC}"
  # Start the server in the background
  go run main.go --db-type=sqlite --sqlite-path=test.db --init-db &
  SERVER_PID=$!
  
  # Wait for the server to start
  echo -e "${YELLOW}Waiting for the server to start...${NC}"
  sleep 3
  
  # Check if the server started successfully
  if ! curl -s "$API_URL" > /dev/null; then
    echo -e "${RED}Failed to start the server${NC}"
    exit 1
  fi
  
  echo -e "${GREEN}Server started successfully${NC}"
  
  # Register a trap to kill the server when the script exits
  trap "kill $SERVER_PID" EXIT
fi

# Function to test an endpoint
test_endpoint() {
  local method=$1
  local endpoint=$2
  local payload=$3
  local description=$4

  echo -e "\n${GREEN}Testing: $description${NC}"
  echo "Method: $method"
  echo "Endpoint: $endpoint"
  
  if [ -n "$payload" ]; then
    echo "Payload: $payload"
    response=$(curl -s -X $method -H "Content-Type: application/json" -d "$payload" "$API_URL$endpoint")
  else
    response=$(curl -s -X $method "$API_URL$endpoint")
  fi

  echo "Response: $response"
  
  if [ -z "$response" ]; then
    echo -e "${RED}Error: No response received${NC}"
  fi
}

# Test the root endpoint
test_endpoint "GET" "/" "" "Root endpoint"

# Test the release endpoint
test_endpoint "POST" "/release" '{
  "application": "auth-service",
  "environment": "production",
  "region": "us-east-1",
  "version": "v2.4.1",
  "deployed_by": "test@example.com"
}' "Release endpoint"

# Test the rollback endpoint
test_endpoint "POST" "/rollback" '{
  "application": "auth-service",
  "environment": "production",
  "region": "us-east-1",
  "deployed_by": "test@example.com"
}' "Rollback endpoint"

# Test the deployments endpoint
test_endpoint "GET" "/deployments?application=auth-service" "" "Deployments endpoint"

# Test the history endpoint
test_endpoint "GET" "/history?application=auth-service&environment=production&region=us-east-1" "" "History endpoint"

# Test the all-deployments endpoint for frontend compatibility
test_endpoint "GET" "/all-deployments" "" "All deployments endpoint for frontend"

echo -e "\n${GREEN}All tests completed${NC}"