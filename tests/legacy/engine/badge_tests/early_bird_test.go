package badge_tests

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/engine"
	"github.com/badge-assignment-system/internal/logging"
	"github.com/badge-assignment-system/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Define test badge and event structures to match the JSON files
type BadgeDefinition struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	ImageURL       string                 `json:"image_url"`
	FlowDefinition map[string]interface{} `json:"flow_definition"`
}

type EventTypeDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
}

// TestEarlyBirdBadge is a real-world integration test for the Early Bird badge
// This test requires a PostgreSQL database to be available
func TestEarlyBirdBadge(t *testing.T) {
	// Skip this test if in short mode or if no DB connection is available
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create DB connection for test
	db, err := setupTestDB()
	if err != nil {
		t.Skipf("Skipping test due to database connection error: %v", err)
		return
	}
	defer cleanupTestDB(db)

	// Create a test user
	testUserID := "test-user-123"

	// Clean up any existing events for this user
	dbWrapper := &models.DB{DB: db}
	cleanupUserEvents(dbWrapper, testUserID)

	// Load badge definition
	badgeDef, err := loadBadgeDefinition("../../../../badges/early-bird.json")
	if err != nil {
		t.Fatalf("Failed to load badge definition: %v", err)
	}

	// Load event type definition
	eventTypeDef, err := loadEventTypeDefinition("../../../../cmd/badgecli/badges/event_types/check-in.json")
	if err != nil {
		t.Fatalf("Failed to load event type definition: %v", err)
	}

	// Instead of creating a new event type, use the existing one
	eventType, err := dbWrapper.GetEventTypeByName("Check In")
	if err != nil {
		// Only create a new event type if it doesn't exist
		newEventType := &models.EventType{
			Name:        eventTypeDef.Name,
			Description: eventTypeDef.Description,
			Schema:      models.JSONB(eventTypeDef.Schema),
		}
		if err := dbWrapper.CreateEventType(newEventType); err != nil {
			t.Fatalf("Failed to create event type: %v", err)
		}
		eventType = *newEventType
		t.Logf("Created event type: %s (ID: %d)", eventType.Name, eventType.ID)
	} else {
		t.Logf("Using existing event type: %s (ID: %d)", eventType.Name, eventType.ID)
	}

	// Create rule engine
	ruleEngine := engine.NewRuleEngine(dbWrapper)
	ruleEngine.SetLogLevel(logging.LogLevelDebug)

	// Print event type info for debugging
	t.Logf("Event type in test: ID=%d, Name='%s'", eventType.ID, eventType.Name)

	// Verify event type can be retrieved correctly
	verifyEventType, err := dbWrapper.GetEventTypeByName(eventType.Name)
	if err != nil {
		t.Fatalf("Failed to retrieve event type by name: %v", err)
	}
	t.Logf("Retrieved event type: ID=%d, Name='%s'", verifyEventType.ID, verifyEventType.Name)

	// Test Case 1: User checks in before 9 AM for 5 consecutive days (should earn badge)
	t.Run("EarlyCheckInsForFiveDays", func(t *testing.T) {
		// Create a badge for this test case
		badge := &models.Badge{
			Name:        badgeDef.Name + " - Test 1",
			Description: badgeDef.Description,
			ImageURL:    badgeDef.ImageURL,
			Active:      true,
		}

		// Create a simplified flow definition that just checks the time field
		flowDefCopy := make(map[string]interface{})

		// Use the exact event type name from the database
		// This ensures we use "Check In" instead of "check-in"
		flowDefCopy["event"] = eventType.Name

		flowDefCopy["criteria"] = map[string]interface{}{
			"time": map[string]interface{}{
				"$lt": "09:00:00",
			},
		}

		// Print the final flow definition for debugging
		flowJSON, _ := json.Marshal(flowDefCopy)
		t.Logf("Badge flow definition: %s", string(flowJSON))

		criteria := &models.BadgeCriteria{
			FlowDefinition: models.JSONB(flowDefCopy),
		}

		if err := dbWrapper.CreateBadge(badge, criteria); err != nil {
			t.Fatalf("Failed to create badge: %v", err)
		}
		t.Logf("Created badge: %s (ID: %d)", badge.Name, badge.ID)

		// Clear previous events
		cleanupUserEvents(dbWrapper, testUserID)

		// Create check-in events: 5 days, all before 9 AM
		baseTime := time.Now().Truncate(24 * time.Hour) // Start of today

		for i := 0; i < 5; i++ {
			// Create a time at 8:00 AM for each day, going back 5 days
			checkInDay := baseTime.AddDate(0, 0, -i) // go back i days
			checkInTime := time.Date(
				checkInDay.Year(),
				checkInDay.Month(),
				checkInDay.Day(),
				8, // 8 hours = 8:00 AM
				0,
				0,
				0,
				checkInDay.Location(),
			)

			// Verify the time is what we expect
			t.Logf("Preparing check-in at %s (%s)",
				checkInTime.Format("2006-01-02 15:04:05"),
				checkInTime.Format("15:04:05"))

			// Create event
			event := &models.Event{
				EventTypeID: eventType.ID,
				UserID:      testUserID,
				OccurredAt:  checkInTime,
				Payload: models.JSONB{
					"user_id":  testUserID,
					"time":     checkInTime.Format("15:04:05"),
					"date":     checkInTime.Format("2006-01-02"),
					"location": "Test Location",
				},
			}

			if err := dbWrapper.CreateEvent(event); err != nil {
				t.Fatalf("Failed to create event: %v", err)
			}

			t.Logf("Created check-in event at %s with ID %d (Time: %s)",
				checkInTime.Format("2006-01-02 15:04:05"),
				event.ID,
				checkInTime.Format("15:04:05"))
		}

		// Evaluate badge criteria
		t.Logf("Evaluating badge criteria for badge ID %d", badge.ID)

		// Get the badge criteria for debugging
		badgeWithCriteria, err := dbWrapper.GetBadgeWithCriteria(badge.ID)
		if err != nil {
			t.Fatalf("Failed to get badge with criteria: %v", err)
		}
		criteriaJSON, _ := json.Marshal(badgeWithCriteria.Criteria.FlowDefinition)
		t.Logf("Badge criteria to evaluate: %s", string(criteriaJSON))

		awarded, metadata, err := ruleEngine.EvaluateBadgeCriteria(badge.ID, testUserID)
		if err != nil {
			t.Fatalf("Failed to evaluate badge criteria: %v", err)
		}

		// Should be awarded since all events are before 9 AM
		if !awarded {
			t.Errorf("Expected badge to be awarded, but it wasn't. Metadata: %v", metadata)
		} else {
			t.Logf("Badge correctly awarded! Metadata: %v", metadata)
		}
	})

	// Test Case 2: User checks in but some days are after 9 AM (should not earn badge)
	t.Run("MixedCheckInTimes", func(t *testing.T) {
		// Create a badge for this test case with count criteria
		badge := &models.Badge{
			Name:        badgeDef.Name + " - Test 2",
			Description: badgeDef.Description,
			ImageURL:    badgeDef.ImageURL,
			Active:      true,
		}

		// Create a flow definition that requires 5 early check-ins
		flowDefCopy := make(map[string]interface{})

		// Use the exact event type name from the database
		// This ensures we use "Check In" instead of "check-in"
		flowDefCopy["event"] = eventType.Name

		// To enforce applying the time filter first before counting,
		// we'll create a nested $and structure with a count embedded
		flowDefCopy["criteria"] = map[string]interface{}{
			"$and": []interface{}{
				// First filter by time
				map[string]interface{}{
					"time": map[string]interface{}{
						"$lt": "09:00:00",
					},
				},
				// Then apply count criteria to the filtered events
				map[string]interface{}{
					"count": map[string]interface{}{
						"$gte": 5,
					},
				},
			},
		}

		// Print the final flow definition for debugging
		flowJSON, _ := json.Marshal(flowDefCopy)
		t.Logf("Badge flow definition: %s", string(flowJSON))

		criteria := &models.BadgeCriteria{
			FlowDefinition: models.JSONB(flowDefCopy),
		}

		if err := dbWrapper.CreateBadge(badge, criteria); err != nil {
			t.Fatalf("Failed to create badge: %v", err)
		}
		t.Logf("Created badge: %s (ID: %d)", badge.Name, badge.ID)

		// Clear previous events
		cleanupUserEvents(dbWrapper, testUserID)

		// Create check-in events: 5 days, but 2 days after 9 AM
		baseTime := time.Now().Truncate(24 * time.Hour) // Start of today

		// 3 days before 9 AM
		for i := 0; i < 3; i++ {
			// Create a time at 8:00 AM for 3 days
			checkInDay := baseTime.AddDate(0, 0, -i) // go back i days
			checkInTime := time.Date(
				checkInDay.Year(),
				checkInDay.Month(),
				checkInDay.Day(),
				8, // 8 hours = 8:00 AM
				0,
				0,
				0,
				checkInDay.Location(),
			)

			// Create event
			event := &models.Event{
				EventTypeID: eventType.ID,
				UserID:      testUserID,
				OccurredAt:  checkInTime,
				Payload: models.JSONB{
					"user_id":  testUserID,
					"time":     checkInTime.Format("15:04:05"),
					"date":     checkInTime.Format("2006-01-02"),
					"location": "Test Location",
				},
			}

			if err := dbWrapper.CreateEvent(event); err != nil {
				t.Fatalf("Failed to create event: %v", err)
			}

			t.Logf("Created check-in event at %s with ID %d (Time: %s)",
				checkInTime.Format("2006-01-02 15:04:05"),
				event.ID,
				checkInTime.Format("15:04:05"))
		}

		// 2 days after 9 AM
		for i := 3; i < 5; i++ {
			// Create a time at 10:00 AM for 2 days
			checkInDay := baseTime.AddDate(0, 0, -i) // go back i days
			checkInTime := time.Date(
				checkInDay.Year(),
				checkInDay.Month(),
				checkInDay.Day(),
				10, // 10 hours = 10:00 AM
				0,
				0,
				0,
				checkInDay.Location(),
			)

			// Create event
			event := &models.Event{
				EventTypeID: eventType.ID,
				UserID:      testUserID,
				OccurredAt:  checkInTime,
				Payload: models.JSONB{
					"user_id":  testUserID,
					"time":     checkInTime.Format("15:04:05"),
					"date":     checkInTime.Format("2006-01-02"),
					"location": "Test Location",
				},
			}

			if err := dbWrapper.CreateEvent(event); err != nil {
				t.Fatalf("Failed to create event: %v", err)
			}

			t.Logf("Created check-in event at %s with ID %d (Time: %s)",
				checkInTime.Format("2006-01-02 15:04:05"),
				event.ID,
				checkInTime.Format("15:04:05"))
		}

		// Evaluate badge criteria
		t.Logf("Evaluating badge criteria for badge ID %d", badge.ID)

		// Get the badge criteria for debugging
		badgeWithCriteria, err := dbWrapper.GetBadgeWithCriteria(badge.ID)
		if err != nil {
			t.Fatalf("Failed to get badge with criteria: %v", err)
		}
		criteriaJSON, _ := json.Marshal(badgeWithCriteria.Criteria.FlowDefinition)
		t.Logf("Badge criteria to evaluate: %s", string(criteriaJSON))

		awarded, metadata, err := ruleEngine.EvaluateBadgeCriteria(badge.ID, testUserID)
		if err != nil {
			t.Fatalf("Failed to evaluate badge criteria: %v", err)
		}

		// Should not be awarded since only 3 events are before 9 AM, but we need 5
		if awarded {
			t.Errorf("Expected badge not to be awarded, but it was. Metadata: %v", metadata)
		} else {
			t.Logf("Badge correctly not awarded. Metadata: %v", metadata)
		}
	})
}

// Helper function to create a check-in event
func createCheckInEvent(t *testing.T, db *models.DB, eventTypeID int, userID string, checkInTime time.Time) {
	event := &models.Event{
		EventTypeID: eventTypeID,
		UserID:      userID,
		OccurredAt:  checkInTime,
		Payload: models.JSONB{
			"user_id":  userID,
			"time":     checkInTime.Format("15:04:05"),
			"date":     checkInTime.Format("2006-01-02"),
			"location": "Test Location",
		},
	}

	if err := db.CreateEvent(event); err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}
	t.Logf("Created check-in event at %s", checkInTime.Format("2006-01-02 15:04:05"))
}

// Helper function to load badge definition from file
func loadBadgeDefinition(filePath string) (*BadgeDefinition, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var badge BadgeDefinition
	if err := json.Unmarshal(fileData, &badge); err != nil {
		return nil, err
	}

	return &badge, nil
}

// Helper function to load event type definition from file
func loadEventTypeDefinition(filePath string) (*EventTypeDefinition, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var eventType EventTypeDefinition
	if err := json.Unmarshal(fileData, &eventType); err != nil {
		return nil, err
	}

	return &eventType, nil
}

// setupTestDB creates a database connection for tests
func setupTestDB() (*sqlx.DB, error) {
	// Load environment variables from .env file if available
	err := godotenv.Load("../../../../.env")
	if err != nil {
		log.Println("Warning: .env file not found, using default database settings")
	}

	// Get database connection details from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "badge_system") // Use the same database as the application
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Open database connection
	return sqlx.Connect("postgres", connStr)
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to clean up test database
func cleanupTestDB(db *sqlx.DB) {
	// Could run cleanup SQL here if needed
	db.Close()
}

// Helper function to clean up events for a user
func cleanupUserEvents(db *models.DB, userID string) {
	// Delete all events for this user
	db.Exec("DELETE FROM events WHERE user_id = $1", userID)
}
