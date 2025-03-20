# Badge Criteria Format Specification

This document provides a comprehensive guide to the badge criteria format used in the Badge Assignment System. It serves as the authoritative reference for all badge criteria structures, validation requirements, and best practices.

## Table of Contents

1. [Basic Structure](#basic-structure)
2. [Criteria Types](#criteria-types)
   - [Comparison Operators](#comparison-operators)
   - [Pattern Criteria](#pattern-criteria)
   - [Time-Based Criteria](#time-based-criteria)
   - [Logical Operators](#logical-operators)
3. [Type Requirements](#type-requirements)
   - [API Calls vs. Go Code](#api-calls-vs-go-code)
4. [Common Templates](#common-templates)
5. [Complete Badge Assignment Flow](#complete-badge-assignment-flow)
   - [Flow Diagram](#flow-diagram)
   - [API Request/Response Examples](#api-requestresponse-examples)
6. [Validation](#validation)
7. [Troubleshooting](#troubleshooting)
8. [Legacy vs. Current Formats](#legacy-vs-current-formats)
9. [Visual Diagrams](#visual-diagrams)

For a visual representation of the concepts in this document, see the [Badge Criteria Diagrams](BADGE_CRITERIA_DIAGRAM.md).

## Basic Structure

Badge criteria follow this general structure:

```json
{
  "event": "event_type_name",
  "criteria": {
    "field_name": {
      "$operator": value
    }
  }
}
```

- The `event` field must match the name of an existing event type in the system
- The `criteria` field contains a map of field names to conditions
- Conditions can use direct comparison or operators

## Criteria Types

### Comparison Operators

Comparison operators are used directly without any wrapper:

```json
"score": {
  "$gte": 90
}
```

The system supports the following comparison operators:

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

### Pattern Criteria

Pattern criteria require a `$pattern` wrapper:

```json
"$pattern": {
  "pattern": "consistent",
  "periodType": "day",
  "minPeriods": 7,
  "maxDeviation": 0.15
}
```

Pattern types include:

1. **Consistent Pattern**: Looks for consistent activity levels across periods
   ```json
   "$pattern": {
     "pattern": "consistent",
     "periodType": "day",
     "minPeriods": 7,
     "maxDeviation": 0.15
   }
   ```

2. **Increasing Pattern**: Looks for a growing trend in activity
   ```json
   "$pattern": {
     "pattern": "increasing",
     "periodType": "week",
     "minPeriods": 4,
     "minIncreasePct": 10.0
   }
   ```

3. **Decreasing Pattern**: Looks for a declining trend in activity
   ```json
   "$pattern": {
     "pattern": "decreasing",
     "periodType": "week",
     "minPeriods": 4,
     "maxDecreasePct": 15.0
   }
   ```

### Time-Based Criteria

Time-based criteria have specialized formats:

1. **Time Period**: Counts activity across time periods
   ```json
   "$timePeriod": {
     "periodType": "day",
     "periodCount": { "$gte": 5 },
     "excludeWeekends": true,
     "excludeHolidays": true
   }
   ```

2. **Time Window**: Define a specific time window for activities
   ```json
   "$timeWindow": {
     "start": "2023-01-01T00:00:00Z",
     "end": "2023-01-31T23:59:59Z",
     "flow": {
       "event": "test_event",
       "criteria": {
         "score": { "$gte": 80 }
       }
     }
   }
   ```

3. **Duration**: Measures time between events
   ```json
   "$duration": {
     "startEvent": {
       "type": "session_start"
     },
     "endEvent": {
       "type": "session_end"
     },
     "matchProperty": "session_id",
     "duration": { "$gte": 30 },
     "unit": "minutes"
   }
   ```

### Logical Operators

Logical operators allow combining multiple criteria:

1. **AND Operator**: All conditions must be true
   ```json
   "$and": [
     {
       "event": "test_event",
       "criteria": {
         "score": { "$gte": 80 }
       }
     },
     {
       "event": "test_event",
       "criteria": {
         "completed": true
       }
     }
   ]
   ```

2. **OR Operator**: At least one condition must be true
   ```json
   "$or": [
     {
       "event": "test_event",
       "criteria": {
         "score": { "$gte": 90 }
       }
     },
     {
       "event": "test_event",
       "criteria": {
         "level": "expert"
       }
     }
   ]
   ```

3. **NOT Operator**: Condition must be false
   ```json
   "$not": {
     "event": "test_event",
     "criteria": {
       "completed": false
     }
   }
   ```

## Type Requirements

### API Calls vs. Go Code

**Important Distinction**:

- **For JSON API Calls**: When sending criteria via JSON in API calls, you don't need to explicitly handle type conversion. The JSON specification naturally defines numbers, and Go's JSON decoder will automatically unmarshal these values as `float64` when deserializing into a `map[string]interface{}`.

  ```json
  {
    "criteria": {
      "score": {
        "$gte": 90
      }
    }
  }
  ```

- **For Go Code**: When writing Go code directly (such as in tests or when building criteria programmatically), you must use explicit `float64()` conversion. This is because Go's compiler interprets numeric literals as `int` by default.

  ```go
  // In Go code, you need the explicit conversion
  criteria := map[string]interface{}{
      "score": map[string]interface{}{
          "$gte": float64(90),  // Required in Go code
      },
  }
  ```

**Why This Matters**: 
The rule engine expects numeric values to be of type `float64`. Without explicit conversion in Go code, the engine may receive an `int` instead of a `float64`, which can cause type assertion errors or unexpected behavior. The validation script specifically checks Go source code, not JSON data.

### Numeric Values in Go Code

In Go code, all numeric values must use explicit `float64()` type conversion:

```go
// Correct
"$gte": float64(90)

// Incorrect - will cause validation errors
"$gte": 90
```

String values and boolean values don't require special type conversion:

```go
"status": "completed"
"active": true
```

## Common Templates

### Simple Comparison

```go
criteria := map[string]interface{}{
    "event": "test_event",
    "criteria": map[string]interface{}{
        "score": map[string]interface{}{
            "$gte": float64(90),
        },
    },
}
```

### Pattern Detection

```go
criteria := map[string]interface{}{
    "$pattern": map[string]interface{}{
        "pattern":      "consistent",
        "periodType":   "day",
        "minPeriods":   float64(7),
        "maxDeviation": float64(0.15),
    },
}
```

### Logical Combination (AND)

```go
criteria := map[string]interface{}{
    "$and": []interface{}{
        map[string]interface{}{
            "event": "test_event",
            "criteria": map[string]interface{}{
                "score": map[string]interface{}{
                    "$gte": float64(80),
                },
            },
        },
        map[string]interface{}{
            "event": "test_event",
            "criteria": map[string]interface{}{
                "completed": true,
            },
        },
    },
}
```

### Time Window

```go
criteria := map[string]interface{}{
    "$timeWindow": map[string]interface{}{
        "start": "2023-01-01T00:00:00Z",
        "end":   "2023-01-31T23:59:59Z",
        "flow": map[string]interface{}{
            "event": "test_event",
            "criteria": map[string]interface{}{
                "score": map[string]interface{}{
                    "$gte": float64(80),
                },
            },
        },
    },
}
```

## Complete Badge Assignment Flow

This section demonstrates the complete flow of creating event types, badges, submitting events, and verifying badge assignments.

### Flow Diagram

```
+-------------------------+      +-------------------------+
|  1. Create Event Type   |      |      Event Schema      |
|  (Name, Description,    +----->+  (Field Definitions,   |
|   Schema)               |      |   Validation Rules)    |
+----------+--------------+      +-------------------------+
           |
           v
+----------+--------------+      +-------------------------+
|  2. Create Badge        |      |    Flow Definition     |
|  (Name, Description,    +----->+  (Event Name, Criteria |
|   Flow Definition)      |      |   with Value Checks)   |
+----------+--------------+      +-------------------------+
           |
           v
+----------+--------------+
|  3. Submit Event        |
|  (Event Type, User ID,  |
|   Data, Payload)        |
+----------+--------------+
           |
           v
+----------+--------------+      +-------------------------+
|  4. Badge Assignment    |      |      Badge Award       |
|  (Automatic based on    +----->+  (User ID, Badge ID,   |
|   criteria match)       |      |   Award Timestamp)     |
+----------+--------------+      +-------------------------+
           |
           v
+----------+--------------+
|  5. Verify Assignment   |
|  (API or Database       |
|   Query)                |
+-------------------------+
```

### API Request/Response Examples

Below are complete API request and response examples for each step in the badge assignment flow, based on our integration testing.

#### 1. Creating an Event Type

**Request:**
```json
POST /api/v1/admin/event-types
{
  "name": "simple_event_1742134382840",
  "description": "Simple test event",
  "schema": {
    "type": "object",
    "properties": {
      "value": {
        "type": "number"
      }
    }
  }
}
```

**Response:**
```json
{
  "id": 50,
  "name": "simple_event_1742134382840",
  "description": "Simple test event",
  "schema": {
    "type": "object",
    "properties": {
      "value": {
        "type": "number"
      }
    }
  }
}
```

#### 2. Creating a Badge with Criteria

**Request:**
```json
POST /api/v1/admin/badges
{
  "name": "Simple Badge_1742134382840",
  "description": "Simple test badge",
  "image_url": "https://example.com/badge.png",
  "flow_definition": {
    "event": "simple_event_1742134382840",
    "criteria": {
      "value": {
        "$gte": 50
      }
    }
  },
  "is_active": true
}
```

**Response:**
```json
{
  "badge": {
    "id": 106,
    "name": "Simple Badge_1742134382840",
    "description": "Simple test badge",
    "image_url": "https://example.com/badge.png",
    "active": true,
    "created_at": "2025-03-16T19:43:02.844034Z",
    "updated_at": "2025-03-16T19:43:02.844034Z"
  },
  "criteria": {
    "id": 106,
    "badge_id": 106,
    "flow_definition": {
      "criteria": {
        "value": {
          "$gte": 50
        }
      },
      "event": "simple_event_1742134382840"
    },
    "created_at": "2025-03-16T19:43:02.844034Z",
    "updated_at": "2025-03-16T19:43:02.844034Z"
  }
}
```

#### 3. Submitting an Event

**Request:**
```json
POST /api/v1/events
{
  "event_type": "simple_event_1742134382840",
  "user_id": "test_user_1",
  "timestamp": "2025-03-16T19:43:02.850000Z",
  "payload": {
    "value": 75
  }
}
```

**Response:**
```json
{
  "success": true,
  "event_id": 500
}
```

#### 4. Retrieving User Badges

**Request:**
```
GET /api/v1/users/test_user_1/badges
```

**Response:**
```json
[
  {
    "awarded_at": "2025-03-16T19:43:02.863425Z",
    "description": "Simple test badge",
    "id": 106,
    "image_url": "https://example.com/badge.png",
    "metadata": "eyJsYXN0X2V2ZW50X2lkIjogNTAwLCAiZmlyc3RfZXZlbnRfaWQiOiA1MDAsICJmaWx0ZXJlZF9ldmVudF9jb3VudCI6IDF9",
    "name": "Simple Badge_1742134382840"
  }
]
```

### Important API Field Requirements

Based on integration testing, here are crucial requirements for API requests:

#### Badge Creation Requirements

1. **Use `flow_definition` not `criteria`**
   - Badge creation requires a `flow_definition` field, not a direct `criteria` field
   - The criteria is nested inside the flow_definition

   ```json
   // CORRECT
   {
     "flow_definition": {
       "event": "event_name",
       "criteria": {
         "field": { "$operator": value }
       }
     }
   }
   
   // INCORRECT
   {
     "criteria": {
       "event": "event_name",
       "field": { "$operator": value }
     }
   }
   ```

2. **Required Badge Fields**
   - `name`: String - Unique badge name
   - `description`: String - Badge description
   - `image_url`: String - URL to badge image
   - `flow_definition`: Object - Contains event name and criteria
   - `is_active`: Boolean - Whether the badge is active

#### Event Submission Requirements

1. **Fields are required**

   ```json
   {
     "event_type": "example_event",
     "user_id": "user123",
     "timestamp": "2023-01-01T12:00:00Z",
     "payload": {
       "value": 75
     }
   }
   ```

2. **Required Event Fields**
   - `event_type`: String - Must match an existing event type name
   - `user_id`: String - ID of the user performing the event
   - `timestamp`: String - ISO 8601 format timestamp (optional, defaults to current time)
   - `payload`: Object - Event data for storage and criteria evaluation

## Validation

### Validation Rules

Our validation script in `tools/validate_test_criteria.sh` checks for:

1. Pattern criteria must use `$pattern` wrappers
2. All numeric values must use explicit `float64()` type conversion
3. Criteria structure must follow the specified format

### Running Validation

To validate your badge criteria in test files:

```bash
./tools/validate_test_criteria.sh
```

### Common Validation Errors

| Error | Description | Solution |
|-------|-------------|----------|
| `pattern criteria without $pattern wrapper` | Pattern criteria missing the `$pattern` wrapper | Add the `$pattern` wrapper |
| `numeric values without explicit type conversion` | Numeric values without `float64()` | Add explicit `float64()` conversion |
| `potential inconsistent criteria nesting` | Possible issues with criteria structure | Check criteria structure against templates |
| `flow definition is required` | Missing flow_definition in badge creation | Use flow_definition field for badge criteria |
| `event type is required` | Missing or invalid event_type in event submission | Ensure event_type field is present and valid |

### False Positives

Some validation warnings may be false positives in special cases:

- Legacy code that follows older formats
- Intentional nesting that deviates from standard formats
- Files explicitly excluded from validation

If the validation script identifies files you believe are correctly formatted, check the script documentation for exclusion options.

## Troubleshooting

### Common Issues

1. **Criteria Not Evaluating as Expected**
   - Check that event type names match exactly (case-sensitive)
   - Verify all numeric values use `float64()` type conversion
   - Ensure pattern criteria use the `$pattern` wrapper

2. **Type Conversion Errors**
   - JavaScript/JSON numbers are automatically parsed as `float64` in Go
   - Explicit conversion is still required in Go code for consistency

3. **Missing Fields**
   - Verify all required fields are present for each criteria type
   - Check for typos in field names (field names are case-sensitive)
   - Badge creation requires `flow_definition` field, not a direct `criteria` field
   - Event submission requires the `payload` field

4. **Logical Operator Issues**
   - Verify each item in a logical operator array is a complete, valid criteria
   - Check that arrays are properly formatted as `[]interface{}`

5. **Badge Creation Fails**
   - Ensure you're using `flow_definition` not `criteria` at the top level
   - Verify your event name matches an existing event type exactly

6. **Event Submission Fails**
   - Ensure the `payload` field is included
   - Check that the event_type matches an existing event type exactly
   - Verify timestamp format is correct (ISO 8601)

## Legacy vs. Current Formats

### Current Format (v2.x+)

All badge criteria should follow the specifications outlined in this document. Key points:

- Pattern criteria require a `$pattern` wrapper
- Comparison operators are used directly without wrappers
- All numeric values require explicit `float64()` type conversion
- Event types must match existing event type names exactly
- Badge creation uses `flow_definition` field to contain criteria
- Event submission requires the `payload` field

### Legacy Format (v1.x)

Legacy formats use different structures that are maintained for backward compatibility but should not be used for new development. Notable differences:

- Some legacy tests may not use `$pattern` wrappers
- Legacy tests may use direct numeric values without type conversion
- Some legacy formats may use different field names

### Handling Legacy Code

When working with legacy code:

- Do not modify the format of existing criteria unless explicitly updating to the new format
- Use the validation script with appropriate exclusions for legacy files
- When updating legacy code, convert to the current format

## Badge Criteria Structure Diagram

```
                +----------------+
                |  Badge Criteria |
                +----------------+
                        |
        +---------------+---------------+
        |                               |
+----------------+              +-----------------+
| Event Criteria |              | Logical Operators |
+----------------+              +-----------------+
        |                               |
+----------------+              +-----------------+
| Comparison Ops |              | $and, $or, $not |
+----------------+              +-----------------+
        |                               |
        |                     +---------+---------+
        |                     |                   |
+----------------+    +----------------+  +----------------+
| Direct Fields  |    | Pattern Criteria |  | Time Criteria |
+----------------+    +----------------+  +----------------+
```

## Appendix: Default Values

When certain optional parameters are omitted, the system uses these defaults:

### Pattern Criteria Defaults
- `minPeriods`: 3 periods
- `minIncreasePct`: 5.0% (for increasing pattern)
- `maxDecreasePct`: 5.0% (for decreasing pattern)
- `maxDeviation`: 0.2 (20% deviation for consistent pattern)

### Time Period Criteria Defaults
- If no `periodCount` is specified, the criterion is met if there's at least one period with activity
- `excludeWeekends` and `excludeHolidays` default to `false` 
