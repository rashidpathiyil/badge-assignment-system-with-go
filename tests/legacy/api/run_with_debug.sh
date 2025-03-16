#!/bin/bash

# Script to restart the server with debug logging and run the tests

# Set paths
PROJECT_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
LOG_FILE="$PROJECT_ROOT/server_debug.log"
ENV_BACKUP="$PROJECT_ROOT/.env.backup"

echo "======================================================================"
echo "  BADGE ASSIGNMENT SYSTEM TEST WITH SERVER DEBUG LOGGING"
echo "======================================================================"
echo "This script will:"
echo "1. Backup your current .env file"
echo "2. Add debug logging settings"
echo "3. Restart the server"
echo "4. Run the API tests"
echo "5. Restore your original .env file"
echo ""
echo "Server logs will be saved to: $LOG_FILE"
echo "======================================================================" 
echo ""

# Check if server is already running
if ! curl -s http://localhost:8080/health > /dev/null; then
  echo "Error: API server is not running at http://localhost:8080"
  echo "Please start the server first"
  exit 1
fi

# Backup current .env file
echo "Backing up current .env file..."
cp "$PROJECT_ROOT/.env" "$ENV_BACKUP"

# Add debug settings to .env
echo "Adding debug logging settings to .env..."
cat >> "$PROJECT_ROOT/.env" << EOF

# Debug settings for badge assignment testing
LOG_LEVEL=trace
ENGINE_LOG_LEVEL=trace
ENABLE_BADGE_DEBUG=true
BADGE_PROCESSING_LOG=true
DETAILED_EVENT_LOGS=true
EOF

# Find the server process and restart it
echo "Restarting the server with debug logging..."
SERVER_PID=$(ps aux | grep "./bin/server" | grep -v grep | awk '{print $2}')

if [ -n "$SERVER_PID" ]; then
  echo "Stopping existing server (PID: $SERVER_PID)..."
  kill $SERVER_PID
  sleep 2
else
  echo "No running server found, will start a new one"
fi

# Start the server with logging to file
cd "$PROJECT_ROOT"
echo "Starting server with debug logging to $LOG_FILE..."
./bin/server > "$LOG_FILE" 2>&1 &
SERVER_PID=$!

# Wait for server to start
echo "Waiting for server to start..."
for i in {1..10}; do
  if curl -s http://localhost:8080/health > /dev/null; then
    echo "Server is running!"
    break
  fi
  if [ $i -eq 10 ]; then
    echo "Server failed to start within the expected time"
    # Restore original .env file
    mv "$ENV_BACKUP" "$PROJECT_ROOT/.env"
    exit 1
  fi
  echo "Waiting... ($i/10)"
  sleep 1
done

# Run the tests
echo "Running API tests with debug logging..."
cd "$PROJECT_ROOT/tests/api"
./run_tests.sh

# Run the debug test for more detailed info
echo "Running detailed debug test..."
RUN_DEBUG_TESTS=true go test -v ./tests/api -run TestDebugBadgeAssignment

# Restore original .env file
echo "Restoring original .env file..."
mv "$ENV_BACKUP" "$PROJECT_ROOT/.env"

# Restart the server with original settings
echo "Restarting the server with original settings..."
kill $SERVER_PID
sleep 2
cd "$PROJECT_ROOT"
./bin/server > /dev/null 2>&1 &

echo "Done!"
echo "Check $LOG_FILE for detailed server logs" 
