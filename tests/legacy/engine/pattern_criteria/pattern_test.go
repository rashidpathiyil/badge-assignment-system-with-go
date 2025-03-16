package pattern_criteria

import (
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/engine"
	"github.com/badge-assignment-system/internal/logging"
	"github.com/badge-assignment-system/internal/models"
)

// MockLogger is a simplified logger for testing
type MockLogger struct{}

// Debug logs a debug message
func (l *MockLogger) Debug(format string, args ...interface{}) {}

// Info logs an info message
func (l *MockLogger) Info(format string, args ...interface{}) {}

// Warn logs a warning message
func (l *MockLogger) Warn(format string, args ...interface{}) {}

// Error logs an error message
func (l *MockLogger) Error(format string, args ...interface{}) {}

// createEventsWithPattern generates events with specific pattern characteristics
func createEventsWithPattern(pattern string, periodType string, periodCount int, startTime time.Time, baseCountPerPeriod int, variation float64) []models.Event {
	var events []models.Event
	currentID := 1

	for period := 0; period < periodCount; period++ {
		var periodStart time.Time
		var periodEnd time.Time

		// Calculate period boundaries
		switch periodType {
		case "day":
			periodStart = startTime.AddDate(0, 0, period)
			periodEnd = startTime.AddDate(0, 0, period+1)
		case "week":
			periodStart = startTime.AddDate(0, 0, period*7)
			periodEnd = startTime.AddDate(0, 0, (period+1)*7)
		case "month":
			periodStart = startTime.AddDate(0, period, 0)
			periodEnd = startTime.AddDate(0, period+1, 0)
		}

		// Determine number of events for this period based on pattern
		countForPeriod := baseCountPerPeriod
		switch pattern {
		case "consistent":
			// Add slight random variation around the base count
			variationRange := int(float64(baseCountPerPeriod) * variation)
			if variationRange > 0 {
				countForPeriod = baseCountPerPeriod + rand.Intn(2*variationRange+1) - variationRange
			}
		case "increasing":
			// Increase by specified percentage each period
			increaseFactor := 1.0 + (variation * float64(period) / float64(periodCount))
			countForPeriod = int(float64(baseCountPerPeriod) * increaseFactor)
		case "decreasing":
			// Decrease by specified percentage each period
			decreaseFactor := 1.0 - (variation * float64(period) / float64(periodCount))
			countForPeriod = int(float64(baseCountPerPeriod) * decreaseFactor)
		}

		// Ensure at least one event per period
		if countForPeriod < 1 {
			countForPeriod = 1
		}

		// Generate events for this period
		for i := 0; i < countForPeriod; i++ {
			// Calculate a random time within the period
			periodDuration := periodEnd.Sub(periodStart)
			randomOffset := time.Duration(rand.Int63n(int64(periodDuration)))
			eventTime := periodStart.Add(randomOffset)

			// Create and append the event
			event := models.Event{
				ID:          currentID,
				UserID:      "test-user",
				EventTypeID: 1,
				OccurredAt:  eventTime,
				Payload:     models.JSONB{"value": float64(currentID * 10)},
			}

			events = append(events, event)
			currentID++
		}
	}

	return events
}

// setupTestRuleEngine creates a RuleEngine for testing with a silent logger
func setupTestRuleEngine() *engine.RuleEngine {
	logger := &MockLogger{}
	return &engine.RuleEngine{
		Logger: logger,
	}
}

// TestUserEngagementPattern tests a scenario where a user's app engagement pattern is analyzed
func TestUserEngagementPattern(t *testing.T) {
	// Create a RuleEngine instance for testing
	re := setupTestRuleEngine()

	// Scenario: We're analyzing daily app usage to see if it's consistent over time
	// A consistent pattern indicates a habit has formed and the user should receive a "Daily User" badge
	startTime := time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC)
	dailyUsageEvents := createEventsWithPattern(
		"consistent", // Pattern type
		"day",        // Period type
		30,           // 30 days of data
		startTime,    // Starting Jan 1, 2023
		5,            // Base of ~5 sessions per day
		0.15,         // With ±15% variation
	)

	// Configure the criteria for a "Daily User" badge
	criteria := map[string]interface{}{
		"$pattern": map[string]interface{}{
			"pattern":      "consistent", // Looking for consistent usage
			"periodType":   "day",        // Daily pattern
			"minPeriods":   float64(20),  // Need at least 20 days of data
			"maxDeviation": float64(0.20), // Allow up to 20% deviation
		},
	}

	// Evaluate the pattern criteria
	metadata := make(map[string]interface{})
	result, err := re.EvaluatePatternCriteria(criteria, dailyUsageEvents, metadata)

	// Assertions
	if err != nil {
		t.Errorf("Error evaluating consistent usage pattern: %v", err)
	}

	if !result {
		t.Errorf("Expected user to qualify for Daily User badge, but criteria was not met. Metadata: %v", metadata)
	}

	// Verify that we have some period data in the metadata
	if periods, ok := metadata["period_keys"].([]string); !ok || len(periods) < 20 {
		t.Errorf("Expected at least 20 period keys in metadata, got: %v", metadata["period_keys"])
	}

	t.Logf("Daily User pattern detected with consistency metrics: %v", metadata)
}

// TestFitnessProgressPattern tests a scenario where a user's fitness activity is increasing over time
func TestFitnessProgressPattern(t *testing.T) {
	// Create a RuleEngine instance for testing
	re := setupTestRuleEngine()

	// Scenario: An exercise app tracks workout counts per week, looking for steady improvement
	// If workout frequency increases over time, the user earns a "Fitness Growth" badge
	startTime := time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC)
	workoutEvents := createEventsWithPattern(
		"increasing", // Pattern type
		"week",       // Weekly data
		8,            // 8 weeks of data
		startTime,    // Starting Jan 1, 2023
		3,            // Starting with ~3 workouts per week
		0.2,          // Increasing by ~20% each week
	)

	// Configure the criteria for a "Fitness Growth" badge
	criteria := map[string]interface{}{
		"$pattern": map[string]interface{}{
			"pattern":        "increasing", // Looking for increasing pattern
			"periodType":     "week",       // Weekly pattern
			"minPeriods":     float64(6),   // Need at least 6 weeks of data
			"minIncreasePct": float64(10.0), // At least 10% average increase
		},
	}

	// Evaluate the pattern criteria
	metadata := make(map[string]interface{})
	result, err := re.EvaluatePatternCriteria(criteria, workoutEvents, metadata)

	// Assertions
	if err != nil {
		t.Errorf("Error evaluating fitness progress pattern: %v", err)
	}

	if !result {
		t.Errorf("Expected user to qualify for Fitness Growth badge, but criteria was not met. Metadata: %v", metadata)
	}

	// Verify we have the expected increase percentage in the metadata
	if avgIncrease, ok := metadata["average_percent_increase"].(float64); !ok || avgIncrease < 10.0 {
		t.Errorf("Expected average_percent_increase to be at least 10%%, got: %v", metadata["average_percent_increase"])
	}

	t.Logf("Fitness Growth pattern detected with growth metrics: %v", metadata)
}

// TestLearningPatternDecline tests a scenario where a user's learning platform usage declines over time
func TestLearningPatternDecline(t *testing.T) {
	// Create a RuleEngine instance for testing
	re := setupTestRuleEngine()

	// Scenario: An educational platform tracks course engagement, looking for declining usage
	// If course views decrease over time, the user gets a "Re-engagement" badge to encourage them
	startTime := time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC)
	learningEvents := createEventsWithPattern(
		"decreasing", // Pattern type
		"week",       // Weekly data
		8,            // 8 weeks of data
		startTime,    // Starting Jan 1, 2023
		10,           // Starting with ~10 course views per week
		0.15,         // Decreasing by ~15% each week
	)

	// Configure the criteria for a "Re-engagement Needed" badge
	criteria := map[string]interface{}{
		"$pattern": map[string]interface{}{
			"pattern":        "decreasing", // Looking for decreasing pattern
			"periodType":     "week",       // Weekly pattern
			"minPeriods":     float64(6),   // Need at least 6 weeks of data
			"maxDecreasePct": float64(15.0), // Maximum 15% decrease per period (looking for gradual decline)
		},
	}

	// Evaluate the pattern criteria
	metadata := make(map[string]interface{})
	result, err := re.EvaluatePatternCriteria(criteria, learningEvents, metadata)

	// Assertions
	if err != nil {
		t.Errorf("Error evaluating learning pattern decline: %v", err)
	}

	if !result {
		t.Errorf("Expected user to qualify for Re-engagement badge, but criteria was not met. Metadata: %v", metadata)
	}

	// Verify we have the expected decrease percentage in the metadata
	if avgDecrease, ok := metadata["average_percent_decrease"].(float64); !ok || avgDecrease < 5.0 {
		t.Errorf("Expected average_percent_decrease to be at least 5%%, got: %v", metadata["average_percent_decrease"])
	}

	t.Logf("Learning decline pattern detected with metrics: %v", metadata)
}

// TestMixedPatterns tests detection of different patterns in the same dataset
func TestMixedPatterns(t *testing.T) {
	// Create a RuleEngine instance for testing
	re := setupTestRuleEngine()

	// Scenario: A dataset with multiple potential patterns
	startTime := time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC)
	mixedEvents := []models.Event{}

	// Weekly consistent data, but with 10% growth trend
	for week := 0; week < 10; week++ {
		weekBase := 5 + int(float64(week)*0.1*5) // Start at 5, grow 10% per week
		weekEvents := createEventsWithPattern(
			"consistent",
			"day",
			7, // 7 days per week
			startTime.AddDate(0, 0, week*7),
			weekBase,
			0.2, // 20% daily variation
		)
		mixedEvents = append(mixedEvents, weekEvents...)
	}

	// Test weekly growth pattern detection
	weeklyGrowthCriteria := map[string]interface{}{
		"$pattern": map[string]interface{}{
			"pattern":        "increasing",
			"periodType":     "week",
			"minPeriods":     float64(8),
			"minIncreasePct": float64(5.0),
		},
	}

	weeklyMetadata := make(map[string]interface{})
	weeklyResult, err := re.EvaluatePatternCriteria(weeklyGrowthCriteria, mixedEvents, weeklyMetadata)

	if err != nil {
		t.Errorf("Error evaluating weekly growth pattern: %v", err)
	}

	if !weeklyResult {
		t.Errorf("Expected weekly growth pattern to be detected, but it wasn't. Metadata: %v", weeklyMetadata)
	}

	// Test daily consistency pattern detection
	dailyConsistencyCriteria := map[string]interface{}{
		"$pattern": map[string]interface{}{
			"pattern":      "consistent",
			"periodType":   "day",
			"minPeriods":   float64(30),
			"maxDeviation": float64(0.25), // Allow up to 25% deviation for weekly consistency
		},
	}

	dailyMetadata := make(map[string]interface{})
	dailyResult, err := re.EvaluatePatternCriteria(dailyConsistencyCriteria, mixedEvents, dailyMetadata)

	if err != nil {
		t.Errorf("Error evaluating daily consistency pattern: %v", err)
	}

	// Test a more demanding criteria that shouldn't be met
	toughCriteria := map[string]interface{}{
		"$pattern": map[string]interface{}{
			"pattern":        "increasing",
			"periodType":     "week",
			"minPeriods":     float64(7),
			"minIncreasePct": float64(20.0), // Require 20% growth - more than our 10% simulation
		},
	}

	toughMetadata := make(map[string]interface{})
	toughResult, err := re.EvaluatePatternCriteria(toughCriteria, mixedEvents, toughMetadata)

	if err != nil {
		t.Errorf("Error evaluating tough growth pattern: %v", err)
	}

	if toughResult {
		t.Errorf("Expected tough criteria NOT to be met, but it was. Metadata: %v", toughMetadata)
	}

	t.Logf("Mixed pattern detection results: Weekly growth: %v, Daily consistency: %v, Tough criteria: %v",
		weeklyResult, dailyResult, toughResult)
}

// TestEdgeCasePatterns tests boundary conditions and edge cases
func TestEdgeCasePatterns(t *testing.T) {
	// Create a RuleEngine instance for testing
	re := setupTestRuleEngine()

	// Test Case 1: Exactly at the boundary of consistency
	t.Run("BoundaryConsistency", func(t *testing.T) {
		startTime := time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC)

		// Create events with exactly 15% deviation (boundary test)
		events := []models.Event{
			{ID: 1, UserID: "test-user", EventTypeID: 1, OccurredAt: startTime.AddDate(0, 0, 0), Payload: models.JSONB{"value": float64(100)}},
			{ID: 2, UserID: "test-user", EventTypeID: 1, OccurredAt: startTime.AddDate(0, 0, 1), Payload: models.JSONB{"value": float64(115)}}, // +15%
			{ID: 3, UserID: "test-user", EventTypeID: 1, OccurredAt: startTime.AddDate(0, 0, 2), Payload: models.JSONB{"value": float64(85)}},  // -15%
			{ID: 4, UserID: "test-user", EventTypeID: 1, OccurredAt: startTime.AddDate(0, 0, 3), Payload: models.JSONB{"value": float64(100)}},
		}

		criteria := map[string]interface{}{
			"$pattern": map[string]interface{}{
				"pattern":      "consistent",
				"periodType":   "day",
				"minPeriods":   float64(4),
				"maxDeviation": float64(0.15), // Exactly at our test boundary
			},
		}

		metadata := make(map[string]interface{})
		result, err := re.EvaluatePatternCriteria(criteria, events, metadata)

		if err != nil {
			t.Errorf("Error evaluating boundary consistency pattern: %v", err)
		}

		// This should pass since we're exactly at the boundary
		if !result {
			t.Errorf("Expected boundary consistency criteria to be met, but it wasn't. Metadata: %v", metadata)
		}
	})

	// Test Case 2: Just over the boundary of consistency
	t.Run("ExceedingConsistency", func(t *testing.T) {
		startTime := time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC)

		// Create events with 15.1% deviation (just over boundary)
		events := []models.Event{
			{ID: 1, UserID: "test-user", EventTypeID: 1, OccurredAt: startTime.AddDate(0, 0, 0), Payload: models.JSONB{"value": float64(100)}},
			{ID: 2, UserID: "test-user", EventTypeID: 1, OccurredAt: startTime.AddDate(0, 0, 1), Payload: models.JSONB{"value": float64(115.1)}}, // +15.1%
			{ID: 3, UserID: "test-user", EventTypeID: 1, OccurredAt: startTime.AddDate(0, 0, 2), Payload: models.JSONB{"value": float64(100)}},
		}

		criteria := map[string]interface{}{
			"$pattern": map[string]interface{}{
				"pattern":      "consistent",
				"periodType":   "day",
				"minPeriods":   float64(3),
				"maxDeviation": float64(0.15), // Below our 15.1% test value
			},
		}

		metadata := make(map[string]interface{})
		result, err := re.EvaluatePatternCriteria(criteria, events, metadata)

		if err != nil {
			t.Errorf("Error evaluating exceeding consistency pattern: %v", err)
		}

		// This should fail since we're over the boundary
		if result {
			t.Errorf("Expected exceeding consistency criteria NOT to be met, but it was. Metadata: %v", metadata)
		}
	})

	// Test Case 3: Minimum number of periods
	t.Run("MinimumPeriods", func(t *testing.T) {
		startTime := time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC)

		// Create events with consistent values but one fewer period than required
		events := []models.Event{
			{ID: 1, UserID: "test-user", EventTypeID: 1, OccurredAt: startTime.AddDate(0, 0, 0), Payload: models.JSONB{"value": float64(100)}},
			{ID: 2, UserID: "test-user", EventTypeID: 1, OccurredAt: startTime.AddDate(0, 0, 1), Payload: models.JSONB{"value": float64(100)}},
			{ID: 3, UserID: "test-user", EventTypeID: 1, OccurredAt: startTime.AddDate(0, 0, 2), Payload: models.JSONB{"value": float64(100)}},
		}

		criteria := map[string]interface{}{
			"$pattern": map[string]interface{}{
				"pattern":      "consistent",
				"periodType":   "day",
				"minPeriods":   float64(4), // Require 4 periods, but we only have 3
				"maxDeviation": float64(0.1),
			},
		}

		metadata := make(map[string]interface{})
		result, err := re.EvaluatePatternCriteria(criteria, events, metadata)

		if err != nil {
			t.Errorf("Error evaluating minimum periods pattern: %v", err)
		}

		// This should fail since we don't have enough periods
		if result {
			t.Errorf("Expected minimum periods criteria NOT to be met, but it was. Metadata: %v", metadata)
		}
	})

	// Test Case 4: No events
	t.Run("NoEvents", func(t *testing.T) {
		// Empty events array
		events := []models.Event{}

		criteria := map[string]interface{}{
			"$pattern": map[string]interface{}{
				"pattern":      "consistent",
				"periodType":   "day",
				"minPeriods":   float64(3),
				"maxDeviation": float64(0.3),
			},
		}

		metadata := make(map[string]interface{})
		result, err := re.EvaluatePatternCriteria(criteria, events, metadata)

		if err != nil {
			t.Errorf("Error evaluating no events pattern: %v", err)
		}

		// This should fail since we have no events
		if result {
			t.Errorf("Expected no events criteria NOT to be met, but it was. Metadata: %v", metadata)
		}
	})
}
