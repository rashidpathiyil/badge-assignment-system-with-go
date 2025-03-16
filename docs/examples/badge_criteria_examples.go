package examples

// This file contains examples of properly formatted badge criteria

import (
	"encoding/json"
	"fmt"
	"time"
)

// GetBasicBadgeCriteria returns a simple badge criteria using comparison operators
func GetBasicBadgeCriteria() map[string]interface{} {
	return map[string]interface{}{
		"event": "test_event",
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": float64(90), // ALWAYS use float64() for numbers
			},
			"completed": true, // Booleans don't need conversion
		},
	}
}

// GetPatternBadgeCriteria returns a badge criteria using pattern detection
func GetPatternBadgeCriteria() map[string]interface{} {
	return map[string]interface{}{
		"$pattern": map[string]interface{}{ // $pattern wrapper is REQUIRED
			"pattern":        "increasing",
			"periodType":     "week",
			"minPeriods":     float64(4),    // ALWAYS use float64() for numbers
			"minIncreasePct": float64(10.0), // ALWAYS use float64() for numbers
		},
	}
}

// GetLogicalAndBadgeCriteria returns a badge criteria using AND logical operator
func GetLogicalAndBadgeCriteria() map[string]interface{} {
	return map[string]interface{}{
		"$and": []interface{}{ // Use []interface{} for arrays
			map[string]interface{}{
				"event": "test_event",
				"criteria": map[string]interface{}{
					"score": map[string]interface{}{
						"$gte": float64(80), // ALWAYS use float64() for numbers
					},
				},
			},
			map[string]interface{}{
				"event": "test_event",
				"criteria": map[string]interface{}{
					"completed": true, // Booleans don't need conversion
				},
			},
		},
	}
}

// GetLogicalOrBadgeCriteria returns a badge criteria using OR logical operator
func GetLogicalOrBadgeCriteria() map[string]interface{} {
	return map[string]interface{}{
		"$or": []interface{}{ // Use []interface{} for arrays
			map[string]interface{}{
				"event": "test_event",
				"criteria": map[string]interface{}{
					"score": map[string]interface{}{
						"$gte": float64(90), // ALWAYS use float64() for numbers
					},
				},
			},
			map[string]interface{}{
				"event": "test_event",
				"criteria": map[string]interface{}{
					"level": "expert", // Strings don't need conversion
				},
			},
		},
	}
}

// GetTimeWindowBadgeCriteria returns a badge criteria using time window
func GetTimeWindowBadgeCriteria() map[string]interface{} {
	startDate := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)
	endDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	return map[string]interface{}{
		"$timeWindow": map[string]interface{}{
			"start": startDate, // RFC3339 format for dates
			"end":   endDate,
			"flow": map[string]interface{}{
				"event": "test_event",
				"criteria": map[string]interface{}{
					"score": map[string]interface{}{
						"$gte": float64(75), // ALWAYS use float64() for numbers
					},
				},
			},
		},
	}
}

// GetComplexBadgeCriteria returns a complex badge criteria combining multiple operators
func GetComplexBadgeCriteria() map[string]interface{} {
	return map[string]interface{}{
		"$and": []interface{}{
			map[string]interface{}{
				"event": "test_event",
				"criteria": map[string]interface{}{
					"score": map[string]interface{}{
						"$gte": float64(80), // ALWAYS use float64() for numbers
					},
				},
			},
			map[string]interface{}{
				"$or": []interface{}{
					map[string]interface{}{
						"event": "test_event",
						"criteria": map[string]interface{}{
							"duration": map[string]interface{}{
								"$lte": float64(300), // ALWAYS use float64() for numbers
							},
						},
					},
					map[string]interface{}{
						"event": "test_event",
						"criteria": map[string]interface{}{
							"completed": map[string]interface{}{
								"$eq": true, // Booleans inside operators
							},
						},
					},
				},
			},
		},
	}
}

// GetMultipleOperatorsBadgeCriteria returns criteria with multiple comparison operators
func GetMultipleOperatorsBadgeCriteria() map[string]interface{} {
	return map[string]interface{}{
		"event": "test_event",
		"criteria": map[string]interface{}{
			"score": map[string]interface{}{
				"$gte": float64(75), // Greater than or equal
				"$lte": float64(95), // Less than or equal - multiple operators for same field
			},
			"tags": map[string]interface{}{
				"$in": []string{"beginner", "intermediate"}, // Array values
			},
		},
	}
}

// CreateBadgeRequest shows the proper API request format for badge creation
// IMPORTANT: flow_definition field is REQUIRED
func CreateBadgeRequest() map[string]interface{} {
	// For API requests, put criteria inside flow_definition
	return map[string]interface{}{
		"name":            "Example Badge",
		"description":     "A badge for achievement",
		"image_url":       "https://example.com/badge.png",
		"flow_definition": GetBasicBadgeCriteria(), // Criteria must be in flow_definition
		"is_active":       true,
	}
}

// CreateEventRequest shows the proper format for event submission
// IMPORTANT: only the payload field is REQUIRED
func CreateEventRequest() map[string]interface{} {
	return map[string]interface{}{
		"event_type": "test_event",
		"user_id":    "user123",
		"timestamp":  time.Now().Format(time.RFC3339),
		"payload": map[string]interface{}{ // Required for both storage and criteria evaluation
			"score":    float64(95),
			"attempts": float64(2),
		},
	}
}

// PrintBadgeCriteriaExample shows how to use these example functions
func PrintBadgeCriteriaExample() {
	// Format for badge creation API request
	badgeRequest := CreateBadgeRequest()

	// Convert to JSON string for viewing
	badgeJSON, _ := json.MarshalIndent(badgeRequest, "", "  ")
	fmt.Println("Badge Creation Request:")
	fmt.Println(string(badgeJSON))

	// Format for event submission
	eventRequest := CreateEventRequest()

	// Convert to JSON string for viewing
	eventJSON, _ := json.MarshalIndent(eventRequest, "", "  ")
	fmt.Println("\nEvent Submission Request:")
	fmt.Println(string(eventJSON))
}
