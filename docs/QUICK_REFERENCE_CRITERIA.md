# Badge Criteria Quick Reference Guide

This is a concise reference for the most common badge criteria formats and requirements.

For detailed visual diagrams, see [Badge Criteria Diagrams](BADGE_CRITERIA_DIAGRAM.md).

## JSON vs. Go Code

**Important:** The `float64()` conversion requirement only applies to Go code (tests, etc.), not to JSON API calls.

- **In JSON API calls:** Numbers are handled automatically
  ```json
  { "$gte": 90 }  /* This is fine in JSON */
  ```

- **In Go code:** Always use explicit `float64()` conversion
  ```go
  "$gte": float64(90)  // Required in Go code
  ```

## Most Common Formats

### Basic Criteria (Direct Field Comparison)

```go
criteria := map[string]interface{}{
    "event": "test_event",
    "criteria": map[string]interface{}{
        "score": map[string]interface{}{
            "$gte": float64(90), // ALWAYS use float64() for numbers
        },
    },
}
```

### Pattern Criteria (MUST use $pattern wrapper)

```go
criteria := map[string]interface{}{
    "$pattern": map[string]interface{}{ // $pattern wrapper is REQUIRED
        "pattern":      "consistent",
        "periodType":   "day",
        "minPeriods":   float64(7),     // ALWAYS use float64() for numbers
        "maxDeviation": float64(0.15),  // ALWAYS use float64() for numbers
    },
}
```

### Logical AND Operator

```go
criteria := map[string]interface{}{
    "$and": []interface{}{ // Use []interface{} for arrays
        map[string]interface{}{
            "event": "test_event",
            "criteria": map[string]interface{}{
                "score": map[string]interface{}{
                    "$gte": float64(80), // ALWAYS use float64() for numbers
                },
            },
        },
        map[string]interface{}{
            "event": "test_event",
            "criteria": map[string]interface{}{
                "completed": true, // Booleans don't need conversion
            },
        },
    },
}
```

## Complete Badge Assignment Flow

The badge assignment process involves these key steps:

1. Create an event type (schema defines data structure)
2. Create a badge with criteria (using flow_definition)
3. Submit events (meeting badge criteria)
4. System automatically awards badges
5. Query awarded badges via API or database

This diagram summarizes the flow:

```
Event Type → Badge Creation → Event Submission → Badge Assignment → Verification
```

## API Request Structure

### Creating a Badge (flow_definition is REQUIRED)

```json
POST /api/v1/admin/badges
{
  "name": "Example Badge",
  "description": "A badge for achievement",
  "image_url": "https://example.com/badge.png",
  "flow_definition": {  // MUST use flow_definition, not criteria at top level
    "event": "test_event",
    "criteria": {
      "value": {
        "$gte": 50
      }
    }
  },
  "is_active": true
}
```

### Submitting an Event (payload field is REQUIRED)

```json
POST /api/v1/events
{
  "event_type": "test_event",  // Must match an existing event type
  "user_id": "user123",
  "timestamp": "2023-01-01T12:00:00Z",
  "payload": {  // REQUIRED for database storage and criteria evaluation
    "value": 75
  }
}
```

## Critical Requirements

1. **ALWAYS use `float64()` for numeric values in Go code (not needed in JSON)**
   ```go
   "$gte": float64(90)  // CORRECT in Go code
   "$gte": 90           // INCORRECT in Go code - will fail validation
   ```

2. **Pattern criteria MUST use `$pattern` wrapper**
   ```go
   "$pattern": map[string]interface{}{  // CORRECT
      "pattern": "consistent",
      // ...
   }
   
   "pattern": "consistent",  // INCORRECT - missing $pattern wrapper
   // ...
   ```

3. **Event names MUST match existing event types exactly**
   ```go
   "event": "test_event"  // Must match an existing event type name (case-sensitive)
   ```

4. **Badge creation MUST use `flow_definition` field**
   ```json
   // CORRECT
   "flow_definition": {
     "event": "test_event",
     "criteria": { ... }
   }
   
   // INCORRECT
   "criteria": {
     "event": "test_event",
     ...
   }
   ```

5. **Event submission requires the `payload` field**
   ```json
   // CORRECT
   {
     "payload": { "value": 75 }
   }
   
   // INCORRECT
   {
     // Missing payload field
   }
   ```

## Validation

Run this command to validate your criteria:

```bash
./tools/validate_test_criteria.sh
```

## Common Errors

| Error | Fix |
|-------|-----|
| `numeric values without explicit type conversion` | Add `float64()` to all numbers in Go code |
| `pattern criteria without $pattern wrapper` | Add the `$pattern` wrapper |
| `event type not found` | Check the event name matches exactly |
| `flow definition is required` | Use flow_definition field in badge creation |
| `event type is required` | Ensure event_type field is present and valid |
| `failed to save event: pq: null value...` | Include both data and payload fields |

## Response Formats

### Badge Creation Response

The response includes both badge and criteria objects:

```json
{
  "badge": {
    "id": 106,
    "name": "Simple Badge",
    "description": "Simple test badge",
    "image_url": "https://example.com/badge.png",
    "active": true,
    "created_at": "2025-03-16T19:43:02.844034Z",
    "updated_at": "2025-03-16T19:43:02.844034Z"
  },
  "criteria": {
    "id": 106,
    "badge_id": 106,
    "flow_definition": { ... },
    "created_at": "2025-03-16T19:43:02.844034Z",
    "updated_at": "2025-03-16T19:43:02.844034Z"
  }
}
```

### User Badges Response

The response is an array of badge objects:

```json
[
  {
    "awarded_at": "2025-03-16T19:43:02.863425Z",
    "description": "Simple test badge",
    "id": 106,
    "image_url": "https://example.com/badge.png",
    "metadata": "eyJsYXN0X...",
    "name": "Simple Badge"
  }
]
```

For complete documentation, see [BADGE_CRITERIA_FORMAT.md](BADGE_CRITERIA_FORMAT.md) 
