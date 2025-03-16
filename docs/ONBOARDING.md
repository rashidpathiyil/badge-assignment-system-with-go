# Developer Onboarding Guide: Badge Assignment System

Welcome to the Badge Assignment System! This guide is designed to help new developers get up to speed quickly with our codebase, architecture, and development workflow.

## Quick Start

### 1. Environment Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/badge-assignment-system.git
cd badge-assignment-system

# Install dependencies
go mod download

# Set up the database
# You'll need PostgreSQL installed
createdb badge_system

# Run migrations
migrate -database "postgres://postgres:postgres@localhost:5432/badge_system?sslmode=disable" -path db/migrations up

# Create configuration
cp .env.example .env
# Edit .env with your settings

# Build and run
go build -o bin/server cmd/server/main.go
./bin/server
```

### 2. Run Tests

```bash
# Run all tests
go test ./...

# Run specific test suites
cd internal/engine
go test -v

# Run real-world scenario tests
cd internal/engine/tests/pattern_criteria
go test -v
```

## System Architecture

The Badge Assignment System is designed as a microservice that handles dynamic badge criteria evaluation. Here's a high-level overview:

```
┌─────────────────┐    ┌────────────────┐    ┌────────────────┐
│   API Layer     │───▶│  Rule Engine   │───▶│  Data Access   │
│  (HTTP Routes)  │◀───│ (Evaluations)  │◀───│    Layer       │
└─────────────────┘    └────────────────┘    └────────────────┘
         │                                           │
         │                                           │
         ▼                                           ▼
┌─────────────────┐                       ┌────────────────┐
│     Cache       │                       │   PostgreSQL   │
│     Layer       │                       │   Database     │
└─────────────────┘                       └────────────────┘
```

### Key Components

1. **API Layer** (`internal/api`): HTTP handlers and routes for badge management and event processing
2. **Rule Engine** (`internal/engine`): Core logic for evaluating badge criteria against events
3. **Data Access Layer** (`internal/models`): Database models and queries
4. **Cache Layer** (`internal/cache`): Performance optimization through caching
5. **Service Layer** (`internal/service`): Business logic and coordination

## The Rule Engine

The heart of our system is the rule engine, which evaluates badge criteria against user events. 

### Data Flow

1. Events are ingested through the API
2. The rule engine evaluates events against badge criteria
3. When criteria are met, badges are awarded to users
4. Results are stored in the database
5. Badge awards are returned via API endpoints

### Badge Criteria Types

Our system supports various criteria types for badge awards, including:

1. **Basic Criteria**: Simple conditions on event properties
2. **Logical Operators**: Combining conditions with `$and`, `$or`, `$not`
3. **Time-Based Criteria**: Evaluating patterns and frequencies over time

## Pattern Detection Algorithms

Our system features sophisticated pattern detection algorithms that analyze user behavior over time to identify meaningful patterns.

### Pattern Types

1. **Consistent Pattern**: Identifies stable user behavior with small, acceptable variations
2. **Increasing Pattern**: Detects growing engagement or improvement over time
3. **Decreasing Pattern**: Identifies declining usage patterns for intervention opportunities

### How Pattern Detection Works

Pattern detection operates through these steps:

1. **Event Grouping**: Events are grouped by time periods (day, week, month)
2. **Count Aggregation**: The number of events per period is calculated
3. **Pattern Analysis**: Algorithms analyze these counts for specific patterns
4. **Metadata Collection**: Detailed metrics are gathered for analysis
5. **Criteria Evaluation**: The pattern is compared against badge criteria

### Key Algorithms

#### Consistent Pattern Detection

```go
// Algorithm overview:
// 1. Calculate average event count across periods
// 2. Calculate standard deviation and coefficient of variation
// 3. Check for outliers that could skew results
// 4. Determine if variations are within acceptable bounds
// 5. Generate detailed metadata about consistency
```

Key metrics:
- Maximum deviation
- Coefficient of variation
- Standard deviation
- Outlier detection

#### Increasing Pattern Detection

```go
// Algorithm overview:
// 1. Calculate percentage increases between consecutive periods
// 2. Determine average percentage increase
// 3. Calculate trend strength using statistical methods
// 4. Check if pattern meets minimum criteria
// 5. Generate detailed growth metadata
```

Key metrics:
- Average percent increase
- Consecutive growth periods
- Trend strength
- Period-by-period breakdown

#### Decreasing Pattern Detection

```go
// Algorithm overview:
// 1. Calculate percentage decreases between consecutive periods
// 2. Apply chronological ordering and possible corrections
// 3. Determine if decrease is gradual rather than sudden
// 4. Check if pattern meets maximum decrease criteria
// 5. Generate detailed metadata on the decline pattern
```

Key metrics:
- Average percent decrease
- Consecutive decrease periods
- Rate of decline
- Trend analysis

## Development Workflow

### Adding a New Badge Criterion

1. Define the new criterion type in `models/criteria.go`
2. Implement the evaluation logic in the rule engine
3. Add test cases to validate the behavior
4. Update documentation to include the new criterion

### Adding a New Pattern Type

1. Define the pattern type in `models/criteria.go`
2. Implement the pattern detection algorithm in `internal/engine/time_utils.go`
3. Ensure the algorithm provides detailed metadata
4. Add real-world scenario tests in `internal/engine/tests/pattern_criteria/`
5. Update the pattern detection documentation

## Testing Strategy

Our testing strategy combines unit tests and real-world scenario tests:

### Unit Tests

Focused on testing individual components in isolation.

### Real-World Scenario Tests

Located in `internal/engine/tests/pattern_criteria/`, these tests simulate actual user behavior to validate pattern detection:

1. **User Engagement**: Testing consistent daily usage patterns
2. **Fitness Progress**: Testing increasing workout activity
3. **Learning Decline**: Testing decreasing engagement patterns
4. **Seasonal Usage**: Testing varying patterns across time scales
5. **Mixed Patterns**: Testing multiple criteria against the same data

## Debugging Tips

### Pattern Detection Issues

If pattern detection tests are failing:

1. Check the metadata output for detailed information
2. Verify that the test data has the expected pattern
3. Review algorithm parameters like `maxDeviation` or `minIncreasePct`
4. Add print statements to track event grouping and counting
5. Troubleshoot using the test metadata:

```go
t.Logf("Pattern detection metadata: %v", metadata)
```

## Common Challenges and Solutions

### 1. Understanding Time Period Grouping

Events are grouped into periods based on their timestamp and the `periodType` parameter. Understanding how this grouping works is key to debugging pattern detection issues.

### 2. Event Count Variations

Small variations in event counts can sometimes cause unexpected test failures. Be sure to check the `maxDeviation` parameter and understand how outliers are handled.

### 3. Edge Cases

Be aware of these common edge cases:
- Empty event sets
- Insufficient periods
- Periods with zero events
- Large outliers in otherwise consistent patterns

## Resources and References

- Complete API documentation: See [API.md](API.md)
- Badge criteria types: See [README_BADGE_CRITERIA.md](README_BADGE_CRITERIA.md)
- Testing documentation: See [Testing.md](Testing.md)

## Get Help

If you're stuck or have questions:
1. Check the existing documentation
2. Review test cases for examples
3. Look at the implementation code with detailed comments
4. Reach out to the team for support

---

We hope this guide helps you get started with the Badge Assignment System. Welcome aboard! 
