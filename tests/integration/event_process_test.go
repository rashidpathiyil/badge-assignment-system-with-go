package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/testutil"
	"github.com/stretchr/testify/assert"
)

// TestProcessEvent tests processing an event and triggering badge assignment
func TestProcessEvent(t *testing.T) {
	SetupTest()

	if EventTypeID == 0 || BadgeID == 0 {
		t.Skip("Event type ID or badge ID not set, skipping test")
	}

	// Get the event type name
	var eventTypeName string
	resp := testutil.MakeRequest("GET", fmt.Sprintf("/api/v1/admin/event-types/%d", EventTypeID), nil)
	testutil.AssertSuccess(t, resp)
	var eventTypeResp EventTypeResponse
	err := testutil.ParseResponse(resp, &eventTypeResp)
	assert.NoError(t, err)
	eventTypeName = eventTypeResp.Name

	// Create event data that meets the badge criteria
	eventPayload := map[string]interface{}{
		"score":     float64(95),                     // Above the 90 threshold
		"duration":  float64(120.5),                  // Example duration
		"completed": true,                            // Required field
		"timestamp": time.Now().Format(time.RFC3339), // Current time
	}

	// Create event request
	eventReq := EventRequest{
		EventType: eventTypeName,
		UserID:    TestUserID,
		Payload:   eventPayload,
		Timestamp: time.Now(),
	}

	// Make the API request
	resp = testutil.MakeRequest("POST", "/api/v1/events", eventReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var eventResp EventResponse
	err = testutil.ParseResponse(resp, &eventResp)
	assert.NoError(t, err)

	// Assert response fields
	assert.NotZero(t, eventResp.ID)
	assert.Equal(t, EventTypeID, eventResp.EventTypeID)
	assert.Equal(t, TestUserID, eventResp.UserID)
	assert.NotNil(t, eventResp.Data)

	// Wait a moment for badge processing to complete
	time.Sleep(1 * time.Second)
}

// TestGetUserBadges tests retrieving badges awarded to a user
func TestGetUserBadges(t *testing.T) {
	if BadgeID == 0 {
		t.Skip("Badge ID not set, skipping test")
	}

	// Make the API request
	endpoint := fmt.Sprintf("/api/v1/users/%s/badges", TestUserID)
	resp := testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var userBadgesResp UserBadgesResponse
	err := testutil.ParseResponse(resp, &userBadgesResp)
	assert.NoError(t, err)

	// Assert response fields
	assert.Equal(t, TestUserID, userBadgesResp.UserID)

	// Check if our badge was awarded
	badgeFound := false
	for _, badge := range userBadgesResp.Badges {
		if badge.BadgeID == BadgeID {
			badgeFound = true
			assert.NotEmpty(t, badge.BadgeName)
			assert.NotEmpty(t, badge.Description)
			assert.NotEmpty(t, badge.ImageURL)
			assert.False(t, badge.AwardedAt.IsZero(), "AwardedAt should be a valid timestamp")
			break
		}
	}

	assert.True(t, badgeFound, "The test badge should have been awarded to the user")
}
