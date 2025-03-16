#!/bin/bash

# migrate_tests.sh - Script to help migrate existing tests to the new structure
#
# Usage: ./tools/migrate_tests.sh
#
# This script will:
# 1. Create the new test directories if they don't exist
# 2. Move API tests to the new integration test structure
# 3. Update import paths in the moved files
# 4. Create a backup of the original tests

set -e

echo "Starting test migration..."

# Create directories if they don't exist
mkdir -p tests/integration/badge
mkdir -p tests/integration/event
mkdir -p tests/integration/condition
mkdir -p tests/integration/api

# Function to copy a test file to the new structure with import updates
migrate_test_file() {
  local src=$1
  local dest=$2
  local file_name=$(basename "$src")
  
  echo "Migrating $file_name to $dest"
  
  # Create a backup of the original file
  cp "$src" "${src}.bak"
  
  # Copy the file to the new location
  mkdir -p "$(dirname "$dest")"
  cp "$src" "$dest"
  
  # Update imports
  sed -i '' 's|"github.com/badge-assignment-system/tests/api/utils"|"github.com/badge-assignment-system/internal/testutil"|g' "$dest"
  
  echo "  ✓ Migrated $file_name"
}

# Migrate badge-related tests
if [ -f "tests/api/badge_test.go" ]; then
  migrate_test_file "tests/api/badge_test.go" "tests/integration/badge/badge_test.go"
fi

if [ -f "tests/api/badge_debug_test.go" ]; then
  migrate_test_file "tests/api/badge_debug_test.go" "tests/integration/badge/badge_debug_test.go"
fi

# Migrate event-related tests
if [ -f "tests/api/event_type_test.go" ]; then
  migrate_test_file "tests/api/event_type_test.go" "tests/integration/event/event_type_test.go"
fi

if [ -f "tests/api/event_processing_test.go" ]; then
  migrate_test_file "tests/api/event_processing_test.go" "tests/integration/event/event_processing_test.go"
fi

# Migrate condition-related tests
if [ -f "tests/api/condition_type_test.go" ]; then
  migrate_test_file "tests/api/condition_type_test.go" "tests/integration/condition/condition_type_test.go"
fi

# Migrate operator tests
if [ -f "tests/api/logical_operator_test.go" ]; then
  migrate_test_file "tests/api/logical_operator_test.go" "tests/integration/api/logical_operator_test.go"
fi

if [ -f "tests/api/operators_test.go" ]; then
  migrate_test_file "tests/api/operators_test.go" "tests/integration/api/operators_test.go"
fi

# Migrate format tests
if [ -f "tests/api/correct_badge_format_test.go" ]; then
  migrate_test_file "tests/api/correct_badge_format_test.go" "tests/integration/badge/format_test.go"
fi

if [ -f "tests/api/negative_criteria_test.go" ]; then
  migrate_test_file "tests/api/negative_criteria_test.go" "tests/integration/condition/negative_criteria_test.go"
fi

# Migrate issue tracking tests
if [ -f "tests/api/issue_tracking_test.go" ]; then
  migrate_test_file "tests/api/issue_tracking_test.go" "tests/integration/api/issue_tracking_test.go"
fi

# Main integration test
if [ -f "tests/api/api_integration_test.go" ]; then
  migrate_test_file "tests/api/api_integration_test.go" "tests/integration/api/api_integration_test.go"
fi

echo "Migration complete!"
echo ""
echo "Next steps:"
echo "1. Review the migrated files to ensure they work with the new structure"
echo "2. Update package declarations in the migrated files"
echo "3. Run the tests with 'make test-integration' to verify they work"
echo "4. Once verified, you can remove the .bak files with 'find tests -name \"*.bak\" -delete'"
echo ""
echo "See TESTING.md for more details on the new test structure." 
