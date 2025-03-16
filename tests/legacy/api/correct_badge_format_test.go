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

// TestCorrectBadgeFormat demonstrates the correct format for badge criteria
// that is compatible with the rule engine implementation
func TestCorrectBadgeFormat(t *testing.T) {
	// Create log file
	logFile, err := os.Create("correct_format_test.log")
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	log := func(format string, args ...interface{}) {
		msg := fmt.Sprintf(format, args...)
		t.Log(msg)
		fmt.Fprintln(logFile, msg)
	}

	log("Starting test with correct badge criteria format")

	// Create a unique timestamp for this test to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000

	// Step 1: Create an Event Type with a clear schema
	log("Step 1: Creating event type")
	eventTypeName := fmt.Sprintf("score_event_%d", timestamp)
	eventTypeReq := utils.EventTypeRequest{
		Name:        eventTypeName,
		Description: "Event with score for testing",
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

	// Step 2: Create a Condition Type with simple logic
	log("Step 2: Creating condition type")
	conditionTypeReq := utils.ConditionTypeRequest{
		Name:        fmt.Sprintf("score_condition_%d", timestamp),
		Description: "Condition that checks if score meets threshold",
		EvaluationLogic: `
			function evaluate(event, params) {
				console.log("DEBUG: Evaluating event with score:", event.payload.score);
				console.log("DEBUG: Against threshold:", params.threshold);
				return event.payload.score >= params.threshold;
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

	// Step 3: Create a Badge with criteria in the correct format
	log("Step 3: Creating badge with correct criteria format")

	// FORMAT 1: Event-based criteria (direct format)
	correctFlowDefinition := map[string]interface{}{
		"event": eventTypeName,
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": 10, // Threshold value
			},
		},
	}

	// Log the flow definition for debugging
	flowJSON, _ := json.MarshalIndent(correctFlowDefinition, "", "  ")
	log("Flow definition: %s", string(flowJSON))

	badgeReq := utils.BadgeRequest{
		Name:           fmt.Sprintf("Correct Format Badge %d", timestamp),
		Description:    "Badge with criteria in correct format",
		ImageURL:       "https://example.com/badges/correct.png",
		FlowDefinition: correctFlowDefinition,
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
	userID := fmt.Sprintf("correct-format-user-%d", timestamp)
	log("Test user ID: %s", userID)

	// Step 5: Submit an event that should trigger the badge
	log("Step 5: Submitting event with score 50")
	eventReq := utils.EventRequest{
		EventType: eventTypeName,
		UserID:    userID,
		Payload: map[string]interface{}{
			"score": 50, // Well above our threshold of 10
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

	// Step 6: Wait and check for badge assignment
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

	if !badgeFound {
		t.Errorf("No badges were assigned to user %s", userID)
		log("DEBUG SUMMARY: Badge award failed despite event submission")
	} else {
		log("SUCCESS: Badge was awarded properly")
	}

	// Now try with a logical operator format ($and)
	log("Creating second badge with $and logical operator format")

	logicalOpFlowDefinition := map[string]interface{}{
		"$and": []map[string]interface{}{
			{
				"event": eventTypeName,
				"criteria": map[string]interface{}{
					"score": map[string]interface{}{
						"$gte": 20, // Different threshold
					},
				},
			},
		},
	}

	// Log the flow definition for debugging
	logicalOpJSON, _ := json.MarshalIndent(logicalOpFlowDefinition, "", "  ")
	log("Logical operator flow definition: %s", string(logicalOpJSON))

	logicalBadgeReq := utils.BadgeRequest{
		Name:           fmt.Sprintf("Logical Op Badge %d", timestamp),
		Description:    "Badge with logical operator format",
		ImageURL:       "https://example.com/badges/logical.png",
		FlowDefinition: logicalOpFlowDefinition,
	}

	response = utils.MakeRequest(http.MethodPost, "/api/v1/admin/badges", logicalBadgeReq)
	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		log("Failed to create logical operator badge: %s", string(response.Body))
		t.Fatalf("Failed to create logical operator badge: %s", string(response.Body))
	}

	var logicalBadgeWithCriteria utils.BadgeWithCriteria
	err = utils.ParseResponse(response, &logicalBadgeWithCriteria)
	assert.NoError(t, err)
	log("Created logical operator badge with ID: %d and name: %s", logicalBadgeWithCriteria.Badge.ID, logicalBadgeWithCriteria.Badge.Name)

	// Create a different user for the logical operator test
	logicalUserID := fmt.Sprintf("logical-format-user-%d", timestamp)
	log("Logical operator test user ID: %s", logicalUserID)

	// Submit an event that should trigger the badge
	log("Submitting event with score 50 for logical operator test")
	logicalEventReq := utils.EventRequest{
		EventType: eventTypeName,
		UserID:    logicalUserID,
		Payload: map[string]interface{}{
			"score": 50, // Well above our threshold of 20
		},
	}

	response = utils.MakeRequest(http.MethodPost, "/api/v1/events", logicalEventReq)
	if response.StatusCode != http.StatusOK {
		log("Failed to submit event for logical operator test: %s", string(response.Body))
		t.Fatalf("Failed to submit event for logical operator test: %s", string(response.Body))
	}
	log("Event response for logical operator test: %s", string(response.Body))

	// Wait and check for badge assignment
	logicalBadgeFound := false

	log("Checking for logical operator badge assignment")
	for attempt := 1; attempt <= 5; attempt++ {
		log("Checking for logical operator badges (attempt %d/5)", attempt)
		time.Sleep(time.Duration(attempt) * time.Second)

		response = utils.MakeRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/badges", logicalUserID), nil)
		log("Get logical user badges response code: %d", response.StatusCode)
		log("Get logical user badges response body: %s", string(response.Body))

		if response.StatusCode != http.StatusOK {
			log("Failed to get logical user badges (status: %d): %s", response.StatusCode, string(response.Body))
			continue
		}

		var logicalUserBadges []utils.UserBadge
		err = utils.ParseResponse(response, &logicalUserBadges)
		if err != nil {
			log("Failed to parse logical user badges: %v", err)
			continue
		}

		if len(logicalUserBadges) > 0 {
			logicalBadgeFound = true
			badgeDetails, _ := json.MarshalIndent(logicalUserBadges, "", "  ")
			log("Logical user badges found: %s", string(badgeDetails))
			break
		}

		log("No logical badges found yet, waiting longer...")
	}

	if !logicalBadgeFound {
		t.Errorf("No badges were assigned to logical operator user %s", logicalUserID)
		log("LOGICAL DEBUG SUMMARY: Badge award failed for logical operator format")
	} else {
		log("LOGICAL SUCCESS: Badge with logical operator format was awarded properly")
	}
}
