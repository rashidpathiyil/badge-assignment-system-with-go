#!/bin/bash

# Kill any previous server processes
pkill -f "go run" || true
echo "Starting server with debug output capturing..."

# Set environment variables for debugging
export DEBUG_RULE_ENGINE=true
export LOG_LEVEL=debug

# Create a logs directory if it doesn't exist
mkdir -p mongo_op_logs

# Start the server in the background with logs redirected
go run cmd/server/main.go > mongo_op_logs/server_logs.txt 2>&1 &
SERVER_PID=$!

echo "Server PID: $SERVER_PID"

# Give the server time to start up
echo "Waiting for server to start up..."
sleep 5

# Run the MongoDB-style operator test
echo "Running MongoDB-style operator test..."
go test -v ./tests/api -run TestMongoOperatorFormat

# Kill the server
echo "Shutting down server..."
kill $SERVER_PID

echo "Debug complete. Check mongo_operator_test.log for test logs and mongo_op_logs/server_logs.txt for server logs." 
