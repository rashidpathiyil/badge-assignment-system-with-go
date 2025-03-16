package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/badge-assignment-system/tests/api/utils"
	"github.com/stretchr/testify/assert"
)

// TestDebugBadgeAssignment performs a more thorough test of the badge assignment process
// with detailed logging to help debug issues
func TestDebugBadgeAssignment(t *testing.T) {
	// Skip in normal test runs unless explicitly requested
	if os.Getenv("RUN_DEBUG_TESTS") != "true" {
		t.Skip("Skipping debug test. Set RUN_DEBUG_TESTS=true to run this test.")
	}

	logFile, err := os.Create("badge_debug.log")
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	log := func(format string, args ...interface{}) {
		msg := fmt.Sprintf(format, args...)
		t.Log(msg)
		fmt.Fprintln(logFile, msg)
	}

	log("Starting badge assignment debug test")

	// Step 1: Create Event Type
	log("Step 1: Creating event type")
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

	timestamp := time.Now().UnixNano() / 1000000
	eventTypeName := fmt.Sprintf("debug_event_%d", timestamp)

	eventTypeReq := utils.EventTypeRequest{
		Name:        eventTypeName,
		Description: "Debug event for badge testing",
		Schema:      schema,
	}

	response := utils.MakeRequest(http.MethodPost, "/api/v1/admin/event-types", eventTypeReq)
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		log("Failed to create event type: %s", string(response.Body))
		t.Fatalf("Failed to create event type: %s", string(response.Body))
	}

	var eventType utils.EventType
	err = utils.ParseResponse(response, &eventType)
	assert.NoError(t, err)
	debug_eventTypeID := eventType.ID
	log("Created event type with ID: %d and name: %s", debug_eventTypeID, eventType.Name)

	// Step 2: Create Condition Type
	log("Step 2: Creating condition type")
	conditionTypeReq := utils.ConditionTypeRequest{
		Name:        fmt.Sprintf("debug_condition_%d", timestamp),
		Description: "Debug condition for badge testing",
		EvaluationLogic: `
			function evaluate(event, params) {
				console.log("Evaluating event:", JSON.stringify(event));
				console.log("With params:", JSON.stringify(params));
				return event.payload.score >= params.threshold;
			}
		`,
	}

	response = utils.MakeRequest(http.MethodPost, "/api/v1/admin/condition-types", conditionTypeReq)
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		log("Failed to create condition type: %s", string(response.Body))
		t.Fatalf("Failed to create condition type: %s", string(response.Body))
	}

	var conditionType utils.ConditionType
	err = utils.ParseResponse(response, &conditionType)
	assert.NoError(t, err)
	debug_conditionTypeID := conditionType.ID
	log("Created condition type with ID: %d and name: %s", debug_conditionTypeID, conditionType.Name)

	// Step 3: Create Badge with simple criteria
	log("Step 3: Creating badge with criteria")
	flowDefinition := map[string]interface{}{
		"steps": []map[string]interface{}{
			{
				"type":            "condition",
				"conditionTypeId": debug_conditionTypeID,
				"eventTypeIds":    []int{debug_eventTypeID},
				"conditionParams": map[string]interface{}{
					"threshold": 50, // Lower threshold to make it easier to trigger
				},
			},
		},
	}

	badgeReq := utils.BadgeRequest{
		Name:           fmt.Sprintf("Debug Badge %d", timestamp),
		Description:    "Debug badge for testing",
		ImageURL:       "https://example.com/badges/debug.png",
		FlowDefinition: flowDefinition,
	}

	response = utils.MakeRequest(http.MethodPost, "/api/v1/admin/badges", badgeReq)
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		log("Failed to create badge: %s", string(response.Body))
		t.Fatalf("Failed to create badge: %s", string(response.Body))
	}

	var badgeWithCriteria utils.BadgeWithCriteria
	err = utils.ParseResponse(response, &badgeWithCriteria)
	assert.NoError(t, err)
	debug_badgeID := badgeWithCriteria.Badge.ID
	log("Created badge with ID: %d and name: %s", debug_badgeID, badgeWithCriteria.Badge.Name)

	// Log the badge criteria details
	criteriaJSON, _ := json.MarshalIndent(badgeWithCriteria.Criteria, "", "  ")
	log("Badge criteria: %s", string(criteriaJSON))

	// Step 4: Create multiple events to trigger the badge
	log("Step 4: Submitting events to trigger badge assignment")

	// Generate a unique user ID for this test to avoid conflicts
	debug_userID := fmt.Sprintf("debug-user-%d", timestamp)

	// Submit 3 different events with different scores
	for i, score := range []int{60, 75, 90} {
		eventReq := utils.EventRequest{
			EventType: eventType.Name,
			UserID:    debug_userID,
			Payload: map[string]interface{}{
				"score": score,
				"level": "debug-level",
			},
		}

		log("Submitting event %d with score %d", i+1, score)
		response = utils.MakeRequest(http.MethodPost, "/api/v1/events", eventReq)
		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
			log("Failed to submit event: %s", string(response.Body))
			t.Fatalf("Failed to submit event: %s", string(response.Body))
		}

		log("Event %d response: %s", i+1, string(response.Body))

		// Wait a bit between events
		time.Sleep(500 * time.Millisecond)
	}

	// Step 5: Wait and check for badge assignment
	log("Step 5: Waiting for badge processing to complete...")

	// Try checking for badges multiple times with delays
	var userBadges []utils.UserBadge
	badgeFound := false

	for attempt := 1; attempt <= 5; attempt++ {
		log("Checking for badges (attempt %d/5)", attempt)
		time.Sleep(time.Duration(attempt) * time.Second)

		response = utils.MakeRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/badges", debug_userID), nil)
		if response.StatusCode != http.StatusOK {
			log("Failed to get user badges: %s", string(response.Body))
			continue
		}

		err = utils.ParseResponse(response, &userBadges)
		if err != nil {
			log("Failed to parse user badges: %v", err)
			continue
		}

		log("User badges response: %s", string(response.Body))

		if len(userBadges) > 0 {
			badgeFound = true
			break
		}

		log("No badges found yet, waiting longer...")
	}

	if badgeFound {
		log("Success! User has been awarded badges:")
		for _, badge := range userBadges {
			log("  - Badge ID: %d, Awarded At: %s", badge.BadgeID, badge.AwardedAt)
		}
	} else {
		log("Failed to find any badges for the user after multiple attempts")

		// Try to get more diagnostic information
		log("Checking server endpoints for additional information:")

		// Check if there's an events endpoint
		response = utils.MakeRequest(http.MethodGet, "/api/v1/events", nil)
		log("Events endpoint response code: %d", response.StatusCode)
		if response.StatusCode == http.StatusOK {
			log("Events data: %s", string(response.Body))
		}

		// Check if there's an endpoint for badge criteria evaluation status
		response = utils.MakeRequest(http.MethodGet, fmt.Sprintf("/api/v1/admin/badges/%d/evaluations", debug_badgeID), nil)
		log("Badge evaluations response code: %d", response.StatusCode)
		if response.StatusCode == http.StatusOK {
			log("Badge evaluations data: %s", string(response.Body))
		}

		t.Errorf("No badges were assigned to user %s", debug_userID)
	}

	log("Debug test completed")
}
