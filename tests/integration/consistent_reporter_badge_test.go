package integration

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/testutil"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

// TestConsistentReporterBadge tests badge assignment using $timePeriod with periodCount:
// 1. Create an event type for issue reporting
// 2. Create a badge that counts unique days with reported issues
// 3. Submit events on multiple different days
// 4. Verify the badge is awarded after events on 3 different days
func TestConsistentReporterBadge(t *testing.T) {
	SetupTest()

	// Create a DB connection for direct database verification
	dbConnStr := "postgres://rashidpathiyil:2426@localhost:5432/badge_system?sslmode=disable"
	db, err := sqlx.Connect("postgres", dbConnStr)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Verify database connection
	err = db.Ping()
	if err != nil {
		t.Fatalf("Could not ping database: %v", err)
	}
	t.Log("Successfully connected to test database")

	// Generate unique test identifiers
	timestamp := time.Now().UnixNano() / 1000000
	eventTypeName := fmt.Sprintf("issue_reported_time_%d", timestamp)
	badgeName := fmt.Sprintf("Consistent Reporter Badge_%d", timestamp)
	testUserID := "test_user_consistent" // Consistent user ID for this test

	// STEP 1: Create an event type for issue reporting
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"issue_id": map[string]interface{}{
				"type": "string",
			},
			"description": map[string]interface{}{
				"type": "string",
			},
			"severity": map[string]interface{}{
				"type": "string",
				"enum": []string{"low", "medium", "high"},
			},
		},
		"required": []string{"issue_id", "severity"},
	}

	eventTypeReq := map[string]interface{}{
		"name":        eventTypeName,
		"description": "Issue reporting event for time-based test",
		"schema":      schema,
	}

	resp := testutil.MakeRequest("POST", "/api/v1/admin/event-types", eventTypeReq)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("Failed to create event type: %s", string(resp.Body))
	}

	// Extract event type ID
	var eventTypeResp map[string]interface{}
	err = testutil.ParseResponse(resp, &eventTypeResp)
	assert.NoError(t, err)

	eventTypeID := int(eventTypeResp["id"].(float64))
	t.Logf("Created event type with ID: %d", eventTypeID)

	// STEP 2: Create a badge with criteria using $timePeriod
	// Define badge criteria that places $timePeriod at the top level, matching the docs
	flowDefinition := map[string]interface{}{
		"$timePeriod": map[string]interface{}{
			"periodType": "day",
			"periodCount": map[string]interface{}{
				"$gte": float64(3), // Must report issues on at least 3 different days
			},
		},
	}

	badgeReq := map[string]interface{}{
		"name":            badgeName,
		"description":     "Awarded for reporting issues on at least 3 different days",
		"image_url":       "https://example.com/consistent_reporter_badge.png",
		"flow_definition": flowDefinition,
		"is_active":       true,
	}

	resp = testutil.MakeRequest("POST", "/api/v1/admin/badges", badgeReq)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("Failed to create badge: %s", string(resp.Body))
	}

	// Parse badge response to extract badge ID
	var badgeWithCriteria struct {
		Badge struct {
			ID int `json:"id"`
		} `json:"badge"`
		Criteria struct {
			ID int `json:"id"`
		} `json:"criteria"`
	}

	err = json.Unmarshal(resp.Body, &badgeWithCriteria)
	if err != nil {
		t.Fatalf("Failed to parse badge response: %v", err)
	}

	badgeID := badgeWithCriteria.Badge.ID
	t.Logf("Created badge with ID: %d", badgeID)

	// STEP 3: Submit events for different days
	// Create clearly distinguishable dates with timezone information
	now := time.Now().UTC()
	t.Logf("Current time (UTC): %s", now.Format(time.RFC3339))

	// Use fixed dates rather than relative dates to ensure they're interpreted correctly
	dayTimeStamps := []time.Time{
		time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 2, 10, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 3, 10, 0, 0, 0, time.UTC),
		time.Date(2023, 1, 4, 10, 0, 0, 0, time.UTC),
	}

	// Submit events for each day (with detailed logging)
	for i, dayTimestamp := range dayTimeStamps {
		// Format the day in the expected format (matches getPeriodKey implementation)
		dayKey := dayTimestamp.Format("2006-01-02")

		// For the first day, submit 2 events to demonstrate that multiple events on the same day
		// only count as one period
		numEvents := 1
		if i == 0 {
			numEvents = 2
		}

		for j := 1; j <= numEvents; j++ {
			issueID := fmt.Sprintf("TIME-ISSUE-%d-%d", i, j)

			// Create the event with explicit timestamp
			formattedTime := dayTimestamp.Format(time.RFC3339)
			t.Logf("Submitting event with timestamp: %s (day key: %s)", formattedTime, dayKey)

			eventReq := map[string]interface{}{
				"event_type": eventTypeName,
				"user_id":    testUserID,
				"timestamp":  formattedTime,
				"payload": map[string]interface{}{
					"issue_id":    issueID,
					"description": fmt.Sprintf("Test issue for day %s", dayKey),
					"severity":    "medium",
				},
			}

			resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)

			// Log the full response for debugging
			t.Logf("Event creation response: %s", string(resp.Body))

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				t.Fatalf("Failed to create event for day %d: %s", i, string(resp.Body))
			}

			t.Logf("Successfully submitted event for day with key: %s", dayKey)

			// Add a small delay between events to ensure they are processed in order
			time.Sleep(300 * time.Millisecond)
		}
	}

	// Allow more time for badge processing to complete
	t.Log("Waiting for badge processing...")
	time.Sleep(4 * time.Second)

	// STEP 4: Verify badge was awarded
	// First check via API
	userBadgesEndpoint := fmt.Sprintf("/api/v1/users/%s/badges", testUserID)
	t.Logf("Checking user badges at endpoint: %s", userBadgesEndpoint)
	badgeResp := testutil.MakeRequest("GET", userBadgesEndpoint, nil)

	if badgeResp.StatusCode < 200 || badgeResp.StatusCode >= 300 {
		t.Fatalf("Failed to get user badges: %s", string(badgeResp.Body))
	}

	// Log the full response for debugging
	t.Logf("User badges response: %s", string(badgeResp.Body))

	// Parse badge response to find our badge
	var badges []map[string]interface{}
	err = json.Unmarshal(badgeResp.Body, &badges)
	if err != nil {
		t.Fatalf("Failed to parse badge check response: %v", err)
	}

	// Look for our badge in the list
	badgeFound := false
	var awardedAt string

	// Check the format of the badges response and adapt our lookup
	if len(badges) > 0 {
		// Print the keys of the first badge to debug
		for k := range badges[0] {
			t.Logf("Badge response has key: %s", k)
		}

		// Look for our badge ID in the list
		for _, badge := range badges {
			// Try different possible key names for the badge ID
			var badgeID int
			if id, ok := badge["badge_id"].(float64); ok {
				badgeID = int(id)
			} else if id, ok := badge["id"].(float64); ok {
				badgeID = int(id)
			}

			if badgeID == badgeWithCriteria.Badge.ID {
				badgeFound = true
				if at, ok := badge["awarded_at"].(string); ok {
					awardedAt = at
					t.Logf("Badge found in user badges! Awarded at: %s", awardedAt)
				} else {
					t.Logf("Badge found in user badges but no awarded_at timestamp")
				}

				// If there's metadata, log it too
				if metadata, ok := badge["metadata"].(string); ok {
					t.Logf("Badge metadata: %s", metadata)
				}

				break
			}
		}
	} else {
		t.Logf("No badges found in user badges response")
	}

	// Query for badge directly in the database as a fallback
	var dbBadgeCount int
	err = db.Get(&dbBadgeCount, "SELECT COUNT(*) FROM user_badges WHERE user_id = $1 AND badge_id = $2",
		testUserID, badgeWithCriteria.Badge.ID)

	if err != nil {
		t.Logf("Error querying database for badge: %v", err)
	} else {
		t.Logf("Database shows %d badges assigned to user", dbBadgeCount)
		if dbBadgeCount > 0 {
			badgeFound = true
		}
	}

	assert.True(t, badgeFound, "Expected to find the badge in the user's badges")

	// Also verify badge assignment directly in the database
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM user_badges WHERE user_id = $1 AND badge_id = $2", testUserID, badgeID)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}

	assert.Equal(t, 1, count, "Expected one badge to be assigned to the user")
	t.Logf("Verified badge assignment in database: badge_id=%d was assigned to user=%s", badgeID, testUserID)

	// Verify the total number of events and confirm unique day count
	var eventCount int
	err = db.Get(&eventCount,
		"SELECT COUNT(*) FROM events WHERE user_id = $1 AND event_type_id = $2",
		testUserID, eventTypeID)

	if err != nil {
		t.Fatalf("Failed to query events: %v", err)
	}

	t.Logf("Total events for this user and event type: %d", eventCount)
	assert.Equal(t, 5, eventCount, "Expected 5 events to be recorded (2 on first day + 1 each on other 3 days)")

	// Check the events in the database to verify their timestamps
	type EventInfo struct {
		ID         int       `db:"id"`
		UserID     string    `db:"user_id"`
		OccurredAt time.Time `db:"occurred_at"`
	}

	var events []EventInfo
	err = db.Select(&events,
		"SELECT id, user_id, occurred_at FROM events WHERE user_id = $1 AND event_type_id = $2 ORDER BY occurred_at",
		testUserID, eventTypeID)

	if err != nil {
		t.Logf("Error querying event timestamps: %v", err)
	} else {
		// Print each event's timestamp to verify correct storage
		t.Logf("Events in database:")
		uniqueDays := make(map[string]bool)

		for _, event := range events {
			dayKey := event.OccurredAt.Format("2006-01-02")
			uniqueDays[dayKey] = true
			t.Logf("  Event ID %d: occurred_at=%s, day_key=%s",
				event.ID, event.OccurredAt.Format(time.RFC3339), dayKey)
		}

		t.Logf("Number of unique days in database: %d", len(uniqueDays))
		assert.GreaterOrEqual(t, len(uniqueDays), 3, "Expected at least 3 unique days in the database")
	}
}
