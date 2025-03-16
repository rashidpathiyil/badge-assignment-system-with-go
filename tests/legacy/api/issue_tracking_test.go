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

// TestIssueBadgeAssignment tests awarding a badge when a user fixes 3 high-priority issues
func TestIssueBadgeAssignment(t *testing.T) {
	// Create log file
	logFile, err := os.Create("issue_badge_test.log")
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	log := func(format string, args ...interface{}) {
		msg := fmt.Sprintf(format, args...)
		t.Log(msg)
		fmt.Fprintln(logFile, msg)
	}

	log("Starting issue badge test - validating a badge is awarded after fixing 3 high-priority issues")

	// Create a unique timestamp for this test to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000

	// Step 1: Create an Event Type for issue fixes with priority
	log("Step 1: Creating issue_fixed event type")
	eventTypeName := fmt.Sprintf("issue_fixed_%d", timestamp)
	eventTypeReq := utils.EventTypeRequest{
		Name:        eventTypeName,
		Description: "Event when an issue is fixed",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"issue_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the fixed issue",
				},
				"priority": map[string]interface{}{
					"type":        "string",
					"description": "The priority of the issue (high, medium, low)",
					"enum":        []string{"high", "medium", "low"},
				},
				"component": map[string]interface{}{
					"type":        "string",
					"description": "The component where the issue was fixed",
				},
			},
			"required": []string{"issue_id", "priority"},
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

	// Step 2: Create a Badge for fixing 3 high-priority issues
	log("Step 2: Creating High Priority Fixer badge")

	// For counting events with a specific field value, we'll use the $and operator
	// to filter for events of the specified type and with priority=high
	badgeFlowDefinition := map[string]interface{}{
		"$count": map[string]interface{}{
			"target": 3,
			"event":  eventTypeName,
			"filter": map[string]interface{}{
				"priority": map[string]interface{}{
					"$eq": "high",
				},
			},
		},
	}

	// Log the flow definition for debugging
	flowJSON, _ := json.MarshalIndent(badgeFlowDefinition, "", "  ")
	log("Flow definition: %s", string(flowJSON))

	badgeReq := utils.BadgeRequest{
		Name:           fmt.Sprintf("High Priority Fixer %d", timestamp),
		Description:    "Awarded for fixing 3 high-priority issues",
		ImageURL:       "https://example.com/badges/high-priority-fixer.png",
		FlowDefinition: badgeFlowDefinition,
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

	// Step 3: Generate a unique user ID for this test
	userID := fmt.Sprintf("issue-fixer-user-%d", timestamp)
	log("Test user ID: %s", userID)

	// Step 4: Submit events for fixing issues
	log("Step 4: Submitting events for issue fixes")

	// First, submit 2 high-priority fixes - should NOT trigger the badge yet
	for i := 1; i <= 2; i++ {
		submitIssueFixedEvent(t, log, eventTypeName, userID, fmt.Sprintf("HIGH-%d", i), "high", fmt.Sprintf("component-%d", i))
	}

	// Check that badge is NOT awarded yet
	log("Checking badges after 2 high-priority fixes (should be none)")
	userBadges := getUserBadges(t, log, userID)

	// Verify user has not received the badge yet
	if hasBadge(userBadges, badgeWithCriteria.Badge.ID) {
		t.Errorf("User incorrectly received the badge after only 2 high-priority fixes")
		log("TEST FAILED: Badge was incorrectly awarded after only 2 high-priority fixes")
	} else {
		log("CORRECT: No badge awarded after only 2 high-priority fixes")
	}

	// Submit a medium-priority fix - should still NOT trigger the badge
	submitIssueFixedEvent(t, log, eventTypeName, userID, "MED-1", "medium", "component-3")

	// Check that badge is still NOT awarded
	log("Checking badges after 2 high + 1 medium priority fixes (should be none)")
	userBadges = getUserBadges(t, log, userID)

	// Verify user still has not received the badge
	if hasBadge(userBadges, badgeWithCriteria.Badge.ID) {
		t.Errorf("User incorrectly received the badge after 2 high + 1 medium priority fixes")
		log("TEST FAILED: Badge was incorrectly awarded after 2 high + 1 medium priority fixes")
	} else {
		log("CORRECT: No badge awarded after 2 high + 1 medium priority fixes")
	}

	// Submit a third high-priority fix - NOW should trigger the badge
	log("Submitting the 3rd high-priority fix - should trigger the badge")
	submitIssueFixedEvent(t, log, eventTypeName, userID, "HIGH-3", "high", "component-4")

	// Wait for badge processing
	log("Waiting for badge processing...")
	time.Sleep(3 * time.Second)

	// Check that badge is awarded
	log("Checking badges after 3 high-priority fixes (should be awarded)")
	userBadges = getUserBadges(t, log, userID)

	// Verify user has received the badge
	if !hasBadge(userBadges, badgeWithCriteria.Badge.ID) {
		t.Errorf("User did not receive the High Priority Fixer badge after 3 high-priority fixes")
		log("TEST FAILED: Badge was not awarded after 3 high-priority fixes")
	} else {
		log("SUCCESS: Badge correctly awarded after 3 high-priority fixes")
	}

	log("Test completed successfully")
}

// Helper function to submit an issue_fixed event
func submitIssueFixedEvent(t *testing.T, log func(string, ...interface{}), eventTypeName, userID, issueID, priority, component string) {
	eventReq := utils.EventRequest{
		EventType: eventTypeName,
		UserID:    userID,
		Payload: map[string]interface{}{
			"issue_id":  issueID,
			"priority":  priority,
			"component": component,
		},
	}

	// Log the event request for debugging
	eventJSON, _ := json.MarshalIndent(eventReq, "", "  ")
	log("Submitting issue fix event: %s", string(eventJSON))

	response := utils.MakeRequest(http.MethodPost, "/api/v1/events", eventReq)
	if response.StatusCode != http.StatusOK {
		log("Failed to submit event: %s", string(response.Body))
		t.Fatalf("Failed to submit event: %s", string(response.Body))
	}
	log("Event response: %s", string(response.Body))
}

// Helper function to get user badges
func getUserBadges(t *testing.T, log func(string, ...interface{}), userID string) []utils.UserBadge {
	response := utils.MakeRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s/badges", userID), nil)
	log("Get user badges response code: %d", response.StatusCode)
	log("Get user badges response body: %s", string(response.Body))

	if response.StatusCode != http.StatusOK {
		log("Failed to get user badges: %s", string(response.Body))
		t.Fatalf("Failed to get user badges: %s", string(response.Body))
	}

	var userBadges []utils.UserBadge
	err := utils.ParseResponse(response, &userBadges)
	if err != nil {
		log("Failed to parse user badges: %v", err)
		t.Fatalf("Failed to parse user badges: %v", err)
	}

	badgeDetails, _ := json.MarshalIndent(userBadges, "", "  ")
	log("Retrieved user badges: %s", string(badgeDetails))

	return userBadges
}

// Helper function to check if a specific badge is in a list of badges
func hasBadge(badges []utils.UserBadge, badgeID int) bool {
	for _, badge := range badges {
		if badge.ID == badgeID {
			return true
		}
	}
	return false
}
