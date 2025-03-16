# Condition Types API

This document provides comprehensive documentation for the Condition Types API endpoints in the Badge Assignment System.

## Table of Contents
- [Create Condition Type](#create-condition-type)
- [List Condition Types](#list-condition-types)
- [Get Condition Type Details](#get-condition-type-details)
- [Update Condition Type](#update-condition-type)
- [Delete Condition Type](#delete-condition-type)

## Overview

Condition Types define reusable conditions that can be referenced in badge criteria. They provide a way to create complex criteria that can be reused across multiple badges.

## Create Condition Type

Creates a new condition type in the system.

**Endpoint:** `POST /api/v1/admin/condition-types`

**Request Body:**
```json
{
  "name": "WorkHoursCondition",
  "description": "Checks if a user has logged a minimum number of work hours",
  "evaluation_logic": "function evaluateCondition(events, context) { return events.reduce((sum, event) => sum + event.payload.hours, 0) >= 40; }"
}
```

**Required Fields:**
- `name`: Name of the condition type
- `evaluation_logic`: JavaScript function that will be evaluated to determine if the condition is met

**Optional Fields:**
- `description`: Description of the condition type

**Response:**
```json
{
  "id": 1,
  "name": "WorkHoursCondition",
  "description": "Checks if a user has logged a minimum number of work hours",
  "evaluation_logic": "function evaluateCondition(events, context) { return events.reduce((sum, event) => sum + event.payload.hours, 0) >= 40; }",
  "created_at": "2023-06-14T09:00:00Z",
  "updated_at": "2023-06-14T09:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid condition type data
- `409 Conflict`: Condition type with the same name already exists

## List Condition Types

Retrieves a list of all condition types in the system.

**Endpoint:** `GET /api/v1/admin/condition-types`

**Response:**
```json
[
  {
    "id": 1,
    "name": "WorkHoursCondition",
    "description": "Checks if a user has logged a minimum number of work hours",
    "evaluation_logic": "function evaluateCondition(events, context) { return events.reduce((sum, event) => sum + event.payload.hours, 0) >= 40; }",
    "created_at": "2023-06-14T09:00:00Z",
    "updated_at": "2023-06-14T09:00:00Z"
  },
  {
    "id": 2,
    "name": "ConsistentAttendance",
    "description": "Checks if a user has attended consistently over a period",
    "evaluation_logic": "function evaluateCondition(events, context) { /* logic for consistent attendance */ }",
    "created_at": "2023-06-14T09:05:00Z",
    "updated_at": "2023-06-14T09:05:00Z"
  }
]
```

## Get Condition Type Details

Retrieves details about a specific condition type.

**Endpoint:** `GET /api/v1/admin/condition-types/{condition_type_id}`

**Path Parameters:**
- `condition_type_id`: The ID of the condition type to retrieve

**Response:**
```json
{
  "id": 1,
  "name": "WorkHoursCondition",
  "description": "Checks if a user has logged a minimum number of work hours",
  "evaluation_logic": "function evaluateCondition(events, context) { return events.reduce((sum, event) => sum + event.payload.hours, 0) >= 40; }",
  "created_at": "2023-06-14T09:00:00Z",
  "updated_at": "2023-06-14T09:00:00Z"
}
```

**Error Responses:**
- `404 Not Found`: Condition type with the specified ID does not exist

## Update Condition Type

Updates an existing condition type.

**Endpoint:** `PUT /api/v1/admin/condition-types/{condition_type_id}`

**Path Parameters:**
- `condition_type_id`: The ID of the condition type to update

**Request Body:**
```json
{
  "name": "WorkHoursCondition",
  "description": "Updated description for work hours condition",
  "evaluation_logic": "function evaluateCondition(events, context) { return events.reduce((sum, event) => sum + event.payload.hours, 0) >= 35; }"
}
```

**Response:** Same format as Get Condition Type Details

**Error Responses:**
- `404 Not Found`: Condition type with the specified ID does not exist
- `400 Bad Request`: Invalid condition type data

## Delete Condition Type

Deletes a condition type from the system.

**Endpoint:** `DELETE /api/v1/admin/condition-types/{condition_type_id}`

**Path Parameters:**
- `condition_type_id`: The ID of the condition type to delete

**Response:** HTTP 200 OK

**Error Responses:**
- `404 Not Found`: Condition type with the specified ID does not exist
- `409 Conflict`: Condition type cannot be deleted because it is referenced by existing badges
