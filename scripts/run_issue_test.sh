#!/bin/bash

# Clear terminal
clear

# Create logs directory if it doesn't exist
mkdir -p issue_test_logs

# Clean previous logs
rm -f issue_test_logs/*

# Start the server in the background and capture its logs
echo "Starting server in the background..."
go run cmd/server/main.go > issue_test_logs/server_logs.txt 2>&1 &
SERVER_PID=$!

# Give the server time to start
echo "Waiting for server to initialize..."
sleep 5

# Run the test
echo "Running high-priority issue fixer badge test..."
go test -v ./tests/api/issue_tracking_test.go -count=1 | tee issue_test_logs/test_logs.txt

# Kill the server
echo "Shutting down server (PID: $SERVER_PID)..."
kill $SERVER_PID

# Print logs summary
echo ""
echo "=== Server Logs Summary ==="
tail -n 20 issue_test_logs/server_logs.txt

echo ""
echo "Complete logs available in the issue_test_logs directory" 
