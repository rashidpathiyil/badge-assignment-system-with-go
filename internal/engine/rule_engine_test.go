package engine

import (
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/models"
	"github.com/badge-assignment-system/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// We need to adapt our test to match the actual RuleEngine implementation

// TestEvaluateBadgeCriteria tests the EvaluateBadgeCriteria method
func TestEvaluateBadgeCriteria(t *testing.T) {
	// Create a mock DB
	mockDB := testutil.NewMockDB()

	// Create a test badge with criteria
	badgeID := 1
	// Update the flow definition to match the format expected by the rule engine
	// The engine looks for event+criteria structure
	testCriteria := map[string]interface{}{
		"event": "test_event",
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": float64(90),
			},
		},
	}

	mockBadge := testutil.CreateTestBadgeWithCriteria(badgeID, "Test Badge", testCriteria)

	// Configure the mock to return our badge
	mockDB.On("GetBadgeWithCriteria", badgeID).Return(mockBadge, nil)

	// Mock GetEventTypeByName as the flow evaluation will need it
	mockEventType := models.EventType{
		ID:   1,
		Name: "test_event",
	}
	mockDB.On("GetEventTypeByName", "test_event").Return(mockEventType, nil)

	// Mock GetUserEventsByType to return empty events
	mockDB.On("GetUserEventsByType", "test-user", mockEventType.ID).Return([]models.Event{}, nil)

	// Create an instance of the rule engine with the mock
	engine := NewRuleEngine(mockDB)

	// Call the method under test
	result, metadata, err := engine.EvaluateBadgeCriteria(badgeID, "test-user")

	// Verify the result
	assert.NoError(t, err)
	// The expected result is false since there are no events that match the criteria
	assert.False(t, result)
	assert.NotNil(t, metadata)

	// Verify that the mock was called as expected
	mockDB.AssertExpectations(t)
}

// TestProcessEvents tests the ProcessEvents method
func TestProcessEvents(t *testing.T) {
	// Create a mock DB
	mockDB := testutil.NewMockDB()

	// Create test badges
	badge1 := models.Badge{
		ID:          1,
		Name:        "Test Badge 1",
		Description: "Description for test badge 1",
		ImageURL:    "https://example.com/badge1.png",
		Active:      true,
	}
	badge2 := models.Badge{
		ID:          2,
		Name:        "Test Badge 2",
		Description: "Description for test badge 2",
		ImageURL:    "https://example.com/badge2.png",
		Active:      true,
	}
	badges := []models.Badge{badge1, badge2}

	// Mock DB to return active badges
	mockDB.On("GetActiveBadges").Return(badges, nil)

	// Mock DB to return no existing badges for the user
	mockDB.On("GetUserBadges", "test-user").Return([]models.UserBadge{}, nil)

	// Create test badge with criteria
	badgeID1 := 1
	testCriteria1 := map[string]interface{}{
		"event": "test_event",
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": float64(90),
			},
		},
	}
	mockBadge1 := testutil.CreateTestBadgeWithCriteria(badgeID1, "Test Badge 1", testCriteria1)

	// Mock GetBadgeWithCriteria for the first badge
	mockDB.On("GetBadgeWithCriteria", badgeID1).Return(mockBadge1, nil)

	// Second badge
	badgeID2 := 2
	testCriteria2 := map[string]interface{}{
		"event": "test_event_2",
		"criteria": map[string]interface{}{
			"count": map[string]interface{}{
				"$gte": 5,
			},
		},
	}
	mockBadge2 := testutil.CreateTestBadgeWithCriteria(badgeID2, "Test Badge 2", testCriteria2)

	// Mock GetBadgeWithCriteria for the second badge
	mockDB.On("GetBadgeWithCriteria", badgeID2).Return(mockBadge2, nil)

	// For badge 1, set up mocks for the event type and user events
	mockEventType1 := models.EventType{
		ID:   1,
		Name: "test_event",
	}
	mockDB.On("GetEventTypeByName", "test_event").Return(mockEventType1, nil)

	// Create a test event with a score of 95
	testEvent := models.Event{
		ID:          1,
		EventTypeID: 1,
		UserID:      "test-user",
		OccurredAt:  time.Now(),
		Payload: models.JSONB{
			"score": 95,
		},
	}

	// Mock GetUserEventsByType to return our test event
	mockDB.On("GetUserEventsByType", "test-user", mockEventType1.ID).Return([]models.Event{testEvent}, nil)

	// For badge 2, set up mocks for the event type
	mockEventType2 := models.EventType{
		ID:   2,
		Name: "test_event_2",
	}
	mockDB.On("GetEventTypeByName", "test_event_2").Return(mockEventType2, nil)

	// Mock GetUserEventsByType to return empty events for the second badge (criteria not met)
	mockDB.On("GetUserEventsByType", "test-user", mockEventType2.ID).Return([]models.Event{}, nil)

	// Mock AwardBadgeToUser for the first badge which will meet the criteria
	mockDB.On("AwardBadgeToUser", mock.MatchedBy(func(badge *models.UserBadge) bool {
		return badge.BadgeID == 1 && badge.UserID == "test-user"
	})).Return(nil)

	// Create an instance of the rule engine with the mock
	engine := NewRuleEngine(mockDB)

	// Call the method under test
	err := engine.ProcessEvents("test-user")

	// Verify the result
	assert.NoError(t, err)

	// Verify that the mock was called as expected
	mockDB.AssertExpectations(t)
}

// Here's a demonstration of table-driven tests for a hypothetical method
func TestHypotheticalEvaluationMethod(t *testing.T) {
	t.Skip("This is a placeholder test demonstrating test patterns")

	// Commenting out engine creation since it's not used and causes a linter warning
	// engine := NewRuleEngine(&models.DB{})

	testCases := []struct {
		name     string
		criteria map[string]interface{}
		expected bool
	}{
		{
			name: "Simple threshold criteria",
			criteria: map[string]interface{}{
				"type": "threshold",
				"parameters": map[string]interface{}{
					"field":     "score",
					"threshold": float64(90),
				},
			},
			expected: true,
		},
		{
			name: "Complex criteria",
			criteria: map[string]interface{}{
				"type":     "composite",
				"operator": "and",
				"rules": []map[string]interface{}{
					{
						"type": "threshold",
						"parameters": map[string]interface{}{
							"field":     "score",
							"threshold": float64(90),
						},
					},
					{
						"type": "count",
						"parameters": map[string]interface{}{
							"minimum": float64(5),
						},
					},
				},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This would be a call to the actual method under test
			// result := engine.SomeEvaluationMethod(tc.criteria)
			// assert.Equal(t, tc.expected, result)
		})
	}
}

// This is a utility function to demonstrate a more complex test pattern
func createTestEngine() (*RuleEngine, *testutil.MockDB) {
	mockDB := testutil.NewMockDB()
	engine := NewRuleEngine(mockDB)
	return engine, mockDB
}

// BenchmarkRuleEvaluation benchmarks the rule evaluation performance
func BenchmarkRuleEvaluation(b *testing.B) {
	// Create a mock DB
	mockDB := testutil.NewMockDB()

	// Create a test badge with criteria
	badgeID := 1
	testCriteria := map[string]interface{}{
		"event": "test_event",
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": float64(90),
			},
		},
	}
	mockBadge := testutil.CreateTestBadgeWithCriteria(badgeID, "Benchmark Badge", testCriteria)

	// Configure the mock to return our badge
	mockDB.On("GetBadgeWithCriteria", badgeID).Return(mockBadge, nil)

	// Mock GetEventTypeByName
	mockEventType := models.EventType{
		ID:   1,
		Name: "test_event",
	}
	mockDB.On("GetEventTypeByName", "test_event").Return(mockEventType, nil)

	// Create test events
	testEvents := []models.Event{
		{
			ID:          1,
			EventTypeID: 1,
			UserID:      "bench-user",
			OccurredAt:  time.Now(),
			Payload: models.JSONB{
				"score": 95,
			},
		},
	}

	// Mock GetUserEventsByType to return our test events
	mockDB.On("GetUserEventsByType", "bench-user", mockEventType.ID).Return(testEvents, nil)

	// Create an instance of the rule engine with the mock
	engine := NewRuleEngine(mockDB)

	// Reset the timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		result, _, _ := engine.EvaluateBadgeCriteria(badgeID, "bench-user")
		// Ensure the result is used to prevent compiler optimizations
		if !result {
			b.Fail()
		}
	}
}
