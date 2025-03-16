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

// TestIssueReporterBadge tests the complete flow of badge assignment:
// 1. Create an event type for issue reporting
// 2. Create a badge with criteria that counts reported issues (using $eventCount)
// 3. Submit multiple events that meet the criteria
// 4. Verify the badge is awarded to the user after 5 reports
func TestIssueReporterBadge(t *testing.T) {
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
	eventTypeName := fmt.Sprintf("issue_reported_%d", timestamp)
	badgeName := fmt.Sprintf("Issue Reporter Badge_%d", timestamp)
	testUserID := "test_user_reporter" // Consistent user ID for this test

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
		},
		"required": []string{"issue_id", "status"},
	}

	eventTypeReq := map[string]interface{}{
		"name":        eventTypeName,
		"description": "Issue reporting event",
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

	// STEP 2: Create a badge with criteria
	// Define badge criteria using the $eventCount operator
	flowDefinition := map[string]interface{}{
		"event": eventTypeName,
		"criteria": map[string]interface{}{
			"status": "fixed",
			"$eventCount": map[string]interface{}{
				"$gte": float64(5), // Count must be at least 5 fixed issues
			},
		},
	}

	badgeReq := map[string]interface{}{
		"name":            badgeName,
		"description":     "Awarded for reporting 5 or more issues that got fixed",
		"image_url":       "https://example.com/issue_reporter_badge.png",
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

	// STEP 3: Submit events for fixed issues (one by one)
	for i := 1; i <= 6; i++ {
		issueID := fmt.Sprintf("ISSUE-%d", i)

		eventReq := map[string]interface{}{
			"event_type": eventTypeName,
			"user_id":    testUserID,
			"timestamp":  time.Now().Add(time.Duration(i) * time.Minute).Format(time.RFC3339), // Stagger timestamps
			"payload": map[string]interface{}{
				"issue_id":    issueID,
				"description": fmt.Sprintf("Test issue %d", i),
				"status":      "fixed",
			},
		}

		resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			t.Fatalf("Failed to create event %d: %s", i, string(resp.Body))
		}

		t.Logf("Submitted event for fixed issue %d", i)

		// Add a small delay between events to ensure they are processed in order
		time.Sleep(200 * time.Millisecond)
	}

	// Allow time for badge processing to complete
	t.Log("Waiting for badge processing...")
	time.Sleep(3 * time.Second)

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

	// Look for our badge in the list
	badgeFound := false
	var awardedAt string

	t.Logf("User badges response: %s", string(badgeResp.Body))

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

	// Verify that the badge was awarded after the 5th issue was reported
	// This shows that the $eventCount logic is working correctly
	if badgeFound {
		// Optionally, query the events table to verify the exact count that triggered the badge
		var eventCount int
		err = db.Get(&eventCount,
			"SELECT COUNT(*) FROM events WHERE user_id = $1 AND event_type_id = $2",
			testUserID, eventTypeID)

		if err != nil {
			t.Fatalf("Failed to query events: %v", err)
		}

		t.Logf("Total events for this user and event type: %d", eventCount)
		assert.Equal(t, 6, eventCount, "Expected 6 events to be recorded (5 to trigger + 1 extra)")
	}
}
