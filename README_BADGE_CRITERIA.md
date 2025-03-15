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
    "count": { "$gte": 5 },
    "excludeWeekends": true,
    "excludeHolidays": true,
    "holidays": ["2023-12-25", "2023-01-01"]
  }
}
```

**Parameters:**
- `periodType`: The type of time period to consider ("day", "week", "month")
- `count`: Criteria for the number of periods with activity
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
    "function": "avg",
    "property": "score",
    "result": { "$gte": 90 }
  }
}
```

**Parameters:**
- `function`: Aggregation function to apply ("min", "max", "avg", "sum", "count")
- `property`: The event property to aggregate
- `result`: Criteria for the aggregation result

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
        { "event": "login", "criteria": { "count": { "$gte": 20 } } },
        { "event": "purchase", "criteria": { "count": { "$gte": 5 } } }
      ]
    }
  }
}
```

**Parameters:**
- `start`: Start datetime of the window (ISO 8601 format)
- `end`: End datetime of the window (ISO 8601 format)
- `flow`: A nested criteria flow to evaluate within the time window

**Use Case:** Award a badge based on activity specifically during January 2023.

## Combining Criteria

These criteria can be combined using logical operators (`$and`, `$or`, `$not`) to create complex badge requirements.

**Example of a Combined Criteria:**
```json
{
  "$and": [
    {
      "$timePeriod": {
        "periodType": "day",
        "count": { "$gte": 5 }
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

## Testing

Each time-based criteria evaluator has corresponding unit tests in `internal/engine/time_evaluators_test.go`. 
These tests provide examples of how to use each criteria type and validate their functionality. 
