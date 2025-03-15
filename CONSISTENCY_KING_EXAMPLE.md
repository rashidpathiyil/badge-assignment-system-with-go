# Consistency King Badge Example

This document provides a detailed example of how to create and test the "Consistency King" badge using time-based criteria from the Badge Assignment System.

## Badge Definition

The "Consistency King" badge is awarded to users who check in and out consistently every weekday for a rolling 28-day period, without missing any workdays.

## Implementation Using Time-Based Criteria

### Badge Configuration

The Consistency King badge can be created with the following criteria:

```json
{
  "name": "Consistency King",
  "description": "Checked in & out without missing a weekday for a 28-day period.",
  "image_url": "https://example.com/badges/consistency-king.png",
  "flow_definition": {
    "$timePeriod": {
      "periodType": "day",
      "count": { "$gte": 20 },
      "excludeWeekends": true,
      "lookbackDays": 28,
      "events": [
        {"type": "check-in", "required": true},
        {"type": "check-out", "required": true}
      ]
    }
  }
}
```

This configuration uses the `$timePeriod` criterion to:
- Look at daily periods
- Require at least 20 unique days (workdays in a 28-day period)
- Exclude weekends from the count
- Require both check-in and check-out events on each counted day

### API Example

Here's how to create this badge using the system's API:

```bash
curl -X POST http://localhost:8080/api/v1/admin/badges \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Consistency King",
    "description": "Checked in & out without missing a weekday for a 28-day period.",
    "image_url": "https://example.com/badges/consistency-king.png",
    "flow_definition": {
      "$timePeriod": {
        "periodType": "day",
        "count": { "$gte": 20 },
        "excludeWeekends": true,
        "lookbackDays": 28,
        "events": [
          {"type": "check-in", "required": true},
          {"type": "check-out", "required": true}
        ]
      }
    }
  }'
```

## Testing the Badge

### Prerequisites

Before creating the badge, you need to:

1. Create the `check-in` and `check-out` event types:

```bash
# Create check-in event type
curl -X POST http://localhost:8080/api/v1/admin/event-types \
  -H "Content-Type: application/json" \
  -d '{
    "name": "check-in",
    "description": "User check-in event",
    "schema": {
      "type": "object",
      "properties": {
        "timestamp": {
          "type": "string",
          "format": "date-time"
        }
      }
    }
  }'

# Create check-out event type
curl -X POST http://localhost:8080/api/v1/admin/event-types \
  -H "Content-Type: application/json" \
  -d '{
    "name": "check-out",
    "description": "User check-out event",
    "schema": {
      "type": "object",
      "properties": {
        "timestamp": {
          "type": "string",
          "format": "date-time"
        }
      }
    }
  }'
```

### Creating Test Data

To test the badge, you need to create check-in and check-out events for a user spanning multiple weekdays:

```bash
# Generate check-ins and check-outs for the last 28 days, skipping weekends
for day in {0..27}; do 
  date_stamp=$(date -v-${day}d "+%Y-%m-%d")
  day_of_week=$(date -v-${day}d "+%u")
  
  # Only create events for weekdays (1-5)
  if [ "$day_of_week" -lt 6 ]; then 
    check_in="${date_stamp}T08:30:00Z"
    check_out="${date_stamp}T17:30:00Z"
    
    # Create check-in event
    curl -s -X POST http://localhost:8080/api/v1/events \
      -H "Content-Type: application/json" \
      -d "{\"event_type\":\"check-in\",\"user_id\":\"test-user\",\"payload\":{\"timestamp\":\"$check_in\"},\"timestamp\":\"$check_in\"}"
      
    # Create check-out event
    curl -s -X POST http://localhost:8080/api/v1/events \
      -H "Content-Type: application/json" \
      -d "{\"event_type\":\"check-out\",\"user_id\":\"test-user\",\"payload\":{\"timestamp\":\"$check_out\"},\"timestamp\":\"$check_out\"}"
  fi
done
```

### Checking Badge Award

To check if the user has been awarded the badge:

```bash
curl http://localhost:8080/api/v1/users/test-user/badges
```

## Unit Testing Time Period Criteria

Here's how the time period criteria could be unit tested specifically for the Consistency King badge:

```go
func TestConsistencyKingBadgeCriteria(t *testing.T) {
    // Create a RuleEngine instance
    re := new(RuleEngine)
    
    // Create events for 20 weekdays
    startTime := time.Now().AddDate(0, 0, -28)
    var events []models.Event
    
    // Event types
    checkInTypeID := 1
    checkOutTypeID := 2
    
    // Create 28 days of events
    for day := 0; day < 28; day++ {
        currentDate := startTime.AddDate(0, 0, day)
        
        // Skip weekends
        if currentDate.Weekday() == time.Saturday || currentDate.Weekday() == time.Sunday {
            continue
        }
        
        // Morning check-in
        checkInTime := time.Date(
            currentDate.Year(), 
            currentDate.Month(), 
            currentDate.Day(), 
            8, 30, 0, 0, 
            time.UTC,
        )
        
        // Evening check-out
        checkOutTime := time.Date(
            currentDate.Year(), 
            currentDate.Month(), 
            currentDate.Day(), 
            17, 30, 0, 0, 
            time.UTC,
        )
        
        // Add check-in event
        events = append(events, models.Event{
            ID:          len(events) + 1,
            UserID:      "test-user",
            EventTypeID: checkInTypeID,
            OccurredAt:  checkInTime,
            Payload:     models.JSONB{"timestamp": checkInTime.Format(time.RFC3339)},
        })
        
        // Add check-out event
        events = append(events, models.Event{
            ID:          len(events) + 1,
            UserID:      "test-user",
            EventTypeID: checkOutTypeID,
            OccurredAt:  checkOutTime,
            Payload:     models.JSONB{"timestamp": checkOutTime.Format(time.RFC3339)},
        })
    }
    
    // Define the Consistency King criteria
    criteria := map[string]interface{}{
        "periodType":      "day",
        "count":           map[string]interface{}{"$gte": float64(20)},
        "excludeWeekends": true,
        "lookbackDays":    float64(28),
        "events": []interface{}{
            map[string]interface{}{
                "type":     "check-in",
                "required": true,
            },
            map[string]interface{}{
                "type":     "check-out",
                "required": true,
            },
        },
    }
    
    // Run the test
    metadata := make(map[string]interface{})
    result, err := re.evaluateTimePeriodCriteria(criteria, events, metadata)
    
    // Verify results
    if err != nil {
        t.Errorf("Error evaluating time period criteria: %v", err)
    }
    
    if !result {
        t.Error("Expected criteria to be met, but it wasn't")
    }
    
    // Check metadata
    uniqueDays, ok := metadata["unique_period_count"].(int)
    if !ok || uniqueDays < 20 {
        t.Errorf("Expected at least 20 unique days with events, got %v", uniqueDays)
    }
}
```

## Troubleshooting

If the Consistency King badge is not being awarded as expected, check:

1. **Event Data**: Ensure there are sufficient check-ins and check-outs on weekdays
2. **Criteria Format**: Verify the criteria format matches the expected structure
3. **Parameter Types**: Ensure numeric values are properly typed as floats
4. **Time Zone Handling**: Be aware of potential time zone issues affecting day calculations

## Extensions

The Consistency King badge concept can be extended with additional criteria:

- Add a requirement for check-ins before a certain time
- Require a minimum duration between check-in and check-out
- Allow for a limited number of "grace days" where check-ins or check-outs can be missed
- Implement seasonality calculations for different requirements based on the time of year 
