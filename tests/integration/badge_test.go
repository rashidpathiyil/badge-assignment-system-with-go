package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/testutil"
	"github.com/stretchr/testify/assert"
)

// Global variable to store the badge ID created during tests
var testBadgeID int

// TestCreateBadge tests the creation of a badge
func TestCreateBadge(t *testing.T) {
	SetupTest()

	if EventTypeID == 0 {
		t.Skip("Event type ID not set, skipping test")
	}

	// Create a unique name with timestamp to avoid conflicts
	timestamp := time.Now().UnixNano() / 1000000
	badgeName := fmt.Sprintf("Test Badge_%d", timestamp)

	// Define the criteria for the badge
	criteria := map[string]interface{}{
		"event": "test_event",
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": float64(90),
			},
		},
	}

	// Create the badge request
	badgeReq := BadgeRequest{
		Name:           badgeName,
		Description:    "A badge for achieving high scores",
		ImageURL:       "https://example.com/badges/high-score.png",
		FlowDefinition: criteria,
		IsActive:       true,
	}

	// Make the API request
	resp := testutil.MakeRequest("POST", "/api/v1/admin/badges", badgeReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badgeResp BadgeResponse
	err := testutil.ParseResponse(resp, &badgeResp)
	assert.NoError(t, err)

	// Assert response fields
	assert.NotZero(t, badgeResp.ID)
	assert.Equal(t, badgeName, badgeResp.Name)
	assert.Equal(t, "A badge for achieving high scores", badgeResp.Description)
	assert.Equal(t, "https://example.com/badges/high-score.png", badgeResp.ImageURL)
	assert.True(t, badgeResp.IsActive)

	// Store the badge ID for subsequent tests
	testBadgeID = badgeResp.ID
	BadgeID = testBadgeID // Update global variable for other tests
}

// TestGetBadge tests retrieving a badge by ID
func TestGetBadge(t *testing.T) {
	if testBadgeID == 0 {
		t.Skip("Badge ID not set, skipping test")
	}

	// Make the API request
	endpoint := fmt.Sprintf("/api/v1/badges/%d", testBadgeID)
	resp := testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badgeResp BadgeResponse
	err := testutil.ParseResponse(resp, &badgeResp)
	assert.NoError(t, err)

	// Assert response fields
	assert.Equal(t, testBadgeID, badgeResp.ID)
	assert.NotEmpty(t, badgeResp.Name)
	assert.NotEmpty(t, badgeResp.Description)
	assert.NotEmpty(t, badgeResp.ImageURL)
}

// TestGetBadgeWithCriteria tests retrieving a badge with its criteria
func TestGetBadgeWithCriteria(t *testing.T) {
	if testBadgeID == 0 {
		t.Skip("Badge ID not set, skipping test")
	}

	// Make the API request
	endpoint := fmt.Sprintf("/api/v1/admin/badges/%d/criteria", testBadgeID)
	resp := testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badgeCriteriaResp BadgeCriteriaResponse
	err := testutil.ParseResponse(resp, &badgeCriteriaResp)
	assert.NoError(t, err)

	// Assert response fields
	assert.Equal(t, testBadgeID, badgeCriteriaResp.ID)
	assert.NotEmpty(t, badgeCriteriaResp.Name)
	assert.NotEmpty(t, badgeCriteriaResp.Description)
	assert.NotEmpty(t, badgeCriteriaResp.ImageURL)
	assert.NotNil(t, badgeCriteriaResp.Criteria)

	// Validate criteria structure - properly handle the type assertion
	criteria := badgeCriteriaResp.Criteria
	assert.NotNil(t, criteria, "Criteria should not be nil")

	// Verify that the criteria contains the correct structure
	assert.Contains(t, criteria, "event")

	// Check if criteria contains complex criteria with operators
	if criteriaObj, ok := criteria["criteria"].(map[string]interface{}); ok {
		assert.Contains(t, criteriaObj, "score", "Criteria should contain score field")
	}
}
