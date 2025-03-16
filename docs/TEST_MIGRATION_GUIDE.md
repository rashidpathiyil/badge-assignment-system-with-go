# Test Migration Guide

This document provides instructions for migrating your existing tests to the new well-organized test structure.

## Overview of Changes

We've reorganized the test structure to follow Go best practices:

1. Unit tests are now placed alongside the code they're testing (`*_test.go` files in the same package)
2. Integration tests are organized by feature in the `tests/integration/` directory
3. Test utilities have been centralized in `internal/testutil/`
4. A Makefile has been added to simplify running tests

## Migration Steps

### 1. Run the Migration Script

We've provided a script to help automate the migration of existing tests:

```bash
chmod +x tools/migrate_tests.sh
./tools/migrate_tests.sh
```

This script will:
- Create the new directory structure
- Copy existing tests to appropriate locations in the new structure
- Update import paths in the migrated files
- Create backups of the original files with `.bak` extensions

### 2. Update Package Declarations

For each migrated test file, you'll need to update the package declaration at the top to match its new location:

- Files in `tests/integration/badge/` should use `package badge`
- Files in `tests/integration/event/` should use `package event`
- Files in `tests/integration/condition/` should use `package condition`
- Files in `tests/integration/api/` should use `package api`

### 3. Convert Utils Import

The migration script updates imports from:
```go
import "github.com/badge-assignment-system/tests/api/utils"
```

To:
```go
import "github.com/badge-assignment-system/internal/testutil"
```

You'll need to update any references from `utils.X` to `testutil.X` in your test files.

### 4. Fix Test Dependencies

Update your tests to use the centralized `integration` package for shared functionality:

```go
import (
    "github.com/badge-assignment-system/internal/testutil"
    "github.com/badge-assignment-system/tests/integration"
)

func TestSomething(t *testing.T) {
    integration.SetupTest()
    
    // Use shared IDs
    if integration.BadgeID == 0 {
        t.Skip("Badge ID not set")
    }
    
    // Use testutil helpers
    response := testutil.MakeRequest(...)
}
```

### 5. Create Unit Tests

For components that don't have unit tests yet, use the template in `tools/templates/unit_test_template.go.txt` as a starting point.

Example for creating a unit test for `internal/engine/rule_engine.go`:

```bash
cp tools/templates/unit_test_template.go.txt internal/engine/rule_engine_test.go
```

Then modify the test file to match your implementation.

### 6. Run Tests

Use the Makefile to run your tests:

```bash
# Run all tests
make test-all

# Run only integration tests
make test-integration

# Run only unit tests
make test-unit

# Run tests for a specific feature
make test-badge
```

### 7. Add Missing Tests

Now that the structure is in place, identify areas with insufficient test coverage and add tests following the patterns in the example files:

- Unit tests should use the table-driven pattern where appropriate
- Integration tests should be organized by feature
- Tests should be independent and not rely on the order of execution

## Best Practices

1. **Use Table-Driven Tests**: For functions with multiple test cases:
   ```go
   testCases := []struct{
       name     string
       input    string
       expected string
   }{
       {"Case1", "input1", "expected1"},
       {"Case2", "input2", "expected2"},
   }
   
   for _, tc := range testCases {
       t.Run(tc.name, func(t *testing.T) {
           result := FunctionToTest(tc.input)
           assert.Equal(t, tc.expected, result)
       })
   }
   ```

2. **Mock Dependencies**: Use testify/mock for unit testing components with dependencies:
   ```go
   mockDep := new(MockDependency)
   mockDep.On("Method", "arg").Return("result", nil)
   
   sut := NewSystemUnderTest(mockDep)
   // Test sut...
   
   mockDep.AssertExpectations(t)
   ```

3. **Skip if Prerequisites Missing**: Use t.Skip for integration tests that depend on previous steps:
   ```go
   if integration.BadgeID == 0 {
       t.Skip("Badge ID not set, skipping test")
   }
   ```

4. **Clean Up After Tests**: Ensure tests clean up resources they create:
   ```go
   t.Cleanup(func() {
       // Clean up resources
   })
   ```

## Further Information

See the [TESTING.md](TESTING.md) document for comprehensive testing guidelines. 
