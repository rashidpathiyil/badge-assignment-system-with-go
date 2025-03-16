package tests

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/engine"
	"github.com/badge-assignment-system/internal/logging"
	"github.com/badge-assignment-system/internal/models"
	"github.com/jmoiron/sqlx"
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
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanupTestDB(db)

	// Load badge definition
	badgeDef, err := loadBadgeDefinition("../../../badges/early-bird.json")
	if err != nil {
		t.Fatalf("Failed to load badge definition: %v", err)
	}

	// Load event type definition
	eventTypeDef, err := loadEventTypeDefinition("../../../cmd/badgecli/badges/event_types/check-in.json")
	if err != nil {
		t.Fatalf("Failed to load event type definition: %v", err)
	}

	// Create models.DB wrapper
	dbWrapper := &models.DB{DB: db}

	// Create event type
	eventType := &models.EventType{
		Name:        eventTypeDef.Name,
		Description: eventTypeDef.Description,
		Schema:      models.JSONB(eventTypeDef.Schema),
	}
	if err := dbWrapper.CreateEventType(eventType); err != nil {
		t.Fatalf("Failed to create event type: %v", err)
	}
	t.Logf("Created event type: %s (ID: %d)", eventType.Name, eventType.ID)

	// Create badge
	badge := &models.Badge{
		Name:        badgeDef.Name,
		Description: badgeDef.Description,
		ImageURL:    badgeDef.ImageURL,
		Active:      true,
	}
	criteria := &models.BadgeCriteria{
		FlowDefinition: models.JSONB(badgeDef.FlowDefinition),
	}
	if err := dbWrapper.CreateBadge(badge, criteria); err != nil {
		t.Fatalf("Failed to create badge: %v", err)
	}
	t.Logf("Created badge: %s (ID: %d)", badge.Name, badge.ID)

	// Create a test user
	testUserID := "test-user-123"

	// Create rule engine
	ruleEngine := engine.NewRuleEngine(dbWrapper)
	ruleEngine.SetLogLevel(logging.LogLevelDebug)

	// Test Case 1: User checks in before 9 AM for 5 consecutive days (should earn badge)
	t.Run("EarlyCheckInsForFiveDays", func(t *testing.T) {
		// Create check-in events: 5 days, all before 9 AM
		baseTime := time.Now().Truncate(24 * time.Hour)
		for i := 0; i < 5; i++ {
			checkInTime := baseTime.AddDate(0, 0, -i).Add(8 * time.Hour) // 8:00 AM each day for past 5 days
			createCheckInEvent(t, dbWrapper, eventType.ID, testUserID, checkInTime)
		}

		// Evaluate badge criteria
		awarded, metadata, err := ruleEngine.EvaluateBadgeCriteria(badge.ID, testUserID)
		if err != nil {
			t.Fatalf("Failed to evaluate badge criteria: %v", err)
		}

		// Should be awarded
		if !awarded {
			t.Errorf("Expected badge to be awarded, but it wasn't. Metadata: %v", metadata)
		} else {
			t.Logf("Badge correctly awarded! Metadata: %v", metadata)
		}
	})

	// Clean up events from first test
	cleanupUserEvents(dbWrapper, testUserID)

	// Test Case 2: User checks in but some days are after 9 AM (should not earn badge)
	t.Run("MixedCheckInTimes", func(t *testing.T) {
		// Create check-in events: 5 days, but 2 days after 9 AM
		baseTime := time.Now().Truncate(24 * time.Hour)

		// 3 days before 9 AM
		for i := 0; i < 3; i++ {
			checkInTime := baseTime.AddDate(0, 0, -i).Add(8 * time.Hour) // 8:00 AM
			createCheckInEvent(t, dbWrapper, eventType.ID, testUserID, checkInTime)
		}

		// 2 days after 9 AM
		for i := 3; i < 5; i++ {
			checkInTime := baseTime.AddDate(0, 0, -i).Add(10 * time.Hour) // 10:00 AM
			createCheckInEvent(t, dbWrapper, eventType.ID, testUserID, checkInTime)
		}

		// Evaluate badge criteria
		awarded, metadata, err := ruleEngine.EvaluateBadgeCriteria(badge.ID, testUserID)
		if err != nil {
			t.Fatalf("Failed to evaluate badge criteria: %v", err)
		}

		// Should not be awarded (only 3 days before 9 AM)
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

// Helper function to set up test database
func setupTestDB() (*sqlx.DB, error) {
	// Get DB connection from environment or use default test connection
	dbConnStr := os.Getenv("TEST_DB_CONNECTION")
	if dbConnStr == "" {
		dbConnStr = "postgres://postgres:postgres@localhost:5432/badge_system_test?sslmode=disable"
	}

	// Connect to DB
	db, err := sqlx.Connect("postgres", dbConnStr)
	if err != nil {
		return nil, err
	}

	return db, nil
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
