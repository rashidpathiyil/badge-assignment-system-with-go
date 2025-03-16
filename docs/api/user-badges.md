# User Badges API

This document provides comprehensive documentation for the User Badges API endpoints in the Badge Assignment System.

## Table of Contents
- [Get User Badges](#get-user-badges)
- [Evaluate User for Badges](#evaluate-user-for-badges) (Planned Feature)

## Get User Badges

Retrieves all badges that have been awarded to a specific user.

**Endpoint:** `GET /api/v1/users/{user_id}/badges`

**Path Parameters:**
- `user_id`: ID of the user whose badges to retrieve

**Query Parameters:**
- `limit`: Maximum number of badges to return (default: 20)
- `offset`: Pagination offset (default: 0)

**Response:**
```json
[
  {
    "id": 123,
    "name": "Early Bird",
    "description": "Checked in before 9 AM for 5 consecutive days",
    "image_url": "https://example.com/badges/early-bird.png",
    "awarded_at": "2023-06-20T08:50:00Z",
    "metadata": {
      "qualifying_events": 5,
      "consecutive_days": 5
    }
  }
]
```

**Response Fields:**
- `id`: ID of the badge
- `name`: Name of the badge
- `description`: Description of the badge
- `image_url`: URL to the badge image
- `awarded_at`: Timestamp when the badge was awarded
- `metadata`: Additional information about how the badge was awarded (varies by badge type)

**Error Responses:**
- `404 Not Found`: User with the specified ID does not exist

## Evaluate User for Badges

> **Note:** This endpoint is documented as a planned feature and has not been implemented in the current API version.

Manually triggers the evaluation process to check if a user qualifies for any badges.

**Endpoint:** `POST /api/v1/users/{user_id}/evaluate-badges`

**Path Parameters:**
- `user_id`: ID of the user to evaluate

**Response:**
```json
{
  "awarded_badges": [
    {
      "id": 123,
      "name": "Early Bird",
      "awarded_at": "2023-06-20T08:50:00Z"
    }
  ],
  "progress": [
    {
      "badge_id": 124,
      "name": "Workaholic",
      "progress": 0.6,
      "requirements": {
        "completed": 24,
        "total": 40
      }
    }
  ]
}
```

**Response Fields:**
- `awarded_badges`: List of badges newly awarded to the user
- `progress`: Progress information for badges the user has not yet earned

**Error Responses:**
- `404 Not Found`: User with the specified ID does not exist 
