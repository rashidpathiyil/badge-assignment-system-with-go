# Event Types API

This document provides comprehensive documentation for the Event Types API endpoints in the Badge Assignment System.

## Table of Contents
- [Create Event Type](#create-event-type)
- [List Event Types](#list-event-types)
- [Get Event Type Details](#get-event-type-details)
- [Update Event Type](#update-event-type)
- [Delete Event Type](#delete-event-type)

## Create Event Type

Creates a new event type in the system.

**Endpoint:** `POST /api/v1/admin/event-types`

**Request Body:**
```json
{
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
}
```

**Required Fields:**
- `name`: Name of the event type
- `schema`: JSON Schema defining the structure of the event payload

**Optional Fields:**
- `description`: Description of the event type

**Response:**
```json
{
  "id": 1,
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
  },
  "created_at": "2023-06-14T09:00:00Z",
  "updated_at": "2023-06-14T09:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid event type data
- `409 Conflict`: Event type with the same name already exists

## List Event Types

Retrieves a list of all event types in the system.

**Endpoint:** `GET /api/v1/admin/event-types`

**Response:**
```json
[
  {
    "id": 1,
    "name": "check-in",
    "description": "User check-in event, used for attendance tracking",
    "schema": { "..." },
    "created_at": "2023-06-14T09:00:00Z",
    "updated_at": "2023-06-14T09:00:00Z"
  },
  {
    "id": 2,
    "name": "check-out",
    "description": "User check-out event, used for attendance tracking",
    "schema": { "..." },
    "created_at": "2023-06-14T09:05:00Z",
    "updated_at": "2023-06-14T09:05:00Z"
  }
]
```

## Get Event Type Details

Retrieves details about a specific event type.

**Endpoint:** `GET /api/v1/admin/event-types/{event_type_id}`

**Path Parameters:**
- `event_type_id`: The ID of the event type to retrieve

**Response:**
```json
{
  "id": 1,
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
  },
  "created_at": "2023-06-14T09:00:00Z",
  "updated_at": "2023-06-14T09:00:00Z"
}
```

**Error Responses:**
- `404 Not Found`: Event type with the specified ID does not exist

## Update Event Type

Updates an existing event type.

**Endpoint:** `PUT /api/v1/admin/event-types/{event_type_id}`

**Path Parameters:**
- `event_type_id`: The ID of the event type to update

**Request Body:**
```json
{
  "name": "check-in",
  "description": "Updated description for check-in events",
  "schema": {
    "type": "object",
    "properties": {
      "user_id": {
        "type": "string"
      },
      "time": {
        "type": "string",
        "format": "time"
      },
      "date": {
        "type": "string",
        "format": "date"
      },
      "location": {
        "type": "string"
      },
      "device_id": {
        "type": "string",
        "description": "ID of the device used for check-in"
      }
    },
    "required": ["user_id", "time", "date"]
  }
}
```

**Response:** Same format as Get Event Type Details

**Error Responses:**
- `404 Not Found`: Event type with the specified ID does not exist
- `400 Bad Request`: Invalid event type data

## Delete Event Type

Deletes an event type from the system.

**Endpoint:** `DELETE /api/v1/admin/event-types/{event_type_id}`

**Path Parameters:**
- `event_type_id`: The ID of the event type to delete

**Response:** HTTP 200 OK

**Error Responses:**
- `404 Not Found`: Event type with the specified ID does not exist
- `409 Conflict`: Event type cannot be deleted because it is referenced by existing badges 
