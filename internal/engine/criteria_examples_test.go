package engine

import (
	"testing"
	"time"
)

// TestCriteriaExamples demonstrates correct criteria formats
// Note: This test is skipped as it's meant for documentation purposes only
func TestCriteriaExamples(t *testing.T) {
	t.Skip("This test is for documentation purposes only")

	// Sample event-based criteria
	eventCriteria := map[string]interface{}{
		"event": "test_event",
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": float64(90),
			},
		},
	}

	// Sample pattern criteria - consistent
	consistentPatternCriteria := map[string]interface{}{
		"$pattern": map[string]interface{}{
			"pattern":      "consistent",
			"periodType":   "day",
			"minPeriods":   float64(7),
			"maxDeviation": float64(0.15),
		},
	}

	// Sample pattern criteria - increasing
	increasingPatternCriteria := map[string]interface{}{
		"$pattern": map[string]interface{}{
			"pattern":        "increasing",
			"periodType":     "week",
			"minPeriods":     float64(4),
			"minIncreasePct": float64(10.0),
		},
	}

	// Sample time window criteria
	timeWindowCriteria := map[string]interface{}{
		"$timeWindow": map[string]interface{}{
			"start": time.Now().AddDate(0, 0, -30).Format(time.RFC3339),
			"end":   time.Now().Format(time.RFC3339),
			"flow": map[string]interface{}{
				"event": "test_event",
				"criteria": map[string]interface{}{
					"value": map[string]interface{}{
						"$gte": float64(0),
					},
				},
			},
		},
	}

	// Sample logical operator - AND
	andCriteria := map[string]interface{}{
		"$and": []interface{}{
			map[string]interface{}{
				"event": "event_type_1",
				"criteria": map[string]interface{}{
					"field_1": map[string]interface{}{
						"$gte": float64(10),
					},
				},
			},
			map[string]interface{}{
				"event": "event_type_2",
				"criteria": map[string]interface{}{
					"field_2": map[string]interface{}{
						"$eq": "value",
					},
				},
			},
		},
	}

	// These are here to avoid compiler warnings - not meant to be executed
	_ = eventCriteria
	_ = consistentPatternCriteria
	_ = increasingPatternCriteria
	_ = timeWindowCriteria
	_ = andCriteria
}
