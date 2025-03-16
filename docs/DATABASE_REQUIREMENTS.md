# Database Tables Requirements Documentation

## Overview

This document provides details about the database tables used in the Badge Assignment System, their dependencies, and which tables are mandatory for the application to function correctly.

## Database Schema

The Badge Assignment System uses the following tables:

1. **event_types** - Defines event types that can be processed by the system
2. **condition_types** - Defines reusable conditions for badge criteria
3. **badges** - Stores information about available badges
4. **badge_criteria** - Defines criteria for badge awards using the rule engine
5. **user_badges** - Records which badges have been awarded to users
6. **events** - Stores event data submitted to the system

## Mandatory vs. Optional Tables

### Mandatory Tables

The following tables are **required** for the core functionality of the Badge Assignment System:

1. **event_types** - Required for processing events; each event must have a valid event type
2. **badges** - Required for the badge assignment functionality
3. **badge_criteria** - Required for defining when badges should be awarded
4. **events** - Required for storing event data that triggers badge awards
5. **user_badges** - Required for tracking which badges have been awarded to users

### Optional Tables

The following table is **optional** and not required for core functionality:

1. **condition_types** - Used for creating reusable conditions, but not essential for the badge assignment process

## PostgreSQL Version Requirements

The system requires **PostgreSQL 12 or later** due to the following features:
- Advanced JSONB functionality for storing and querying dynamic badge criteria
- JSONB indexing capabilities for performance optimization
- Transaction handling for atomic badge awards
- Support for complex JOIN operations used in the rule engine

If using an earlier PostgreSQL version, you may encounter issues with:
- JSONB query operators not working as expected
- Performance degradation when querying JSONB fields
- Unable to use some of the more advanced query patterns

## Table Dependencies

The database tables have the following dependencies:

1. **badge_criteria** → **badges** (Foreign key: `badge_id`)
   - Each badge criterion must be associated with a valid badge
   - If a badge is deleted, all associated criteria are automatically deleted (CASCADE)

2. **user_badges** → **badges** (Foreign key: `badge_id`)
   - Each user badge must reference a valid badge
   - If a badge is deleted, all associated user badges are automatically deleted (CASCADE)

3. **events** → **event_types** (Foreign key: `event_type_id`)
   - Each event must reference a valid event type
   - If an event type is deleted, the reference in events is set to NULL (SET NULL)

## Minimal Configuration

To run the Badge Assignment System with minimal configuration:

1. Create all mandatory tables:
   - event_types
   - badges
   - badge_criteria
   - user_badges
   - events

2. Ensure the following minimal data is present:
   - At least one event_type record for processing events
   - At least one badge with associated criteria for the rule engine to evaluate

3. The condition_types table can be empty if you're not using reusable conditions.

## Table Relationships Diagram

```
+--------------+       +----------------+       +-----------------+
| event_types  |<------| events         |       | badges          |
+--------------+       +----------------+       +-----------------+
      ^                       |                        ^
      |                       |                        |
      |                       |                        |
      |                       v                        |
+----------------+      +----------------+      +----------------+
| condition_types|      | rule_engine    |----->| badge_criteria |
+----------------+      | (evaluations)  |      +----------------+
                        +----------------+             |
                               |                       |
                               v                       v
                        +------------------------------------+
                        |           user_badges              |
                        +------------------------------------+
```

## Performance Considerations

1. All tables include appropriate indexes to improve query performance.
2. The events table may grow rapidly in a production environment. Consider implementing a data retention policy or archiving strategy for historical events.
3. For high-traffic applications, consider additional indexes on frequently queried fields in the events table.

### Database Indexes

The following indexes are created automatically during schema migration:

1. **events_user_id_idx**: `CREATE INDEX idx_events_user_id ON events(user_id)`
   - Improves performance when looking up events for a specific user

2. **events_event_type_id_idx**: `CREATE INDEX idx_events_event_type_id ON events(event_type_id)`
   - Improves performance when filtering events by type

3. **events_occurred_at_idx**: `CREATE INDEX idx_events_occurred_at ON events(occurred_at)`
   - Improves performance for time-based queries and filtering events by date ranges

4. **user_badges_user_id_idx**: `CREATE INDEX idx_user_badges_user_id ON user_badges(user_id)`
   - Improves performance when retrieving badges for a specific user

5. **user_badges_badge_id_idx**: `CREATE INDEX idx_user_badges_badge_id ON user_badges(badge_id)`
   - Improves performance when looking up users who have earned a specific badge

6. **badge_criteria_badge_id_idx**: `CREATE INDEX idx_badge_criteria_badge_id ON badge_criteria(badge_id)`
   - Improves performance when retrieving criteria for a specific badge

### JSONB Performance Considerations

For tables with JSONB fields (event_types, badge_criteria, user_badges, events):

1. **Query Optimization**:
   - Use specific path queries rather than full document scans
   - Example: `payload->>'field'` instead of scanning the entire payload

2. **Index Considerations**:
   - For frequently queried JSONB paths, consider adding GIN indexes:
     ```sql
     CREATE INDEX idx_events_payload ON events USING GIN (payload jsonb_path_ops);
     ```

3. **Size Limitations**:
   - While PostgreSQL can handle large JSONB documents, try to keep payload sizes reasonable
   - Consider extracting frequently queried fields to separate columns if they're used in WHERE clauses

### Uniqueness Constraints

While not explicitly enforced by the database schema, the application enforces the following uniqueness rules:

1. **Event Type Names**: Each event type should have a unique name (enforced by application code)
2. **Condition Type Names**: Each condition type should have a unique name (enforced by application code)
3. **Badge Names**: Badge names should be unique (enforced by application code)
4. **User Badges**: A user should not receive the same badge twice (enforced by application code)

For production deployments, consider adding explicit uniqueness constraints:

```sql
ALTER TABLE event_types ADD CONSTRAINT event_types_name_unique UNIQUE (name);
ALTER TABLE condition_types ADD CONSTRAINT condition_types_name_unique UNIQUE (name);
ALTER TABLE badges ADD CONSTRAINT badges_name_unique UNIQUE (name);
ALTER TABLE user_badges ADD CONSTRAINT user_badges_user_badge_unique UNIQUE (user_id, badge_id);
```

## Database Integrity

The database schema enforces integrity through:

1. **Primary Keys**: Each table has a unique identifier field
2. **Foreign Keys**: Proper references between tables with appropriate actions (CASCADE, SET NULL)
3. **NOT NULL Constraints**: On critical fields like names and required relationships
4. **Default Values**: For creation timestamps and other fields with sensible defaults

## Schema Evolution and Migrations

The system uses database migrations to manage schema changes. When upgrading:

1. **Migration Tool**: Use the golang-migrate tool for applying migrations:
   ```
   migrate -database "postgres://postgres:postgres@localhost:5432/badge_system?sslmode=disable" -path db/migrations up
   ```

2. **Backward Compatibility**: New migrations should maintain backward compatibility with existing data
   - Avoid dropping columns that may contain important data
   - When adding required fields, provide default values
   - When changing field types, ensure data can be properly converted

3. **Migration Testing**: Always test migrations on a copy of production data before applying to production

4. **Rollback Plan**: Each migration should have a corresponding down migration to allow for rollbacks

## Table Details

### event_types

```sql
CREATE TABLE event_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    schema JSONB,             -- JSON schema defining the expected event payload structure
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### condition_types (Optional)

```sql
CREATE TABLE condition_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    evaluation_logic TEXT,    -- Reference or JSON definition for dynamic mapping
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### badges

```sql
CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### badge_criteria

```sql
CREATE TABLE badge_criteria (
    id SERIAL PRIMARY KEY,
    badge_id INTEGER REFERENCES badges(id) ON DELETE CASCADE,
    flow_definition JSONB NOT NULL,  -- Dynamic rule definition using JSON operators
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### user_badges

```sql
CREATE TABLE user_badges (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,   -- External user identifier
    badge_id INTEGER REFERENCES badges(id) ON DELETE CASCADE,
    awarded_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB                   -- Additional details (e.g., event IDs, evaluation metrics)
);
```

### events

```sql
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    event_type_id INTEGER REFERENCES event_types(id) ON DELETE SET NULL,
    user_id VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,          -- The raw event data
    occurred_at TIMESTAMP DEFAULT NOW()
);
```

## Conclusion

Understanding the database requirements is essential for maintaining and scaling the Badge Assignment System. While the condition_types table is optional, the other five tables (event_types, badges, badge_criteria, user_badges, and events) form the core of the application and must be present for proper functionality. 
