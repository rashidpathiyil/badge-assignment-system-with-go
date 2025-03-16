# Understanding Count Operators in the Badge Assignment System

The Badge Assignment System implements two distinct types of count operators that serve different purposes. This document explains the difference between them and when to use each one.

## 1. Event Count (`$eventCount`)

The `$eventCount` operator counts the raw number of events that match specified criteria.

### Example

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

### Key Characteristics

- **What it counts**: Raw events
- **Use case**: When you want to count the total number of times something happened
- **Duplicates**: Each event is counted individually, even if they occur on the same day
- **Metadata field**: `event_count`

### When to Use `$eventCount`

- Counting total number of activities (e.g., "Complete 10 workouts")
- Measuring volume of actions (e.g., "Make 5 purchases")
- When timing doesn't matter, only the total count

### Implementation Details

- `$eventCount` must be placed within the `criteria` field of an event-specific definition
- It always operates on events of a specific type (specified by the `event` field)
- The operator can be combined with other field-specific criteria (e.g., filtering by status)

## 2. Period Count (`periodCount` in `$timePeriod`)

The `periodCount` parameter within the `$timePeriod` operator counts unique time periods (days, weeks, months) that have any activity.

### Example

```json
{
  "$timePeriod": {
    "periodType": "day",
    "periodCount": {
      "$gte": 5
    },
    "excludeWeekends": true
  }
}
```

### Key Characteristics

- **What it counts**: Unique time periods (days, weeks, months) with any activity
- **Use case**: When you want to measure consistency over time
- **Duplicates**: Multiple events on the same day/week/month only count as one period
- **Metadata field**: `unique_period_count`

### When to Use `periodCount`

- Measuring consistency (e.g., "Be active for 5 different days")
- Streak-based badges (e.g., "Use the app every week for a month")
- When distribution of activity over time is important

### Implementation Details

- `$timePeriod` **must be placed at the top level** of the flow definition, not within a `criteria` field
- It automatically considers all events for the user, regardless of event type
- The time periods are determined by the `periodType` parameter according to ISO calendar standards
- For `day` type, a period key is generated using the format "2006-01-02" (year-month-day)

## Example Illustrating the Difference

Consider a user who:
- Completes 3 workouts on Monday
- Completes 2 workouts on Tuesday
- Completes 0 workouts on Wednesday
- Completes 4 workouts on Thursday

With `$eventCount`:
```json
{
  "event": "workout_completed",
  "criteria": {
    "$eventCount": {
      "$gte": 8
    }
  }
}
```
This would be satisfied because there are 9 total workout events (3+2+0+4).

With `$timePeriod` and `periodCount`:
```json
{
  "$timePeriod": {
    "periodType": "day",
    "periodCount": {
      "$gte": 3
    }
  }
}
```
This would be satisfied because there are 3 different days with activity (Monday, Tuesday, Thursday).

## Combining Both Approaches

You can create sophisticated badges by combining both types of counting:

```json
{
  "$and": [
    {
      "$timePeriod": {
        "periodType": "week",
        "periodCount": {
          "$gte": 4
        }
      }
    },
    {
      "event": "workout_completed",
      "criteria": {
        "$eventCount": {
          "$gte": 20
        }
      }
    }
  ]
}
```

This would require both being active for at least 4 different weeks AND completing at least 20 total workouts.

### Implementation Note for Combined Operators

When combining both operators using logical operators like `$and`:

1. Each operator must maintain its correct structure:
   - `$eventCount` stays within the `criteria` field of an event-specific definition
   - `$timePeriod` remains at the top level of its branch in the flow definition

2. When both criteria are met, the badge metadata will contain:
   - `event_count`: The total number of events counted
   - `unique_period_count`: The number of unique time periods with activity

3. Example metadata for a combined badge (in JSON format):
   ```json
   {
     "event_count": 10, 
     "unique_period_count": 3
   }
   ```

## Implementation Note

The distinction between these two count types is intentional in the API design:
- `$eventCount` is prefixed with `$` because it's an operator (following MongoDB syntax style)
- `periodCount` doesn't have a `$` prefix because it's a configuration parameter for the `$timePeriod` operator

This naming convention helps make the different roles clear in the JSON configuration. 

## Common Pitfalls to Avoid

1. **Incorrect Placement**: Placing `$timePeriod` inside a `criteria` field instead of at the top level will cause the badge to not work properly.

2. **Forgetting Event Type**: When using `$eventCount`, always specify the event type with the `event` field.

3. **Timestamp Format**: When creating events with specific timestamps for testing, use the RFC3339 format and be aware that the system parses these timestamps to group events by their appropriate periods.

4. **Timezone Issues**: Be aware that the system uses the timestamp's timezone when determining which period an event belongs to.

## Alternative Approach

If you need time-based filtering, consider using the `timestamp` field directly in your criteria:

```json
{
  "event": "purchase_event",
  "criteria": {
    "$eventCount": {
      "$gte": 3
    },
    "timestamp": {
      "$gte": "2023-12-01T00:00:00Z",
      "$lte": "2023-12-31T23:59:59Z"
    }
  }
}
```

This approach has been tested and confirmed to work correctly. See our [Timestamp Filter Test](/tests/integration/timestamp_filter_badge_test.go) for a complete example.

## Example Implementation

For detailed examples of implementing both count operators in various scenarios, see:
- [Count Operators Example JSON](/docs/examples/count_operators_examples.json)
- [Test Implementation in Go](/tests/integration/consistent_reporter_badge_test.go) - Real-world test showing how to use `$timePeriod`
- [Combined Operators Test](/tests/integration/super_reporter_badge_test.go) - Test that combines both count operators

### Known Issues

- The `$timeWindow` operator with the `last` parameter (e.g., `"last": "30d"`) may not correctly filter events.
- Using fixed dates with `start` and `end` parameters also showed inconsistent results in our testing.
- Further investigation is needed to fully resolve these issues.
- For a working alternative, see the [Timestamp Filter Test](/tests/integration/timestamp_filter_badge_test.go) which demonstrates using the `timestamp` field directly.
