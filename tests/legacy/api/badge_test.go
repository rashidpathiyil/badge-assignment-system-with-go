package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/badge-assignment-system/tests/api/utils"
	"github.com/stretchr/testify/assert"
)

// TestCreateBadge tests the creation of a badge with criteria
func TestCreateBadge(t *testing.T) {
	if eventTypeID == 0 || conditionTypeID == 0 {
		t.Skip("Event type ID or condition type ID not set, skipping test")
	}

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	badgeName := fmt.Sprintf("Master Challenger_%d", timestamp)

	// Define the flow definition for the badge criteria
	flowDefinition := map[string]interface{}{
		"steps": []map[string]interface{}{
			{
				"type":            "condition",
				"conditionTypeId": conditionTypeID,
				"eventTypeIds":    []int{eventTypeID},
				"conditionParams": map[string]interface{}{
					"threshold": 90,
				},
			},
		},
	}

	// Create the badge request
	badgeReq := utils.BadgeRequest{
		Name:           badgeName,
		Description:    "Awarded for completing challenges with a high score",
		ImageURL:       "https://example.com/badges/master_challenger.png",
		FlowDefinition: flowDefinition,
	}

	// Make the API request
	response := utils.MakeRequest(http.MethodPost, "/api/v1/admin/badges", badgeReq)
	utils.AssertSuccess(t, response)

	// Parse the response
	var badgeWithCriteria utils.BadgeWithCriteria
	err := utils.ParseResponse(response, &badgeWithCriteria)
	assert.NoError(t, err)

	// Validate response data
	assert.NotZero(t, badgeWithCriteria.Badge.ID)
	assert.Equal(t, badgeReq.Name, badgeWithCriteria.Badge.Name)
	assert.Equal(t, badgeReq.Description, badgeWithCriteria.Badge.Description)
	assert.Equal(t, badgeReq.ImageURL, badgeWithCriteria.Badge.ImageURL)
	assert.NotZero(t, badgeWithCriteria.Criteria.ID)

	// Store the badge ID for later use
	badgeID = badgeWithCriteria.Badge.ID
	t.Logf("Created badge with ID: %d", badgeID)
}

// TestGetBadge tests retrieving a badge
func TestGetBadge(t *testing.T) {
	if badgeID == 0 {
		t.Skip("Badge ID not set, skipping test")
	}

	// Make the API request
	response := utils.MakeRequest(http.MethodGet, fmt.Sprintf("/api/v1/badges/%d", badgeID), nil)
	utils.AssertSuccess(t, response)

	// Parse the response
	var badge utils.Badge
	err := utils.ParseResponse(response, &badge)
	assert.NoError(t, err)

	// Validate response data
	assert.Equal(t, badgeID, badge.ID)
	assert.NotEmpty(t, badge.Name)
}

// TestGetBadgeWithCriteria tests retrieving a badge with criteria
func TestGetBadgeWithCriteria(t *testing.T) {
	if badgeID == 0 {
		t.Skip("Badge ID not set, skipping test")
	}

	// Make the API request
	response := utils.MakeRequest(http.MethodGet, fmt.Sprintf("/api/v1/admin/badges/%d/criteria", badgeID), nil)
	utils.AssertSuccess(t, response)

	// Parse the response
	var badgeWithCriteria utils.BadgeWithCriteria
	err := utils.ParseResponse(response, &badgeWithCriteria)
	assert.NoError(t, err)

	// Validate response data
	assert.Equal(t, badgeID, badgeWithCriteria.Badge.ID)
	assert.NotEmpty(t, badgeWithCriteria.Badge.Name)
	assert.NotZero(t, badgeWithCriteria.Criteria.ID)
}
