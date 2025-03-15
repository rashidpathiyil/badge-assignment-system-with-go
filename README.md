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

## API Documentation

### Public APIs

#### List all badges
```
GET /api/v1/badges
```

#### List all active badges
```
GET /api/v1/badges/active
```

#### Get badge details
```
GET /api/v1/badges/:id
```

#### Get user badges
```
GET /api/v1/users/:id/badges
```

#### Process an event
```
POST /api/v1/events
```
Example payload:
```json
{
  "event_type": "meeting-attendance",
  "user_id": "user123",
  "payload": {
    "meeting_id": "m456",
    "duration_minutes": 60
  },
  "timestamp": "2023-04-01T09:00:00Z"
}
```

### Admin APIs

#### Event Type Management
```
POST /api/v1/admin/event-types
GET /api/v1/admin/event-types
GET /api/v1/admin/event-types/:id
PUT /api/v1/admin/event-types/:id
DELETE /api/v1/admin/event-types/:id
```

#### Badge Management
```
POST /api/v1/admin/badges
GET /api/v1/admin/badges/:id/criteria
PUT /api/v1/admin/badges/:id
DELETE /api/v1/admin/badges/:id
```

#### Condition Type Management
```
POST /api/v1/admin/condition-types
GET /api/v1/admin/condition-types
GET /api/v1/admin/condition-types/:id
PUT /api/v1/admin/condition-types/:id
DELETE /api/v1/admin/condition-types/:id
```

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

The badge system includes a suite of unit tests for the time-based criteria evaluators in `internal/engine/time_evaluators_test.go`. These tests verify that the enhanced badge criteria functionality works correctly.

#### Running the Tests

1. Navigate to the `internal/engine` directory:
   ```bash
   cd internal/engine
   ```

2. Run all time evaluator tests:
   ```bash
   go test -v
   ```

3. Run a specific test:
   ```bash
   go test -v -run TestTimePeriodCriteria
   ```

#### Test Descriptions

- **TestTimePeriodCriteria**: Tests the evaluation of time period-based criteria (days, weeks, months)
- **TestPatternCriteria**: Tests pattern recognition for consistent, increasing, or decreasing event frequencies
- **TestGapCriteria**: Tests evaluation of time gaps between events
- **TestDurationCriteria**: Tests measurement of durations between paired events (e.g., start/end)
- **TestAggregationCriteria**: Tests computation of aggregated values across events

#### Troubleshooting Failed Tests

If tests are failing:

1. Check the test output for specific error messages
2. Verify that the time_utils.go and time_evaluators.go files have consistent parameter names
3. Ensure that the models referenced in the tests match the current model definitions

#### Extending the Tests

To add new test cases:

1. Review the existing test patterns in time_evaluators_test.go
2. Create test events with the appropriate timestamps and payload data
3. Define criteria that test edge cases or new functionality
4. Verify that metadata is correctly populated during evaluation

## License

This project is licensed under the MIT License. 
