# Contributing to Integration Tests

This guide explains how to effectively use and contribute to the integration test suite.

## Test Structure Overview

The integration tests are organized in a flat structure:

```
tests/integration/
├── README.md               # Overview documentation
├── CONTRIBUTING.md         # This guide
├── setup.go                # Contains setup code and global variables
├── integration_test.go     # Base test to verify setup
├── badge_sample_test.go    # Sample test showing the correct pattern
└── [your_feature]_test.go  # Your new tests following the pattern
```

## Adding New Tests

When adding new tests:

1. Create a new file named `feature_test.go` in the `tests/integration` directory
2. Use the `integration` package (not subpackages)
3. Follow the pattern shown in `badge_sample_test.go`
4. Utilize the global variables defined in `setup.go`

## Best Practices

1. **Proper Imports**:
   ```go
   import (
      "testing"
      
      "github.com/stretchr/testify/assert"
   )
   ```

2. **Test Function Structure**:
   ```go
   func TestFeature(t *testing.T) {
      // Skip if needed
      // SkipIfNotReady(t, "Not ready yet")
      
      // Setup
      
      // Test logic
      
      // Assertions
      
      // Cleanup
   }
   ```

3. **Use Global Variables**:
   ```go
   // Use these global variables from setup.go
   EventTypeID     // Default: 1
   ConditionTypeID // Default: 1
   BadgeID         // Default: 1
   TestUserID      // Default: "test-user-integration"
   ```

4. **Helper Functions**:
   ```go
   // Use these helper functions from setup.go
   SetupTest()          // Sets up the test environment
   SkipIfNotReady(t, "reason") // Skips a test with reason
   ```

## Running Tests

To run all integration tests:

```bash
make test-integration
```

To run a specific test:

```bash
go test -v ./tests/integration -run TestFeature
```

## Debugging Tests

1. Add detailed assertions to identify failures:
   ```go
   assert.Equal(t, expected, actual, "Description of what this checks")
   ```

2. Use `t.Logf()` for debugging info:
   ```go
   t.Logf("Debug info: %v", someValue)
   ```

3. Environment variables can be used to control test behavior:
   ```go
   // In setup.go
   TEST_API_TIMEOUT  // Timeout for API availability checks
   TEST_API_INTERVAL // Interval between API checks
   ```

## Reporting Issues

If you encounter issues with the integration tests:

1. Create a GitHub issue
2. Include the specific test that fails
3. Provide the complete error message
4. Document steps to reproduce the issue 
