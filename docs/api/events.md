# Events API

This document provides comprehensive documentation for the Events API endpoints in the Badge Assignment System.

## Table of Contents
- [Create Event](#create-event)
- [Get User Events](#get-user-events)

## Create Event

Records a new event in the system. When an event is recorded, the system evaluates if the user qualifies for any badges based on this event.

**Endpoint:** `POST /api/v1/events`

**Request Body:**
```json
{
  "event_type": "check-in",
  "user_id": "user123",
  "payload": {
    "time": "08:45:00",
    "date": "2023-06-20",
    "location": "Main Office"
  },
  "timestamp": "2023-06-20T08:45:00Z"
}
```

**Required Fields:**
- `event_type`: The type of event (must be a registered event type name)
- `user_id`: Identifier for the user who performed the action
- `payload`: JSON object containing event data according to the event type schema

**Optional Fields:**
- `timestamp`: When the event occurred (ISO 8601 format, defaults to current time)

**Response:**
```json
{
  "id": 500,
  "event_type_id": 1,
  "user_id": "user123",
  "payload": {
    "time": "08:45:00",
    "date": "2023-06-20",
    "location": "Main Office"
  },
  "occurred_at": "2023-06-20T08:45:00Z"
}
```

**Response Fields:**
- `id`: Unique identifier for the event
- `event_type_id`: ID of the event type
- `user_id`: ID of the user who performed the action
- `payload`: Event data
- `occurred_at`: Timestamp when the event occurred

**Error Responses:**
- `400 Bad Request`: Invalid event data
- `404 Not Found`: Specified event type does not exist
- `422 Unprocessable Entity`: Event payload doesn't match the schema for the event type

## Get User Events

Retrieves events for a specific user.

**Endpoint:** `GET /api/v1/users/{user_id}/events`

**Path Parameters:**
- `user_id`: ID of the user whose events to retrieve

**Query Parameters:**
- `event_type`: Filter by event type ID or name
- `start_date`: Filter events after this date (ISO 8601 format)
- `end_date`: Filter events before this date (ISO 8601 format)
- `limit`: Maximum number of events to return (default: 50)
- `offset`: Pagination offset (default: 0)

**Response:**
```json
[
  {
    "id": 500,
    "event_type_id": 1,
    "event_type_name": "check-in",
    "user_id": "user123",
    "payload": {
      "time": "08:45:00",
      "date": "2023-06-20",
      "location": "Main Office"
    },
    "occurred_at": "2023-06-20T08:45:00Z"
  },
  {
    "id": 501,
    "event_type_id": 2,
    "event_type_name": "check-out",
    "user_id": "user123",
    "payload": {
      "time": "17:30:00",
      "date": "2023-06-20",
      "location": "Main Office"
    },
    "occurred_at": "2023-06-20T17:30:00Z"
  }
]
```

**Error Responses:**
- `404 Not Found`: User with the specified ID does not exist 
