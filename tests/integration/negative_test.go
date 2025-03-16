package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/testutil"
	"github.com/stretchr/testify/assert"
)

// Global variables for negative tests
var (
	negativeTestBadgeID int
)

// TestNegativeScenarios tests that events that don't match criteria don't earn badges
func TestNegativeScenarios(t *testing.T) {
	SetupTest()

	if EventTypeID == 0 {
		t.Skip("Event type ID not set, skipping test")
	}

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	badgeName := fmt.Sprintf("Negative Test Badge_%d", timestamp)

	// Define criteria requiring high score
	criteria := map[string]interface{}{
		"event": "test_event",
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": float64(90),
			},
		},
	}

	// Create badge request
	badgeReq := BadgeRequest{
		Name:           badgeName,
		Description:    "Badge that requires high score",
		ImageURL:       "https://example.com/badges/high-score.png",
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
	negativeTestBadgeID = badgeResp.ID

	// Get the event type name
	var eventTypeName string
	resp = testutil.MakeRequest("GET", fmt.Sprintf("/api/v1/admin/event-types/%d", EventTypeID), nil)
	testutil.AssertSuccess(t, resp)
	var eventTypeResp EventTypeResponse
	err = testutil.ParseResponse(resp, &eventTypeResp)
	assert.NoError(t, err)
	eventTypeName = eventTypeResp.Name

	// Create event that does NOT meet the criteria (score < 90)
	eventPayload := map[string]interface{}{
		"score":     float64(85), // Less than required 90
		"completed": true,
	}

	// Create event request
	eventReq := EventRequest{
		EventType: eventTypeName,
		UserID:    fmt.Sprintf("%s_negative", TestUserID),
		Payload:   eventPayload,
		Timestamp: time.Now(),
	}

	// Make the API request to process the event
	resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Wait for badge processing
	time.Sleep(1 * time.Second)

	// Check if badge was awarded (it should NOT be)
	endpoint := fmt.Sprintf("/api/v1/users/%s_negative/badges", TestUserID)
	resp = testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var userBadgesResp UserBadgesResponse
	err = testutil.ParseResponse(resp, &userBadgesResp)
	assert.NoError(t, err)

	// Verify badge was NOT awarded
	badgeFound := false
	for _, badge := range userBadgesResp.Badges {
		if badge.BadgeID == negativeTestBadgeID {
			badgeFound = true
			break
		}
	}

	assert.False(t, badgeFound, "The badge should NOT have been awarded to the user")
}

// TestInactiveBadge tests that inactive badges do not get awarded
func TestInactiveBadge(t *testing.T) {
	SetupTest()

	if EventTypeID == 0 {
		t.Skip("Event type ID not set, skipping test")
	}

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	badgeName := fmt.Sprintf("Inactive Badge_%d", timestamp)

	// Define simple criteria
	criteria := map[string]interface{}{
		"event": "test_event",
		"criteria": map[string]interface{}{
			"completed": map[string]interface{}{
				"$eq": true,
			},
		},
	}

	// Create badge request with isActive = false
	badgeReq := BadgeRequest{
		Name:           badgeName,
		Description:    "Badge that is inactive",
		ImageURL:       "https://example.com/badges/inactive.png",
		FlowDefinition: criteria,
		IsActive:       false, // Important: badge is inactive
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
	inactiveBadgeID := badgeResp.ID

	// Get the event type name
	var eventTypeName string
	resp = testutil.MakeRequest("GET", fmt.Sprintf("/api/v1/admin/event-types/%d", EventTypeID), nil)
	testutil.AssertSuccess(t, resp)
	var eventTypeResp EventTypeResponse
	err = testutil.ParseResponse(resp, &eventTypeResp)
	assert.NoError(t, err)
	eventTypeName = eventTypeResp.Name

	// Create event that would meet the criteria if badge was active
	eventPayload := map[string]interface{}{
		"completed": true,
	}

	// Create event request
	eventReq := EventRequest{
		EventType: eventTypeName,
		UserID:    fmt.Sprintf("%s_inactive", TestUserID),
		Payload:   eventPayload,
		Timestamp: time.Now(),
	}

	// Make the API request to process the event
	resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Wait for badge processing
	time.Sleep(1 * time.Second)

	// Check if badge was awarded (it should NOT be because it's inactive)
	endpoint := fmt.Sprintf("/api/v1/users/%s_inactive/badges", TestUserID)
	resp = testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var userBadgesResp UserBadgesResponse
	err = testutil.ParseResponse(resp, &userBadgesResp)
	assert.NoError(t, err)

	// Verify badge was NOT awarded
	badgeFound := false
	for _, badge := range userBadgesResp.Badges {
		if badge.BadgeID == inactiveBadgeID {
			badgeFound = true
			break
		}
	}

	assert.False(t, badgeFound, "The inactive badge should NOT have been awarded to the user")
}

// TestInvalidBadgeRequest tests invalid badge requests
func TestInvalidBadgeRequest(t *testing.T) {
	SetupTest()

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	badgeName := fmt.Sprintf("Invalid Badge_%d", timestamp)

	// Create an invalid badge request (missing required fields)
	invalidBadgeReq := BadgeRequest{
		Name:           badgeName,
		Description:    "This badge request is invalid",
		ImageURL:       "https://example.com/badges/invalid.png",
		FlowDefinition: nil, // Missing flow definition
		IsActive:       true,
	}

	// Make the API request
	resp := testutil.MakeRequest("POST", "/api/v1/admin/badges", invalidBadgeReq)

	// Assert response
	testutil.AssertError(t, resp, "Badge request is invalid")
}

// TestInvalidCriteriaFormat tests invalid badge criteria formats
func TestInvalidCriteriaFormat(t *testing.T) {
	SetupTest()

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	badgeName := fmt.Sprintf("Invalid Criteria_%d", timestamp)

	// Create a badge request with nil flow definition
	invalidCriteriaReq := BadgeRequest{
		Name:           badgeName,
		Description:    "This badge has invalid criteria format",
		ImageURL:       "https://example.com/badges/invalid-criteria.png",
		FlowDefinition: nil, // Nil flow definition should cause an error
		IsActive:       true,
	}

	// Make the API request
	resp := testutil.MakeRequest("POST", "/api/v1/admin/badges", invalidCriteriaReq)

	// Assert response
	testutil.AssertError(t, resp, "Badge request is invalid")
}

// TestInvalidEventRequest tests invalid event requests
func TestInvalidEventRequest(t *testing.T) {
	SetupTest()

	// Attempt to create an event with invalid data
	invalidEventReq := EventRequest{
		EventType: "non_existent_event_type", // Invalid event type
		UserID:    "",                        // Missing user ID
		Payload:   nil,                       // Missing payload
		Timestamp: time.Now(),
	}

	// Make the API request
	resp := testutil.MakeRequest("POST", "/api/v1/events", invalidEventReq)

	// Assert response
	testutil.AssertError(t, resp, "Event request is invalid")

	// Attempt to create an event with invalid event type but valid user ID
	invalidEventTypeReq := EventRequest{
		EventType: "non_existent_event", // Event type that doesn't exist
		UserID:    TestUserID,
		Payload: map[string]interface{}{
			"score": float64(95),
		},
		Timestamp: time.Now(),
	}

	// Make the API request
	resp = testutil.MakeRequest("POST", "/api/v1/events", invalidEventTypeReq)

	// Assert response
	testutil.AssertError(t, resp, "Event request is invalid")
}
