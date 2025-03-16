package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/badge-assignment-system/tests/api/utils"
	"github.com/stretchr/testify/assert"
)

// TestCreateEventType tests the creation of an event type
func TestCreateEventType(t *testing.T) {
	// Define the event type schema
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"score": map[string]interface{}{
				"type":        "integer",
				"description": "The score achieved",
			},
			"level": map[string]interface{}{
				"type":        "string",
				"description": "The difficulty level",
			},
		},
		"required": []string{"score", "level"},
	}

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	eventTypeName := fmt.Sprintf("challenge_completed_%d", timestamp)

	// Create the event type request
	eventTypeReq := utils.EventTypeRequest{
		Name:        eventTypeName,
		Description: "Triggered when a user completes a challenge",
		Schema:      schema,
	}

	// Make the API request
	response := utils.MakeRequest(http.MethodPost, "/api/v1/admin/event-types", eventTypeReq)
	utils.AssertSuccess(t, response)

	// Parse the response
	var eventType utils.EventType
	err := utils.ParseResponse(response, &eventType)
	assert.NoError(t, err)

	// Validate response data
	assert.NotZero(t, eventType.ID)
	assert.Equal(t, eventTypeReq.Name, eventType.Name)
	assert.Equal(t, eventTypeReq.Description, eventType.Description)

	// Store the event type ID for later use
	eventTypeID = eventType.ID
	t.Logf("Created event type with ID: %d", eventTypeID)
}

// TestGetEventType tests retrieving an event type
func TestGetEventType(t *testing.T) {
	if eventTypeID == 0 {
		t.Skip("Event type ID not set, skipping test")
	}

	// Make the API request
	response := utils.MakeRequest(http.MethodGet, fmt.Sprintf("/api/v1/admin/event-types/%d", eventTypeID), nil)
	utils.AssertSuccess(t, response)

	// Parse the response
	var eventType utils.EventType
	err := utils.ParseResponse(response, &eventType)
	assert.NoError(t, err)

	// Validate response data
	assert.Equal(t, eventTypeID, eventType.ID)
	assert.NotEmpty(t, eventType.Name)
}
