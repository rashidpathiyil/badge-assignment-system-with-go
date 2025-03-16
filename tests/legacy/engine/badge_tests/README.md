# Badge System API Testing

This directory contains API tests for the Badge Assignment System. These tests validate badge criteria functionality using the system's API endpoints.

## Overview

The Badge Assignment System supports RESTful API endpoints for:
- Creating event types
- Processing events
- Checking user badges
- Managing badges and their criteria

The API tests in this directory focus on validating the "Early Bird" badge, which is awarded to users who check in before 9 AM for at least 5 consecutive days.

## Test Files

- `early_bird_api_test.sh`: Bash script for testing the Early Bird badge via API calls
- `early_bird_api_test.go`: Go equivalent of the bash script, integrated with Go's testing framework

## Prerequisites

Before running the tests, make sure:

1. The Badge Assignment System database is running (PostgreSQL)
2. The API server is running (`go run cmd/server/main.go`)
3. The "Early Bird" badge is defined in the system
4. The "Check In" event type is defined in the system

## Running the Tests

### Shell Script Test

To run the shell script test:

```bash
cd internal/engine/tests/badge_tests
chmod +x early_bird_api_test.sh
./early_bird_api_test.sh
```

### Go Test

To run the Go test:

```bash
cd internal/engine/tests/badge_tests
go test -run TestEarlyBirdBadgeAPI -v
```

## Test Cases

Both test files implement the same test cases:

1. **Test Case 1: Early Check-ins for Five Days**
   - Creates a user with 5 check-in events before 9 AM
   - Verifies the user is awarded the Early Bird badge

2. **Test Case 2: Mixed Check-in Times**
   - Creates a user with 3 check-in events before 9 AM and 2 after 9 AM
   - Verifies the user is NOT awarded the Early Bird badge

## Expected Output

### Successful Tests

For the shell script, you should see:
```
=== Early Bird Badge API Test ===
...
All tests completed successfully!
```

For the Go test, you should see:
```
=== RUN   TestEarlyBirdBadgeAPI
=== RUN   TestEarlyBirdBadgeAPI/EarlyCheckInsForFiveDaysAPI
    early_bird_api_test.go:61: User successfully awarded the Early Bird badge
=== RUN   TestEarlyBirdBadgeAPI/MixedCheckInTimesAPI
    early_bird_api_test.go:97: User correctly NOT awarded the Early Bird badge
--- PASS: TestEarlyBirdBadgeAPI
```

## Troubleshooting

If the tests fail:

1. **API server not running**: Ensure the server is running on the expected URL/port.
2. **Badge definition issue**: Verify the Early Bird badge is correctly defined with criteria requiring check-ins before 9 AM for at least 5 consecutive days.
3. **Event type issue**: Ensure the "Check In" event type is defined and matches exactly what is expected in the tests (case-sensitive).
4. **Database connection**: Check if the server can connect to the database.

## Customization

- You can modify `API_URL` in the shell script or `apiURL` in the Go test to point to a different server.
- User IDs can be changed to test with different test users.
- Check-in times can be modified to test different scenarios. 
