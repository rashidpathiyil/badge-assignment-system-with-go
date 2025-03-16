#!/bin/bash

# Kill any previous server processes
pkill -f "go run" || true
echo "Starting server with debug output capturing..."

# Set environment variable to enable the debug test
export RUN_DEBUG_TESTS=true

# Set environment variable to enable debug logging in the rule engine (if supported)
export DEBUG_RULE_ENGINE=true
export LOG_LEVEL=debug

# Create a logs directory if it doesn't exist
mkdir -p badge_debug_logs

# Start the server in the background with logs redirected
go run cmd/server/main.go > badge_debug_logs/server_logs.txt 2>&1 &
SERVER_PID=$!

echo "Server PID: $SERVER_PID"

# Give the server time to start up
echo "Waiting for server to start up..."
sleep 5

# Run the debug test
echo "Running badge assignment debug test..."
# Use -v flag for verbose output and run only the debug test
go test -v ./tests/api -run TestDebugBadgeAssignment

# Kill the server
echo "Shutting down server..."
kill $SERVER_PID

echo "Debug complete. Check badge_debug.log for test logs and badge_debug_logs/server_logs.txt for server logs." 
