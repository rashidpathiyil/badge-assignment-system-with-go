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

// TestMongoOperatorFormat tests badge assignment using the MongoDB-style operator format
func TestMongoOperatorFormat(t *testing.T) {
	// Create log file
	logFile, err := os.Create("mongo_operator_test.log")
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	log := func(format string, args ...interface{}) {
		msg := fmt.Sprintf(format, args...)
		t.Log(msg)
		fmt.Fprintln(logFile, msg)
	}

	log("Starting MongoDB-style operator format test")

	// Create a unique timestamp for this test to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000

	// Step 1: Create an Event Type with a very clear schema
	log("Step 1: Creating event type")
	eventTypeReq := utils.EventTypeRequest{
		Name:        fmt.Sprintf("mongo_event_%d", timestamp),
		Description: "Event with score for testing MongoDB-style operators",
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
		Name:        fmt.Sprintf("mongo_condition_%d", timestamp),
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

	// Step 3: Create a Badge with MongoDB-style operator format criteria
	log("Step 3: Creating badge with MongoDB-style operator format criteria")

	// Use a simple MongoDB-style query format
	// 1. Try with direct field comparison
	mongoFlowDefinition := map[string]interface{}{
		"payload.score": map[string]interface{}{
			"$gte": 10, // Threshold value
		},
		"event_type_id": eventType.ID,
	}

	// Log the flow definition for debugging
	flowJSON, _ := json.MarshalIndent(mongoFlowDefinition, "", "  ")
	log("Flow definition: %s", string(flowJSON))

	badgeReq := utils.BadgeRequest{
		Name:           fmt.Sprintf("Mongo Op Badge %d", timestamp),
		Description:    "Badge with MongoDB-style operator format criteria",
		ImageURL:       "https://example.com/badges/mongo.png",
		FlowDefinition: mongoFlowDefinition,
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
	userID := fmt.Sprintf("mongo-user-%d", timestamp)
	log("Test user ID: %s", userID)

	// Step 5: Submit an event that should trigger the badge
	log("Step 5: Submitting event with score 50")
	eventReq := utils.EventRequest{
		EventType: eventType.Name,
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
}
