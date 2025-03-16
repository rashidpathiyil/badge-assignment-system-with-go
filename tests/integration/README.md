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
