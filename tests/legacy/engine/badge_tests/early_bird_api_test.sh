#!/bin/bash

# Early Bird Badge API Test
# This script tests the Early Bird badge functionality via the API
# It submits check-in events and verifies badge assignment

# Configuration
API_URL="http://localhost:8080"
TEST_USER_ID="test-user-api-123"
EVENT_TYPE="Check In"

# Color outputs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Early Bird Badge API Test ===${NC}"
echo -e "${BLUE}API URL: ${API_URL}${NC}"
echo -e "${BLUE}Test User: ${TEST_USER_ID}${NC}"

# Function to check if the API server is running
check_server() {
  echo -e "${YELLOW}Checking if API server is running...${NC}"
  
  server_response=$(curl -s -o /dev/null -w "%{http_code}" "${API_URL}/health")
  
  if [ "$server_response" = "200" ]; then
    echo -e "${GREEN}API server is running.${NC}"
    return 0
  else
    echo -e "${RED}API server is not running. Please start the server using 'go run cmd/server/main.go'${NC}"
    return 1
  fi
}

# Function to submit a check-in event
submit_checkin() {
  local user_id=$1
  local date=$2
  local time=$3
  
  echo -e "${YELLOW}Submitting check-in for ${user_id} on ${date} at ${time}${NC}"
  
  response=$(curl -s -X POST "${API_URL}/api/v1/events" \
    -H "Content-Type: application/json" \
    -d '{
      "event_type": "'"${EVENT_TYPE}"'",
      "user_id": "'"${user_id}"'",
      "payload": {
        "user_id": "'"${user_id}"'",
        "time": "'"${time}"'",
        "date": "'"${date}"'",
        "location": "API Test Location"
      },
      "timestamp": "'"${date}T${time}Z"'"
    }')
  
  echo "$response" | grep -q "Event processed successfully"
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}Check-in event submitted successfully.${NC}"
  else
    echo -e "${RED}Failed to submit check-in event: ${response}${NC}"
    return 1
  fi
}

# Function to check user badges
check_user_badges() {
  local user_id=$1
  
  echo -e "${YELLOW}Checking badges for user ${user_id}...${NC}"
  
  badges=$(curl -s "${API_URL}/api/v1/users/${user_id}/badges")
  
  if echo "$badges" | grep -q "Early Bird"; then
    echo -e "${GREEN}SUCCESS: User has been awarded the Early Bird badge!${NC}"
    return 0
  else
    echo -e "${RED}FAIL: User has not been awarded the Early Bird badge.${NC}"
    echo -e "${YELLOW}API Response: ${badges}${NC}"
    return 1
  fi
}

# Main test function
run_test() {
  # Check if server is running
  check_server || return 1
  
  echo -e "${BLUE}\nTest 1: User checks in before 9 AM for 5 consecutive days${NC}"
  
  # Get today's date and format it
  today=$(date +%Y-%m-%d)
  
  # Create 5 check-ins before 9 AM for 5 consecutive days
  for i in {0..4}; do
    # Calculate date (starting from today, going backward)
    test_date=$(date -v-${i}d +%Y-%m-%d 2>/dev/null || date -d "${today} -${i} days" +%Y-%m-%d)
    
    # Submit check-in at 8:00 AM
    submit_checkin "${TEST_USER_ID}" "${test_date}" "08:00:00" || return 1
  done
  
  # Wait a bit for badge processing
  echo -e "${YELLOW}Waiting for badge processing...${NC}"
  sleep 2
  
  # Check if user has received the badge
  check_user_badges "${TEST_USER_ID}" || return 1
  
  echo -e "${BLUE}\nTest 2: User with mixed check-in times (should not earn badge)${NC}"
  
  # Change test user ID for the second test
  TEST_USER_ID="test-user-api-456"
  
  # Create 3 check-ins before 9 AM
  for i in {0..2}; do
    test_date=$(date -v-${i}d +%Y-%m-%d 2>/dev/null || date -d "${today} -${i} days" +%Y-%m-%d)
    submit_checkin "${TEST_USER_ID}" "${test_date}" "08:00:00" || return 1
  done
  
  # Create 2 check-ins after 9 AM
  for i in {3..4}; do
    test_date=$(date -v-${i}d +%Y-%m-%d 2>/dev/null || date -d "${today} -${i} days" +%Y-%m-%d)
    submit_checkin "${TEST_USER_ID}" "${test_date}" "10:00:00" || return 1
  done
  
  # Wait a bit for badge processing
  echo -e "${YELLOW}Waiting for badge processing...${NC}"
  sleep 2
  
  # User should NOT have the badge in this case
  badges=$(curl -s "${API_URL}/api/v1/users/${TEST_USER_ID}/badges")
  
  if echo "$badges" | grep -q "Early Bird"; then
    echo -e "${RED}FAIL: User was incorrectly awarded the Early Bird badge.${NC}"
    echo -e "${YELLOW}API Response: ${badges}${NC}"
    return 1
  else
    echo -e "${GREEN}SUCCESS: User was correctly NOT awarded the Early Bird badge.${NC}"
  fi
  
  echo -e "${GREEN}\nAll tests completed successfully!${NC}"
  return 0
}

# Run the tests
run_test
exit_code=$?

if [ $exit_code -eq 0 ]; then
  echo -e "${GREEN}\nEarly Bird Badge API Test: SUCCESS${NC}"
else
  echo -e "${RED}\nEarly Bird Badge API Test: FAILED${NC}"
fi

exit $exit_code 
