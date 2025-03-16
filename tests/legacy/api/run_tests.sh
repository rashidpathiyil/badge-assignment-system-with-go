#!/bin/bash

# Script to run the API integration tests

# Set default API URL if not provided
export API_TEST_URL=${API_TEST_URL:-http://localhost:8080}

# Set log file path
LOG_FILE="./server_test_logs.txt"
CURL_LOG_FILE="./curl_debug.txt"

# Print test configuration
echo "Running API tests against: $API_TEST_URL"
echo "Server logs will be saved to: $LOG_FILE"

# Check if the API is running
if ! curl -s "$API_TEST_URL/health" > /dev/null; then
  echo "Error: API server is not running at $API_TEST_URL"
  echo "Please start the server first or specify a different URL with API_TEST_URL env variable"
  exit 1
fi

# Capture server status and configuration
echo "=== SERVER INFORMATION ===" > $LOG_FILE
echo "Timestamp: $(date)" >> $LOG_FILE
echo "API URL: $API_TEST_URL" >> $LOG_FILE
echo "" >> $LOG_FILE

# Fetch badges before test
echo "=== EXISTING BADGES BEFORE TEST ===" >> $LOG_FILE
curl -s "$API_TEST_URL/api/v1/badges" | jq . >> $LOG_FILE 2>&1
echo "" >> $LOG_FILE

# Fetch user badges before test
echo "=== USER BADGES BEFORE TEST ===" >> $LOG_FILE
curl -s "$API_TEST_URL/api/v1/users/test-user-12345/badges" | jq . >> $LOG_FILE 2>&1
echo "" >> $LOG_FILE

# Ask for server logs if possible
echo "=== SERVER LOGS BEFORE TEST ===" >> $LOG_FILE
curl -s "$API_TEST_URL/admin/logs" >> $LOG_FILE 2>&1 || echo "No admin log endpoint available" >> $LOG_FILE
echo "" >> $LOG_FILE

# Run the tests with curl debugging for critical requests
echo "=== RUNNING TESTS ===" >> $LOG_FILE

# Define a function to run tests
run_tests() {
  cd "$(dirname "$0")/../.." # Move to project root
  go test -v ./tests/api
  return $?
}

# Run the tests and capture output
run_tests | tee -a $LOG_FILE

# Save the exit code
TEST_EXIT_CODE=${PIPESTATUS[0]}

# After tests, capture event processing details with verbose curl
echo "" >> $LOG_FILE
echo "=== EVENT PROCESSING DEBUG INFO ===" >> $LOG_FILE

# Get detailed information about the most recent event
echo "Capturing detailed event information..." >> $LOG_FILE
curl -v "$API_TEST_URL/api/v1/users/test-user-12345/badges" > $CURL_LOG_FILE 2>&1
cat $CURL_LOG_FILE >> $LOG_FILE
echo "" >> $LOG_FILE

# Fetch badges after test
echo "=== EXISTING BADGES AFTER TEST ===" >> $LOG_FILE
curl -s "$API_TEST_URL/api/v1/badges" | jq . >> $LOG_FILE 2>&1
echo "" >> $LOG_FILE

# Try to get event processing status if available
echo "=== EVENT PROCESSING STATUS ===" >> $LOG_FILE
curl -s "$API_TEST_URL/api/v1/events/latest" >> $LOG_FILE 2>&1 || echo "No event status endpoint available" >> $LOG_FILE
echo "" >> $LOG_FILE

# Ask for server logs again
echo "=== SERVER LOGS AFTER TEST ===" >> $LOG_FILE
curl -s "$API_TEST_URL/admin/logs" >> $LOG_FILE 2>&1 || echo "No admin log endpoint available" >> $LOG_FILE
echo "" >> $LOG_FILE

echo "Test logs and server information saved to $LOG_FILE"

# Check if tests passed
if [ $TEST_EXIT_CODE -eq 0 ]; then
  echo -e "\n✅ API tests passed successfully"
else
  echo -e "\n❌ API tests failed"
  exit 1
fi 
