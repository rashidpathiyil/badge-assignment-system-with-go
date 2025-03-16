# API Integration Tests

This directory contains integration tests for the Badge Assignment System API.

## Overview

The tests verify the full badge assignment flow:

1. Creating an event type
2. Creating a condition type
3. Creating a badge with criteria
4. Processing an event for a user
5. Checking that the user has been awarded the badge

## Running Tests

To run these tests, you need to have the Badge Assignment System API running:

```bash
# From the project root
go test -v ./tests/api
```

By default, the tests connect to `http://localhost:8080`. You can override this by setting the `API_TEST_URL` environment variable:

```bash
API_TEST_URL=http://your-server:8080 go test -v ./tests/api
```

## Test Structure

- `api_integration_test.go`: The main test file that runs through the full flow
- `utils/`: Helper functions and models for the tests
  - `test_helpers.go`: Common test utilities
  - `models.go`: Models for the API request/response data

## Test Flow

The tests run in a specific order to simulate the real-world usage of the badge system:

1. **Create Event Type**: Creates a new event type for tracking challenge completions
2. **Create Condition Type**: Creates a condition type to evaluate event data
3. **Create Badge**: Creates a badge with criteria using the condition type
4. **Process Event**: Submits an event that should trigger badge assignment
5. **Check User Badges**: Verifies that the user has been awarded the badge

## Debugging

If tests fail, look at the error message and response contents to understand what went wrong.

Common issues:
- Server not running
- Database not properly initialized
- Missing permissions
- Invalid request payloads 
