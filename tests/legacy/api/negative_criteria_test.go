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

// TestNegativeCriteria verifies that badges are NOT awarded when criteria aren't met
func TestNegativeCriteria(t *testing.T) {
	// Create log file
	logFile, err := os.Create("negative_criteria_test.log")
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	log := func(format string, args ...interface{}) {
		msg := fmt.Sprintf(format, args...)
		t.Log(msg)
		fmt.Fprintln(logFile, msg)
	}

	log("Starting negative criteria test - validating badges are NOT awarded when criteria aren't met")

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

	// Step 2: Create a Badge with criteria that requires a high score
	log("Step 2: Creating badge with high threshold criteria")

	// Setting a high threshold that won't be met
	highThresholdFlowDefinition := map[string]interface{}{
		"event": eventTypeName,
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": 90, // High threshold that won't be met
			},
		},
	}

	// Log the flow definition for debugging
	flowJSON, _ := json.MarshalIndent(highThresholdFlowDefinition, "", "  ")
	log("Flow definition with high threshold: %s", string(flowJSON))

	badgeReq := utils.BadgeRequest{
		Name:           fmt.Sprintf("High Threshold Badge %d", timestamp),
		Description:    "Badge with criteria requiring a high score",
		ImageURL:       "https://example.com/badges/high-score.png",
		FlowDefinition: highThresholdFlowDefinition,
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

	// Step 3: Create a Badge with criteria that requires a low score (control case)
	log("Step 3: Creating control badge with low threshold criteria")

	// Setting a low threshold that will be met
	lowThresholdFlowDefinition := map[string]interface{}{
		"event": eventTypeName,
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": 10, // Low threshold that will be met
			},
		},
	}

	controlBadgeReq := utils.BadgeRequest{
		Name:           fmt.Sprintf("Low Threshold Badge %d", timestamp),
		Description:    "Badge with criteria requiring a low score (control)",
		ImageURL:       "https://example.com/badges/low-score.png",
		FlowDefinition: lowThresholdFlowDefinition,
	}

	response = utils.MakeRequest(http.MethodPost, "/api/v1/admin/badges", controlBadgeReq)
	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		log("Failed to create control badge: %s", string(response.Body))
		t.Fatalf("Failed to create control badge: %s", string(response.Body))
	}

	var controlBadgeWithCriteria utils.BadgeWithCriteria
	err = utils.ParseResponse(response, &controlBadgeWithCriteria)
	assert.NoError(t, err)
	log("Created control badge with ID: %d and name: %s", controlBadgeWithCriteria.Badge.ID, controlBadgeWithCriteria.Badge.Name)

	// Step 4: Generate a unique user ID for this test
	userID := fmt.Sprintf("negative-test-user-%d", timestamp)
	log("Test user ID: %s", userID)

	// Step 5: Submit an event with a medium score that meets control criteria but not high threshold
	log("Step 5: Submitting event with score 50 (meets control criteria but not high threshold)")
	eventReq := utils.EventRequest{
		EventType: eventTypeName,
		UserID:    userID,
		Payload: map[string]interface{}{
			"score": 50, // Above low threshold (10) but below high threshold (90)
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

	// Step 6: Wait and check that only the control badge was awarded
	var userBadges []utils.UserBadge
	log("Step 6: Checking for badge assignment")

	// Wait for badge processing
	time.Sleep(3 * time.Second)

	response = utils.MakeRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/badges", userID), nil)
	log("Get user badges response code: %d", response.StatusCode)
	log("Get user badges response body: %s", string(response.Body))

	if response.StatusCode != http.StatusOK {
		log("Failed to get user badges: %s", string(response.Body))
		t.Fatalf("Failed to get user badges: %s", string(response.Body))
	}

	err = utils.ParseResponse(response, &userBadges)
	if err != nil {
		log("Failed to parse user badges: %v", err)
		t.Fatalf("Failed to parse user badges: %v", err)
	}

	// Verify that user received only the control badge (low threshold) and not the high threshold badge
	gotHighThresholdBadge := false
	gotLowThresholdBadge := false

	for _, badge := range userBadges {
		if badge.ID == badgeWithCriteria.Badge.ID {
			gotHighThresholdBadge = true
			log("ERROR: User incorrectly received the high threshold badge")
		}
		if badge.ID == controlBadgeWithCriteria.Badge.ID {
			gotLowThresholdBadge = true
			log("User correctly received the control badge (low threshold)")
		}
	}

	// The test passes if:
	// 1. User did NOT get the high threshold badge (criteria not met)
	// 2. User DID get the low threshold badge (criteria met - control case)
	if gotHighThresholdBadge {
		t.Errorf("User incorrectly received the high threshold badge (score 50 < required 90)")
		log("TEST FAILED: User incorrectly received the high threshold badge")
	} else {
		log("TEST PASSED: User correctly did NOT receive the high threshold badge")
	}

	if !gotLowThresholdBadge {
		t.Errorf("User did not receive the control badge (score 50 > required 10)")
		log("CONTROL TEST FAILED: User did not receive the control badge")
	} else {
		log("CONTROL TEST PASSED: User correctly received the control badge")
	}

	// Log the complete summary
	if !gotHighThresholdBadge && gotLowThresholdBadge {
		log("OVERALL TEST PASSED: Badge system correctly awarded only badges with met criteria")
	} else {
		log("OVERALL TEST FAILED: Badge system did not behave as expected")
	}
}
