package integration

import (
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

	if EventTypeID == 0 {
		t.Skip("Event type ID not set, skipping test")
	}

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	badgeName := fmt.Sprintf("Logical Op Badge_%d", timestamp)

	// Define criteria with logical operators ($and and $or)
	criteria := map[string]interface{}{
		"$and": []interface{}{
			map[string]interface{}{
				"event": "test_event",
				"criteria": map[string]interface{}{
					"score": map[string]interface{}{
						"$gte": float64(80),
					},
				},
			},
			map[string]interface{}{
				"$or": []interface{}{
					map[string]interface{}{
						"event": "test_event",
						"criteria": map[string]interface{}{
							"duration": map[string]interface{}{
								"$lte": float64(300),
							},
						},
					},
					map[string]interface{}{
						"event": "test_event",
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
	resp := testutil.MakeRequest("POST", "/api/v1/admin/badges", badgeReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badgeResp BadgeResponse
	err := testutil.ParseResponse(resp, &badgeResp)
	assert.NoError(t, err)

	// Store the badge ID
	logicalOpBadgeID = badgeResp.ID

	// Get the event type name
	var eventTypeName string
	resp = testutil.MakeRequest("GET", fmt.Sprintf("/api/v1/admin/event-types/%d", EventTypeID), nil)
	testutil.AssertSuccess(t, resp)
	var eventTypeResp EventTypeResponse
	err = testutil.ParseResponse(resp, &eventTypeResp)
	assert.NoError(t, err)
	eventTypeName = eventTypeResp.Name

	// Create event that meets the logical criteria
	eventPayload := map[string]interface{}{
		"score":     float64(85),  // >=80 ✓
		"duration":  float64(250), // <=300 ✓
		"completed": true,         // Already true, but this fulfills the $or condition as well
	}

	// Create event request
	eventReq := EventRequest{
		EventType: eventTypeName,
		UserID:    fmt.Sprintf("%s_logical", TestUserID),
		Payload:   eventPayload,
		Timestamp: time.Now(),
	}

	// Make the API request to process the event
	resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Wait for badge processing
	time.Sleep(1 * time.Second)

	// Check if badge was awarded
	endpoint := fmt.Sprintf("/api/v1/users/%s_logical/badges", TestUserID)
	resp = testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var userBadgesResp UserBadgesResponse
	err = testutil.ParseResponse(resp, &userBadgesResp)
	assert.NoError(t, err)

	// Verify badge was awarded
	badgeFound := false
	for _, badge := range userBadgesResp.Badges {
		if badge.BadgeID == logicalOpBadgeID {
			badgeFound = true
			break
		}
	}

	assert.True(t, badgeFound, "The logical operator badge should have been awarded to the user")
}

// TestTimeBasedCriteria tests creating badges with time-based criteria and verifying they work
func TestTimeBasedCriteria(t *testing.T) {
	SetupTest()

	if EventTypeID == 0 {
		t.Skip("Event type ID not set, skipping test")
	}

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	badgeName := fmt.Sprintf("Time-Based Badge_%d", timestamp)

	// Define time window criteria (events must occur within a specific time frame)
	startDate := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	endDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	criteria := map[string]interface{}{
		"$timeWindow": map[string]interface{}{
			"start": startDate,
			"end":   endDate,
			"flow": map[string]interface{}{
				"event": "test_event",
				"criteria": map[string]interface{}{
					"score": map[string]interface{}{
						"$gte": float64(75),
					},
				},
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
	resp := testutil.MakeRequest("POST", "/api/v1/admin/badges", badgeReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badgeResp BadgeResponse
	err := testutil.ParseResponse(resp, &badgeResp)
	assert.NoError(t, err)

	// Store the badge ID
	timeBasedCriteriaBadgeID = badgeResp.ID

	// Get the event type name
	var eventTypeName string
	resp = testutil.MakeRequest("GET", fmt.Sprintf("/api/v1/admin/event-types/%d", EventTypeID), nil)
	testutil.AssertSuccess(t, resp)
	var eventTypeResp EventTypeResponse
	err = testutil.ParseResponse(resp, &eventTypeResp)
	assert.NoError(t, err)
	eventTypeName = eventTypeResp.Name

	// Create event that meets the time-based criteria
	eventPayload := map[string]interface{}{
		"score":     float64(80), // >=75 ✓
		"completed": true,
	}

	// Create event request with timestamp inside the time window
	eventReq := EventRequest{
		EventType: eventTypeName,
		UserID:    fmt.Sprintf("%s_timewindow", TestUserID),
		Payload:   eventPayload,
		Timestamp: time.Now(), // Current time is within the window
	}

	// Make the API request to process the event
	resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Wait for badge processing
	time.Sleep(1 * time.Second)

	// Check if badge was awarded
	endpoint := fmt.Sprintf("/api/v1/users/%s_timewindow/badges", TestUserID)
	resp = testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var userBadgesResp UserBadgesResponse
	err = testutil.ParseResponse(resp, &userBadgesResp)
	assert.NoError(t, err)

	// Verify badge was awarded
	badgeFound := false
	for _, badge := range userBadgesResp.Badges {
		if badge.BadgeID == timeBasedCriteriaBadgeID {
			badgeFound = true
			break
		}
	}

	assert.True(t, badgeFound, "The time-based criteria badge should have been awarded to the user")
}
