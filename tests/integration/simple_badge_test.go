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

// TestSimpleBadgeAssignment tests the complete flow of badge assignment:
// 1. Create an event type
// 2. Create a badge with criteria based on the event type
// 3. Submit an event that meets the criteria
// 4. Verify the badge is awarded to the user
func TestSimpleBadgeAssignment(t *testing.T) {
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
	eventTypeName := fmt.Sprintf("simple_event_%d", timestamp)
	badgeName := fmt.Sprintf("Simple Badge_%d", timestamp)
	testUserID := "test_user_1" // Consistent user ID for this test

	// STEP 1: Create an event type
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"value": map[string]interface{}{
				"type": "number",
			},
		},
	}

	eventTypeReq := map[string]interface{}{
		"name":        eventTypeName,
		"description": "Simple test event",
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
	// Define badge criteria using the $gte operator
	flowDefinition := map[string]interface{}{
		"event": eventTypeName,
		"criteria": map[string]interface{}{
			"value": map[string]interface{}{
				"$gte": float64(50), // Value must be greater than or equal to 50
			},
		},
	}

	badgeReq := map[string]interface{}{
		"name":            badgeName,
		"description":     "Simple test badge",
		"image_url":       "https://example.com/badge.png",
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

	// STEP 3: Submit an event that meets the badge criteria
	eventReq := map[string]interface{}{
		"event_type": eventTypeName,
		"user_id":    testUserID,
		"timestamp":  time.Now().Format(time.RFC3339),
		"payload": map[string]interface{}{
			"value": float64(75), // Value exceeds the threshold in criteria
		},
	}

	resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("Failed to create event: %s", string(resp.Body))
	}

	// Allow time for badge processing to complete
	t.Log("Waiting for badge processing...")
	time.Sleep(2 * time.Second)

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
	for _, badge := range badges {
		if int(badge["id"].(float64)) == badgeID {
			badgeFound = true
			t.Logf("Badge found in user badges! Awarded at: %s", badge["awarded_at"].(string))
			break
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
}
