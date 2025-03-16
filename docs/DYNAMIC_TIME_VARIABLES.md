# Dynamic Time Variables Proposal

## Overview

This document proposes enhancing the Badge Assignment System with dynamic time variables to improve time-based badge criteria. The current implementation requires hard-coded timestamps, which limits flexibility and creates maintenance overhead.

## Problem Statement

Currently, time-based badge criteria have the following limitations:

1. The `$timeWindow` operator has proven unreliable in testing
2. Using direct timestamp filtering requires hard-coding timestamps:
   ```json
   "timestamp": {
     "$gte": "2023-12-01T00:00:00Z"
   }
   ```
3. For rolling time windows (e.g., "last 30 days"), badge criteria must be manually updated or recreated periodically
4. Time zones and daylight saving time add further complexity

## Proposed Solution: Dynamic Time Variables

We propose adding support for dynamic time variables that are evaluated at runtime:

### Core Variables

1. **`$NOW`**: Evaluates to the current timestamp when the badge criteria is checked
2. **`$NOW(<adjustment>)`**: Evaluates to the current timestamp adjusted by the specified duration

**Important Note**: Dynamic time variables can be used in **any field** that expects a date/time value, not just the `timestamp` field. This provides flexible time-based filtering across all temporal attributes of your events and users.

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

```json
{
  "event": "purchase",
  "criteria": {
    "$eventCount": { "$gte": 3 },
    "timestamp": {
      "$gte": "$NOW(-1y)",  // From 1 year ago
      "$lte": "$NOW"        // To now
    }
  }
}
```

### Multi-Field Example

Dynamic time variables work with any date field in your criteria, allowing for complex temporal conditions:

```json
{
  "event": "user_session",
  "criteria": {
    "$eventCount": { "$gte": 10 },
    "timestamp": {
      "$gte": "$NOW(-30d)"  // Recent activity (last 30 days)
    },
    "user": {
      "created_at": {
        "$lte": "$NOW(-1y)"  // Account at least 1 year old
      },
      "subscription": {
        "expires_at": {
          "$gte": "$NOW"  // Active subscription
        }
      }
    },
    "last_purchase": {
      "$gte": "$NOW(-90d)"  // Purchased something in last 90 days
    }
  }
}
```

### Supported Time Adjustments

| Unit | Description | Example |
|------|-------------|---------|
| `s`  | Seconds     | `$NOW(-30s)` |
| `m`  | Minutes     | `$NOW(-15m)` |
| `h`  | Hours       | `$NOW(-6h)` |
| `d`  | Days        | `$NOW(-30d)` |
| `w`  | Weeks       | `$NOW(-2w)` |
| `M`  | Months      | `$NOW(-3M)` |
| `y`  | Years       | `$NOW(-1y)` |

Multiple adjustments could also be supported:
```
$NOW(-1y-3M)  // 1 year and 3 months ago
```

## Benefits

1. **Truly Dynamic Criteria**: Badge criteria remain valid indefinitely without manual updates
2. **Reduced Maintenance**: No need to periodically regenerate badge definitions
3. **Improved Readability**: Badge definitions clearly express time-based logic
4. **Consistent Time Handling**: Centralized time calculation ensures consistency
5. **Field Flexibility**: Can be used in any date-related field, not just event timestamps
6. **Complex Temporal Conditions**: Enables sophisticated criteria combining multiple time dimensions:
   - Recent activity (timestamp)
   - Account age (created_at)
   - Subscription status (expires_at)
   - Purchase recency (last_purchase_date)
   - Event frequency within custom time periods
7. **User Journey Mapping**: Create badges that reward specific user journeys across time dimensions

## Implementation Approach

1. **Parser Enhancement**: Update the criteria parser to recognize `$NOW` syntax
2. **Variable Resolution**: Add a time variable resolver in the rule engine
3. **Timestamp Calculation**: Implement adjustment calculations for different time units
4. **Cache Mechanism**: Implement a per-evaluation cache to ensure consistent timestamp usage

## Implementation Considerations

1. **Timezone Handling**: All time calculations should respect a configurable timezone
2. **Performance Impact**: Minimal, as variable resolution happens once per badge evaluation
3. **Backward Compatibility**: Existing criteria with hard-coded timestamps will continue to work
4. **Testing Strategy**: Add comprehensive tests for various time variable scenarios

### Real-World Examples: Multi-Field Time Criteria

#### Loyal Active Customer Badge

This badge rewards users who:
- Have been customers for at least 1 year
- Have been active in the last 30 days
- Have a valid subscription
- Have made at least 5 purchases within the last 6 months

```json
{
  "name": "Loyal Active Customer",
  "description": "Awarded to long-time active customers with recent purchases",
  "criteria": {
    "$and": [
      {
        "user.created_at": {
          "$lte": "$NOW(-1y)"  // Account is at least 1 year old
        }
      },
      {
        "user.subscription.status": "active"
      },
      {
        "user.subscription.expires_at": {
          "$gte": "$NOW"  // Subscription is not expired
        }
      },
      {
        "event": "user_activity",
        "criteria": {
          "timestamp": {
            "$gte": "$NOW(-30d)"  // Activity in last 30 days
          },
          "$eventCount": { "$gte": 1 }
        }
      },
      {
        "event": "purchase",
        "criteria": {
          "timestamp": {
            "$gte": "$NOW(-6M)"  // Purchases in last 6 months
          },
          "$eventCount": { "$gte": 5 }
        }
      }
    ]
  }
}
```

This badge definition will always evaluate correctly regardless of when it's run, thanks to the dynamic time variables automatically adjusting to the current date.

## Additional Future Enhancements

Beyond `$NOW`, we could consider other time variables:

1. **`$STARTOFDAY`**, **`$STARTOFMONTH`**, **`$STARTOFYEAR`**: For period-specific calculations
2. **`$EVENT_TIMESTAMP`**: Reference the timestamp of the triggering event
3. **`$USER_CREATED_AT`**: Reference when the user was created for tenure-based badges

## Next Steps

1. RFC Review and Feedback
2. Implementation Plan
3. Prototype Development
4. Testing Framework
   - Test basic variable resolution
   - Test complex nested criteria
   - Test multiple date fields in the same criteria
   - Test timezone handling
   - Test with real-world badge scenarios
5. Documentation Update
   - Add code examples for various date fields
   - Document best practices for complex time criteria
6. Rollout Strategy
   - Progressive deployment
   - Migration guide for existing hard-coded timestamps 
