# Time-Based Criteria in the Badge Assignment System

This document summarizes our findings and recommendations for implementing time-based criteria in the badge assignment system.

## Overview

Time-based criteria are essential for creating badges that reward user behavior within specific time periods, such as:
- Recent activity badges (e.g., "Active in the last 30 days")
- Seasonal badges (e.g., "Holiday Shopper")
- Time-limited challenges (e.g., "Weekend Warrior")

## Approaches

### 1. Using the `$timeWindow` Operator

The system includes a `$timeWindow` operator designed to evaluate criteria within specific time windows.

**Example Implementation:**
```json
{
  "event": "purchase",
  "criteria": {
    "$timeWindow": {
      "start": "2023-12-01T00:00:00Z",
      "end": "2023-12-31T23:59:59Z",
      "criteria": {
        "$eventCount": {
          "$gte": 3
        }
      }
    }
  }
}
```

**Known Issues:**
- Inconsistent behavior when combined with other operators
- May not correctly evaluate events with the `last` parameter (e.g., `"last": "30d"`)
- Fixed date parameters (`start` and `end`) also showed inconsistent results in testing

### 2. Using the `timestamp` Field (Recommended)

Our testing has confirmed that using the `timestamp` field directly in criteria provides reliable time-based filtering.

**Example Implementation:**
```json
{
  "event": "user_activity",
  "criteria": {
    "$eventCount": {
      "$gte": 5
    },
    "timestamp": {
      "$gte": "2023-12-01T00:00:00Z",
      "$lte": "2023-12-31T23:59:59Z"
    }
  }
}
```

**Benefits:**
- Simpler implementation
- More predictable behavior
- Confirmed to work correctly in testing
- Compatible with all other operators

## Testing Results

We conducted several tests to verify the behavior of time-based criteria:

1. **Holiday Shopper Badge Test** - Used `$timeWindow` with fixed dates (Dec 1-31, 2023)
   - Result: Badge was not awarded despite meeting criteria
   - See: `tests/integration/holiday_shopper_badge_test.go`

2. **Recent Activity Badge Test** - Used `$timeWindow` with `last` parameter (30 days)
   - Result: Badge was not awarded despite meeting criteria
   - See: `tests/integration/recent_activity_badge_test.go`

3. **Basic Event Count Test** - Used simple `$eventCount` without time constraints
   - Result: Badge was awarded correctly
   - See: `tests/integration/time_window_badge_test.go`

4. **Timestamp Filter Test** - Used `timestamp` field directly for filtering
   - Result: Badge was awarded correctly
   - See: `tests/integration/timestamp_filter_badge_test.go`

## Implementation Recommendations

1. **Use the `timestamp` field directly** for all time-based criteria
2. Ensure proper timezone handling in your timestamp values
3. Combine with `$eventCount` or other operators as needed
4. Refer to the timestamp filter test for a complete working example

## Future Work

- Investigate and fix issues with the `$timeWindow` operator
- Add more comprehensive tests for various time-based scenarios
- Consider adding timezone support for more user-friendly time windows

## Proposed Enhancement: Dynamic Time Variables

To address the limitations of the current time-based filtering approaches, we propose adding support for dynamic time variables in badge criteria:

### Core Time Variables

1. **`$NOW`**: Evaluates to the current timestamp when the badge criteria is checked
2. **`$NOW(<adjustment>)`**: Evaluates to the current timestamp adjusted by the specified duration

### Example Usage

```json
{
  "event": "user_activity",
  "criteria": {
    "$eventCount": { "$gte": 5 },
    "timestamp": {
      "$gte": "$NOW(-30d)"  // Events from the last 30 days
    }
  }
}
```

This approach would:
- Eliminate the need to manually update timestamp criteria
- Provide consistent, reliable time-based filtering
- Improve readability and maintainability of badge definitions

For more details on this proposal, see [Dynamic Time Variables Proposal](DYNAMIC_TIME_VARIABLES.md).

## References

- [Count Operators Documentation](COUNT_OPERATORS.md)
- [Timestamp Filter Test](/tests/integration/timestamp_filter_badge_test.go)
- [Rule Engine Implementation](/internal/engine/rule_engine.go)
- [Dynamic Time Variables Proposal](DYNAMIC_TIME_VARIABLES.md) 
