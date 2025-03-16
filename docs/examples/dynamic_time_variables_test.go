package examples

import (
	"testing"
	"time"
)

func TestTimeVariableResolver(t *testing.T) {
	// Use a fixed time for predictable test results
	fixedTime := time.Date(2023, 12, 15, 12, 0, 0, 0, time.UTC)
	resolver := NewTimeVariableResolverWithFixedTime(fixedTime)

	// Test cases for different time variables
	testCases := []struct {
		name     string
		variable string
		expected string
	}{
		{
			name:     "Simple $NOW",
			variable: "$NOW",
			expected: "2023-12-15T12:00:00Z",
		},
		{
			name:     "30 days ago",
			variable: "$NOW(-30d)",
			expected: "2023-11-15T12:00:00Z",
		},
		{
			name:     "1 year in future",
			variable: "$NOW(+1y)",
			expected: "2024-12-15T12:00:00Z",
		},
		{
			name:     "6 hours ago",
			variable: "$NOW(-6h)",
			expected: "2023-12-15T06:00:00Z",
		},
		{
			name:     "2 weeks ago",
			variable: "$NOW(-2w)",
			expected: "2023-12-01T12:00:00Z",
		},
		{
			name:     "3 months ago",
			variable: "$NOW(-3M)",
			expected: "2023-09-15T12:00:00Z",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a simple criteria with the test variable
			criteria := map[string]interface{}{
				"timestamp": tc.variable,
			}

			// Resolve the criteria
			resolved := resolver.ResolveCriteria(criteria)

			// Check the result
			if timestamp, ok := resolved["timestamp"].(string); ok {
				if timestamp != tc.expected {
					t.Errorf("Expected %s, got %s", tc.expected, timestamp)
				}
			} else {
				t.Errorf("Expected string timestamp, got %T", resolved["timestamp"])
			}
		})
	}
}

func TestComplexCriteriaResolution(t *testing.T) {
	// Use a fixed time for predictable test results
	fixedTime := time.Date(2023, 12, 15, 12, 0, 0, 0, time.UTC)
	resolver := NewTimeVariableResolverWithFixedTime(fixedTime)

	// Create a complex nested criteria
	criteria := map[string]interface{}{
		"$eventCount": map[string]interface{}{
			"$gte": float64(5),
		},
		"timestamp": map[string]interface{}{
			"$gte": "$NOW(-30d)",
			"$lte": "$NOW",
		},
		"nested": map[string]interface{}{
			"created_at": "$NOW(-1M)",
			"updated_at": "$NOW(-1d)",
			"deeper": map[string]interface{}{
				"time": "$NOW(-1h)",
			},
		},
	}

	// Resolve the criteria
	resolved := resolver.ResolveCriteria(criteria)

	// Expected resolved values
	expectedValues := map[string]string{
		"timestamp.$gte":     "2023-11-15T12:00:00Z", // 30 days ago
		"timestamp.$lte":     "2023-12-15T12:00:00Z", // now
		"nested.created_at":  "2023-11-15T12:00:00Z", // 1 month ago
		"nested.updated_at":  "2023-12-14T12:00:00Z", // 1 day ago
		"nested.deeper.time": "2023-12-15T11:00:00Z", // 1 hour ago
	}

	// Verify timestamp.$gte
	timestampCriteria, ok := resolved["timestamp"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected timestamp to be a map, got %T", resolved["timestamp"])
	}

	if gte, ok := timestampCriteria["$gte"].(string); ok {
		if gte != expectedValues["timestamp.$gte"] {
			t.Errorf("Expected timestamp.$gte to be %s, got %s",
				expectedValues["timestamp.$gte"], gte)
		}
	} else {
		t.Errorf("Expected timestamp.$gte to be string, got %T", timestampCriteria["$gte"])
	}

	// Verify timestamp.$lte
	if lte, ok := timestampCriteria["$lte"].(string); ok {
		if lte != expectedValues["timestamp.$lte"] {
			t.Errorf("Expected timestamp.$lte to be %s, got %s",
				expectedValues["timestamp.$lte"], lte)
		}
	} else {
		t.Errorf("Expected timestamp.$lte to be string, got %T", timestampCriteria["$lte"])
	}

	// Verify nested values
	nestedCriteria, ok := resolved["nested"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected nested to be a map, got %T", resolved["nested"])
	}

	// Check created_at
	if createdAt, ok := nestedCriteria["created_at"].(string); ok {
		if createdAt != expectedValues["nested.created_at"] {
			t.Errorf("Expected nested.created_at to be %s, got %s",
				expectedValues["nested.created_at"], createdAt)
		}
	} else {
		t.Errorf("Expected nested.created_at to be string, got %T", nestedCriteria["created_at"])
	}

	// Check updated_at
	if updatedAt, ok := nestedCriteria["updated_at"].(string); ok {
		if updatedAt != expectedValues["nested.updated_at"] {
			t.Errorf("Expected nested.updated_at to be %s, got %s",
				expectedValues["nested.updated_at"], updatedAt)
		}
	} else {
		t.Errorf("Expected nested.updated_at to be string, got %T", nestedCriteria["updated_at"])
	}

	// Check deeper.time
	deeperCriteria, ok := nestedCriteria["deeper"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected nested.deeper to be a map, got %T", nestedCriteria["deeper"])
	}

	if time, ok := deeperCriteria["time"].(string); ok {
		if time != expectedValues["nested.deeper.time"] {
			t.Errorf("Expected nested.deeper.time to be %s, got %s",
				expectedValues["nested.deeper.time"], time)
		}
	} else {
		t.Errorf("Expected nested.deeper.time to be string, got %T", deeperCriteria["time"])
	}
}

// Test demonstrating how this would be used for a "last 30 days" badge
func TestLastThirtyDaysBadge(t *testing.T) {
	// Use a fixed time for the test
	fixedTime := time.Date(2023, 12, 15, 12, 0, 0, 0, time.UTC)
	resolver := NewTimeVariableResolverWithFixedTime(fixedTime)

	// Badge criteria with dynamic time window
	flowDefinition := map[string]interface{}{
		"event": "user_activity",
		"criteria": map[string]interface{}{
			"$eventCount": map[string]interface{}{
				"$gte": float64(5),
			},
			"timestamp": map[string]interface{}{
				"$gte": "$NOW(-30d)", // Dynamic "last 30 days" window
			},
		},
	}

	// Resolve the time variables
	resolvedFlow := resolver.ResolveCriteria(flowDefinition)

	// Extract the resolved criteria for verification
	criteria, ok := resolvedFlow["criteria"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected criteria to be a map, got %T", resolvedFlow["criteria"])
	}

	timestamp, ok := criteria["timestamp"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected timestamp to be a map, got %T", criteria["timestamp"])
	}

	// Verify the resolved timestamp
	expectedTimestamp := "2023-11-15T12:00:00Z" // 30 days before our fixed time
	if gte, ok := timestamp["$gte"].(string); ok {
		if gte != expectedTimestamp {
			t.Errorf("Expected timestamp.$gte to be %s, got %s",
				expectedTimestamp, gte)
		}
	} else {
		t.Errorf("Expected timestamp.$gte to be string, got %T", timestamp["$gte"])
	}

	// Verify that other parts of the criteria remain unchanged
	eventCount, ok := criteria["$eventCount"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected $eventCount to be a map, got %T", criteria["$eventCount"])
	}

	if gte, ok := eventCount["$gte"].(float64); ok {
		if gte != float64(5) {
			t.Errorf("Expected $eventCount.$gte to be 5, got %v", gte)
		}
	} else {
		t.Errorf("Expected $eventCount.$gte to be float64, got %T", eventCount["$gte"])
	}
}
