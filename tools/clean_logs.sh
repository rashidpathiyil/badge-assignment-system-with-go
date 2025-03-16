#!/bin/bash

# Script to clean up log directories in the root
echo "Cleaning up log directories in the root..."

# Create logs directory structure
mkdir -p tests/logs/badge_debug
mkdir -p tests/logs/correct_format
mkdir -p tests/logs/issue_test
mkdir -p tests/logs/logical_op
mkdir -p tests/logs/mongo_op
mkdir -p tests/logs/negative_test

# Move log directories to tests/logs
echo "Moving log directories to tests/logs..."

# Function to move logs
move_logs() {
    src_dir=$1
    dest_dir=$2
    
    if [ -d "$src_dir" ]; then
        echo "Moving $src_dir to $dest_dir"
        if [ "$(ls -A "$src_dir" 2>/dev/null)" ]; then
            # Directory has content
            cp -r "$src_dir"/* "$dest_dir"/
        fi
        rm -rf "$src_dir"
    fi
}

# Move each log directory
move_logs "badge_debug_logs" "tests/logs/badge_debug"
move_logs "correct_format_logs" "tests/logs/correct_format"
move_logs "issue_test_logs" "tests/logs/issue_test"
move_logs "logical_op_logs" "tests/logs/logical_op"
move_logs "mongo_op_logs" "tests/logs/mongo_op"
move_logs "negative_test_logs" "tests/logs/negative_test"

# Update README.md in logs directory
cat > "tests/logs/README.md" << 'EOF'
# Logs Directory

This directory contains logs from test runs. These files are not tracked in version control.

## Structure

- `badge_debug/`: Logs from badge debug tests
- `correct_format/`: Logs from format correctness tests
- `issue_test/`: Logs from issue tracking tests
- `logical_op/`: Logs from logical operator tests
- `mongo_op/`: Logs from MongoDB operator tests
- `negative_test/`: Logs from negative criteria tests

## Notes

Log files should never be committed to version control.
EOF

echo "Log directories cleanup complete!"
echo "All logs are now organized in the tests/logs directory." 
