# Dynamic Time Variables Usage Guide

This document explains how to use the newly implemented dynamic time variables in the Badge Assignment System. Dynamic time variables allow you to create badge criteria with relative time references that automatically adjust based on when they are evaluated.

## Basic Usage

Dynamic time variables are evaluated at runtime when a badge criteria is checked. The system currently supports the `$NOW` variable and its adjusted forms.

### `$NOW` Variable

The `$NOW` variable represents the current date and time when the badge criteria is being evaluated.

**Example:**

```json
{
  "event": "user_activity",
  "criteria": {
    "timestamp": {
      "$lte": "$NOW"  // Events up to now
    }
  }
}
```

### Time Adjustments

You can add time adjustments to the `$NOW` variable using the format `$NOW(<adjustment>)`, where adjustment is a combination of a sign and a time unit.

**Example:**

```json
{
  "event": "user_activity",
  "criteria": {
    "timestamp": {
      "$gte": "$NOW(-30d)"  // Events from the last 30 days
    }
  }
}
```

### Supported Time Units

| Unit | Description | Example |
|------|-------------|---------|
| `s`  | Seconds     | `$NOW(-30s)` |
| `m`  | Minutes     | `$NOW(-15m)` |
| `h`  | Hours       | `$NOW(-6h)` |
| `d`  | Days        | `$NOW(-30d)` |
| `w`  | Weeks       | `$NOW(-2w)` |
| `M`  | Months      | `$NOW(-3M)` |
| `y`  | Years       | `$NOW(-1y)` |

You can also combine multiple time adjustments:

```json
{
  "event": "user_activity",
  "criteria": {
    "timestamp": {
      "$gte": "$NOW(-1y-3M)"  // Events from 1 year and 3 months ago until now
    }
  }
}
```

## Advanced Usage

### Time Windows

Dynamic time variables can be used to create flexible time windows:

```json
{
  "event": "purchase",
  "criteria": {
    "timestamp": {
      "$gte": "$NOW(-30d)",  // From 30 days ago
      "$lte": "$NOW"         // To now
    }
  }
}
```

### Usage in Multiple Fields

Dynamic time variables can be used in any field that expects a date/time value, not just the `timestamp` field:

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

## Real-World Examples

### Active Users Badge

Award a badge to users who have been active in the last 30 days:

```json
{
  "event": "user_activity",
  "criteria": {
    "timestamp": {
      "$gte": "$NOW(-30d)"
    },
    "$eventCount": { "$gte": 1 }
  }
}
```

### Loyal Customer Badge

Award a badge to users who have been customers for at least 1 year and made purchases in the last 3 months:

```json
{
  "criteria": {
    "$and": [
      {
        "user.created_at": {
          "$lte": "$NOW(-1y)"
        }
      },
      {
        "event": "purchase",
        "criteria": {
          "timestamp": {
            "$gte": "$NOW(-3M)"
          },
          "$eventCount": { "$gte": 1 }
        }
      }
    ]
  }
}
```

### Active Subscriber Badge

Award a badge to users with an active subscription:

```json
{
  "criteria": {
    "user.subscription.expires_at": {
      "$gte": "$NOW"
    }
  }
}
```

### Power User Badge

Award a badge to users who have been active for at least 10 days in the last 30 days:

```json
{
  "event": "user_activity",
  "criteria": {
    "timestamp": {
      "$gte": "$NOW(-30d)"
    },
    "$timePeriod": {
      "periodType": "day",
      "periodCount": {
        "$gte": 10
      }
    }
  }
}
```

## Best Practices

1. **Use for Rolling Windows**: Dynamic time variables are ideal for criteria that need to represent rolling time windows rather than hard-coded dates.

2. **Consistent Time Zone Handling**: All time calculations are performed in UTC to ensure consistency.

3. **Cache Consistency**: The system uses a time variable cache to ensure that multiple `$NOW` references within a single badge evaluation will use the same timestamp.

4. **Testing**: When testing badge criteria with dynamic time variables, remember that the system will use the current time for evaluation. For testing specific scenarios, you can adjust your criteria accordingly.

5. **Documentation**: When defining badge criteria with dynamic time variables, include clear comments to explain the time-based logic.

## Limitations

1. Dynamic time variables are currently only supported for the `$NOW` variable. Future enhancements may include additional time variables.

2. Time adjustments are limited to the supported time units listed above.

3. The system does not currently support more complex time expressions or calculations.

## Migration from Hard-Coded Timestamps

If you have existing badge criteria with hard-coded timestamps, you can migrate them to use dynamic time variables for improved maintenance. For example:

**Before:**
```json
{
  "timestamp": {
    "$gte": "2023-12-01T00:00:00Z"
  }
}
```

**After:**
```json
{
  "timestamp": {
    "$gte": "$NOW(-30d)"
  }
}
```

This change eliminates the need to periodically update badge criteria with new timestamps. 
