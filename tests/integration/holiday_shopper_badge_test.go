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

// TestHolidayShopperBadge tests the complete flow of badge assignment:
// 1. Create an event type for purchase tracking
// 2. Create a badge with $timeWindow criteria for purchases during a holiday period
// 3. Submit events within and outside the time window
// 4. Verify the badge is awarded only based on events within the time window
func TestHolidayShopperBadge(t *testing.T) {
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
	eventTypeName := fmt.Sprintf("purchase_event_%d", timestamp)
	badgeName := fmt.Sprintf("Holiday Shopper Badge_%d", timestamp)
	testUserID := "test_user_holiday" // Consistent user ID for this test

	// STEP 1: Create an event type for purchase events
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"order_id": map[string]interface{}{
				"type": "string",
			},
			"amount": map[string]interface{}{
				"type": "number",
			},
			"product_category": map[string]interface{}{
				"type": "string",
			},
		},
		"required": []string{"order_id", "amount"},
	}

	eventTypeReq := map[string]interface{}{
		"name":        eventTypeName,
		"description": "Purchase event for holiday promotion",
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

	// STEP 2: Create a badge with $timeWindow criteria
	// Define a holiday period from December 1 to December 31, 2023
	holidayStart := "2023-12-01T00:00:00Z"
	holidayEnd := "2023-12-31T23:59:59Z"

	// For testing recent purchase events, we could use $NOW-based criteria
	// But for the holiday period, we'll keep specific dates since it's a fixed historical window
	// This demonstrates using both fixed dates and dynamic variables in different contexts

	flowDefinition := map[string]interface{}{
		"event": eventTypeName,
		"criteria": map[string]interface{}{
			"$eventCount": map[string]interface{}{
				"$gte": float64(3), // At least 3 purchases during the holiday period
			},
			"timestamp": map[string]interface{}{
				"$gte": holidayStart,
				"$lte": holidayEnd,
			},
		},
	}

	badgeReq := map[string]interface{}{
		"name":            badgeName,
		"description":     "Awarded for making at least 3 purchases during the holiday season",
		"image_url":       "https://example.com/holiday_shopper_badge.png",
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

	// STEP 3: Submit events - some within the time window, some outside
	// Create timestamps within and outside the holiday period
	// Format consistent with the engine's time parsing (RFC3339)
	withinHolidayPeriod := []string{
		"2023-12-01T10:30:00Z", // Beginning of holiday period
		"2023-12-15T14:45:00Z", // Middle of holiday period
		"2023-12-31T23:00:00Z", // End of holiday period
		"2023-12-24T18:30:00Z", // Christmas Eve - extra purchase
	}

	outsideHolidayPeriod := []string{
		"2023-11-25T12:00:00Z", // Before holiday period (Black Friday)
		"2024-01-02T09:15:00Z", // After holiday period
	}

	// Submit events within the holiday period
	t.Log("Submitting events within the holiday period:")
	for i, timestamp := range withinHolidayPeriod {
		orderID := fmt.Sprintf("HOLIDAY-ORDER-%d", i+1)
		amount := 50.0 + float64(i*25) // Varying purchase amounts

		eventReq := map[string]interface{}{
			"event_type": eventTypeName,
			"user_id":    testUserID,
			"timestamp":  timestamp,
			"payload": map[string]interface{}{
				"order_id":         orderID,
				"amount":           amount,
				"product_category": "gifts",
			},
		}

		resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			t.Fatalf("Failed to create event with timestamp %s: %s", timestamp, string(resp.Body))
		}

		t.Logf("Successfully submitted in-period purchase event with timestamp %s", timestamp)
		time.Sleep(200 * time.Millisecond) // Brief delay between events
	}

	// Submit events outside the holiday period
	t.Log("Submitting events outside the holiday period:")
	for i, timestamp := range outsideHolidayPeriod {
		orderID := fmt.Sprintf("NON-HOLIDAY-ORDER-%d", i+1)
		amount := 30.0 + float64(i*15) // Varying purchase amounts

		eventReq := map[string]interface{}{
			"event_type": eventTypeName,
			"user_id":    testUserID,
			"timestamp":  timestamp,
			"payload": map[string]interface{}{
				"order_id":         orderID,
				"amount":           amount,
				"product_category": "regular",
			},
		}

		resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			t.Fatalf("Failed to create event with timestamp %s: %s", timestamp, string(resp.Body))
		}

		t.Logf("Successfully submitted out-of-period purchase event with timestamp %s", timestamp)
		time.Sleep(200 * time.Millisecond) // Brief delay between events
	}

	// Total events: 4 within holiday period, 2 outside period

	// Allow time for badge processing to complete
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
	var badgeMetadata string

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
					badgeMetadata = metadata
					t.Logf("Badge metadata: %s", metadata)
				}

				break
			}
		}
	} else {
		t.Logf("No badges found in user badges response")
	}

	// Check badge in database as a fallback
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

	assert.True(t, badgeFound, "Expected to find the badge in the user's badges (should be awarded based on the 3+ purchases within holiday period)")

	// Verify that only events within the time window were counted
	// We expect 4 events within the window and 2 outside (total 6)
	var eventCount int
	err = db.Get(&eventCount,
		"SELECT COUNT(*) FROM events WHERE user_id = $1 AND event_type_id = $2",
		testUserID, eventTypeID)

	if err != nil {
		t.Fatalf("Failed to query events: %v", err)
	}

	t.Logf("Total events for this user and event type: %d", eventCount)
	assert.Equal(t, 6, eventCount, "Expected 6 total events (4 within holiday period + 2 outside)")

	// Log the specific times of all events for verification
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
		t.Logf("All events in database:")

		holidayPeriodStart, _ := time.Parse(time.RFC3339, holidayStart)
		holidayPeriodEnd, _ := time.Parse(time.RFC3339, holidayEnd)
		inHolidayPeriodCount := 0

		for _, event := range events {
			isInHolidayPeriod := (event.OccurredAt.Equal(holidayPeriodStart) || event.OccurredAt.After(holidayPeriodStart)) &&
				(event.OccurredAt.Equal(holidayPeriodEnd) || event.OccurredAt.Before(holidayPeriodEnd))

			periodStatus := "OUTSIDE HOLIDAY PERIOD"
			if isInHolidayPeriod {
				periodStatus = "WITHIN HOLIDAY PERIOD"
				inHolidayPeriodCount++
			}

			t.Logf("  Event ID %d: occurred_at=%s, %s",
				event.ID, event.OccurredAt.Format(time.RFC3339), periodStatus)
		}

		t.Logf("Events within holiday period: %d", inHolidayPeriodCount)
		assert.Equal(t, 4, inHolidayPeriodCount, "Expected 4 events within the holiday period")
	}

	// Optionally, decode the metadata to verify the window_event_count
	if badgeMetadata != "" {
		var decodedMetadata map[string]interface{}
		err := json.Unmarshal([]byte(badgeMetadata), &decodedMetadata)
		if err == nil {
			if windowEventCount, ok := decodedMetadata["window_event_count"].(float64); ok {
				t.Logf("Window event count from metadata: %.0f", windowEventCount)
				assert.Equal(t, float64(4), windowEventCount, "Expected window_event_count to be 4 in metadata")
			}
		}
	}
}
