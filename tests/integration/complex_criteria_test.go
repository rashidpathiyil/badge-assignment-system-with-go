package integration

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/testutil"
	"github.com/stretchr/testify/assert"
)

// Global variables for complex criteria tests
var (
	logicalOpBadgeID         int
	timeBasedCriteriaBadgeID int
)

// TestLogicalOperators tests creating badges with logical operators and verifying they work
func TestLogicalOperators(t *testing.T) {
	SetupTest()

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	eventTypeName := fmt.Sprintf("test_event_%d", timestamp)
	badgeName := fmt.Sprintf("Logical Op Badge_%d", timestamp)
	testUserID := fmt.Sprintf("test_user_logical_%d", timestamp)

	// STEP 1: Create an event type
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"score": map[string]interface{}{
				"type": "number",
			},
			"duration": map[string]interface{}{
				"type": "number",
			},
			"completed": map[string]interface{}{
				"type": "boolean",
			},
		},
	}

	eventTypeReq := map[string]interface{}{
		"name":        eventTypeName,
		"description": "Test event type for logical operators",
		"schema":      schema,
	}

	resp := testutil.MakeRequest("POST", "/api/v1/admin/event-types", eventTypeReq)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("Failed to create event type: %s", string(resp.Body))
	}

	var eventTypeCreateResp map[string]interface{}
	err := testutil.ParseResponse(resp, &eventTypeCreateResp)
	assert.NoError(t, err)
	eventTypeID := int(eventTypeCreateResp["id"].(float64))
	t.Logf("Created event type with ID: %d", eventTypeID)

	// STEP 2: Define criteria with logical operators ($and and $or)
	criteria := map[string]interface{}{
		"$and": []interface{}{
			map[string]interface{}{
				"event": eventTypeName,
				"criteria": map[string]interface{}{
					"score": map[string]interface{}{
						"$gte": float64(80),
					},
				},
			},
			map[string]interface{}{
				"$or": []interface{}{
					map[string]interface{}{
						"event": eventTypeName,
						"criteria": map[string]interface{}{
							"duration": map[string]interface{}{
								"$lte": float64(300),
							},
						},
					},
					map[string]interface{}{
						"event": eventTypeName,
						"criteria": map[string]interface{}{
							"completed": map[string]interface{}{
								"$eq": true,
							},
						},
					},
				},
			},
		},
	}

	// Create badge request
	badgeReq := BadgeRequest{
		Name:           badgeName,
		Description:    "Badge with logical operators criteria",
		ImageURL:       "https://example.com/badges/logical-op.png",
		FlowDefinition: criteria,
		IsActive:       true,
	}

	// Make the API request
	resp = testutil.MakeRequest("POST", "/api/v1/admin/badges", badgeReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badgeResp struct {
		Badge struct {
			ID          int       `json:"id"`
			Name        string    `json:"name"`
			Description string    `json:"description"`
			ImageURL    string    `json:"image_url"`
			Active      bool      `json:"active"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		} `json:"badge"`
	}
	err = testutil.ParseResponse(resp, &badgeResp)
	assert.NoError(t, err)

	// Store the badge ID
	logicalOpBadgeID = badgeResp.Badge.ID

	// Get the event type name
	resp = testutil.MakeRequest("GET", fmt.Sprintf("/api/v1/admin/event-types/%d", eventTypeID), nil)
	testutil.AssertSuccess(t, resp)
	var eventTypeGetResp map[string]interface{}
	err = testutil.ParseResponse(resp, &eventTypeGetResp)
	assert.NoError(t, err)
	// Use the event type name we created earlier instead of trying to get it from the response
	// This is more reliable and avoids potential issues with the API

	// Create event that meets the logical criteria
	eventPayload := map[string]interface{}{
		"score":     float64(85),  // >=80 ✓
		"duration":  float64(250), // <=300 ✓
		"completed": true,         // Already true, but this fulfills the $or condition as well
	}

	// Create event request
	eventReq := EventRequest{
		EventType: eventTypeName,
		UserID:    testUserID,
		Payload:   eventPayload,
		Timestamp: time.Now(),
	}

	// Make the API request to process the event
	resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Wait for badge processing
	t.Log("Waiting for badge processing...")
	time.Sleep(4 * time.Second)

	// Check if badge was awarded
	endpoint := fmt.Sprintf("/api/v1/users/%s/badges", testUserID)
	resp = testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badges []map[string]interface{}
	err = json.Unmarshal(resp.Body, &badges)
	if err != nil {
		t.Fatalf("Failed to parse badge check response: %v", err)
	}

	// Verify badge was awarded
	badgeFound := false
	for _, badge := range badges {
		if badgeIDFloat, ok := badge["id"].(float64); ok && int(badgeIDFloat) == logicalOpBadgeID {
			badgeFound = true
			break
		}
	}

	assert.True(t, badgeFound, "The logical operator badge should have been awarded to the user")
}

// TestTimeBasedCriteria tests creating badges with time-based criteria and verifying they work
func TestTimeBasedCriteria(t *testing.T) {
	SetupTest()

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	eventTypeName := fmt.Sprintf("test_event_%d", timestamp)
	badgeName := fmt.Sprintf("Time-Based Badge_%d", timestamp)
	testUserID := fmt.Sprintf("test_user_timewindow_%d", timestamp)

	// STEP 1: Create an event type
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"score": map[string]interface{}{
				"type": "number",
			},
			"completed": map[string]interface{}{
				"type": "boolean",
			},
		},
	}

	eventTypeReq := map[string]interface{}{
		"name":        eventTypeName,
		"description": "Test event type for time-based criteria",
		"schema":      schema,
	}

	resp := testutil.MakeRequest("POST", "/api/v1/admin/event-types", eventTypeReq)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("Failed to create event type: %s", string(resp.Body))
	}

	var eventTypeCreateResp map[string]interface{}
	err := testutil.ParseResponse(resp, &eventTypeCreateResp)
	assert.NoError(t, err)
	eventTypeID := int(eventTypeCreateResp["id"].(float64))
	t.Logf("Created event type with ID: %d", eventTypeID)

	// STEP 2: Define time-based criteria using timestamp range instead of $timeWindow
	criteria := map[string]interface{}{
		"event": eventTypeName,
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": float64(75),
			},
			"timestamp": map[string]interface{}{
				"$gte": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				"$lte": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
		},
	}

	// Create badge request
	badgeReq := BadgeRequest{
		Name:           badgeName,
		Description:    "Badge with time-based criteria",
		ImageURL:       "https://example.com/badges/time-based.png",
		FlowDefinition: criteria,
		IsActive:       true,
	}

	// Make the API request
	resp = testutil.MakeRequest("POST", "/api/v1/admin/badges", badgeReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badgeResp struct {
		Badge struct {
			ID          int       `json:"id"`
			Name        string    `json:"name"`
			Description string    `json:"description"`
			ImageURL    string    `json:"image_url"`
			Active      bool      `json:"active"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		} `json:"badge"`
	}
	err = testutil.ParseResponse(resp, &badgeResp)
	assert.NoError(t, err)

	// Store the badge ID
	timeBasedCriteriaBadgeID = badgeResp.Badge.ID

	// Create event that meets the time-based criteria
	eventPayload := map[string]interface{}{
		"score":     float64(80), // >=75 ✓
		"completed": true,
	}

	// Create event request with timestamp inside the time window
	eventReq := EventRequest{
		EventType: eventTypeName,
		UserID:    testUserID,
		Payload:   eventPayload,
		Timestamp: time.Now(), // Current time is within the window
	}

	// Make the API request to process the event
	resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Wait for badge processing
	t.Log("Waiting for badge processing...")
	time.Sleep(4 * time.Second)

	// Check if badge was awarded
	endpoint := fmt.Sprintf("/api/v1/users/%s/badges", testUserID)
	resp = testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badges []map[string]interface{}
	err = json.Unmarshal(resp.Body, &badges)
	if err != nil {
		t.Fatalf("Failed to parse badge check response: %v", err)
	}

	// Verify badge was awarded
	badgeFound := false
	for _, badge := range badges {
		if badgeIDFloat, ok := badge["id"].(float64); ok && int(badgeIDFloat) == timeBasedCriteriaBadgeID {
			badgeFound = true
			break
		}
	}

	assert.True(t, badgeFound, "The time-based criteria badge should have been awarded to the user")
}
