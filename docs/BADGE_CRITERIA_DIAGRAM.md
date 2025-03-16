# Badge Criteria Structure Diagrams

This document provides visual diagrams to help understand the badge criteria structure and badge assignment flow.

## Badge Creation Structure

```
+-----------------------------------+
|          Badge Object             |
+-----------------------------------+
| - name: string                    |
| - description: string             |
| - image_url: string               |
| - is_active: boolean              |
| - flow_definition: Object ------+ |
+-----------------------------------+
                                  |
                                  v
+-----------------------------------+
|        Flow Definition            |
+-----------------------------------+
| - event: string (event_type_name) |
| - criteria: Object -------------+ |
+-----------------------------------+
                                  |
                                  v
+-----------------------------------+
|           Criteria                |
+-----------------------------------+
|                                   |
|  +---------------------------+    |
|  | Direct Field Comparisons  |    |
|  | field: { $operator: val } |    |
|  +---------------------------+    |
|                                   |
|  +---------------------------+    |
|  | Logical Operators         |    |
|  | $and, $or, $not: [...]    |    |
|  +---------------------------+    |
|                                   |
|  +---------------------------+    |
|  | Pattern Detection         |    |
|  | $pattern: {...}           |    |
|  +---------------------------+    |
|                                   |
|  +---------------------------+    |
|  | Time-Based Criteria       |    |
|  | $timeWindow, $timePeriod  |    |
|  +---------------------------+    |
|                                   |
+-----------------------------------+
```

## Badge Assignment Flow

```
+-------------------+       +-------------------+
| Create Event Type | ----> | Define Data Schema|
+--------+----------+       +-------------------+
         |
         v
+--------+----------+       +-----------------------+
| Create Badge      | ----> | Set Flow Definition   |
| with Criteria     |       | (Event + Value Checks)|
+--------+----------+       +-----------------------+
         |
         v
+--------+----------+
| Submit Events     |
| (payload)         |
+--------+----------+
         |
         v
+--------+----------+       +----------------+
| Automatic Badge   | ----> | Store in DB    |
| Assignment        |       | User-Badge Link|
+--------+----------+       +----------------+
         |
         v
+--------+----------+
| Badge Retrieval   |
| (API or DB query) |
+-------------------+
```

## API Requests Structure

### Badge Creation Request

```
POST /api/v1/admin/badges
{
  name: "Example Badge",
  description: "...",
  image_url: "...",
  is_active: true,
  flow_definition: {      <-- Required field
    event: "event_name",  <-- Must match existing event type
    criteria: {           <-- Criteria for badge assignment
      field: {
        $operator: value
      }
    }
  }
}
```

### Event Submission Request

```
POST /api/v1/events
{
  event_type: "event_name",  <-- Must match existing event type
  user_id: "user123",
  timestamp: "2023-01-01T12:00:00Z",
  data: {                    <-- Used for criteria evaluation
    field: value
  },
  payload: {                 <-- Required for storage
    field: value
  }
}
```

## Common Error Patterns

```
+-------------------------------+       +---------------------------+
| Error: flow definition        | ----> | Solution: Use proper      |
| is required                   |       | flow_definition structure |
+-------------------------------+       +---------------------------+

+-------------------------------+       +---------------------------+
| Error: event type             | ----> | Solution: Ensure event    |
| is required                   |       | type exists and is used   |
+-------------------------------+       +---------------------------+

+-------------------------------+       +---------------------------+
| Error: null value in          | ----> | Solution: Include both    |
| column "payload"              |       | data and payload fields   |
+-------------------------------+       +---------------------------+

+-------------------------------+       +---------------------------+
| Error: pattern criteria       | ----> | Solution: Always use      |
| without $pattern wrapper      |       | $pattern wrapper          |
+-------------------------------+       +---------------------------+

+-------------------------------+       +---------------------------+
| Error: numeric values without | ----> | Solution: Use float64()   |
| explicit type conversion      |       | for all numeric values    |
+-------------------------------+       +---------------------------+
```

See the main [Badge Criteria Format Specification](BADGE_CRITERIA_FORMAT.md) for detailed documentation. 
