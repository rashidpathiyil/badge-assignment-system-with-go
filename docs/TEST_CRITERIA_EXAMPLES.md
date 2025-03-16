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
