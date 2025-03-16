#!/bin/bash

# Script to organize the codebase structure
echo "Organizing codebase structure..."

# Create organized directory structure
echo "Creating organized directory structure..."
mkdir -p tests/legacy/api
mkdir -p tests/logs
mkdir -p tests/debug

# Move old test files to legacy directory
echo "Moving old API tests to legacy directory..."
mv tests/api/*.go tests/legacy/api/
mv tests/api/*.sh tests/legacy/api/
mv tests/api/README.md tests/legacy/api/

# Move logs to dedicated logs directory
echo "Moving log files to logs directory..."
mv tests/api/*.log tests/logs/
mv tests/api/*.txt tests/logs/
mv tests/api/badge_debug_logs tests/logs/

# Move utils to a common location
echo "Moving utils to common location..."
mkdir -p tests/common/utils
mv tests/api/utils/* tests/common/utils/
rmdir tests/api/utils

# Clean up empty directories
echo "Cleaning up empty directories..."
rmdir tests/api 2>/dev/null || true

# Move clean_tests.sh to archive
echo "Archiving unneeded scripts..."
mkdir -p tools/archive
mv tools/clean_tests.sh tools/archive/
mv tools/cleanup_codebase.sh tools/archive/

# Create a README in the tools directory
cat > tools/README.md << 'EOF'
# Tools Directory

This directory contains utility scripts and tools for the badge assignment system.

## Available Tools

- `final_cleanup.sh`: Removes temporary files and cleans up the codebase
- `migrate_tests.sh`: Script used to migrate tests to the new structure
- `organize_codebase.sh`: Organizes the codebase structure for better maintainability

## Archived Tools

Tools that were used during refactoring and are kept for reference:

- `archive/clean_tests.sh`: Used to clean up tests during refactoring
- `archive/cleanup_codebase.sh`: Used to remove redundant test files

## Database Tools

See `db-utils/README.md` for information about database utilities.
EOF

# Create a README in the legacy directory
cat > tests/legacy/README.md << 'EOF'
# Legacy Tests

This directory contains legacy test files that were part of the original codebase.
These tests are kept for reference but are not part of the active test suite.

## Structure

- `api/`: Original API tests before refactoring to the new integration test structure

## Notes

When writing new tests, please follow the pattern in `tests/integration` rather than
using these legacy tests as examples.
EOF

echo "Codebase organization complete!"
echo "The codebase is now properly organized with:"
echo "- Integration tests in tests/integration/"
echo "- Legacy tests in tests/legacy/"
echo "- Logs in tests/logs/"
echo "- Common utilities in tests/common/"
echo "- Organized tools in tools/ with appropriate documentation" 
