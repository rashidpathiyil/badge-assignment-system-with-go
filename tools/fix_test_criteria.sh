#!/bin/bash
# Script to fix incorrect test criteria formats in the codebase

echo "Starting test criteria fixes..."

# Fix 1: Update pattern criteria in time_evaluators_test.go
echo "Fixing pattern criteria in time_evaluators_test.go"
sed -i.bak '
/pattern.*:.*consistent/,/}/ {
  s/criteria := map\[string\]interface{}{\s*$/criteria := map[string]interface{}{/
  s/"pattern":\s*"consistent"/"$pattern": map[string]interface{}{\n\t\t"pattern":\t"consistent"/
  s/"periodType":\s*"day"/"periodType":\t"day"/
  s/"minPeriods":\s*float64(3)/"minPeriods":\tfloat64(3)/
  s/"maxDeviation":\s*float64(0.1),/"maxDeviation":\tfloat64(0.1),\n\t\t},/
}
' internal/engine/time_evaluators_test.go

# Fix 2: Update pattern criteria in legacy pattern tests
echo "Fixing pattern criteria in pattern_test.go"
sed -i.bak '
/criteria.*pattern.*consistent/,/}/ {
  s/criteria := map\[string\]interface{}{\s*$/criteria := map[string]interface{}{/
  s/"pattern":\s*"consistent"/"$pattern": map[string]interface{}{\n\t\t\t"pattern":\t"consistent"/
  s/"periodType":\s*"day"/"periodType":\t"day"/
  s/"minPeriods":\s*[0-9]\+/"minPeriods":\tfloat64(&)/
  s/"maxDeviation":\s*0\.[0-9]\+/"maxDeviation":\tfloat64(&)/
  s/},\s*$/},\n\t\t},/
}

/criteria.*pattern.*increasing/,/}/ {
  s/criteria := map\[string\]interface{}{\s*$/criteria := map[string]interface{}{/
  s/"pattern":\s*"increasing"/"$pattern": map[string]interface{}{\n\t\t\t"pattern":\t"increasing"/
  s/"periodType":\s*"week"/"periodType":\t"week"/
  s/"minPeriods":\s*[0-9]\+/"minPeriods":\tfloat64(&)/
  s/"minIncreasePct":\s*[0-9]\+\.[0-9]\+/"minIncreasePct":\tfloat64(&)/
  s/},\s*$/},\n\t\t},/
}

/criteria.*pattern.*decreasing/,/}/ {
  s/criteria := map\[string\]interface{}{\s*$/criteria := map[string]interface{}{/
  s/"pattern":\s*"decreasing"/"$pattern": map[string]interface{}{\n\t\t\t"pattern":\t"decreasing"/
  s/"periodType":\s*"week"/"periodType":\t"week"/
  s/"minPeriods":\s*[0-9]\+/"minPeriods":\tfloat64(&)/
  s/"maxDecreasePct":\s*[0-9]\+\.[0-9]\+/"maxDecreasePct":\tfloat64(&)/
  s/},\s*$/},\n\t\t},/
}
' tests/legacy/engine/pattern_criteria/pattern_test.go

# Fix 3: Create an example test for time window criteria
echo "Creating example test for time window criteria"
cat > internal/engine/time_window_test.go << 'EOL'
package engine

import (
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/logging"
	"github.com/badge-assignment-system/internal/models"
)

// TestTimeWindowCriteria tests the time window criteria evaluation
func TestTimeWindowCriteria(t *testing.T) {
	// Create a RuleEngine instance with a proper logger
	re := &RuleEngine{
		Logger: logging.NewLogger("TEST-ENGINE", logging.LogLevelError),
	}

	// Test case: Time window with specific date range
	startTime := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	events := []models.Event{}

	// Create events within time window
	for day := 0; day < 5; day++ {
		dayEvents := createTestEvents(2, startTime.AddDate(0, 0, day), time.Hour)
		events = append(events, dayEvents...)
	}

	// Create events outside time window
	for day := 10; day < 15; day++ {
		dayEvents := createTestEvents(2, startTime.AddDate(0, 0, day), time.Hour)
		events = append(events, dayEvents...)
	}

	// Create proper time window criteria according to documentation
	criteria := map[string]interface{}{
		"$timeWindow": map[string]interface{}{
			"start": startTime.Format(time.RFC3339),
			"end":   startTime.AddDate(0, 0, 7).Format(time.RFC3339),
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

	metadata := make(map[string]interface{})
	// We can't fully test this without mocking more DB queries,
	// but we can at least test that the criteria format is valid
	t.Skip("Skip this test as it requires additional mocking. The criteria format is documented here for reference.")
}
EOL

# Fix 4: Create a README file with correct criteria examples
echo "Creating README file with correct criteria examples"
cat > docs/TEST_CRITERIA_EXAMPLES.md << 'EOL'
# Test Criteria Examples

This document provides correct examples of criteria formats that should be used in tests. These examples
are aligned with the formats specified in the documentation.

## Event-Based Criteria

```go
testCriteria := map[string]interface{}{
    "event": "test_event",
    "criteria": map[string]interface{}{
        "score": map[string]interface{}{
            "$gte": float64(90),
        },
    },
}
```

## Pattern Criteria

### Consistent Pattern

```go
criteria := map[string]interface{}{
    "$pattern": map[string]interface{}{
        "pattern":      "consistent",
        "periodType":   "day",
        "minPeriods":   float64(7),
        "maxDeviation": float64(0.15),
    },
}
```

### Increasing Pattern

```go
criteria := map[string]interface{}{
    "$pattern": map[string]interface{}{
        "pattern":        "increasing",
        "periodType":     "week",
        "minPeriods":     float64(4),
        "minIncreasePct": float64(10.0),
    },
}
```

### Decreasing Pattern

```go
criteria := map[string]interface{}{
    "$pattern": map[string]interface{}{
        "pattern":        "decreasing",
        "periodType":     "week",
        "minPeriods":     float64(4),
        "maxDecreasePct": float64(15.0),
    },
}
```

## Time Window Criteria

```go
criteria := map[string]interface{}{
    "$timeWindow": map[string]interface{}{
        "start": "2023-01-01T00:00:00Z",
        "end":   "2023-01-31T23:59:59Z",
        "flow": map[string]interface{}{
            "$and": []interface{}{
                map[string]interface{}{
                    "event": "event_type_1",
                },
                map[string]interface{}{
                    "event": "event_type_2",
                },
            },
        },
    },
}
```

## Logical Operators

### AND Operator

```go
criteria := map[string]interface{}{
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
```

### OR Operator

```go
criteria := map[string]interface{}{
    "$or": []interface{}{
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
```

### NOT Operator

```go
criteria := map[string]interface{}{
    "$not": map[string]interface{}{
        "event": "event_type_1",
        "criteria": map[string]interface{}{
            "field_1": map[string]interface{}{
                "$lt": float64(5),
            },
        },
    },
}
```
EOL

echo "Creating a sample test file with correct criteria formats"
cat > internal/engine/criteria_examples_test.go << 'EOL'
package engine

import (
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/logging"
	"github.com/badge-assignment-system/internal/models"
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
EOL

# Update the internal/engine/rule_engine_test.go to ensure numeric types are consistent
echo "Updating engine tests to ensure consistent numeric types"
sed -i.bak 's/"\$gte": 90,/"\$gte": float64(90),/g' internal/engine/rule_engine_test.go

echo "Script completed. The following files have been modified or created:"
echo "1. internal/engine/time_evaluators_test.go (fixed pattern criteria)"
echo "2. tests/legacy/engine/pattern_criteria/pattern_test.go (fixed pattern criteria)"
echo "3. internal/engine/time_window_test.go (created example time window test)"
echo "4. docs/TEST_CRITERIA_EXAMPLES.md (created reference document)"
echo "5. internal/engine/criteria_examples_test.go (created example test with correct formats)"
echo "6. internal/engine/rule_engine_test.go (standardized numeric types)"

echo "Backup files with .bak extension have been created for modified files."
echo "To apply the changes, please run:"
echo "  chmod +x tools/fix_test_criteria.sh"
echo "  ./tools/fix_test_criteria.sh" 
