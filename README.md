# Scalable Badge Assignment System

A flexible, dynamic, and scalable Badge Assignment System built as a standalone microservice. The system supports dynamic badge creation through a JSON-based configuration model, where badge criteria are defined entirely via JSON configurations using MongoDB-like query operators.

## Features

- **Dynamic Rule Engine:** Badge criteria are defined using JSON with MongoDB-like query operators (`$gte`, `$lt`, `$and`, etc.) to express conditions
- **Modular Architecture:** Independent microservice that can interface with multiple external systems
- **Flexible Data Modeling:** PostgreSQL with JSONB for storing dynamic badge criteria
- **RESTful API:** Full API for badge management, event ingestion, and rule evaluation
- **Scalable Design:** Capable of handling high volumes of events with a clean separation of concerns

## System Components

1. **Core System:**
   - Event ingestion and storage
   - Dynamic rule evaluation engine
   - Badge assignment logic

2. **APIs:**
   - Public APIs for users (badge listing, user badges)
   - Admin APIs for badge management
   - Event processing endpoint
   
3. **Database Schema:**
   - Event types
   - Condition types
   - Badges
   - Badge criteria (JSON-based)
   - User badges
   - Events

## Prerequisites

- Go 1.18 or later
- PostgreSQL 12 or later

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/badge-assignment-system.git
   cd badge-assignment-system
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Set up the database:
   - Create a PostgreSQL database named `badge_system`
   - Run the migrations:
     ```
     # You'll need to install the migrate tool first:
     # https://github.com/golang-migrate/migrate/
     
     migrate -database "postgres://postgres:postgres@localhost:5432/badge_system?sslmode=disable" -path db/migrations up
     ```

4. Create a `.env` file:
   ```
   cp .env.example .env
   # Then edit .env with your configuration
   ```

5. Build and run the server:
   ```
   go build -o bin/server cmd/server/main.go
   ./bin/server
   ```

## Documentation

Comprehensive documentation for the Badge Assignment System is available in the `docs` directory:

- [API Documentation](docs/api/README.md): Complete API reference with endpoints, request/response formats, and examples
- [Badge Criteria Documentation](docs/README_BADGE_CRITERIA.md): Detailed explanation of badge criteria syntax
- [Rule Engine Guide](docs/RULE_ENGINE_GUIDE.md): Technical guide to the rule engine
- [Pattern Detection Reference](docs/PATTERN_DETECTION_REFERENCE.md): Documentation for pattern detection features
- [Count Operators Documentation](docs/COUNT_OPERATORS.md): Guide to count operators for badge criteria
- [Database Requirements](docs/DATABASE_REQUIREMENTS.md): Information about database tables, dependencies, and mandatory requirements
- [Testing Documentation](docs/Testing.md): Information on testing the system
- [Onboarding Guide](docs/ONBOARDING.md): Guide for new developers

For a quick overview of the available APIs, see below:

### Public APIs

- `GET /api/v1/badges` - List all badges
- `GET /api/v1/badges/active` - List all active badges
- `GET /api/v1/badges/:id` - Get badge details
- `GET /api/v1/users/:id/badges` - Get user badges
- `POST /api/v1/events` - Process an event

### Admin APIs

- `/api/v1/admin/badges/*` - Badge management
- `/api/v1/admin/event-types/*` - Event type management
- `/api/v1/admin/condition-types/*` - Condition type management

For detailed API documentation, including request/response formats and examples, please refer to the [API Documentation](docs/api/README.md).

## Badge Criteria Examples

### Meeting Maestro Badge
A badge awarded to users who attend at least 10 meetings in a month:
```json
{
  "event": "meeting-attendance",
  "criteria": {
    "$and": [
      { "timestamp": { "$gte": "2023-04-01T00:00:00Z" } },
      { "timestamp": { "$lt": "2023-05-01T00:00:00Z" } }
    ],
    "count": { "$gte": 10 }
  }
}
```

### Punctuality Pro Badge
A badge awarded to users who arrive on time for at least 5 meetings:
```json
{
  "event": "meeting-attendance",
  "criteria": {
    "arrived_on_time": true,
    "count": { "$gte": 5 }
  }
}
```

### Task Champion Badge
A badge awarded to users who complete at least 20 tasks with a high priority:
```json
{
  "event": "task-completion",
  "criteria": {
    "priority": "high",
    "count": { "$gte": 20 }
  }
}
```

## Development

### Project Structure
```
badge-assignment-system/
├── cmd/
│   └── server/             # Main application entry point
├── config/                 # Configuration files
├── db/
│   └── migrations/         # Database migrations
├── internal/
│   ├── api/                # HTTP handlers and routes
│   ├── engine/             # Rule evaluation engine
│   ├── models/             # Database models and queries
│   └── service/            # Business logic
├── pkg/
│   └── utils/              # Shared utilities
├── .env.example           
├── go.mod
├── go.sum
└── README.md
```

### Running Time-Based Criteria Tests

The badge system includes a comprehensive suite of tests for the time-based criteria evaluators. These tests verify that the enhanced badge criteria functionality works correctly, with a focus on real-world usage scenarios.

#### Important Notes on Time Window Operators

Our testing has revealed some limitations with the `$timeWindow` operator:

- The `$timeWindow` operator may not work correctly when combined with other operators like `$eventCount`.
- Using the `last` parameter (e.g., `"last": "30d"`) or fixed dates with `start` and `end` parameters showed inconsistent results.

**Recommended Alternative:** Use the `timestamp` field directly in your criteria for time-based filtering:

```json
{
  "event": "user_activity",
  "criteria": {
    "$eventCount": {
      "$gte": 5
    },
    "timestamp": {
      "$gte": "2023-12-01T00:00:00Z",
      "$lte": "2023-12-31T23:59:59Z"
    }
  }
}
```

**Proposed Enhancement:** We're working on adding support for dynamic time variables that would make time-based criteria more flexible:

```json
{
  "event": "user_activity",
  "criteria": {
    "$eventCount": {
      "$gte": 5
    },
    "timestamp": {
      "$gte": "$NOW(-30d)"  // Dynamic variable for "30 days ago"
    }
  }
}
```

For more details, see:
- [Count Operators Documentation](docs/COUNT_OPERATORS.md)
- [Time-Based Criteria Documentation](docs/TIME_BASED_CRITERIA.md)
- [Dynamic Time Variables Proposal](docs/DYNAMIC_TIME_VARIABLES.md)
- [Timestamp Filter Test](/tests/integration/timestamp_filter_badge_test.go)

#### Running the Tests

1. Run all tests:
   ```bash
   go test ./...
   ```

2. Run time evaluator tests:
   ```bash
   cd internal/engine
   go test -v
   ```

3. Run real-world pattern detection tests:
   ```bash
   cd internal/engine/tests/pattern_criteria
   go test -v
   ```

4. Run all real-world scenario tests with the test runner:
   ```bash
   cd internal/engine/tests
   go run run_tests.go
   ```

#### Test Descriptions

- **TestTimePeriodCriteria**: Tests the evaluation of time period-based criteria (days, weeks, months)
- **TestPatternCriteria**: Tests pattern recognition for consistent, increasing, or decreasing event frequencies
- **TestGapCriteria**: Tests evaluation of time gaps between events
- **TestDurationCriteria**: Tests measurement of durations between paired events (e.g., start/end)
- **TestAggregationCriteria**: Tests computation of aggregated values across events

#### Real-World Scenario Tests

The system includes real-world scenario tests that simulate actual user behavior:

- **User Engagement Pattern**: Tests consistent daily app usage detection
- **Fitness Progress Pattern**: Tests increasing workout counts detection
- **Learning Pattern Decline**: Tests gradual decline in learning activities
- **Seasonal Usage Pattern**: Tests consistent weekly usage with monthly variations
- **Mixed Pattern Detection**: Tests multiple patterns on the same dataset

These tests are located in the `internal/engine/tests/pattern_criteria` directory and use the `createEventsWithPattern` helper function to generate realistic test data.

#### Improved Pattern Detection Algorithms

The system features robust pattern detection algorithms that can identify consistent, increasing, and decreasing patterns in user behavior:

- **Consistent Pattern Detection**: Identifies stable usage patterns with tolerance for isolated anomalies
- **Increasing Pattern Detection**: Detects growing engagement with trend strength calculation
- **Decreasing Pattern Detection**: Identifies declining usage with chronological correction

These algorithms provide detailed metadata about detected patterns, including:
- Consistency metrics (deviation, coefficient of variation)
- Growth/decline percentages
- Trend strength measurements
- Period-by-period breakdowns

#### Troubleshooting Failed Tests

If tests are failing:

1. Check the test output for specific error messages
2. Verify that the time_utils.go and time_evaluators.go files have consistent parameter names
3. Ensure that the models referenced in the tests match the current model definitions

#### Extending the Tests

To add new test cases:

1. Review the existing test patterns in time_evaluators_test.go or pattern_test.go
2. Create test events with the appropriate timestamps and payload data
3. Define criteria that test edge cases or new functionality
4. Verify that metadata is correctly populated during evaluation

## License

This project is licensed under the MIT License.
