# API Error Handling

This document provides comprehensive documentation for error handling in the Badge Assignment System API.

## Error Response Format

All API endpoints return standard HTTP status codes along with a structured JSON error response:

```json
{
  "error": {
    "code": "error_code",
    "message": "A human-readable error message",
    "details": {
      "field": "specific_field",
      "value": "problematic_value"
    }
  }
}
```

## Common HTTP Status Codes

| Status Code | Description | Common Causes |
|-------------|-------------|--------------|
| `200 OK` | The request was successful | - |
| `201 Created` | The resource was created successfully | POST requests that create a new resource |
| `400 Bad Request` | The request was invalid | Missing required fields, invalid data format |
| `401 Unauthorized` | Authentication is required | Missing or invalid authentication token |
| `403 Forbidden` | The client does not have permission | Attempting to access admin endpoints without admin privileges |
| `404 Not Found` | The resource was not found | Requesting a non-existent badge, event type, or user |
| `409 Conflict` | The request conflicts with the current state | Creating a duplicate resource |
| `422 Unprocessable Entity` | The request data failed validation | Event payload doesn't match the schema for the event type |
| `500 Internal Server Error` | An unexpected error occurred | Server-side issues |

## Error Codes

### General Errors

| Error Code | Description | HTTP Status |
|------------|-------------|-------------|
| `invalid_request` | The request format is invalid | 400 |
| `resource_not_found` | The requested resource doesn't exist | 404 |
| `unauthorized` | Authentication is required | 401 |
| `forbidden` | The client doesn't have permission | 403 |
| `internal_error` | An unexpected server error occurred | 500 |

### Badge-Specific Errors

| Error Code | Description | HTTP Status |
|------------|-------------|-------------|
| `invalid_badge_data` | The badge data is invalid | 400 |
| `badge_not_found` | The requested badge doesn't exist | 404 |
| `duplicate_badge` | A badge with the same name already exists | 409 |
| `invalid_criteria` | The badge criteria format is invalid | 400 |

### Event-Specific Errors

| Error Code | Description | HTTP Status |
|------------|-------------|-------------|
| `invalid_event_data` | The event data is invalid | 400 |
| `event_type_not_found` | The specified event type doesn't exist | 404 |
| `invalid_payload` | The event payload doesn't match the schema | 422 |

### Event Type-Specific Errors

| Error Code | Description | HTTP Status |
|------------|-------------|-------------|
| `invalid_event_type_data` | The event type data is invalid | 400 |
| `event_type_not_found` | The requested event type doesn't exist | 404 |
| `duplicate_event_type` | An event type with the same name already exists | 409 |
| `invalid_schema` | The event type schema is invalid | 400 |
| `in_use_event_type` | Cannot delete an event type that is used by badges | 409 |

### Condition Type-Specific Errors

| Error Code | Description | HTTP Status |
|------------|-------------|-------------|
| `invalid_condition_type_data` | The condition type data is invalid | 400 |
| `condition_type_not_found` | The requested condition type doesn't exist | 404 |
| `duplicate_condition_type` | A condition type with the same name already exists | 409 |
| `in_use_condition_type` | Cannot delete a condition type that is used by badges | 409 |

## Examples

### Resource Not Found Example

**Request:**
```
GET /api/v1/badges/9999
```

**Response:**
```
HTTP/1.1 404 Not Found
Content-Type: application/json

{
  "error": {
    "code": "badge_not_found",
    "message": "Badge with ID 9999 not found"
  }
}
```

### Invalid Request Body Example

**Request:**
```
POST /api/v1/admin/badges
Content-Type: application/json

{
  "name": "Early Bird",
  "flow_definition": {
    "invalid_field": "value"
  }
}
```

**Response:**
```
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "error": {
    "code": "invalid_criteria",
    "message": "Badge criteria contains invalid operator",
    "details": {
      "field": "flow_definition.invalid_field",
      "value": "value"
    }
  }
}
```

### Duplicate Resource Example

**Request:**
```
POST /api/v1/admin/event-types
Content-Type: application/json

{
  "name": "Check In",
  "description": "Event type that already exists"
}
```

**Response:**
```
HTTP/1.1 409 Conflict
Content-Type: application/json

{
  "error": {
    "code": "duplicate_event_type",
    "message": "Event type with name 'Check In' already exists"
  }
}
``` 
