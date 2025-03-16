# Rule Engine Documentation

The Badge Assignment System includes a powerful rule engine for evaluating whether users have met the criteria to earn badges.

## Overview

The rule engine evaluates badge criteria against user events to determine if a badge should be awarded. It uses a flexible JSON-based criteria definition format that can express complex conditions.

## Criteria Evaluation Model

Badge criteria are defined using a JSON structure that supports the following components:

### Event-Based Criteria

Event-based criteria check events of a specific type:

```json
{
  "event": "login",
  "criteria": {
    "count": { "$gte": 5 }
  }
}
```

In this example, the criteria check if the user has logged in at least 5 times.

### Logical Operators

The rule engine supports the following logical operators:

- `$and`: All conditions must be true
- `$or`: At least one condition must be true
- `$not`: The condition must be false

Example using logical operators:

```json
{
  "$and": [
    {
      "event": "login",
      "criteria": {
        "count": { "$gte": 5 }
      }
    },
    {
      "event": "project-contribution",
      "criteria": {
        "count": { "$gte": 3 }
      }
    }
  ]
}
```

This criteria requires the user to have both logged in at least 5 times AND made at least 3 project contributions.

### Time-Based Criteria

The rule engine supports various time-based criteria evaluations:

#### Time Period Criteria

Evaluates events over specific time periods:

```json
{
  "event": "login",
  "criteria": {
    "$timePeriod": {
      "periodType": "day",
      "periodCount": { "$gte": 5 },
      "excludeWeekends": true,
      "excludeHolidays": true,
      "holidays": ["2023-12-25", "2024-01-01"]
    }
  }
}
```

Supported period types:
- `day`: Evaluates daily events
- `week`: Evaluates weekly events
- `month`: Evaluates monthly events

#### Pattern Criteria

Detects specific patterns in event occurrences:

```json
{
  "event": "login",
  "criteria": {
    "$pattern": {
      "pattern": "consecutive",
      "count": { "$gte": 5 },
      "ignoreGaps": false,
      "maxGap": 1
    }
  }
}
```

Supported patterns:
- `consecutive`: Events occurring in sequence
- `weekly`: Events occurring on the same day of the week
- `monthly`: Events occurring on the same day of the month
- `custom`: Custom pattern defined by the pattern array

#### Sequence Criteria

Evaluates if events occur in a specific sequence:

```json
{
  "$sequence": [
    { "event": "start-task" },
    { "event": "complete-task" },
    { "event": "review-task" }
  ]
}
```

#### Gap Criteria

Evaluates gaps between events:

```json
{
  "event": "login",
  "criteria": {
    "$gap": {
      "maxGap": "48h",
      "minGap": "24h"
    }
  }
}
```

#### Duration Criteria

Evaluates the duration of time between specific events:

```json
{
  "event": "login",
  "criteria": {
    "$duration": {
      "total": { "$gte": "30d" },
      "from": "first-event",
      "to": "last-event"
    }
  }
}
```

### Aggregation Criteria

Performs calculations over event data:

```json
{
  "event": "work-hours",
  "criteria": {
    "$aggregate": {
      "type": "sum",
      "field": "payload.hours",
      "value": { "$gte": 40 }
    }
  }
}
```

Supported aggregation types:
- `sum`: Calculates the sum of values
- `avg`: Calculates the average value
- `min`: Finds the minimum value
- `max`: Finds the maximum value
- `count`: Counts the number of events

### Evaluation Process

When evaluating badge criteria:

1. The rule engine retrieves all relevant events for the user
2. It applies the criteria definition recursively
3. For time-based criteria, it organizes events chronologically
4. For aggregate criteria, it processes numeric data from the events
5. The engine returns a boolean result indicating if the criteria are met, along with metadata about the evaluation

## Integration with the API

The rule engine is integrated with the Badge Assignment System's API through the following endpoints:

- `POST /api/v1/events`: Submits new events that can trigger criteria evaluation
- `GET /api/v1/users/{user_id}/badges`: Returns badges that a user has earned
- `POST /api/v1/users/{user_id}/evaluate`: (Planned) Manually triggers criteria evaluation for a user

## Advanced Features

The rule engine includes additional features for complex requirements:

- **Progress Tracking**: The engine can track progress toward badge criteria
- **Metadata Storage**: Evaluation details are stored as metadata for diagnostics
- **Conditional Logic**: Complex conditions can be expressed through nested operators

## Examples

### Consistent Learning Badge

```json
{
  "$and": [
    {
      "event": "lesson-completed",
      "criteria": {
        "count": { "$gte": 10 }
      }
    },
    {
      "event": "lesson-completed",
      "criteria": {
        "$pattern": {
          "pattern": "consecutive",
          "count": { "$gte": 5 }
        }
      }
    }
  ]
}
```

This criteria awards a badge when a user completes at least 10 lessons total AND has completed at least 5 lessons consecutively.

### Active Contributor Badge

```json
{
  "$and": [
    {
      "event": "code-commit",
      "criteria": {
        "count": { "$gte": 20 }
      }
    },
    {
      "event": "code-commit",
      "criteria": {
        "$timePeriod": {
          "periodType": "week",
          "periodCount": { "$gte": 3 }
        }
      }
    },
    {
      "event": "pull-request",
      "criteria": {
        "count": { "$gte": 5 }
      }
    }
  ]
}
```

This criteria awards a badge when a user has made at least 20 code commits, contributed in at least 3 different weeks, and submitted at least 5 pull requests. 
