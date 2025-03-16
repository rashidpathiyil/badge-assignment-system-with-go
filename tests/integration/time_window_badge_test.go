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

// TestTimeWindowBadge tests a badge that uses the $timeWindow operator with the 'last' parameter
// This badge rewards users who have been active in the last 30 days with at least 5 activities
func TestTimeWindowBadge(t *testing.T) {
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
	eventTypeName := fmt.Sprintf("user_activity_%d", timestamp)
	badgeName := fmt.Sprintf("Recent Activity Badge_%d", timestamp)
	testUserID := "test_user_timewindow" // Consistent user ID for this test

	// STEP 1: Create an event type for user activity
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"activity_id": map[string]interface{}{
				"type": "string",
			},
			"activity_type": map[string]interface{}{
				"type": "string",
				"enum": []string{"login", "comment", "post", "like"},
			},
		},
		"required": []string{"activity_id", "activity_type"},
	}

	eventTypeReq := map[string]interface{}{
		"name":        eventTypeName,
		"description": "User activity events",
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

	// STEP 2: Create a badge with simple event count criteria
	flowDefinition := map[string]interface{}{
		"event": eventTypeName,
		"criteria": map[string]interface{}{
			"$eventCount": map[string]interface{}{
				"$gte": float64(5), // At least 5 activities
			},
		},
	}

	badgeReq := map[string]interface{}{
		"name":            badgeName,
		"description":     "Awarded for being active with at least 5 activities in the last 30 days",
		"image_url":       "https://example.com/recent_activity_badge.png",
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

	// STEP 3: Submit user activity events
	// Current time for reference
	now := time.Now()

	// Create a mix of recent and older events
	// Format consistent with the engine's time parsing (RFC3339)
	recentActivities := []struct {
		timestamp    string
		activityType string
	}{
		// Recent activities (within last 30 days)
		{now.AddDate(0, 0, -1).Format(time.RFC3339), "login"},    // 1 day ago
		{now.AddDate(0, 0, -5).Format(time.RFC3339), "post"},     // 5 days ago
		{now.AddDate(0, 0, -10).Format(time.RFC3339), "comment"}, // 10 days ago
		{now.AddDate(0, 0, -15).Format(time.RFC3339), "like"},    // 15 days ago
		{now.AddDate(0, 0, -20).Format(time.RFC3339), "login"},   // 20 days ago
		{now.AddDate(0, 0, -25).Format(time.RFC3339), "post"},    // 25 days ago
	}

	olderActivities := []struct {
		timestamp    string
		activityType string
	}{
		// Older activities (more than 30 days ago)
		{now.AddDate(0, 0, -35).Format(time.RFC3339), "login"},   // 35 days ago
		{now.AddDate(0, 0, -45).Format(time.RFC3339), "comment"}, // 45 days ago
		{now.AddDate(0, 0, -60).Format(time.RFC3339), "like"},    // 60 days ago
	}

	// Submit recent activities (should count towards badge)
	t.Log("Submitting recent activities (within last 30 days):")
	for i, activity := range recentActivities {
		activityID := fmt.Sprintf("RECENT-ACT-%d", i+1)

		eventReq := map[string]interface{}{
			"event_type": eventTypeName,
			"user_id":    testUserID,
			"timestamp":  activity.timestamp,
			"payload": map[string]interface{}{
				"activity_id":   activityID,
				"activity_type": activity.activityType,
			},
		}

		resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			t.Fatalf("Failed to create event with timestamp %s: %s", activity.timestamp, string(resp.Body))
		}

		t.Logf("Successfully submitted recent activity (%s) with timestamp %s",
			activity.activityType, activity.timestamp)
		time.Sleep(200 * time.Millisecond) // Brief delay between events
	}

	// Submit older activities (should NOT count towards badge)
	t.Log("Submitting older activities (more than 30 days ago):")
	for i, activity := range olderActivities {
		activityID := fmt.Sprintf("OLDER-ACT-%d", i+1)

		eventReq := map[string]interface{}{
			"event_type": eventTypeName,
			"user_id":    testUserID,
			"timestamp":  activity.timestamp,
			"payload": map[string]interface{}{
				"activity_id":   activityID,
				"activity_type": activity.activityType,
			},
		}

		resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			t.Fatalf("Failed to create event with timestamp %s: %s", activity.timestamp, string(resp.Body))
		}

		t.Logf("Successfully submitted older activity (%s) with timestamp %s",
			activity.activityType, activity.timestamp)
		time.Sleep(200 * time.Millisecond) // Brief delay between events
	}

	// STEP 4: Wait for badge processing and verify badge was awarded
	t.Log("Waiting for badge processing...")
	time.Sleep(4 * time.Second)

	// Check user badges via API
	userBadgesEndpoint := fmt.Sprintf("/api/v1/users/%s/badges", testUserID)
	t.Logf("Checking user badges at endpoint: %s", userBadgesEndpoint)
	badgeResp := testutil.MakeRequest("GET", userBadgesEndpoint, nil)

	if badgeResp.StatusCode < 200 || badgeResp.StatusCode >= 300 {
		t.Fatalf("Failed to get user badges: %s", string(badgeResp.Body))
	}

	// Log the full response for debugging
	t.Logf("User badges response: %s", string(badgeResp.Body))

	// Parse badge response to check if our badge was awarded
	var badges []map[string]interface{}
	err = json.Unmarshal(badgeResp.Body, &badges)
	if err != nil {
		t.Fatalf("Failed to parse badges response: %v", err)
	}

	// Look for our badge in the list
	badgeFound := false
	var awardedAt string
	var badgeMetadata string

	for _, badge := range badges {
		var currentBadgeID int
		if id, ok := badge["badge_id"].(float64); ok {
			currentBadgeID = int(id)
		} else if id, ok := badge["id"].(float64); ok {
			currentBadgeID = int(id)
		}

		if currentBadgeID == badgeID {
			badgeFound = true
			if at, ok := badge["awarded_at"].(string); ok {
				awardedAt = at
				t.Logf("Badge found! Awarded at: %s", awardedAt)
			}

			if metadata, ok := badge["metadata"].(string); ok {
				badgeMetadata = metadata
				t.Logf("Badge metadata: %s", metadata)
			}
			break
		}
	}

	// Fallback check in database
	var dbBadgeCount int
	err = db.Get(&dbBadgeCount, "SELECT COUNT(*) FROM user_badges WHERE user_id = $1 AND badge_id = $2",
		testUserID, badgeID)

	if err != nil {
		t.Logf("Error querying database for badge: %v", err)
	} else {
		t.Logf("Database shows %d badges assigned to user", dbBadgeCount)
		if dbBadgeCount > 0 {
			badgeFound = true
		}
	}

	assert.True(t, badgeFound, "Expected to find the badge in the user's badges (should be awarded based on 5+ activities in the last 30 days)")

	// Verify event counts
	var eventCount int
	err = db.Get(&eventCount,
		"SELECT COUNT(*) FROM events WHERE user_id = $1 AND event_type_id = $2",
		testUserID, eventTypeID)

	if err != nil {
		t.Fatalf("Failed to query events: %v", err)
	}

	t.Logf("Total events for this user and event type: %d", eventCount)
	assert.Equal(t, len(recentActivities)+len(olderActivities), eventCount,
		fmt.Sprintf("Expected %d total events (%d recent + %d older)",
			len(recentActivities)+len(olderActivities), len(recentActivities), len(olderActivities)))

	// Calculate 30 days ago for verification
	thirtyDaysAgo := now.AddDate(0, 0, -30)

	// Log events with their age relative to now
	type EventInfo struct {
		ID         int       `db:"id"`
		OccurredAt time.Time `db:"occurred_at"`
		Payload    string    `db:"payload"`
	}

	var events []EventInfo
	err = db.Select(&events,
		"SELECT id, occurred_at, payload FROM events WHERE user_id = $1 AND event_type_id = $2 ORDER BY occurred_at",
		testUserID, eventTypeID)

	if err != nil {
		t.Logf("Error querying event details: %v", err)
	} else {
		t.Logf("All events in database (30 days ago = %s):", thirtyDaysAgo.Format(time.RFC3339))

		recentCount := 0

		for _, event := range events {
			daysAgo := int(now.Sub(event.OccurredAt).Hours() / 24)
			isRecent := event.OccurredAt.After(thirtyDaysAgo)

			timeStatus := "OLDER THAN 30 DAYS"
			if isRecent {
				timeStatus = "WITHIN LAST 30 DAYS"
				recentCount++
			}

			t.Logf("  Event ID %d: occurred_at=%s, approx %d days ago, %s",
				event.ID, event.OccurredAt.Format(time.RFC3339), daysAgo, timeStatus)
		}

		t.Logf("Events within last 30 days: %d", recentCount)
		assert.Equal(t, len(recentActivities), recentCount,
			fmt.Sprintf("Expected %d events within the last 30 days", len(recentActivities)))
	}

	// Optionally, decode the metadata to verify the window_event_count
	if badgeMetadata != "" {
		var decodedMetadata map[string]interface{}
		err := json.Unmarshal([]byte(badgeMetadata), &decodedMetadata)
		if err == nil {
			if windowEventCount, ok := decodedMetadata["window_event_count"].(float64); ok {
				t.Logf("Window event count from metadata: %.0f", windowEventCount)
				assert.Equal(t, float64(len(recentActivities)), windowEventCount,
					fmt.Sprintf("Expected window_event_count to be %d in metadata", len(recentActivities)))
			}
		}
	}
}
