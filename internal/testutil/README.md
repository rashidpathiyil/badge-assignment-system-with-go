# Test Utilities Package

This package contains utilities for testing the Badge Assignment System.

## Contents

- `models.go` - Test models and data structures
- `helpers.go` - Helper functions for API testing

## Usage

### Making API Requests

```go
import (
    "net/http"
    "testing"
    
    "github.com/badge-assignment-system/internal/testutil"
)

func TestExample(t *testing.T) {
    // Make a GET request
    response := testutil.MakeRequest(http.MethodGet, "/api/v1/badges", nil)
    testutil.AssertSuccess(t, response)
    
    // Parse the response into a struct
    var badges []testutil.Badge
    err := testutil.ParseResponse(response, &badges)
    if err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }
}
```

### Creating Test Models

```go
// Create a badge request
badgeReq := testutil.BadgeRequest{
    Name:        "Test Badge",
    Description: "This is a test badge",
    ImageURL:    "https://example.com/badge.png",
    FlowDefinition: map[string]interface{}{
        // Your flow definition here
    },
}

// Make a POST request with the badge
response := testutil.MakeRequest(http.MethodPost, "/api/v1/admin/badges", badgeReq)
```

### Configuration

The test utilities use environment variables for configuration:

- `API_TEST_URL` - Base URL for API requests (default: http://localhost:8080)

You can set these variables before running the tests:

```bash
export API_TEST_URL=http://localhost:8081
go test ./...
```

## Extending

If you need to add new test utilities:

1. Add them to the appropriate file (`models.go` for data structures, `helpers.go` for functions)
2. Document any new functions or types with comments
3. If adding new configuration options, use environment variables with sensible defaults 
