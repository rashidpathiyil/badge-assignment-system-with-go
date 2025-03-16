#!/bin/bash

# Clear terminal
clear

# Create logs directory if it doesn't exist
mkdir -p negative_test_logs

# Clean previous logs
rm -f negative_test_logs/*

# Start the server in the background and capture its logs
echo "Starting server in the background..."
go run cmd/server/main.go > negative_test_logs/server_logs.txt 2>&1 &
SERVER_PID=$!

# Give the server time to start
echo "Waiting for server to initialize..."
sleep 5

# Run the test
echo "Running negative criteria test..."
go test -v ./tests/api/negative_criteria_test.go -count=1 | tee negative_test_logs/test_logs.txt

# Kill the server
echo "Shutting down server (PID: $SERVER_PID)..."
kill $SERVER_PID

# Print logs summary
echo ""
echo "=== Server Logs Summary ==="
tail -n 20 negative_test_logs/server_logs.txt

echo ""
echo "Complete logs available in the negative_test_logs directory" 
