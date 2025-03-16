# Badge Assignment System Testing Guide

## Table of Contents
- [Overview](#overview)
- [Testing Structure](#testing-structure)
- [Running Tests](#running-tests)
- [Writing Tests](#writing-tests)
- [Best Practices](#best-practices)
- [Mocking Dependencies](#mocking-dependencies)
- [Test Coverage](#test-coverage)

## Overview

This document provides guidelines for testing the Badge Assignment System. Our testing approach includes:

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test interactions between components
- **API Tests**: Test HTTP endpoints and API behavior

## Testing Structure

Our test structure follows Go best practices, organizing tests close to the code they're testing:

```
├── internal/                 # Application code
│   ├── engine/               # Engine implementation
│   │   └── rule_engine.go    # Rule engine implementation
│   │   └── rule_engine_test.go # Unit tests for rule engine
│   ├── cache/
│   │   └── badge_cache.go
│   │   └── badge_cache_test.go
│   └── testutil/             # Test utilities
│       ├── models.go         # Test models
│       └── helpers.go        # Test helper functions
├── tests/                    # Integration tests
│   └── integration/          # Integration test suite
│       ├── setup.go          # Test setup code
│       ├── badge/            # Badge-related integration tests
│       │   └── badge_test.go
│       ├── event/            # Event-related integration tests
│       └── condition/        # Condition-related integration tests
├── Makefile                  # Test commands
└── TESTING.md                # This guide
```

### Key Principles

1. **Unit tests** are placed in the same package as the code they're testing (`*_test.go` alongside implementation)
2. **Integration tests** are organized in the `tests/integration/` directory, grouped by feature
3. **Test utilities** are centralized in `internal/testutil/`

## Running Tests

We use a Makefile to simplify running tests:

```bash
# Run all tests
make test-all

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run specific feature tests
make test-badge
make test-event
make test-condition
make test-api

# Generate test coverage report
make test-coverage
```

### Test Environment Variables

Tests can be configured using environment variables:

- `API_TEST_URL`: URL for the API server in integration tests (default: http://localhost:8080)
- `TEST_API_TIMEOUT`: Timeout for waiting for API (default: 30s)
- `TEST_API_INTERVAL`: Polling interval for API health check (default: 1s)

Example:
```bash
API_TEST_URL=http://localhost:8081 make test-integration
```

## Writing Tests

### Unit Tests

Unit tests should be written in the same package as the code they're testing:

```go
// file: internal/engine/rule_engine_test.go
package engine

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
)

func TestRuleEvaluation(t *testing.T) {
    // Test setup
    engine := NewRuleEngine(...)
    
    // Test cases
    t.Run("ValidRule", func(t *testing.T) {
        result, err := engine.EvaluateRule(...)
        assert.NoError(t, err)
        assert.True(t, result)
    })
}
```

### Integration Tests

Integration tests should be in the `tests/integration/` directory:

```go
// file: tests/integration/badge/badge_test.go
package badge

import (
    "testing"
    "net/http"
    
    "github.com/badge-assignment-system/internal/testutil"
    "github.com/badge-assignment-system/tests/integration"
    "github.com/stretchr/testify/assert"
)

func TestBadgeCreation(t *testing.T) {
    integration.SetupTest()
    
    // Make API request
    response := testutil.MakeRequest(http.MethodPost, "/api/v1/badges", ...)
    testutil.AssertSuccess(t, response)
    
    // Assert expected behavior
    var badge testutil.Badge
    err := testutil.ParseResponse(response, &badge)
    assert.NoError(t, err)
    assert.NotZero(t, badge.ID)
}
```

### Table-Driven Tests

Use table-driven tests for testing multiple cases:

```go
func TestThresholdEvaluation(t *testing.T) {
    testCases := []struct{
        name           string
        scoreValue     float64
        threshold      float64
        expectedResult bool
    }{
        {"AboveThreshold", 95, 90, true},
        {"BelowThreshold", 85, 90, false},
        {"EqualToThreshold", 90, 90, true},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Test implementation using tc.scoreValue, tc.threshold
            // Assert result equals tc.expectedResult
        })
    }
}
```

## Best Practices

1. **Isolation**: Unit tests should not depend on external services or databases
2. **Deterministic**: Tests should be repeatable and produce the same results every time
3. **Fast**: Unit tests should run quickly to enable rapid development
4. **Independent**: Tests should not depend on each other or run in a specific order
5. **Clear**: Test names should clearly describe what is being tested
6. **Comprehensive**: Test both happy paths and error cases

### Naming Conventions

- Test functions should start with `Test`
- Test files should end with `_test.go`
- Benchmark functions should start with `Benchmark`
- Example functions should start with `Example`

## Mocking Dependencies

We use [testify/mock](https://github.com/stretchr/testify#mock-package) for mocking dependencies:

```go
// Define a mock for a dependency
type MockEventStore struct {
    mock.Mock
}

func (m *MockEventStore) GetEventsForUser(userID string) ([]models.Event, error) {
    args := m.Called(userID)
    return args.Get(0).([]models.Event), args.Error(1)
}

func TestWithMockedDependency(t *testing.T) {
    // Create the mock
    mockStore := new(MockEventStore)
    
    // Define expectations
    mockStore.On("GetEventsForUser", "user123").Return([]models.Event{...}, nil)
    
    // Inject the mock
    engine := NewRuleEngine(mockStore)
    
    // Test with the mock
    result, err := engine.EvaluateForUser("user123")
    
    // Verify expectations
    mockStore.AssertExpectations(t)
}
```

## Test Coverage

We aim for high test coverage but prioritize meaningful tests over coverage percentage:

```bash
# Generate coverage report
make test-coverage
```

The coverage report will be available at `coverage.html`.

### Coverage Targets

- Core business logic: 90%+ coverage
- API handlers: 80%+ coverage
- Utility functions: 70%+ coverage

## Continuous Integration

Tests are automatically run in CI pipelines on pull requests to ensure code quality.

## Debugging Tests

- Use `t.Log()` to add debug information
- Run specific tests with `-v` flag for verbose output:
  ```bash
  go test -v ./internal/engine -run TestSpecificFunction
  ```
- Use `-count=1` to disable test caching:
  ```bash
  go test -count=1 ./...
  ``` 
