#!/bin/bash

# Set environment variable to enable the debug test
export RUN_DEBUG_TESTS=true

# Run the debug query test
echo "Running database debug query test..."
go test -v ./tests/api -run TestQueryDebugInfo

echo "Debug complete. Check db_debug.log for detailed database information." 
