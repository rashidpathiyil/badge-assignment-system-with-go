package engine

import (
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/logging"
	"github.com/badge-assignment-system/internal/models"
)

// Helper function for creating test events
func createTestEvents(count int, startTime time.Time, gap time.Duration) []models.Event {
	events := make([]models.Event, count)
	currentTime := startTime

	for i := 0; i < count; i++ {
		events[i] = models.Event{
			ID:          i + 1,
			UserID:      "test-user",
			EventTypeID: 1,
			OccurredAt:  currentTime,
			Payload:     models.JSONB{"value": float64(i * 10)},
		}
		currentTime = currentTime.Add(gap)
	}

	return events
}

// Test time period criteria evaluation
func TestTimePeriodCriteria(t *testing.T) {
	// Create a RuleEngine instance with a proper logger
	re := &RuleEngine{
		Logger: logging.NewLogger("TEST-ENGINE", logging.LogLevelError),
	}

	// Test case: Daily periods
	startTime := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	events := []models.Event{}

	for day := 0; day < 5; day++ {
		dayEvents := createTestEvents(2, startTime.AddDate(0, 0, day), time.Hour)
		events = append(events, dayEvents...)
	}

	criteria := map[string]interface{}{
		"periodType": "day",
		"periodCount": map[string]interface{}{
			"$gte": float64(3),
		},
	}

	metadata := make(map[string]interface{})
	result, err := re.evaluateTimePeriodCriteria(criteria, events, metadata)

	if err != nil {
		t.Errorf("Error evaluating time period criteria: %v", err)
	}

	if !result {
		t.Error("Expected criteria to be met, but it wasn't")
	}

	expectedPeriods := 5
	if actualPeriods, ok := metadata["unique_period_count"].(int); !ok || actualPeriods != expectedPeriods {
		t.Errorf("Expected %d unique periods, got %v", expectedPeriods, metadata["unique_period_count"])
	}
}

// Test pattern criteria evaluation
func TestPatternCriteria(t *testing.T) {
	// Create a RuleEngine instance with a proper logger
	re := &RuleEngine{
		Logger: logging.NewLogger("TEST-ENGINE", logging.LogLevelError),
	}

	// Test case: Consistent pattern
	startTime := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	events := []models.Event{}

	for day := 0; day < 5; day++ {
		// 3 events per day
		dayEvents := createTestEvents(3, startTime.AddDate(0, 0, day), time.Hour)
		events = append(events, dayEvents...)
	}

	criteria := map[string]interface{}{
		"$pattern": map[string]interface{}{
			"pattern":      "consistent",
			"periodType":   "day",
			"minPeriods":   float64(3),
			"maxDeviation": float64(0.1),
		},
	}

	// Need to extract the inner criteria for the evaluatePatternCriteria function
	// since it expects unwrapped criteria (the wrapper is handled at a higher level)
	innerCriteria := criteria["$pattern"].(map[string]interface{})

	metadata := make(map[string]interface{})
	result, err := re.evaluatePatternCriteria(innerCriteria, events, metadata)

	if err != nil {
		t.Errorf("Error evaluating pattern criteria: %v", err)
	}

	if !result {
		t.Error("Expected criteria to be met, but it wasn't")
	}
}

// Test gap criteria evaluation
func TestGapCriteria(t *testing.T) {
	// Create a RuleEngine instance with a proper logger
	re := &RuleEngine{
		Logger: logging.NewLogger("TEST-ENGINE", logging.LogLevelError),
	}

	// Create events with specific gaps
	baseTime := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	events := []models.Event{
		{ID: 1, UserID: "test-user", EventTypeID: 1, OccurredAt: baseTime, Payload: models.JSONB{}},
		{ID: 2, UserID: "test-user", EventTypeID: 1, OccurredAt: baseTime.Add(3 * time.Hour), Payload: models.JSONB{}},
		{ID: 3, UserID: "test-user", EventTypeID: 1, OccurredAt: baseTime.Add(7 * time.Hour), Payload: models.JSONB{}},
		{ID: 4, UserID: "test-user", EventTypeID: 1, OccurredAt: baseTime.Add(12 * time.Hour), Payload: models.JSONB{}},
	}

	criteria := map[string]interface{}{
		"minGapHours": float64(2),
		"maxGapHours": float64(6),
	}

	metadata := make(map[string]interface{})
	result, err := re.evaluateGapCriteria(criteria, events, metadata)

	if err != nil {
		t.Errorf("Error evaluating gap criteria: %v", err)
	}

	if !result {
		t.Error("Expected criteria to be met, but it wasn't")
	}
}

// Test duration criteria evaluation
func TestDurationCriteria(t *testing.T) {
	// Create a RuleEngine instance with a proper logger
	re := &RuleEngine{
		Logger: logging.NewLogger("TEST-ENGINE", logging.LogLevelError),
	}

	// Create start and end events
	baseTime := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	events := []models.Event{
		{ID: 1, UserID: "test-user", EventTypeID: 1, OccurredAt: baseTime, Payload: models.JSONB{"type": "session_start", "session_id": "123"}},
		{ID: 2, UserID: "test-user", EventTypeID: 2, OccurredAt: baseTime.Add(30 * time.Minute), Payload: models.JSONB{"type": "session_end", "session_id": "123"}},
		{ID: 3, UserID: "test-user", EventTypeID: 1, OccurredAt: baseTime.Add(2 * time.Hour), Payload: models.JSONB{"type": "session_start", "session_id": "456"}},
		{ID: 4, UserID: "test-user", EventTypeID: 2, OccurredAt: baseTime.Add(3 * time.Hour), Payload: models.JSONB{"type": "session_end", "session_id": "456"}},
	}

	criteria := map[string]interface{}{
		"startEvent": map[string]interface{}{
			"type": "session_start",
		},
		"endEvent": map[string]interface{}{
			"type": "session_end",
		},
		"matchProperty": "session_id",
		"duration": map[string]interface{}{
			"$gte": float64(30),
		},
		"unit": "minutes",
	}

	metadata := make(map[string]interface{})
	result, err := re.evaluateDurationCriteria(criteria, events, metadata)

	if err != nil {
		t.Errorf("Error evaluating duration criteria: %v", err)
	}

	if !result {
		t.Error("Expected criteria to be met, but it wasn't")
	}
}

// Test aggregation criteria evaluation
func TestAggregationCriteria(t *testing.T) {
	// Create a RuleEngine instance with a proper logger
	re := &RuleEngine{
		Logger: logging.NewLogger("TEST-ENGINE", logging.LogLevelError),
	}

	// Create events with values to aggregate
	baseTime := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	events := []models.Event{
		{ID: 1, UserID: "test-user", EventTypeID: 1, OccurredAt: baseTime, Payload: models.JSONB{"score": float64(10)}},
		{ID: 2, UserID: "test-user", EventTypeID: 1, OccurredAt: baseTime.Add(1 * time.Hour), Payload: models.JSONB{"score": float64(20)}},
		{ID: 3, UserID: "test-user", EventTypeID: 1, OccurredAt: baseTime.Add(2 * time.Hour), Payload: models.JSONB{"score": float64(30)}},
		{ID: 4, UserID: "test-user", EventTypeID: 1, OccurredAt: baseTime.Add(3 * time.Hour), Payload: models.JSONB{"score": float64(40)}},
	}

	criteria := map[string]interface{}{
		"type":  "avg",
		"field": "score",
		"value": map[string]interface{}{
			"$gte": float64(20),
		},
	}

	metadata := make(map[string]interface{})
	result, err := re.evaluateAggregationCriteria(criteria, events, metadata)

	if err != nil {
		t.Errorf("Error evaluating aggregation criteria: %v", err)
	}

	if !result {
		t.Error("Expected criteria to be met, but it wasn't")
	}

	if avg, ok := metadata["avg_score"].(float64); !ok || avg != 25.0 {
		t.Errorf("Expected average score of 25.0, got %v", metadata["avg_score"])
	}
}
