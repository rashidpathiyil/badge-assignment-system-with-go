# Pattern Detection Badge Example

This document provides a detailed example of creating a badge that uses pattern detection criteria to identify increasing user engagement.

## Overview

In this example, we'll create a "Fitness Growth" badge that detects an increasing pattern in workout frequency over time. This badge is awarded to users who show a consistent increase in their workout frequency over several weeks.

## Step 1: Create the Workout Event Type

First, create the "workout-completed" event type:

```bash
curl -X POST "http://localhost:8080/api/v1/admin/event-types" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "workout-completed",
    "description": "Tracks completion of a workout session",
    "schema": {
      "type": "object",
      "properties": {
        "duration_minutes": {
          "type": "number",
          "description": "Duration of the workout in minutes"
        },
        "calories_burned": {
          "type": "number"
        },
        "workout_type": {
          "type": "string"
        }
      },
      "required": ["duration_minutes"]
    }
  }'
```

## Step 2: Create the Fitness Growth Badge

Now, create the badge with pattern detection criteria:

```bash
curl -X POST "http://localhost:8080/api/v1/admin/badges" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Fitness Growth",
    "description": "Awarded for increasing workout frequency",
    "image_url": "https://example.com/badges/fitness-growth.png",
    "flow_definition": {
      "event": "workout-completed",
      "criteria": {
        "$pattern": {
          "type": "increasing",
          "periodType": "week",
          "min_periods": 6,
          "min_percent_increase": 10
        }
      }
    }
  }'
```

This badge requires at least 6 weeks of data with an average increase of 10% in workout frequency.

## Step 3: Simulate User Workout Data

Let's simulate a user who gradually increases their workout frequency over 8 weeks:

```bash
# Week 1: 2 workouts
# Week 2: 3 workouts
# Week 3: 3 workouts
# Week 4: 4 workouts
# Week 5: 4 workouts
# Week 6: 5 workouts
# Week 7: 5 workouts
# Week 8: 6 workouts
```

For each workout, submit an event similar to:

```bash
curl -X POST "http://localhost:8080/api/v1/events" \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "workout-completed",
    "user_id": "fitness-user",
    "payload": {
      "duration_minutes": 30,
      "calories_burned": 150,
      "workout_type": "cardio"
    },
    "timestamp": "2023-06-01T10:00:00Z"
  }'
```

## Step 4: Check if the Badge is Awarded

After submitting all the events, check if the user earned the badge:

```bash
curl -X GET "http://localhost:8080/api/v1/users/fitness-user/badges" | jq
``` 
