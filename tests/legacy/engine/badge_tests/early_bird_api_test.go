package badge_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// TestEarlyBirdBadgeAPI tests the Early Bird badge functionality through the API
func TestEarlyBirdBadgeAPI(t *testing.T) {
	// Skip test if in short mode
	if testing.Short() {
		t.Skip("Skipping API test in short mode")
	}

	// Configuration
	apiURL := "http://localhost:8080"
	testUserID1 := "test-user-api-123"
	testUserID2 := "test-user-api-456"
	eventType := "Check In"

	// Check if API server is running
	if !isServerRunning(apiURL) {
		t.Skip("API server is not running. Please start the server using 'go run cmd/server/main.go'")
	}

	// Test case 1: User checks in before 9 AM for 5 consecutive days (should earn badge)
	t.Run("EarlyCheckInsForFiveDaysAPI", func(t *testing.T) {
		// Get today's date
		today := time.Now()

		// Create 5 check-ins before 9 AM for 5 consecutive days
		for i := 0; i < 5; i++ {
			// Calculate date (starting from today, going backward)
			checkInDay := today.AddDate(0, 0, -i)
			checkInDate := checkInDay.Format("2006-01-02")

			// Submit check-in at 8:00 AM
			err := submitCheckIn(apiURL, eventType, testUserID1, checkInDate, "08:00:00")
			if err != nil {
				t.Fatalf("Failed to submit check-in: %v", err)
			}
		}

		// Wait a bit for badge processing
		time.Sleep(2 * time.Second)

		// Check if user has received the badge
		badges, err := getUserBadges(apiURL, testUserID1)
		if err != nil {
			t.Fatalf("Failed to get user badges: %v", err)
		}

		// Check if Early Bird badge is in the list
		if !checkForBadge(badges, "Early Bird") {
			t.Errorf("Expected user to have the Early Bird badge, but it wasn't awarded. Badges: %v", badges)
		} else {
			t.Logf("User successfully awarded the Early Bird badge")
		}
	})

	// Test case 2: User with mixed check-in times (should not earn badge)
	t.Run("MixedCheckInTimesAPI", func(t *testing.T) {
		// Get today's date
		today := time.Now()

		// Create 3 check-ins before 9 AM
		for i := 0; i < 3; i++ {
			checkInDay := today.AddDate(0, 0, -i)
			checkInDate := checkInDay.Format("2006-01-02")

			err := submitCheckIn(apiURL, eventType, testUserID2, checkInDate, "08:00:00")
			if err != nil {
				t.Fatalf("Failed to submit check-in: %v", err)
			}
		}

		// Create 2 check-ins after 9 AM
		for i := 3; i < 5; i++ {
			checkInDay := today.AddDate(0, 0, -i)
			checkInDate := checkInDay.Format("2006-01-02")

			err := submitCheckIn(apiURL, eventType, testUserID2, checkInDate, "10:00:00")
			if err != nil {
				t.Fatalf("Failed to submit check-in: %v", err)
			}
		}

		// Wait a bit for badge processing
		time.Sleep(2 * time.Second)

		// Check if user has received the badge
		badges, err := getUserBadges(apiURL, testUserID2)
		if err != nil {
			t.Fatalf("Failed to get user badges: %v", err)
		}

		// User should NOT have the badge in this case
		if checkForBadge(badges, "Early Bird") {
			t.Errorf("Expected user NOT to have the Early Bird badge, but it was awarded. Badges: %v", badges)
		} else {
			t.Logf("User correctly NOT awarded the Early Bird badge")
		}
	})
}

// Helper function to check if the API server is running
func isServerRunning(apiURL string) bool {
	resp, err := http.Get(apiURL + "/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close()
	return true
}

// Helper function to submit a check-in event
func submitCheckIn(apiURL, eventType, userID, date, time string) error {
	// Create request payload
	payload := map[string]interface{}{
		"event_type": eventType,
		"user_id":    userID,
		"payload": map[string]interface{}{
			"user_id":  userID,
			"time":     time,
			"date":     date,
			"location": "API Test Location",
		},
		"timestamp": fmt.Sprintf("%sT%sZ", date, time),
	}

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Send request
	resp, err := http.Post(
		apiURL+"/api/v1/events",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
	}

	return nil
}

// Helper function to get user badges
func getUserBadges(apiURL, userID string) ([]map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/users/%s/badges", apiURL, userID))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
	}

	var badges []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&badges); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return badges, nil
}

// Helper function to check if a specific badge is in the list
func checkForBadge(badges []map[string]interface{}, badgeName string) bool {
	for _, badge := range badges {
		badgeObj, ok := badge["badge"].(map[string]interface{})
		if !ok {
			continue
		}

		name, ok := badgeObj["name"].(string)
		if !ok {
			continue
		}

		if name == badgeName || name == badgeName+" - Test 1" || name == badgeName+" - Test 2" {
			return true
		}
	}
	return false
}
