# Legacy Engine Tests

This directory contains legacy test files for the rule engine that were part of the original codebase.
These tests are kept for reference but are not part of the active test suite.

## Structure

- `badge_test.go`: Original badge tests
- `badge_tests/`: Tests for badge assignment functionality
- `pattern_criteria/`: Tests for pattern detection functionality

## Migration Status

These tests have been migrated to proper unit tests:
- `badge_tests/` → `internal/engine/badge/badge_test.go`
- `pattern_criteria/` → `internal/engine/pattern_test.go`

## Notes

When writing new tests, please follow the pattern in `tests/integration` rather than
using these legacy tests as examples.
