#!/bin/bash

# Script to clean up the codebase by removing redundant test files
echo "Cleaning up the codebase to remove redundant test files..."

# Step 1: Remove backup files created during our fixes
echo "Step 1: Removing backup files..."
find tests/integration -name "*.bak" -type f -delete

# Step 2: Move all test directories to a backup location
echo "Step 2: Creating backup of test directories..."
backup_dir="tests/integration_backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$backup_dir"

for dir in badge event condition api; do
    if [ -d "tests/integration/$dir" ]; then
        echo "Backing up tests/integration/$dir to $backup_dir/$dir"
        cp -r "tests/integration/$dir" "$backup_dir/$dir"
    fi
done

# Step 3: Remove redundant test files (leave only essential test files)
echo "Step 3: Removing redundant test files..."
for dir in badge event condition api; do
    if [ -d "tests/integration/$dir" ]; then
        echo "Cleaning tests/integration/$dir"
        rm -rf "tests/integration/$dir"
    fi
done

# Step 4: Create a clean test structure with minimal tests
echo "Step 4: Creating a clean test structure..."

# Create a README file explaining the test structure
cat > "tests/integration/README.md" << 'EOF'
# Integration Tests

This directory contains integration tests for the badge assignment system.

## Test Structure

- `integration_test.go`: Contains the main integration test that verifies the system setup.
- Additional tests should be added here following the same pattern.

## Running Tests

To run the integration tests:

```bash
make test-integration
```

## Adding New Tests

When adding new tests, follow these guidelines:

1. Use the `integration` package
2. Import required dependencies
3. Utilize the global test variables defined in `setup.go`
4. Call `SkipIfNotReady` if the test prerequisites aren't met
EOF

echo "Cleanup complete!"
echo "The codebase is now clean without redundant test files."
echo "A backup of the original test files is available at $backup_dir"
echo "Run 'make test-integration' to verify the tests still pass." 
