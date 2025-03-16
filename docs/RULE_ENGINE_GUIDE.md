# Badge Assignment Rule Engine Guide

This document provides a detailed explanation of the badge assignment rule engine, which is a core component of the system. Understanding how the rule engine works is essential for making changes to the badge assignment logic.

## Overview

The rule engine evaluates incoming events against defined conditions to determine if a badge should be assigned to a user. The process follows these steps:

1. An event is received with a specific event type and associated data
2. The rule engine looks up rules that match the event type
3. For each matching rule, conditions are evaluated
4. If conditions are met, the badge defined in the rule is assigned to the user

## Key Components

### Event

An event represents something that happened in the system, such as:

```go
type Event struct {
    Type      int               // Event type ID
    UserID    string            // User who triggered the event
    Timestamp time.Time         // When the event occurred
    Data      map[string]string // Additional data associated with the event
}
```

### Condition

A condition defines criteria that must be met for a rule to trigger:

```go
type Condition struct {
    ID            int    // Condition ID
    Type          int    // Condition type ID
    Operator      string // Comparison operator (e.g., "eq", "gt", "contains")
    ExpectedValue string // Value to compare against
    Field         string // Field in the event data to evaluate
}
```

### Rule

A rule connects events to badges through conditions:

```go
type Rule struct {
    ID          int         // Rule ID
    EventTypeID int         // Event type this rule applies to
    BadgeID     int         // Badge to assign when conditions are met
    Conditions  []Condition // Conditions that must be satisfied
    Operator    string      // Logical operator between conditions ("and", "or")
}
```

## Rule Evaluation

The rule engine evaluates conditions using these steps:

1. Load rules that match the incoming event type
2. For each rule, evaluate all conditions
3. Apply the logical operator (AND/OR) between conditions
4. If the overall evaluation is true, assign the badge

```go
// Pseudocode for rule evaluation
func EvaluateRules(event Event) {
    rules := GetRulesByEventType(event.Type)
    
    for _, rule := range rules {
        if EvaluateRule(rule, event) {
            AssignBadge(event.UserID, rule.BadgeID)
        }
    }
}
```

## Operators

The rule engine supports various operators for condition evaluation:

| Operator | Description | Example |
|----------|-------------|---------|
| `eq` | Equal | `field eq "value"` |
| `neq` | Not equal | `field neq "value"` |
| `gt` | Greater than | `field gt "5"` |
| `lt` | Less than | `field lt "10"` |
| `contains` | String contains | `field contains "substring"` |
| `startsWith` | String starts with | `field startsWith "prefix"` |
| `endsWith` | String ends with | `field endsWith "suffix"` |
| `regex` | Regular expression match | `field regex "pattern"` |

## Extending the Rule Engine

### Adding a New Operator

To add a new operator:

1. Define the operator function in `internal/rules/operators.go`
2. Register the operator in the `OperatorFunctions` map

```go
// Example: Adding a "between" operator
func between(value, expected string) (bool, error) {
    // Parse expected range (e.g., "5,10")
    ranges := strings.Split(expected, ",")
    if len(ranges) != 2 {
        return false, fmt.Errorf("between operator requires two values separated by comma")
    }
    
    // Implementation details...
    
    return result, nil
}

// Register in the map
OperatorFunctions["between"] = between
```

### Adding a New Rule Type

To add a new rule type:

1. Define the rule logic in `internal/rules/engine.go`
2. Update the database schema if necessary
3. Add API endpoints for managing the new rule type
4. Update the UI to support creating rules of this type

## Testing the Rule Engine

When testing the rule engine, follow these guidelines:

1. Create a mock event with appropriate test data
2. Define test rules with conditions targeting that data
3. Verify badge assignments occur when conditions are met
4. Test boundary conditions and edge cases

```go
// Example test for rule engine
func TestRuleEngine(t *testing.T) {
    // Create a test event
    event := Event{
        Type:   1,
        UserID: "test-user",
        Data:   map[string]string{"score": "100"},
    }
    
    // Define a test rule
    rule := Rule{
        ID:          1,
        EventTypeID: 1,
        BadgeID:     1,
        Conditions:  []Condition{
            {
                ID:            1,
                Operator:      "gt",
                ExpectedValue: "90",
                Field:         "score",
            },
        },
    }
    
    // Test rule evaluation
    result := EvaluateRule(rule, event)
    assert.True(t, result, "Rule should evaluate to true")
}
```

## Common Issues and Solutions

### Performance Considerations

- Rules are evaluated sequentially, so prioritize frequently used rules
- Complex conditions may impact performance, especially with large data sets
- Consider caching rule definitions for faster lookups

### Debugging Rules

- Use debug logs to trace rule evaluation steps
- Inspect the event data to ensure expected fields are present
- Verify operator behavior with test cases for each scenario

## Best Practices

1. **Keep rules simple**: Complex rules are harder to maintain and debug
2. **Document rule logic**: Add comments explaining the purpose of each rule
3. **Test thoroughly**: Cover edge cases and unexpected inputs
4. **Log rule evaluation**: Enable detailed logging during development
5. **Validate inputs**: Ensure events and conditions have valid data

## Core Components

1. **Events**: Actions performed by users that are tracked by the system
2. **Event Types**: Categories of events with specific schemas
3. **Condition Types**: Reusable evaluation logic written in JavaScript
4. **Badges**: Achievements awarded to users
5. **Badge Criteria**: Rules that determine when a badge should be awarded

## Badge Criteria Format

The rule engine supports several formats for defining badge criteria. The criteria are defined in the `FlowDefinition` field when creating a badge.

### Supported Formats

#### 1. Event-Based Criteria (Direct Format)

This is the simplest format, which directly specifies an event type and criteria to match against the event payload:

```json
{
  "event": "event_type_name",
  "criteria": {
    "field_name": {
      "$gte": 10
    }
  }
}
```

In this example:
- `event`: Specifies the event type name to match
- `criteria`: Defines conditions to check against the event payload
  - `field_name`: A field in the event payload to evaluate
  - `$gte`: Greater than or equal operator (supported operators include `$gt`, `$lt`, `$lte`, `$eq`, etc.)

#### 2. Logical Operators

The rule engine supports logical operators to combine multiple conditions:

##### $and Operator

```json
{
  "$and": [
    {
      "event": "event_type_name",
      "criteria": {
        "field_name": {
          "$gte": 10
        }
      }
    },
    {
      "event": "another_event_type",
      "criteria": {
        "another_field": {
          "$eq": "value"
        }
      }
    }
  ]
}
```

The `$and` operator requires all conditions to be true.

##### $or Operator

```json
{
  "$or": [
    {
      "event": "event_type_name",
      "criteria": {
        "field_name": {
          "$gte": 10
        }
      }
    },
    {
      "event": "another_event_type",
      "criteria": {
        "another_field": {
          "$eq": "value"
        }
      }
    }
  ]
}
```

The `$or` operator requires at least one condition to be true.

##### $not Operator

```json
{
  "$not": {
    "event": "event_type_name",
    "criteria": {
      "field_name": {
        "$lt": 5
      }
    }
  }
}
```

The `$not` operator negates the condition it contains.

#### 3. Time-Based Operators

For more complex scenarios, you can use time-based operators:

```json
{
  "$timeWindow": {
    "start": "2023-01-01T00:00:00Z",
    "end": "2023-01-31T23:59:59Z",
    "flow": {
      "$and": [
        {
          "event": "event_type_1"
        },
        {
          "event": "event_type_2"
        }
      ]
    }
  }
}
```

You can also use relative time windows:

```json
{
  "$timeWindow": {
    "last": "30d",  // Supports "d" (days), "w" (weeks), "m" (months), "q" (quarters), "y" (years)
    "businessDaysOnly": true,  // Optional: filter out weekends
    "flow": {
      "event": "login",
      "criteria": {
        "$eventCount": { "$gte": 5 }
      }
    }
  }
}
```

## Unsupported Formats

The following formats are **NOT** supported and will result in errors:

### Incorrect: "steps" Format

```json
{
  "steps": [
    {
      "conditionTypeId": 1,
      "eventTypeIds": [1],
      "conditionParams": {
        "threshold": 50
      }
    }
  ]
}
```

### Incorrect: "condition" Format

```json
{
  "condition": {
    "conditionTypeId": 1,
    "eventTypeIds": [1]
  }
}
```

## Debugging Badge Assignment Issues

If badges are not being assigned as expected, follow these troubleshooting steps:

### 1. Enable Debug Logging

Ensure debug logging is enabled to see detailed information about badge criteria evaluation:

```go
re.Logger.Debug("Evaluating flow definition: %v", flowDefinition)
```

### 2. Check Server Logs

Look for these specific messages in the server logs:

- `"Evaluating flow definition"` - Shows the flow definition being processed
- `"Unsupported operator"` - Indicates that an invalid operator was used
- `"Evaluating criteria against event"` - Shows the evaluation process

### 3. Verify Event Payload Structure

Make sure the event payload structure matches what the criteria expect:

```json
{
  "score": 50
}
```

### 4. Check Field Names and Types

Ensure the field names and data types in your criteria match those in your event payload:

- Field names are case-sensitive
- Data types should match (e.g., numeric comparisons require numeric field values)

## Best Practices

1. **Use Simple Criteria First**: Start with simple event-based criteria before using complex logical operators.
2. **Test with Lower Thresholds**: For testing, use lower threshold values to make it easier to trigger badge assignments.
3. **Include Detailed Logging**: Add logs in your tests to capture the complete flow of events and badge checks.
4. **Validate Event Payloads**: Ensure your event payloads comply with the defined schemas.

## Examples

### Simple Score-Based Badge

```json
{
  "event": "score_event",
  "criteria": {
    "score": {
      "$gte": 10
    }
  }
}
```

### Logical Combination Badge

```json
{
  "$and": [
    {
      "event": "login_event",
      "criteria": {
        "consecutive_days": {
          "$gte": 7
        }
      }
    },
    {
      "event": "purchase_event",
      "criteria": {
        "amount": {
          "$gte": 100
        }
      }
    }
  ]
}
```

### Time Window Badge

```json
{
  "$withinTimeWindow": {
    "window": "24h",
    "conditions": {
      "$and": [
        {
          "event": "visit_page_event",
          "criteria": {
            "page": {
              "$eq": "homepage"
            }
          }
        },
        {
          "event": "visit_page_event",
          "criteria": {
            "page": {
              "$eq": "product_page"
            }
          }
        },
        {
          "event": "purchase_event"
        }
      ]
    }
  }
}
```

## Event Processing Flow

When an event is processed, the following happens:

1. The event is stored in the database
2. The rule engine's `ProcessEvent` method is called
3. The engine retrieves all badges and their criteria
4. For each badge, the criteria are evaluated against user events
5. If criteria are met, the badge is awarded to the user

## Logical Operators

The rule engine supports logical operators for combining conditions:

- **$and**: All conditions must be true
- **$or**: At least one condition must be true
- **$not**: The condition must be false

Example of using the `$and` operator:

```json
{
  "$and": [
    {
      "event": "challenge_completed",
      "criteria": {
        "score": { "$gte": 90 }
      }
    },
    {
      "event": "login",
      "criteria": {
        "streak": { "$gte": 5 }
      }
    }
  ]
}
```

## Time-Based Operators

The engine supports various time-based operators:

- **$timePeriod**: Evaluates events over a period of time
- **$pattern**: Detects patterns in user behavior over time
- **$sequence**: Checks for a sequence of events
- **$gap**: Evaluates time gaps between events
- **$duration**: Checks total duration of activities
- **$timeWindow**: Applies criteria within a time window

## Testing Format

In our tests, we use a simplified format for badge criteria:

```json
{
  "steps": [
    {
      "type": "condition",
      "conditionTypeId": 123,
      "eventTypeIds": [456],
      "conditionParams": {
        "threshold": 90
      }
    }
  ]
}
```

This matches what we confirmed in the codebase, where a badge definition contains `FlowDefinition` with a `steps` array.

## Comparison with Test Structure

Our tests correctly use the expected format for badge criteria, with the proper `steps` array containing condition objects that reference condition types and event types by ID.

The test structure:
1. Creates an event type
2. Creates a condition type with JavaScript evaluation logic
3. Creates a badge with a flow definition referencing these types
4. Processes events that should trigger the badge
5. Checks if the user has been awarded the badge

If badges are not being assigned in tests despite correct setup, potential issues could be:

1. The rule engine's `ProcessEvent` implementation may not be properly evaluating criteria
2. There might be a mismatch between the event payload and what the condition expects
3. The criteria might not be correctly formatted for the rule engine
4. Event processing might be asynchronous and not completing during the test

## Troubleshooting

If badge assignment is not working as expected:

1. Enable debug logging in the rule engine: `re.SetLogLevel(logging.LogLevelDebug)`
2. Check the event payload structure matches what the condition expects
3. Verify the condition's JavaScript logic is evaluating as expected
4. Ensure all IDs are correctly referenced between event types, condition types, and badges
5. Check if any background processing is required for badge assignment

## Summary

The badge assignment system uses a powerful rule engine that evaluates user events against predefined criteria to award badges. The criteria are defined using a flexible flow-based format that can express complex conditions through logical operators and time-based evaluation.

The test structure correctly follows this format, but issue may lie in the rule engine implementation or event processing pipeline. 
