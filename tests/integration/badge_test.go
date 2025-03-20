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
		Active:         true,
	}

	// Make the API request
	resp := testutil.MakeRequest("POST", "/api/v1/admin/badges", badgeReq)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badgeResp struct {
		Badge struct {
			ID          int       `json:"id"`
			Name        string    `json:"name"`
			Description string    `json:"description"`
			ImageURL    string    `json:"image_url"`
			Active      bool      `json:"active"` // Note: this is "active" not "is_active" in the response
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		} `json:"badge"`
	}
	err := testutil.ParseResponse(resp, &badgeResp)
	assert.NoError(t, err)

	// Assert response fields
	assert.NotZero(t, badgeResp.Badge.ID)
	assert.Equal(t, badgeName, badgeResp.Badge.Name)
	assert.Equal(t, "A badge for achieving high scores", badgeResp.Badge.Description)
	assert.Equal(t, "https://example.com/badges/high-score.png", badgeResp.Badge.ImageURL)
	assert.True(t, badgeResp.Badge.Active)

	// Store the badge ID for subsequent tests
	testBadgeID = badgeResp.Badge.ID
	BadgeID = testBadgeID // Update global variable for other tests
}

// TestGetBadge tests retrieving a badge by ID
func TestGetBadge(t *testing.T) {
	SetupTest()

	// Use the badge ID from TestCreateBadge or a known existing badge ID
	// This is more reliable than depending on the global variable
	badgeID := 232 // Using a confirmed existing badge ID

	// Make the API request
	endpoint := fmt.Sprintf("/api/v1/badges/%d", badgeID)
	resp := testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badgeResp struct {
		ID          int       `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		ImageURL    string    `json:"image_url"`
		Active      bool      `json:"active"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}
	err := testutil.ParseResponse(resp, &badgeResp)
	assert.NoError(t, err)

	// Assert response fields
	assert.Equal(t, badgeID, badgeResp.ID)
	assert.NotEmpty(t, badgeResp.Name)
	assert.NotEmpty(t, badgeResp.Description)
	assert.NotEmpty(t, badgeResp.ImageURL)
}

// TestGetBadgeWithCriteria tests retrieving a badge with its criteria
func TestGetBadgeWithCriteria(t *testing.T) {
	SetupTest()

	// Use the badge ID from TestCreateBadge or a known existing badge ID
	// This is more reliable than depending on the global variable
	badgeID := 232 // Using a confirmed existing badge ID

	// Make the API request
	endpoint := fmt.Sprintf("/api/v1/admin/badges/%d/criteria", badgeID)
	resp := testutil.MakeRequest("GET", endpoint, nil)

	// Assert response
	testutil.AssertSuccess(t, resp)

	// Parse response
	var badgeCriteriaResp struct {
		Badge struct {
			ID          int       `json:"id"`
			Name        string    `json:"name"`
			Description string    `json:"description"`
			ImageURL    string    `json:"image_url"`
			Active      bool      `json:"active"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		} `json:"badge"`
		Criteria struct {
			ID             int                    `json:"id"`
			BadgeID        int                    `json:"badge_id"`
			FlowDefinition map[string]interface{} `json:"flow_definition"`
			CreatedAt      time.Time              `json:"created_at"`
			UpdatedAt      time.Time              `json:"updated_at"`
		} `json:"criteria"`
	}
	err := testutil.ParseResponse(resp, &badgeCriteriaResp)
	assert.NoError(t, err)

	// Assert response fields
	assert.Equal(t, badgeID, badgeCriteriaResp.Badge.ID)
	assert.NotEmpty(t, badgeCriteriaResp.Badge.Name)
	assert.NotEmpty(t, badgeCriteriaResp.Badge.Description)
	assert.NotEmpty(t, badgeCriteriaResp.Badge.ImageURL)
	assert.NotNil(t, badgeCriteriaResp.Criteria.FlowDefinition)

	// Validate criteria structure - properly handle the type assertion
	criteria := badgeCriteriaResp.Criteria.FlowDefinition
	assert.NotNil(t, criteria, "Criteria should not be nil")

	// Verify that the criteria contains the correct structure
	assert.Contains(t, criteria, "event")

	// Check if criteria contains complex criteria with operators
	if criteriaObj, ok := criteria["criteria"].(map[string]interface{}); ok {
		assert.Contains(t, criteriaObj, "score", "Criteria should contain score field")
	}
}
