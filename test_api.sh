#!/bin/bash

# Base URL for the API
BASE_URL="http://localhost:3000"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Starting API tests..."

# Test 1: Create a new task
echo -e "\n${GREEN}Test 1: Create a new task${NC}"
response=$(curl -s -X POST "$BASE_URL/tasks" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Task",
    "description": "This is a test task",
    "status": "pending"
  }')

echo "Response: $response"
task_id=$(echo $response | jq -r '.id')
echo "Created task ID: $task_id"

# Test 2: Get all tasks
echo -e "\n${GREEN}Test 2: Get all tasks${NC}"
curl -s -X GET "$BASE_URL/tasks" | jq '.'

# Test 3: Get specific task
echo -e "\n${GREEN}Test 3: Get specific task${NC}"
curl -s -X GET "$BASE_URL/tasks/$task_id" | jq '.'

# Test 4: Update task
echo -e "\n${GREEN}Test 4: Update task${NC}"
curl -s -X PUT "$BASE_URL/tasks/$task_id" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Test Task",
    "description": "This is an updated test task",
    "status": "completed"
  }' | jq '.'

# Test 5: Verify update
echo -e "\n${GREEN}Test 5: Verify update${NC}"
curl -s -X GET "$BASE_URL/tasks/$task_id" | jq '.'

# Test 6: Delete task
echo -e "\n${GREEN}Test 6: Delete task${NC}"
curl -s -X DELETE "$BASE_URL/tasks/$task_id"

# Test 7: Verify deletion
echo -e "\n${GREEN}Test 7: Verify deletion${NC}"
curl -s -X GET "$BASE_URL/tasks/$task_id" | jq '.'

# Test 8: Health check
echo -e "\n${GREEN}Test 8: Health check${NC}"
curl -s -X GET "$BASE_URL/health" | jq '.'

echo -e "\n${GREEN}All tests completed!${NC}" 