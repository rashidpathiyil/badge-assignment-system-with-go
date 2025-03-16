# Badge Creation Example

This document provides a complete example of creating a badge using the Badge Assignment System API.

## Prerequisites

Before creating a badge, you need to ensure that the event types used in the badge criteria exist in the system. For this example, we'll create an "Early Bird" badge that requires the "check-in" event type.

## Step 1: Check if the Event Type Exists

First, check if the "check-in" event type exists:

```bash
curl -X GET "http://localhost:8080/api/v1/admin/event-types" | jq
```

## Step 2: Create the Event Type (if it doesn't exist)

If the event type doesn't exist, create it:

```bash
curl -X POST "http://localhost:8080/api/v1/admin/event-types" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "check-in",
    "description": "User check-in event, used for attendance tracking",
    "schema": {
      "type": "object",
      "properties": {
        "user_id": {
          "type": "string"
        },
        "time": {
          "type": "string",
          "format": "time",
          "description": "The time of check-in in HH:MM:SS format"
        },
        "date": {
          "type": "string",
          "format": "date"
        },
        "location": {
          "type": "string"
        }
      },
      "required": ["user_id", "time", "date"]
    }
  }'
```

## Step 3: Create the Badge

Now that the event type exists, create the "Early Bird" badge:

```bash
curl -X POST "http://localhost:8080/api/v1/admin/badges" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Early Bird",
    "description": "Checked in before 9 AM for 5 consecutive days",
    "image_url": "https://example.com/badges/early-bird.png",
    "flow_definition": {
      "event": "check-in",
      "criteria": {
        "$timePeriod": {
          "periodType": "day",
          "periodCount": { "$gte": 5 },
          "excludeWeekends": false
        },
        "payload": {
          "time": { "$lt": "09:00:00" }
        }
      }
    }
  }'
```

## Step 4: Verify Badge Creation

Verify that the badge was created successfully:

```bash
curl -X GET "http://localhost:8080/api/v1/badges" | jq
```

## Step 5: Trigger the Badge with Events

To test the badge, simulate a user checking in before 9 AM for several days:

```bash
# Day 1
curl -X POST "http://localhost:8080/api/v1/events" \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "check-in",
    "user_id": "user123",
    "payload": {
      "time": "08:45:00",
      "date": "2023-06-20",
      "location": "Main Office"
    },
    "timestamp": "2023-06-20T08:45:00Z"
  }'

# Day 2
curl -X POST "http://localhost:8080/api/v1/events" \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "check-in",
    "user_id": "user123",
    "payload": {
      "time": "08:30:00",
      "date": "2023-06-21",
      "location": "Main Office"
    },
    "timestamp": "2023-06-21T08:30:00Z"
  }'

# Continue for Days 3, 4, and 5...
```

## Step 6: Check if the User Earned the Badge

After submitting events for 5 consecutive days, check if the user earned the badge:

```bash
curl -X GET "http://localhost:8080/api/v1/users/user123/badges" | jq
```

If the user met the criteria, you should see the "Early Bird" badge in the response. 
