# Badge Assignment System API Documentation

## Overview

The Badge Assignment System provides a flexible, rule-based engine for defining and assigning badges to users based on their activities and behaviors. This document outlines the API endpoints and the JSON schema for badge criteria definition.

## Badge Criteria Types

### Basic Criteria

| Operator | Description | Example |
|----------|-------------|---------|
| `$eq` | Equal | `{"field": "score", "operator": "$eq", "value": 100}` |
| `$ne` | Not Equal | `{"field": "attempts", "operator": "$ne", "value": 0}` |
| `$gt` | Greater Than | `{"field": "points", "operator": "$gt", "value": 50}` |
| `$gte` | Greater Than or Equal | `{"field": "level", "operator": "$gte", "value": 5}` |
| `$lt` | Less Than | `{"field": "errors", "operator": "$lt", "value": 3}` |
| `$lte` | Less Than or Equal | `{"field": "time", "operator": "$lte", "value": 60}` |
| `$in` | In Array | `{"field": "category", "operator": "$in", "value": ["sports", "fitness"]}` |
| `$nin` | Not In Array | `{"field": "tags", "operator": "$nin", "value": ["beginner", "tutorial"]}` |

### Logical Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `$and` | All conditions must be true | `{"$and": [{"field": "score", "operator": "$gt", "value": 80}, {"field": "attempts", "operator": "$lt", "value": 3}]}` |
| `$or` | At least one condition must be true | `{"$or": [{"field": "achievements", "operator": "$gte", "value": 5}, {"field": "level", "operator": "$gte", "value": 10}]}` |
| `$not` | Negates the result of the condition | `{"$not": {"field": "status", "operator": "$eq", "value": "inactive"}}` |

### Time-Based Criteria

#### Time Period Criteria

Evaluates events within specific time periods (days, weeks, months).

```json
{
  "$timePeriod": {
    "periodType": "day",
    "periodCount": { "$gte": 5 },
    "excludeWeekends": true,
    "excludeHolidays": true,
    "holidays": ["2023-12-25", "2024-01-01"]
  }
}
```

#### Event Count Criteria

Counts raw events that match specified criteria.

```json
{
  "event": "workout_completed",
  "criteria": {
    "$eventCount": {
      "$gte": 10
    }
  }
}
```

#### Pattern Criteria

Detects patterns in event occurrences over time periods.

```json
{
  "$timePeriod": {
    "period": "day",
    "periodCount": 7,
    "flow": {
      "event": "workout_completed",
      "criteria": {
        "$pattern": {
          "type": "increasing",
          "min_periods": 5,
          "min_percent_increase": 10
        }
      }
    }
  }
}
```

**Parameters:**
- `period`: Time period to group events by. Options: "day", "week", "month", "quarter", "year"
- `periodCount`: Number of periods to evaluate
- `flow`: Criteria to evaluate for each period
- `$pattern`: Pattern to detect
  - `type`: Pattern type. Options: "consistent", "increasing", "decreasing"
  - `min_periods`: Minimum number of periods with events (default: periodCount)
  - `min_percent_increase`: For "increasing" pattern - minimum percent increase between periods
  - `min_percent_decrease`: For "decreasing" pattern - minimum percent decrease between periods
  - `max_deviation`: For "consistent" pattern - maximum allowed deviation from average

**Pattern Type Descriptions:**
- `consistent`: Events occur at a consistent frequency across periods
- `increasing`: Events show an increasing trend over time
- `decreasing`: Events show a decreasing trend over time

**Example Use Cases:**
- Award "Consistent Learner" badge for users with consistent daily learning activity
- Award "Fitness Growth" badge for users with increasing workout frequency
- Detect decreased engagement pattern for potential user re-engagement

#### Sequence Criteria

Checks if events occur in a specific sequence.

```json
{
  "$sequence": {
    "sequence": ["login", "view_course", "complete_quiz"],
    "maxGapSeconds": 3600,
    "requireStrict": true
  }
}
```

#### Gap Criteria

Checks for gaps in event occurrence.

```json
{
  "$gap": {
    "maxGapHours": 24,
    "minGapHours": 1,
    "periodType": "day"
  }
}
```

#### Duration Criteria

Assesses the time duration between related events.

```json
{
  "$duration": {
    "startEvent": { "type": "start_workout" },
    "endEvent": { "type": "complete_workout" },
    "duration": { "$lte": 30 },
    "unit": "minutes"
  }
}
```

#### Aggregation Criteria

Handles calculations over event values.

```json
{
  "$aggregate": {
    "type": "avg",
    "field": "duration",
    "value": { "$gte": 30 }
  }
}
```

#### Time Window Criteria

Filters events within a specific time window.

```json
{
  "$timeWindow": {
    "start": "2023-01-01T00:00:00Z",
    "end": "2023-01-31T23:59:59Z",
    "flow": {
      "event": "login",
      "criteria": {
        "$eventCount": { "$gte": 10 }
      }
    }
  }
}
```

Relative time windows are also supported:

```json
{
  "$timeWindow": {
    "last": "30d",
    "businessDaysOnly": true,
    "flow": {
      "event": "login",
      "criteria": {
        "$eventCount": { "$gte": 5 }
      }
    }
  }
}
```

**Parameters:**
- `start`: Start datetime (ISO 8601)
- `end`: End datetime (ISO 8601)
- `last`: Relative time window (e.g., "30d", "2w", "1m")
   - Supported units: "d" (days), "w" (weeks), "m" (months), "q" (quarters), "y" (years)
- `businessDaysOnly`: When true, excludes weekends from the time window
- `flow`: The criteria to evaluate within the time window

### Payload Field Criteria

Filters events based on specific values in the event payload data.

#### Direct Equality Comparison

Checks if a field in the event payload exactly matches a specified value.

```json
{
  "payload": {
    "status": "fixed",
    "priority": "high"
  }
}
```

This checks if the `status` field equals "fixed" AND the `priority` field equals "high".

#### Operator Comparison

Uses comparison operators for more flexible matching.

```json
{
  "payload": {
    "time": { "$lt": "09:00:00" },
    "score": { "$gte": 80 },
    "tags": { "$in": ["important", "featured"] }
  }
}
```

#### Supported Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `$eq` | Equal | `{"score": {"$eq": 100}}` |
| `$ne` | Not Equal | `{"status": {"$ne": "cancelled"}}` |
| `$gt` | Greater Than | `{"points": {"$gt": 50}}` |
| `$gte` | Greater Than or Equal | `{"level": {"$gte": 5}}` |
| `$lt` | Less Than | `{"time": {"$lt": "09:00:00"}}` |
| `$lte` | Less Than or Equal | `{"errors": {"$lte": 3}}` |
| `$in` | In Array | `{"category": {"$in": ["sports", "fitness"]}}` |
| `$nin` | Not In Array | `{"tags": {"$nin": ["beginner", "tutorial"]}}` |

#### Example Badge with Payload Criteria

```json
{
  "name": "Early Bird",
  "description": "Checked in before 9 AM for 5 consecutive days",
  "imageUrl": "https://example.com/badges/early-bird.png",
  "flow_definition": {
    "event": "check-in",
    "criteria": {
      "$timePeriod": {
        "periodType": "day",
        "periodCount": { "$gte": 5 },
        "excludeWeekends": false
      },
      "payload": {
        "time": { "$lt": "09:00:00" }
      }
    }
  }
}
```

This badge is awarded when a user checks in before 9:00 AM for at least 5 days.

## API Endpoints

### API Overview

The Badge Assignment System API is divided into two main sections:

1. **Public API Endpoints**: These endpoints are accessible to regular users and provide read access to badges and functionality for submitting events.
   - `GET /api/v1/badges` - List all badges
   - `GET /api/v1/badges/:id` - Get badge details
   - `GET /api/v1/users/:id/badges` - Get user badges
   - `POST /api/v1/events` - Submit events

2. **Admin API Endpoints**: These endpoints are for administrative operations and require appropriate authentication.
   - `/api/v1/admin/badges/*` - Badge management (create, update, delete)
   - `/api/v1/admin/event-types/*` - Event type management

### Badge API

#### Create Badge Definition

**Endpoint:** `POST /api/v1/admin/badges`

**Request Body:**

```json
{
  "name": "Early Bird",
  "description": "Checked in before 9 AM for 5 consecutive days",
  "image_url": "https://example.com/badges/early-bird.png",
  "flow_definition": {
    "event": "check-in",
    "criteria": {
      "$timePeriod": {
        "periodType": "day",
        "periodCount": { "$gte": 5 },
        "excludeWeekends": false
      },
      "payload": {
        "time": { "$lt": "09:00:00" }
      }
    }
  }
}
```

**Response:**

```json
{
  "badge": {
    "id": 123,
    "name": "Early Bird",
    "description": "Checked in before 9 AM for 5 consecutive days",
    "image_url": "https://example.com/badges/early-bird.png",
    "active": true,
    "created_at": "2023-06-15T10:30:00Z",
    "updated_at": "2023-06-15T10:30:00Z"
  },
  "criteria": {
    "id": 456,
    "badge_id": 123,
    "flow_definition": {
      "event": "check-in",
      "criteria": {
        "$timePeriod": {
          "periodType": "day",
          "periodCount": { "$gte": 5 },
          "excludeWeekends": false
        },
        "payload": {
          "time": { "$lt": "09:00:00" }
        }
      }
    },
    "created_at": "2023-06-15T10:30:00Z",
    "updated_at": "2023-06-15T10:30:00Z"
  }
}
```

#### Update Badge Definition

**Endpoint:** `PUT /api/v1/admin/badges/{badge_id}`

**Request Body:** Same as Create Badge Definition

**Response:** Same as Create Badge Definition

#### List Badges

**Endpoint:** `GET /api/v1/badges`

**Query Parameters:**
- `limit`: Maximum number of badges to return (default: 20)
- `offset`: Pagination offset (default: 0)
- `active`: Filter by active status (true, false, all)

**Response:**
```json
[
  {
    "id": 123,
    "name": "Early Bird",
    "description": "Checked in before 9 AM for 5 consecutive days",
    "image_url": "https://example.com/badges/early-bird.png",
    "active": true,
    "created_at": "2023-06-15T10:30:00Z",
    "updated_at": "2023-06-15T10:30:00Z"
  },
  {
    "id": 124,
    "name": "Workaholic",
    "description": "Logged 40+ hours in a week",
    "image_url": "https://example.com/badges/workaholic.png",
    "active": true,
    "created_at": "2023-06-16T14:20:00Z",
    "updated_at": "2023-06-16T14:20:00Z"
  }
]
```

#### Get Badge Details

**Endpoint:** `GET /api/v1/badges/{badge_id}`

**Response:**
```json
{
  "id": 123,
  "name": "Early Bird",
  "description": "Checked in before 9 AM for 5 consecutive days",
  "image_url": "https://example.com/badges/early-bird.png",
  "active": true,
  "created_at": "2023-06-15T10:30:00Z",
  "updated_at": "2023-06-15T10:30:00Z"
}
```

#### Get Badge with Criteria

**Endpoint:** `GET /api/v1/admin/badges/{badge_id}/criteria`

**Response:** Same as the response from Create Badge Definition

#### Delete Badge

**Endpoint:** `DELETE /api/v1/admin/badges/{badge_id}`

**Response:** HTTP 200 OK

### Event Type API

#### Create Event Type

**Endpoint:** `POST /api/v1/admin/event-types`

**Request Body:**
```json
{
  "name": "Check In",
  "description": "User check-in event, used for attendance tracking",
  "schema": {
    "type": "object",
    "properties": {
      "user_id": {
        "type": "string"
      },
      "time": {
        "type": "string",
        "format": "time",
        "description": "The time of check-in in HH:MM:SS format"
      },
      "date": {
        "type": "string",
        "format": "date"
      },
      "location": {
        "type": "string"
      }
    },
    "required": ["user_id", "time", "date"]
  }
}
```

**Response:**
```json
{
  "id": 1,
  "name": "Check In",
  "description": "User check-in event, used for attendance tracking",
  "schema": {
    "type": "object",
    "properties": {
      "user_id": {
        "type": "string"
      },
      "time": {
        "type": "string",
        "format": "time",
        "description": "The time of check-in in HH:MM:SS format"
      },
      "date": {
        "type": "string",
        "format": "date"
      },
      "location": {
        "type": "string"
      }
    },
    "required": ["user_id", "time", "date"]
  },
  "created_at": "2023-06-14T09:00:00Z",
  "updated_at": "2023-06-14T09:00:00Z"
}
```

#### Update Event Type

**Endpoint:** `PUT /api/v1/admin/event-types/{event_type_id}`

**Request Body:**
```json
{
  "name": "Check In",
  "description": "Updated description for check-in events",
  "schema": {
    "type": "object",
    "properties": {
      "user_id": {
        "type": "string"
      },
      "time": {
        "type": "string",
        "format": "time"
      },
      "date": {
        "type": "string",
        "format": "date"
      },
      "location": {
        "type": "string"
      },
      "device_id": {
        "type": "string",
        "description": "ID of the device used for check-in"
      }
    },
    "required": ["user_id", "time", "date"]
  }
}
```

**Response:** Updated event type object

#### List Event Types

**Endpoint:** `GET /api/v1/admin/event-types`

**Response:**
```json
[
  {
    "id": 1,
    "name": "Check In",
    "description": "User check-in event, used for attendance tracking",
    "schema": { "..." },
    "created_at": "2023-06-14T09:00:00Z",
    "updated_at": "2023-06-14T09:00:00Z"
  },
  {
    "id": 2,
    "name": "Check Out",
    "description": "User check-out event, used for attendance tracking",
    "schema": { "..." },
    "created_at": "2023-06-14T09:05:00Z",
    "updated_at": "2023-06-14T09:05:00Z"
  }
]
```

#### Get Event Type Details

**Endpoint:** `GET /api/v1/admin/event-types/{event_type_id}`

**Response:** Event type object

#### Delete Event Type

**Endpoint:** `DELETE /api/v1/admin/event-types/{event_type_id}`

**Response:** HTTP 200 OK

### Event API

#### Create Event

**Endpoint:** `POST /api/v1/events`

**Request Body:**
```json
{
  "event_type": "check-in",
  "user_id": "user123",
  "payload": {
    "time": "08:45:00",
    "date": "2023-06-20",
    "location": "Main Office"
  },
  "timestamp": "2023-06-20T08:45:00Z"
}
```

**Response:**
```json
{
"message": "Event processed successfully"
}
```

#### Get User Events

**Endpoint:** `GET /api/v1/users/{user_id}/events`

**Query Parameters:**
- `event_type`: Filter by event type ID or name
- `start_date`: Filter events after this date (ISO 8601 format)
- `end_date`: Filter events before this date (ISO 8601 format)
- `limit`: Maximum number of events to return (default: 50)
- `offset`: Pagination offset (default: 0)

**Response:**
```json
[
  {
    "id": 500,
    "event_type_id": 1,
    "event_type_name": "Check In",
    "user_id": "user123",
    "payload": {
      "time": "08:45:00",
      "date": "2023-06-20",
      "location": "Main Office"
    },
    "occurred_at": "2023-06-20T08:45:00Z"
  },
  {
    "id": 501,
    "event_type_id": 2,
    "event_type_name": "Check Out",
    "user_id": "user123",
    "payload": {
      "time": "17:30:00",
      "date": "2023-06-20",
      "location": "Main Office"
    },
    "occurred_at": "2023-06-20T17:30:00Z"
  }
]
```

### User Badge API

#### Evaluate User for Badges (Not Implemented)

> **Note:** This endpoint is documented as a planned feature but has not been implemented in the current API.

**Endpoint:** `POST /api/v1/users/{user_id}/evaluate-badges`

**Response:**
```json
{
  "awarded_badges": [
    {
      "badge_id": 123,
      "name": "Early Bird",
      "awarded_at": "2023-06-20T08:50:00Z"
    }
  ],
  "progress": [
    {
      "badge_id": 124,
      "name": "Workaholic",
      "progress": 0.6,
      "requirements": {
        "completed": 24,
        "total": 40
      }
    }
  ]
}
```

#### Get User Badges

**Endpoint:** `GET /api/v1/users/{user_id}/badges`

**Query Parameters:**
- `limit`: Maximum number of badges to return (default: 20)
- `offset`: Pagination offset (default: 0)

**Response:**
```json
[
  {
    "id": 1001,
    "badge_id": 123,
    "user_id": "user123",
    "awarded_at": "2023-06-20T08:50:00Z",
    "badge": {
      "name": "Early Bird",
      "description": "Checked in before 9 AM for 5 consecutive days",
      "image_url": "https://example.com/badges/early-bird.png"
    }
  }
]
```

## Badge and Event Type Relationship

When creating badges, it's important to ensure that the event types referenced in the badge criteria exist in the system. Each badge's flow definition references one or more event types and specifies criteria that events of those types must meet for the badge to be awarded.

### Prerequisites

1. **Create Event Types First**: Before creating badges, ensure that all required event types exist
2. **Match Schema Fields**: Any payload fields referenced in badge criteria must exist in the corresponding event type's schema
3. **Type Compatibility**: The data types in badge criteria must be compatible with the schema data types

### Example Workflow

1. Create an event type:
   ```
   POST /api/v1/admin/event-types
   {
     "name": "Check In",
     "description": "User check-in event",
     "schema": {
       "type": "object",
       "properties": {
         "time": {
           "type": "string",
           "format": "time"
         }
       },
       "required": ["time"]
     }
   }
   ```

2. Create a badge that uses this event type:
   ```
   POST /api/v1/admin/badges
   {
     "name": "Early Bird",
     "description": "Checked in before 9 AM for 5 consecutive days",
     "image_url": "https://example.com/badges/early-bird.png",
     "flow_definition": {
       "event": "Check In",
       "criteria": {
         "payload": {
           "time": { "$lt": "09:00:00" }
         }
       }
     }
   }
   ```

3. Submit events that might trigger the badge:
   ```
   POST /api/v1/events
   {
     "event_type": "check-in",
     "user_id": "user123",
     "payload": {
       "time": "08:45:00"
     },
     "timestamp": "2023-06-20T08:45:00Z"
   }
   ```

4. Check if the user was awarded the badge (Note: Manual evaluation endpoint is not implemented):
   ```
   GET /api/v1/users/user123/badges
   ```

## Error Responses

All API endpoints return standard HTTP status codes:

- `200 OK`: The request was successful
- `201 Created`: The resource was created successfully
- `400 Bad Request`: The request was invalid
- `401 Unauthorized`: Authentication is required
- `403 Forbidden`: The client does not have permission
- `404 Not Found`: The resource was not found
- `500 Internal Server Error`: An unexpected error occurred

Error responses include a JSON body with details:

```json
{
  "error": {
    "code": "invalid_criteria",
    "message": "Badge criteria contains invalid operator",
    "details": {
      "field": "criteria.type",
      "value": "unknown_type"
    }
  }
}
``` 
