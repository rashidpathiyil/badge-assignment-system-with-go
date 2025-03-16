# Badge Assignment System: Extended Criteria Capabilities

This document provides an overview of the extended badge criteria capabilities that have been implemented in the Badge Assignment System. These extensions allow for more sophisticated badge awarding based on various temporal and behavioral patterns.

## Time-Based Criteria

### Time Period Criteria (`$timePeriod`)

The Time Period Criteria allows for badge awarding based on user activity across different time periods like days, weeks, or months.

**Sample Configuration:**
```json
{
  "$timePeriod": {
    "periodType": "day",
    "periodCount": { "$gte": 5 },
    "excludeWeekends": true,
    "excludeHolidays": true,
    "holidays": ["2023-12-25", "2023-01-01"]
  }
}
```

**Parameters:**
- `periodType`: The type of time period to consider ("day", "week", "month")
- `periodCount`: Criteria for the number of periods with activity
- `excludeWeekends`: Whether weekend days should be excluded from counting
- `excludeHolidays`: Whether holiday days should be excluded from counting
- `holidays`: A list of specific holiday dates to exclude

**Use Case:** Award a badge if a user is active for at least 5 weekdays, excluding holidays.

### Pattern Criteria (`$pattern`)

The Pattern Criteria allows for badge awarding based on consistent, increasing, or decreasing patterns of activity over time.

**Sample Configuration:**
```json
{
  "$pattern": {
    "pattern": "increasing",
    "periodType": "week",
    "minPeriods": 3,
    "minIncreasePct": 10,
    "maxDeviation": 0.2
  }
}
```

**Parameters:**
- `pattern`: The type of pattern to look for ("consistent", "increasing", "decreasing")
- `periodType`: The time period granularity ("day", "week", "month")
- `minPeriods`: Minimum number of consecutive periods required
- `minIncreasePct`: For increasing pattern, minimum percentage increase required
- `maxDecreasePct`: For decreasing pattern, maximum percentage decrease allowed
- `maxDeviation`: For consistent pattern, maximum deviation percentage allowed

**Implementation Details:**
- **Consistent Pattern**: Checks event counts across periods and allows for a maximum deviation from the average. The algorithm is robust to outliers and can detect consistent patterns even with isolated anomalies.
- **Increasing Pattern**: Identifies growth trends using percentage increases between periods and checks if the average increase meets the minimum threshold. The algorithm calculates trend strength and can identify patterns even with occasional dips.
- **Decreasing Pattern**: Similar to increasing pattern but looks for a gradual decrease in activity. The implementation handles chronological ordering and can detect seasonal decline patterns.

**Metadata Output:**
Each pattern evaluator provides detailed metadata that explains why a pattern was detected or not detected:
- `is_consistent`, `is_increasing`, or `is_decreasing`: Boolean indicating if the pattern was detected
- `average`: Average count across periods (for consistent patterns)
- `max_deviation`: Maximum deviation from average (for consistent patterns)
- `coefficient_var`: Coefficient of variation as a normalized measure of dispersion
- `average_percent_increase/decrease`: Average percentage change between periods
- `trend_strength`: A value between 0 and 1 indicating how strong the trend is
- `period_counts`: Array of counts for each period
- `period_keys`: Array of period keys (dates/weeks/months)

**Use Case:** Award a badge if user activity shows an increasing trend week over week with at least 10% growth.

### Sequence Criteria (`$sequence`)

The Sequence Criteria allows for badge awarding based on specific sequences of events occurring in order.

**Sample Configuration:**
```json
{
  "$sequence": {
    "sequence": ["login", "search", "view", "purchase"],
    "maxGapSeconds": 900,
    "requireStrict": true
  }
}
```

**Parameters:**
- `sequence`: Ordered array of event types to look for
- `maxGapSeconds`: Maximum time (in seconds) allowed between consecutive events
- `requireStrict`: If true, no other events can occur between the sequence events

**Use Case:** Award a badge when a user follows the complete conversion funnel within 15 minutes.

### Gap Criteria (`$gap`)

The Gap Criteria allows for badge awarding based on the time gaps between consecutive events.

**Sample Configuration:**
```json
{
  "$gap": {
    "minGapHours": 2,
    "maxGapHours": 24,
    "periodType": "day"
  }
}
```

**Parameters:**
- `minGapHours`: Minimum gap hours required between events
- `maxGapHours`: Maximum gap hours allowed between events
- `periodType`: Optional period type for special gap calculations

**Use Case:** Award a badge if the user consistently maintains a healthy spacing between activities.

### Duration Criteria (`$duration`)

The Duration Criteria allows for badge awarding based on the duration between paired start/end events.

**Sample Configuration:**
```json
{
  "$duration": {
    "startEvent": { "type": "session_start" },
    "endEvent": { "type": "session_end" },
    "matchProperty": "session_id",
    "duration": { "$gte": 30 },
    "unit": "minutes"
  }
}
```

**Parameters:**
- `startEvent`: Criteria to identify the starting event
- `endEvent`: Criteria to identify the ending event
- `matchProperty`: Property to match to pair start and end events
- `duration`: Criteria for the time duration
- `unit`: Time unit for duration ("seconds", "minutes", "hours", "days")

**Use Case:** Award a badge when a user spends at least 30 minutes in a single session.

### Aggregation Criteria (`$aggregate`)

The Aggregation Criteria allows for badge awarding based on aggregate calculations on event properties.

**Sample Configuration:**
```json
{
  "$aggregate": {
    "type": "avg",
    "field": "score",
    "value": { "$gte": 90 }
  }
}
```

**Parameters:**
- `type`: Aggregation type to apply ("min", "max", "avg", "sum", "count")
- `field`: The event field to aggregate
- `value`: Criteria for the aggregation result

**Use Case:** Award a badge when a user's average score is 90 or higher.

### Time Window Criteria (`$timeWindow`)

The Time Window Criteria allows for evaluating a sub-criteria only within a specific time window.

**Sample Configuration:**
```json
{
  "$timeWindow": {
    "start": "2023-01-01T00:00:00Z",
    "end": "2023-01-31T23:59:59Z",
    "flow": {
      "$or": [
        { "event": "login", "criteria": { "$eventCount": { "$gte": 20 } } },
        { "event": "purchase", "criteria": { "$eventCount": { "$gte": 5 } } }
      ]
    }
  }
}
```

**Relative Time Window Example:**
```json
{
  "$timeWindow": {
    "last": "30d",
    "businessDaysOnly": true,
    "flow": {
      "event": "login",
      "criteria": {
        "$eventCount": { "$gte": 10 }
      }
    }
  }
}
```

**Parameters:**
- `start`: Start datetime of the window (ISO 8601 format)
- `end`: End datetime of the window (ISO 8601 format)
- `last`: Relative time window (e.g., "30d", "2w", "1m", "1q", "1y")
  - Supported units: "d" (days), "w" (weeks), "m" (months), "q" (quarters), "y" (years)
- `businessDaysOnly`: If true, excludes weekends from the time window
- `flow`: A nested criteria flow to evaluate within the time window

**Use Case:** 
- Award a badge based on activity specifically during January 2023
- Award a badge based on activity within the last 30 business days

### Event Count Criteria (`$eventCount`)

The Event Count Criteria allows for badge awarding based on the raw count of matching events.

**Sample Configuration:**
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

**Parameters:**
- `$eventCount`: Criteria for the number of matching events
  - Supports all comparison operators (`$eq`, `$gt`, `$gte`, `$lt`, `$lte`, `$ne`)

**Use Case:** Award a badge when a user completes at least 10 workouts.

**Implementation Details:**
- Counts the raw number of events (not unique time periods)
- Can be combined with other filtering criteria
- Provides detailed metadata including `event_count`
- Must be placed within the `criteria` field of an event-specific definition
- It always operates on events of a specific type (specified by the `event` field)

### Time Period Counting (Unique Days/Weeks/Months)

For counting unique time periods with activity, use the `$timePeriod` operator:

```json
{
  "$timePeriod": {
    "periodType": "day",
    "periodCount": { "$gte": 3 }
  }
}
```

**Important Implementation Note:**
- `$timePeriod` must be placed at the top level of the flow definition, not within a `criteria` field
- It automatically considers all events for the user, regardless of event type
- Adds `unique_period_count` to the badge metadata

## Combining Criteria

These criteria can be combined using logical operators (`$and`, `$or`, `$not`) to create complex badge requirements.

**Example of a Combined Counting Criteria:**
```json
{
  "$and": [
    {
      "event": "workout_completed",
      "criteria": {
        "$eventCount": {
          "$gte": 10
        }
      }
    },
    {
      "$timePeriod": {
        "periodType": "day", 
        "periodCount": { "$gte": 3 }
      }
    }
  ]
}
```

**Note on Combined Operators:**
When combining both count operators, each must maintain its correct structure within the flow definition. When both criteria are met, the badge metadata will contain both `event_count` and `unique_period_count`.

**Example of a Complex Combined Criteria:**
```json
{
  "$and": [
    {
      "$timePeriod": {
        "periodType": "day",
        "periodCount": { "$gte": 5 }
      }
    },
    {
      "$or": [
        {
          "$pattern": {
            "pattern": "increasing",
            "periodType": "day",
            "minPeriods": 3
          }
        },
        {
          "$aggregate": {
            "function": "sum",
            "property": "points",
            "result": { "$gte": 1000 }
          }
        }
      ]
    }
  ]
}
```

## Performance Considerations

For optimal performance when working with time-based criteria:

1. Use the appropriate indexes on timestamp fields in your database
2. Consider implementing a caching layer for frequently accessed criteria results
3. Process badge criteria evaluations asynchronously when possible
4. Aggregate events data at regular intervals for high-volume applications

## Default Values for Criteria Parameters

Most criteria types have default values for optional parameters:

### Pattern Criteria Defaults
- `minPeriods`: 3 periods
- `minIncreasePct`: 5.0% (for increasing pattern)
- `maxDecreasePct`: 5.0% (for decreasing pattern)
- `maxDeviation`: 0.2 (20% deviation for consistent pattern)

### Time Period Criteria Defaults
- If no `periodCount` is specified, the criterion is met if there's at least one period with activity
- `excludeWeekends` and `excludeHolidays` default to `false`

### Metadata Output

Each criteria evaluation produces detailed metadata:

- **Time Period Criteria**: Includes `unique_period_count` (number of unique periods with activity)
- **Pattern Criteria**: Includes pattern-specific metrics like `average`, `max_deviation`, `trend_strength`
- **Event Count Criteria**: Includes `event_count` (total number of matching events)

For more detailed information about the differences between counting events and counting periods, see [COUNT_OPERATORS.md](COUNT_OPERATORS.md).

## Payload Field Criteria

The Payload Field Criteria allows for badge awarding based on specific values in the event payload. Each event in the system contains a payload with event-specific data, and badge criteria can filter events based on this payload data.

### Overview

The `payload` field in a badge criteria flow definition specifies conditions that event payloads must meet. This allows badges to be awarded based on specific values or attributes within events.

**Sample Configuration:**
```json
{
  "criteria": {
    "payload": {
      "time": { "$lt": "09:00:00" }
    }
  },
  "event": "check-in"
}
```

### How It Works

When evaluating badge criteria, the system:
1. Retrieves events of the specified type for the user
2. For each event, checks if its payload contains the specified fields
3. Applies the comparison conditions to the payload field values
4. Filters out events that don't match the payload criteria

### Payload Field Structure

The payload field can be structured in two ways:

1. **Direct equality comparison**: 
   ```json
   "payload": {
     "status": "fixed"
   }
   ```
   This checks if the `status` field in the event payload exactly equals "fixed".

2. **Operator comparison**:
   ```json
   "payload": {
     "time": {
       "$lt": "09:00:00"
     }
   }
   ```
   This uses comparison operators to check if the `time` field is less than "09:00:00".

### Supported Comparison Operators

The following operators are supported for payload field comparisons:

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

### Real-World Examples

1. **Early Bird Badge**:
   ```json
   "payload": {
     "time": {
       "$lt": "09:00:00"
     }
   }
   ```
   Awards the badge when check-in events have a time earlier than 9:00 AM.

2. **Overtime Warrior Badge**:
   ```json
   "payload": {
     "is_overtime": true
   }
   ```
   Awards the badge when work-log events have the overtime flag set to true.

3. **Bug Hunter Badge**:
   ```json
   "payload": {
     "status": "fixed"
   }
   ```
   Awards the badge when bug-report events have a status of "fixed".

### Combining with Other Criteria

Payload criteria can be combined with other criteria types for more complex badge requirements:

```json
{
  "event": "work-log",
  "criteria": {
    "$timePeriod": {
      "periodType": "week",
      "periodCount": { "$gte": 1 }
    },
    "payload": {
      "is_overtime": true
    },
    "$aggregate": {
      "function": "sum",
      "property": "hours",
      "result": { "$gte": 10 }
    }
  }
}
```

This awards a badge when a user logs overtime hours that sum to at least 10 hours in a week.

### Important Considerations

1. **Field validation**: Ensure that payload fields referenced in criteria exist in the corresponding event type schema
2. **Type matching**: Comparison operators require values of compatible types (e.g., comparing numbers with numbers)
3. **Case sensitivity**: Field names and string values are case-sensitive
4. **Null handling**: If a payload field doesn't exist in an event, the event won't match the criteria

## Testing

Each time-based criteria evaluator has corresponding unit tests in `internal/engine/time_evaluators_test.go`. 
These tests provide examples of how to use each criteria type and validate their functionality. 
