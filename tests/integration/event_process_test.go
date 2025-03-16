package integration

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/testutil"
	"github.com/stretchr/testify/assert"
)

// TestProcessEvent tests processing an event and triggering badge assignment
func TestProcessEvent(t *testing.T) {
	SetupTest()

	// Use the existing event type provided by the user
	existingEventTypeID := "issue_reported_super_1742154397309"
	userID := "test_user_super"

	// Create event using the map format (this is how other tests create events successfully)
	eventReq := map[string]interface{}{
		"event_type": existingEventTypeID,
		"user_id":    userID,
		"timestamp":  time.Now().Format(time.RFC3339),
		"payload": map[string]interface{}{
			"score":     float64(95),
			"duration":  float64(120.5),
			"completed": true,
		},
	}

	// Make the API request
	resp := testutil.MakeRequest("POST", "/api/v1/events", eventReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Log the full response for debugging
	t.Logf("Event creation response: %s", string(resp.Body))

	// Parse response and check for success message
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body, &response)
	assert.NoError(t, err)

	// Check for success message
	message, ok := response["message"].(string)
	assert.True(t, ok, "Response should contain a message field")
	assert.Equal(t, "Event processed successfully", message, "Response should indicate successful processing")

	// Wait for badge processing to complete
	time.Sleep(4 * time.Second)
}

// TestGetUserBadges tests retrieving badges awarded to a user
func TestGetUserBadges(t *testing.T) {
	// Use the existing user provided by the user
	userID := "test_user_super"

	// Make the API request
	endpoint := fmt.Sprintf("/api/v1/users/%s/badges", userID)
	resp := testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badges []map[string]interface{}
	err := json.Unmarshal(resp.Body, &badges)
	assert.NoError(t, err)

	// Assert we have at least one badge
	assert.True(t, len(badges) > 0, "The user should have at least one badge")
}
