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

// TestDebugRuleEngine is a focused test to debug the rule engine's badge assignment
func TestDebugRuleEngine(t *testing.T) {
	// Create log file
	logFile, err := os.Create("rule_engine_debug.log")
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	log := func(format string, args ...interface{}) {
		msg := fmt.Sprintf(format, args...)
		t.Log(msg)
		fmt.Fprintln(logFile, msg)
	}

	log("Starting rule engine debug test")

	// Create a unique timestamp for this test to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000

	// Step 1: Create an Event Type with a very clear schema
	log("Step 1: Creating event type")
	eventTypeReq := utils.EventTypeRequest{
		Name:        fmt.Sprintf("score_event_%d", timestamp),
		Description: "Debug event for rule engine testing",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"score": map[string]interface{}{
					"type":        "integer",
					"description": "The score value",
				},
			},
			"required": []string{"score"},
		},
	}

	response := utils.MakeRequest(http.MethodPost, "/api/v1/admin/event-types", eventTypeReq)
	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		log("Failed to create event type: %s", string(response.Body))
		t.Fatalf("Failed to create event type: %s", string(response.Body))
	}

	var eventType utils.EventType
	err = utils.ParseResponse(response, &eventType)
	assert.NoError(t, err)
	log("Created event type with ID: %d and name: %s", eventType.ID, eventType.Name)
	log("Event type schema: %s", eventType.Schema)

	// Step 2: Create a Condition Type with very simple logic and debug logging
	log("Step 2: Creating condition type")
	conditionTypeReq := utils.ConditionTypeRequest{
		Name:        fmt.Sprintf("debug_condition_%d", timestamp),
		Description: "Debug condition with extensive logging",
		EvaluationLogic: `
			function evaluate(event, params) {
				// Log everything for debugging
				console.log("DEBUG: Event object:", JSON.stringify(event));
				console.log("DEBUG: Event type:", event.event_type_id);
				console.log("DEBUG: Event payload:", JSON.stringify(event.payload));
				console.log("DEBUG: Params:", JSON.stringify(params));
				
				// Very basic condition that should always succeed if event has score
				if (event.payload && typeof event.payload.score === 'number') {
					console.log("DEBUG: Score found:", event.payload.score);
					console.log("DEBUG: Threshold:", params.threshold);
					console.log("DEBUG: Evaluation result:", event.payload.score >= params.threshold);
					return event.payload.score >= params.threshold;
				} else {
					console.log("DEBUG: Score not found in payload or not a number");
					console.log("DEBUG: Payload keys:", Object.keys(event.payload || {}));
					return false;
				}
			}
		`,
	}

	response = utils.MakeRequest(http.MethodPost, "/api/v1/admin/condition-types", conditionTypeReq)
	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		log("Failed to create condition type: %s", string(response.Body))
		t.Fatalf("Failed to create condition type: %s", string(response.Body))
	}

	var conditionType utils.ConditionType
	err = utils.ParseResponse(response, &conditionType)
	assert.NoError(t, err)
	log("Created condition type with ID: %d and name: %s", conditionType.ID, conditionType.Name)

	// Step 3: Create a Badge with very simple criteria
	log("Step 3: Creating badge with simple criteria")
	// Use a very low threshold to ensure the condition is met
	flowDefinition := map[string]interface{}{
		"steps": []map[string]interface{}{
			{
				"type":            "condition",
				"conditionTypeId": conditionType.ID,
				"eventTypeIds":    []int{eventType.ID},
				"conditionParams": map[string]interface{}{
					"threshold": 1, // Very low threshold that should be easy to meet
				},
			},
		},
	}

	// Log the flow definition for debugging
	flowJSON, _ := json.MarshalIndent(flowDefinition, "", "  ")
	log("Flow definition: %s", string(flowJSON))

	badgeReq := utils.BadgeRequest{
		Name:           fmt.Sprintf("Debug Badge %d", timestamp),
		Description:    "Debug badge for rule engine testing",
		ImageURL:       "https://example.com/badges/debug.png",
		FlowDefinition: flowDefinition,
	}

	response = utils.MakeRequest(http.MethodPost, "/api/v1/admin/badges", badgeReq)
	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		log("Failed to create badge: %s", string(response.Body))
		t.Fatalf("Failed to create badge: %s", string(response.Body))
	}

	var badgeWithCriteria utils.BadgeWithCriteria
	err = utils.ParseResponse(response, &badgeWithCriteria)
	assert.NoError(t, err)
	log("Created badge with ID: %d and name: %s", badgeWithCriteria.Badge.ID, badgeWithCriteria.Badge.Name)

	// Step 4: Generate a unique user ID for this test
	userID := fmt.Sprintf("debug-user-%d", timestamp)
	log("Test user ID: %s", userID)

	// Step 5: Submit an event with a score that should trigger the badge
	log("Step 5: Submitting event with score 50")
	eventReq := utils.EventRequest{
		EventType: eventType.Name,
		UserID:    userID,
		Payload: map[string]interface{}{
			"score": 50, // Well above the threshold of 1
		},
	}

	// Log the event request for debugging
	eventJSON, _ := json.MarshalIndent(eventReq, "", "  ")
	log("Event request: %s", string(eventJSON))

	response = utils.MakeRequest(http.MethodPost, "/api/v1/events", eventReq)
	if response.StatusCode != http.StatusOK {
		log("Failed to submit event: %s", string(response.Body))
		t.Fatalf("Failed to submit event: %s", string(response.Body))
	}
	log("Event response: %s", string(response.Body))

	// Step 6: Wait and check for badge assignment with multiple attempts
	var userBadges []utils.UserBadge
	badgeFound := false

	log("Step 6: Checking for badge assignment")
	for attempt := 1; attempt <= 5; attempt++ {
		log("Checking for badges (attempt %d/5)", attempt)
		time.Sleep(time.Duration(attempt) * time.Second)

		response = utils.MakeRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/badges", userID), nil)
		log("Get user badges response code: %d", response.StatusCode)
		log("Get user badges response body: %s", string(response.Body))

		if response.StatusCode != http.StatusOK {
			log("Failed to get user badges (status: %d): %s", response.StatusCode, string(response.Body))
			continue
		}

		err = utils.ParseResponse(response, &userBadges)
		if err != nil {
			log("Failed to parse user badges: %v", err)
			continue
		}

		if len(userBadges) > 0 {
			badgeFound = true
			badgeDetails, _ := json.MarshalIndent(userBadges, "", "  ")
			log("User badges found: %s", string(badgeDetails))
			break
		}

		log("No badges found yet, waiting longer...")
	}

	// Step 7: Check specific badge details
	if !badgeFound {
		log("DEBUG: No badges assigned. Checking badge details to verify references...")

		// Check badge criteria in the database
		response = utils.MakeRequest(http.MethodGet, fmt.Sprintf("/api/v1/admin/badges/%d/criteria", badgeWithCriteria.Badge.ID), nil)
		log("Badge criteria response: %s", string(response.Body))

		// Check events in the database
		response = utils.MakeRequest(http.MethodGet, fmt.Sprintf("/api/v1/admin/events/user/%s", userID), nil)
		if response.StatusCode == http.StatusOK {
			log("User events: %s", string(response.Body))
		} else {
			log("Could not get user events, endpoint may not exist: %d", response.StatusCode)
		}

		log("DEBUG SUMMARY: Badge award failed despite event submission")
		t.Errorf("No badges were assigned to user %s", userID)
	} else {
		log("SUCCESS: Badge was awarded properly")
	}
}
