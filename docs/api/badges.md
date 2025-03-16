# Badge API

This document provides comprehensive documentation for the Badge API endpoints in the Badge Assignment System.

## Table of Contents
- [Public Endpoints](#public-endpoints)
  - [List Badges](#list-badges)
  - [List Active Badges](#list-active-badges)
  - [Get Badge Details](#get-badge-details)
- [Admin Endpoints](#admin-endpoints)
  - [Create Badge](#create-badge)
  - [Update Badge](#update-badge)
  - [Get Badge with Criteria](#get-badge-with-criteria)
  - [Delete Badge](#delete-badge)

## Public Endpoints

These endpoints are available to all clients without authentication.

### List Badges

Retrieves a list of all badges in the system.

**Endpoint:** `GET /api/v1/badges`

**Query Parameters:**
- `limit`: Maximum number of badges to return (default: 20)
- `offset`: Pagination offset (default: 0)
- `active`: Filter by active status (true, false, all)

**Response:**
```json
[
  {
    "id": 123,
    "name": "Early Bird",
    "description": "Checked in before 9 AM for 5 consecutive days",
    "image_url": "https://example.com/badges/early-bird.png",
    "active": true,
    "created_at": "2023-06-15T10:30:00Z",
    "updated_at": "2023-06-15T10:30:00Z"
  },
  {
    "id": 124,
    "name": "Workaholic",
    "description": "Logged 40+ hours in a week",
    "image_url": "https://example.com/badges/workaholic.png",
    "active": true,
    "created_at": "2023-06-16T14:20:00Z",
    "updated_at": "2023-06-16T14:20:00Z"
  }
]
```

### List Active Badges

Retrieves a list of all active badges in the system.

**Endpoint:** `GET /api/v1/badges/active`

**Query Parameters:**
- `limit`: Maximum number of badges to return (default: 20)
- `offset`: Pagination offset (default: 0)

**Response:**
Same format as the List Badges endpoint, but only includes badges with `active: true`.

### Get Badge Details

Retrieves details about a specific badge.

**Endpoint:** `GET /api/v1/badges/{badge_id}`

**Path Parameters:**
- `badge_id`: The ID of the badge to retrieve

**Response:**
```json
{
  "id": 123,
  "name": "Early Bird",
  "description": "Checked in before 9 AM for 5 consecutive days",
  "image_url": "https://example.com/badges/early-bird.png",
  "active": true,
  "created_at": "2023-06-15T10:30:00Z",
  "updated_at": "2023-06-15T10:30:00Z"
}
```

**Error Responses:**
- `404 Not Found`: Badge with the specified ID does not exist

## Admin Endpoints

These endpoints are for administrative operations and require appropriate authentication.

### Create Badge

Creates a new badge in the system.

**Endpoint:** `POST /api/v1/admin/badges`

**Request Body:**
```json
{
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
}
```

**Required Fields:**
- `name`: Name of the badge
- `description`: Description of the badge
- `flow_definition`: Rules that determine when this badge is awarded

**Optional Fields:**
- `image_url`: URL to the badge image
- `active`: Badge active status (default: true)

**Response:**
```json
{
  "badge": {
    "id": 123,
    "name": "Early Bird",
    "description": "Checked in before 9 AM for 5 consecutive days",
    "image_url": "https://example.com/badges/early-bird.png",
    "active": true,
    "created_at": "2023-06-15T10:30:00Z",
    "updated_at": "2023-06-15T10:30:00Z"
  },
  "criteria": {
    "id": 456,
    "badge_id": 123,
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
    },
    "created_at": "2023-06-15T10:30:00Z",
    "updated_at": "2023-06-15T10:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid badge data or criteria
- `409 Conflict`: Badge with the same name already exists

### Update Badge

Updates an existing badge.

**Endpoint:** `PUT /api/v1/admin/badges/{badge_id}`

**Path Parameters:**
- `badge_id`: The ID of the badge to update

**Request Body:** Same as Create Badge

**Response:** Same as Create Badge

**Error Responses:**
- `404 Not Found`: Badge with the specified ID does not exist
- `400 Bad Request`: Invalid badge data or criteria

### Get Badge with Criteria

Retrieves a badge along with its criteria.

**Endpoint:** `GET /api/v1/admin/badges/{badge_id}/criteria`

**Path Parameters:**
- `badge_id`: The ID of the badge to retrieve

**Response:** Same as the response from Create Badge

**Error Responses:**
- `404 Not Found`: Badge with the specified ID does not exist

### Delete Badge

Deletes a badge from the system.

**Endpoint:** `DELETE /api/v1/admin/badges/{badge_id}`

**Path Parameters:**
- `badge_id`: The ID of the badge to delete

**Response:** HTTP 200 OK

**Error Responses:**
- `404 Not Found`: Badge with the specified ID does not exist
- `409 Conflict`: Badge cannot be deleted because it is already awarded to users 
