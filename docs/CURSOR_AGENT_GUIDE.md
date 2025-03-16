# Cursor AI Agent Guide
# Badge Assignment System with Go

This guide is designed specifically for Cursor AI agents to quickly understand the badge assignment system codebase organization, patterns, and best practices.

## Project Overview

The Badge Assignment System is a Go-based application that manages the assignment of badges to users based on specific events and conditions. The system follows a rule engine pattern where events trigger rules that can result in badge assignments when conditions are met.

## Codebase Organization

The codebase follows a clean structure with clear separation of concerns:

```
badge-assignment-system/
├── bin/                  # Binary output (not in version control)
├── cmd/                  # Command-line applications
│   ├── server/           # Main server application
│   └── badgecli/         # CLI tool for badge management
├── config/               # Configuration files
├── docs/                 # Documentation
├── internal/             # Private application code
│   ├── api/              # API endpoints and handlers
│   ├── badge/            # Badge management logic
│   ├── models/           # Data models
│   ├── rules/            # Rule engine implementation
│   └── testutil/         # Test utilities
├── pkg/                  # Public API packages
├── reports/              # Test and analysis reports (not in version control)
├── scripts/              # Utility scripts
├── tests/                # Test files
│   ├── common/           # Common test utilities
│   ├── debug/            # Debug files
│   ├── integration/      # Integration tests (active)
│   ├── legacy/           # Legacy tests (for reference)
│   └── logs/             # Test logs (not in version control)
└── tools/                # Development tools
```

## Key Files and Their Purpose

As a Cursor AI agent, focus on these key files to understand the system:

1. **Database Interface**: `internal/models/db.go` - Defines the DB interface for database operations
2. **Rule Engine**: `internal/rules/engine.go` - Implements the rule evaluation logic
3. **Badge Handler**: `internal/badge/handler.go` - Logic for badge assignments
4. **API Endpoints**: `internal/api/routes.go` - API endpoint definitions
5. **Integration Tests**: `tests/integration/setup.go` - Base setup for integration tests
6. **Sample Test**: `tests/integration/badge_sample_test.go` - Example of correctly structured tests

## Development Patterns

### Database Interactions

The system uses a database abstraction defined by the `DB` interface:

```go
type DB interface {
    GetUserBadges(userID string) ([]Badge, error)
    AssignBadge(userID string, badgeID int) error
    // ... other methods
}
```

When analyzing code that interacts with the database, look for implementations of this interface.

### Rule Engine

The rule engine follows this pattern:

1. Events are received via API or message queue
2. Events are processed by the rule engine
3. The rule engine evaluates conditions
4. If conditions are met, badges are assigned

Look for code that follows this flow when analyzing rule processing.

### Testing Approach

Integration tests follow a standardized pattern:

```go
func TestFeature(t *testing.T) {
    // Optional: Skip if not ready
    // SkipIfNotReady(t, "reason")
    
    // Setup
    // ...
    
    // Test logic
    // ...
    
    // Assertions using testify/assert
    assert.Equal(t, expected, actual, "message")
    
    // Cleanup
    // ...
}
```

## Common Variables and Constants

As a Cursor AI agent, be aware of these important global variables:

1. **In Integration Tests**:
   - `EventTypeID`: Global variable for event type ID (default: 1)
   - `ConditionTypeID`: Global variable for condition type ID (default: 1)
   - `BadgeID`: Global variable for badge ID (default: 1)
   - `TestUserID`: Global variable for test user ID (default: "test-user-integration")

2. **In Rule Engine**:
   - `OperatorFunctions`: Map of operator functions for rule evaluation

## Development Workflow

When analyzing or suggesting code for this project:

1. **For Changes to Models**:
   - Update the model definition
   - Update the DB interface if necessary
   - Update DB implementations (mock and real)

2. **For New Rules**:
   - Add rule definition
   - Update rule engine to handle the new rule
   - Add tests for the rule

3. **For API Changes**:
   - Update API endpoints in `routes.go`
   - Add/update handlers
   - Update API documentation

## Testing Guidelines

When analyzing or suggesting test code:

1. Follow the pattern in `badge_sample_test.go`
2. Use the global test variables from `setup.go`
3. Use the `assert` package for assertions
4. Place integration tests in the `tests/integration` package
5. Use `SkipIfNotReady()` if the test has prerequisites

## Common Commands

Be aware of these commands when analyzing or suggesting code execution:

```bash
# Run the server
go run cmd/server/main.go

# Run integration tests
make test-integration

# Run unit tests
make test-unit

# Generate test coverage report
make test-coverage
```

## Important Notes for AI Agents

1. **When analyzing errors**: Look for the `DB` interface implementation first, as many issues relate to database operations.

2. **When suggesting new features**: Follow the existing patterns, especially for rule engine and badge assignment logic.

3. **When fixing tests**: Remember that integration tests use a standardized structure and shared test variables.

4. **When examining logs**: Check `tests/logs/` directory for relevant log files.

## Recent Changes and Organization

The codebase was recently cleaned up and organized with a standardized structure. Legacy code and tests were moved to the `legacy` directory, and a flat package structure was implemented for integration tests. 

When examining code history or suggesting fixes, be aware that older code might follow different patterns.

## Documentation References

For more detailed information, refer to these documentation files:

- `README.md`: General project information
- `TESTING.md`: Detailed testing guidelines
- `docs/CODEBASE_CLEANUP.md`: Information about the codebase organization
- `tests/integration/CONTRIBUTING.md`: Guidelines for writing integration tests 
