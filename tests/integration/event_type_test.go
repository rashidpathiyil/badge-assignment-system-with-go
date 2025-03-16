package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/testutil"
	"github.com/stretchr/testify/assert"
)

// Global variable to store the event type ID created during tests
var testEventTypeID int

// TestCreateEventType tests the creation of an event type
func TestCreateEventType(t *testing.T) {
	SetupTest()

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	eventTypeName := fmt.Sprintf("Test Event Type_%d", timestamp)

	// Define the event type schema
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"score": map[string]interface{}{
				"type":        "integer",
				"description": "The score achieved",
			},
			"duration": map[string]interface{}{
				"type":        "number",
				"description": "Time taken in seconds",
			},
			"completed": map[string]interface{}{
				"type":        "boolean",
				"description": "Whether the activity was completed",
			},
		},
		"required": []string{"score", "completed"},
	}

	// Create the event type request
	eventTypeReq := EventTypeRequest{
		Name:        eventTypeName,
		Description: "An event type for testing purposes",
		Schema:      schema,
	}

	// Make the API request
	resp := testutil.MakeRequest("POST", "/api/v1/admin/event-types", eventTypeReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var eventTypeResp EventTypeResponse
	err := testutil.ParseResponse(resp, &eventTypeResp)
	assert.NoError(t, err)

	// Assert response fields
	assert.NotZero(t, eventTypeResp.ID)
	assert.Equal(t, eventTypeName, eventTypeResp.Name)
	assert.Equal(t, "An event type for testing purposes", eventTypeResp.Description)

	// Store the event type ID for subsequent tests
	testEventTypeID = eventTypeResp.ID
	EventTypeID = testEventTypeID // Update global variable for other tests
}

// TestGetEventType tests retrieving an event type by ID
func TestGetEventType(t *testing.T) {
	if testEventTypeID == 0 {
		t.Skip("Event type ID not set, skipping test")
	}

	// Make the API request
	endpoint := fmt.Sprintf("/api/v1/admin/event-types/%d", testEventTypeID)
	resp := testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var eventTypeResp EventTypeResponse
	err := testutil.ParseResponse(resp, &eventTypeResp)
	assert.NoError(t, err)

	// Assert response fields
	assert.Equal(t, testEventTypeID, eventTypeResp.ID)
	assert.NotEmpty(t, eventTypeResp.Name)
	assert.NotEmpty(t, eventTypeResp.Description)
	assert.NotNil(t, eventTypeResp.Schema)
}
