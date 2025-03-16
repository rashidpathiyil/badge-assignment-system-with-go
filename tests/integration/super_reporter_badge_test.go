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

// TestSuperReporterBadge tests a badge that combines both count operators:
// 1. Create an event type for issue reporting
// 2. Create a badge with criteria that uses both $eventCount AND $timePeriod
// 3. Submit events that will satisfy both criteria
// 4. Verify the badge is awarded
func TestSuperReporterBadge(t *testing.T) {
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
	eventTypeName := fmt.Sprintf("issue_reported_super_%d", timestamp)
	badgeName := fmt.Sprintf("Super Reporter Badge_%d", timestamp)
	testUserID := "test_user_super" // Consistent user ID for this test

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
			"status": map[string]interface{}{
				"type": "string",
				"enum": []string{"reported", "in_progress", "fixed"},
			},
			"severity": map[string]interface{}{
				"type": "string",
				"enum": []string{"low", "medium", "high"},
			},
		},
		"required": []string{"issue_id", "status", "severity"},
	}

	eventTypeReq := map[string]interface{}{
		"name":        eventTypeName,
		"description": "Issue reporting event for combined test",
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

	// STEP 2: Create a badge with criteria using both $eventCount and $timePeriod
	// Define badge criteria to require:
	// 1. At least 10 total fixed issues
	// 2. Activity on at least 3 different days
	flowDefinition := map[string]interface{}{
		"$and": []interface{}{
			map[string]interface{}{
				"event": eventTypeName,
				"criteria": map[string]interface{}{
					"status": "fixed",
					"$eventCount": map[string]interface{}{
						"$gte": float64(10), // Count must be at least 10 fixed issues
					},
				},
			},
			map[string]interface{}{
				"$timePeriod": map[string]interface{}{
					"periodType": "day",
					"periodCount": map[string]interface{}{
						"$gte": float64(3), // Must report issues on at least 3 different days
					},
				},
			},
		},
	}

	badgeReq := map[string]interface{}{
		"name":            badgeName,
		"description":     "Awarded for reporting 10+ issues that got fixed across at least 3 different days",
		"image_url":       "https://example.com/super_reporter_badge.png",
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

	// STEP 3: Submit events across multiple days
	// Create base time (today)
	now := time.Now()

	// Create timestamps for different days
	dayTimeStamps := []time.Time{
		now.AddDate(0, 0, -2), // 2 days ago
		now.AddDate(0, 0, -1), // 1 day ago
		now,                   // today
	}

	// Submit multiple events for each day (different numbers to show both criteria)
	eventsPerDay := []int{2, 3, 6} // Total of 11 events across 3 days

	eventCount := 0
	for i, dayTimestamp := range dayTimeStamps {
		numEvents := eventsPerDay[i]

		for j := 1; j <= numEvents; j++ {
			eventCount++
			issueID := fmt.Sprintf("SUPER-ISSUE-%d", eventCount)

			eventReq := map[string]interface{}{
				"event_type": eventTypeName,
				"user_id":    testUserID,
				"timestamp":  dayTimestamp.Format(time.RFC3339),
				"payload": map[string]interface{}{
					"issue_id":    issueID,
					"description": fmt.Sprintf("Test issue %d on day %d", j, i),
					"status":      "fixed",
					"severity":    "medium",
				},
			}

			resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				t.Fatalf("Failed to create event %d on day %d: %s", j, i, string(resp.Body))
			}

			t.Logf("Submitted event %d for day %d (timestamp: %s)", j, i, dayTimestamp.Format("2006-01-02"))

			// Add a small delay between events to ensure they are processed in order
			time.Sleep(200 * time.Millisecond)
		}
	}

	// Allow time for badge processing to complete
	t.Log("Waiting for badge processing...")
	time.Sleep(5 * time.Second) // Extra time for more complex criteria

	// STEP 4: Verify badge was awarded
	// First check via API
	userBadgesEndpoint := fmt.Sprintf("/api/v1/users/%s/badges", testUserID)
	badgeResp := testutil.MakeRequest("GET", userBadgesEndpoint, nil)

	if badgeResp.StatusCode < 200 || badgeResp.StatusCode >= 300 {
		t.Fatalf("Failed to get user badges: %s", string(badgeResp.Body))
	}

	// Parse badge response to find our badge
	var badges []map[string]interface{}
	err = json.Unmarshal(badgeResp.Body, &badges)
	if err != nil {
		t.Fatalf("Failed to parse badge check response: %v", err)
	}

	// Log full badge response for debugging
	t.Logf("User badges response: %s", string(badgeResp.Body))

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

	assert.True(t, badgeFound, "Expected to find the badge in the user's badges")

	// Also verify badge assignment directly in the database
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM user_badges WHERE user_id = $1 AND badge_id = $2", testUserID, badgeID)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}

	assert.Equal(t, 1, count, "Expected one badge to be assigned to the user")
	t.Logf("Verified badge assignment in database: badge_id=%d was assigned to user=%s", badgeID, testUserID)

	// Query metadata from database for additional debug info
	var metadata string
	err = db.Get(&metadata, "SELECT metadata FROM user_badges WHERE user_id = $1 AND badge_id = $2", testUserID, badgeID)
	if err == nil && metadata != "" {
		t.Logf("Badge metadata from database: %s", metadata)
	}

	// Verify the total number of events
	var eventTotal int
	err = db.Get(&eventTotal,
		"SELECT COUNT(*) FROM events WHERE user_id = $1 AND event_type_id = $2",
		testUserID, eventTypeID)

	if err != nil {
		t.Fatalf("Failed to query events: %v", err)
	}

	t.Logf("Total events for this user and event type: %d", eventTotal)
	assert.Equal(t, 11, eventTotal, "Expected 11 events to be recorded across 3 days")
}
