#!/bin/bash

# Script to debug badge assignment with extensive logging

# Set default API URL if not provided
export API_TEST_URL=${API_TEST_URL:-http://localhost:8080}
export RUN_DEBUG_TESTS=true

# Set log file paths
BADGE_LOG_DIR="./badge_debug_logs"
SERVER_LOG_FILE="$BADGE_LOG_DIR/server_logs.txt"
TEST_LOG_FILE="$BADGE_LOG_DIR/test_run.txt"
API_RESPONSES_DIR="$BADGE_LOG_DIR/api_responses"

# Create log directories
mkdir -p "$BADGE_LOG_DIR"
mkdir -p "$API_RESPONSES_DIR"

# Print debug configuration
echo "Running badge assignment debug against: $API_TEST_URL"
echo "Debug logs will be saved to: $BADGE_LOG_DIR"

# Check if the API is running
if ! curl -s "$API_TEST_URL/health" > /dev/null; then
  echo "Error: API server is not running at $API_TEST_URL"
  echo "Please start the server first or specify a different URL with API_TEST_URL env variable"
  exit 1
fi

echo "=== BADGE ASSIGNMENT DEBUG SESSION ===" > "$SERVER_LOG_FILE"
echo "Timestamp: $(date)" >> "$SERVER_LOG_FILE"
echo "API URL: $API_TEST_URL" >> "$SERVER_LOG_FILE"
echo "" >> "$SERVER_LOG_FILE"

# Collect environment and server information
echo "Collecting server information..."
echo "=== SERVER ENVIRONMENT ===" >> "$SERVER_LOG_FILE"
curl -s "$API_TEST_URL/health" > "$API_RESPONSES_DIR/health.json"
cat "$API_RESPONSES_DIR/health.json" >> "$SERVER_LOG_FILE"
echo "" >> "$SERVER_LOG_FILE"

# Collect existing badges
echo "Collecting existing badges..."
echo "=== EXISTING BADGES ===" >> "$SERVER_LOG_FILE"
curl -s "$API_TEST_URL/api/v1/badges" > "$API_RESPONSES_DIR/badges.json"
cat "$API_RESPONSES_DIR/badges.json" >> "$SERVER_LOG_FILE"
echo "" >> "$SERVER_LOG_FILE"

# Collect existing event types
echo "Collecting event types..."
echo "=== EVENT TYPES ===" >> "$SERVER_LOG_FILE"
curl -s "$API_TEST_URL/api/v1/admin/event-types" > "$API_RESPONSES_DIR/event_types.json"
cat "$API_RESPONSES_DIR/event_types.json" >> "$SERVER_LOG_FILE"
echo "" >> "$SERVER_LOG_FILE"

# Collect existing condition types
echo "Collecting condition types..."
echo "=== CONDITION TYPES ===" >> "$SERVER_LOG_FILE"
curl -s "$API_TEST_URL/api/v1/admin/condition-types" > "$API_RESPONSES_DIR/condition_types.json"
cat "$API_RESPONSES_DIR/condition_types.json" >> "$SERVER_LOG_FILE"
echo "" >> "$SERVER_LOG_FILE"

# Look for server logs endpoint (might not exist)
echo "Checking for server logs..."
echo "=== SERVER LOGS BEFORE TEST ===" >> "$SERVER_LOG_FILE"
curl -s "$API_TEST_URL/admin/logs" > "$API_RESPONSES_DIR/server_logs_before.txt"
cat "$API_RESPONSES_DIR/server_logs_before.txt" >> "$SERVER_LOG_FILE" || echo "No server logs endpoint available" >> "$SERVER_LOG_FILE"
echo "" >> "$SERVER_LOG_FILE"

# Run the debug test
echo "Running badge assignment debug test..."
echo "=== DEBUG TEST OUTPUT ===" >> "$SERVER_LOG_FILE"
cd "$(dirname "$0")/../.." # Move to project root

# Run only the debug test
echo "Running badge debug test..."
go test -v ./tests/api -run TestDebugBadgeAssignment | tee "$TEST_LOG_FILE"
TEST_EXIT_CODE=${PIPESTATUS[0]}

# Merge the test output into the server log
cat "$TEST_LOG_FILE" >> "$SERVER_LOG_FILE"
echo "" >> "$SERVER_LOG_FILE"

# Check for badge_debug.log file created by the test
if [ -f "badge_debug.log" ]; then
  echo "Found detailed debug log"
  cp "badge_debug.log" "$BADGE_LOG_DIR/"
  echo "=== DETAILED DEBUG LOG ===" >> "$SERVER_LOG_FILE"
  cat "badge_debug.log" >> "$SERVER_LOG_FILE"
  echo "" >> "$SERVER_LOG_FILE"
fi

# Collect logs after test
echo "Collecting final data..."

# Check server logs again
echo "=== SERVER LOGS AFTER TEST ===" >> "$SERVER_LOG_FILE"
curl -s "$API_TEST_URL/admin/logs" > "$API_RESPONSES_DIR/server_logs_after.txt"
cat "$API_RESPONSES_DIR/server_logs_after.txt" >> "$SERVER_LOG_FILE" || echo "No server logs endpoint available" >> "$SERVER_LOG_FILE"
echo "" >> "$SERVER_LOG_FILE"

# Summarize results
echo "=== DEBUG SUMMARY ===" >> "$SERVER_LOG_FILE"
if grep -q "Success! User has been awarded badges" "$TEST_LOG_FILE"; then
  echo "✅ Badge assignment succeeded" >> "$SERVER_LOG_FILE"
else
  echo "❌ Badge assignment failed" >> "$SERVER_LOG_FILE"
  
  # Collect any errors from the test log
  echo "Error details:" >> "$SERVER_LOG_FILE"
  grep -i "error\|fail\|fatal\|panic" "$TEST_LOG_FILE" >> "$SERVER_LOG_FILE" || echo "No specific errors found in test log" >> "$SERVER_LOG_FILE"
fi

echo ""
echo "Debug session completed!"
echo "Check $BADGE_LOG_DIR for detailed logs"

# Make the log files more readable
echo "Formatting JSON files for better readability..."
for file in "$API_RESPONSES_DIR"/*.json; do
  if [ -f "$file" ]; then
    # Create a temporary file with formatted JSON
    jq . "$file" > "$file.formatted" 2>/dev/null || true
    # If formatting succeeded, replace the original
    if [ -s "$file.formatted" ]; then
      mv "$file.formatted" "$file"
    else
      rm -f "$file.formatted"
    fi
  fi
done

# Set executable permissions
chmod +x ./tests/api/debug_badge_assignment.sh

# Exit with the test exit code
exit $TEST_EXIT_CODE 
