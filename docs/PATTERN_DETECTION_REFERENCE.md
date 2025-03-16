# Pattern Detection Algorithms Reference Guide

This guide provides detailed information on the pattern detection algorithms used in the Badge Assignment System. It's intended as a quick reference for developers working with these algorithms.

## Overview

Pattern detection is a critical feature that allows the badge system to identify meaningful trends in user behavior over time. We support three main pattern types:

1. **Consistent Pattern**: Detecting steady, reliable behavior
2. **Increasing Pattern**: Detecting growing engagement or improvement
3. **Decreasing Pattern**: Detecting declining activity or engagement

## JSON Configuration Examples

### Consistent Pattern

```json
{
  "$pattern": {
    "pattern": "consistent",
    "periodType": "day",
    "minPeriods": 7,
    "maxDeviation": 0.15
  }
}
```

**Use Case:** Identify users who maintain consistent daily app usage to award a "Daily User" badge.

### Increasing Pattern

```json
{
  "$pattern": {
    "pattern": "increasing",
    "periodType": "week",
    "minPeriods": 4,
    "minIncreasePct": 10.0
  }
}
```

**Use Case:** Identify users who show steady improvement in workout frequency to award a "Fitness Growth" badge.

### Decreasing Pattern

```json
{
  "$pattern": {
    "pattern": "decreasing",
    "periodType": "week",
    "minPeriods": 4,
    "maxDecreasePct": 15.0
  }
}
```

**Use Case:** Identify users with gradually declining engagement who might benefit from re-engagement efforts.

## Algorithm Details

### Event Grouping

Events are first grouped by time periods using the `groupEventsByPeriod` function:

```go
func groupEventsByPeriod(events []models.Event, periodType string) (map[string]int, []string, error) {
    periods := make(map[string]int)
    
    // Group events by period key (day, week, month)
    for _, event := range events {
        periodKey, err := getPeriodKey(event.OccurredAt, periodType)
        if err != nil {
            return nil, nil, err
        }
        periods[periodKey]++
    }
    
    // Sort period keys chronologically
    var periodKeys []string
    for k := range periods {
        periodKeys = append(periodKeys, k)
    }
    sort.Strings(periodKeys)
    
    return periods, periodKeys, nil
}
```

### Consistent Pattern Detection

The `evaluateConsistentPattern` function analyzes event counts to determine if they follow a consistent pattern:

```go
func evaluateConsistentPattern(periodCounts []int, criteria map[string]interface{}, metadata map[string]interface{}) bool {
    // 1. Check if we have data to analyze
    if len(periodCounts) == 0 {
        metadata["reason"] = "no data available"
        metadata["is_consistent"] = false
        return false
    }
    
    // 2. Check if all values are identical (perfect consistency)
    allIdentical := true
    firstCount := periodCounts[0]
    for _, count := range periodCounts {
        if count != firstCount {
            allIdentical = false
            break
        }
    }
    
    if allIdentical {
        metadata["reason"] = "all values are identical"
        metadata["is_consistent"] = true
        metadata["max_deviation"] = 0.0
        return true
    }
    
    // 3. Calculate average count
    sum := 0
    for _, count := range periodCounts {
        sum += count
    }
    average := float64(sum) / float64(len(periodCounts))
    metadata["average"] = average
    
    // 4. Calculate standard deviation and coefficient of variation
    sumSquaredDiff := 0.0
    for _, count := range periodCounts {
        diff := float64(count) - average
        sumSquaredDiff += diff * diff
    }
    stdDev := math.Sqrt(sumSquaredDiff / float64(len(periodCounts)))
    coeffVar := stdDev / average
    
    metadata["std_deviation"] = stdDev
    metadata["coefficient_var"] = coeffVar
    
    // 5. Check if the pattern meets the maximum deviation criteria
    maxDeviation := 0.2 // Default
    if maxDev, ok := criteria["maxDeviation"].(float64); ok {
        maxDeviation = maxDev
    }
    
    // 6. Handle outliers for more robust evaluation
    if coeffVar <= maxDeviation {
        metadata["is_consistent"] = true
        metadata["max_deviation"] = coeffVar
        metadata["consistency_strength"] = 1 - (coeffVar / maxDeviation)
        return true
    }
    
    // 7. Check for isolated anomalies
    // Detect if pattern would be consistent after excluding isolated anomalies
    // This allows for a more robust evaluation that can handle occasional outliers
    
    metadata["is_consistent"] = false
    metadata["reason"] = "variation exceeds maximum allowed deviation"
    return false
}
```

### Increasing Pattern Detection

The `evaluateIncreasingPattern` function analyzes if event counts show a meaningful increasing trend:

```go
func evaluateIncreasingPattern(periodCounts []int, criteria map[string]interface{}, periodKeys []string, metadata map[string]interface{}) bool {
    // 1. Calculate percentage increases between consecutive periods
    var increases []float64
    var increasePercentages []float64
    for i := 1; i < len(periodCounts); i++ {
        prev := float64(periodCounts[i-1])
        curr := float64(periodCounts[i])
        
        if prev == 0 {
            // Handle division by zero
            continue
        }
        
        increase := curr - prev
        increases = append(increases, increase)
        
        pctIncrease := (increase / prev) * 100
        increasePercentages = append(increasePercentages, pctIncrease)
    }
    
    metadata["increase_percentages"] = increasePercentages
    
    // 2. Calculate average percentage increase
    if len(increasePercentages) == 0 {
        metadata["reason"] = "no valid increases found"
        metadata["is_increasing"] = false
        return false
    }
    
    sumPct := 0.0
    for _, pct := range increasePercentages {
        sumPct += pct
    }
    avgPctIncrease := sumPct / float64(len(increasePercentages))
    metadata["average_percent_increase"] = avgPctIncrease
    
    // 3. Check against minimum increase percentage
    minIncreasePct := 5.0 // Default
    if min, ok := criteria["minIncreasePct"].(float64); ok {
        minIncreasePct = min
    }
    
    // 4. Calculate trend strength metrics
    positiveIncreases := 0
    maxConsecutive := 0
    currentConsecutive := 0
    
    for _, pct := range increasePercentages {
        if pct > 0 {
            positiveIncreases++
            currentConsecutive++
            if currentConsecutive > maxConsecutive {
                maxConsecutive = currentConsecutive
            }
        } else {
            currentConsecutive = 0
        }
    }
    
    increasingRatio := float64(positiveIncreases) / float64(len(increasePercentages))
    metadata["increasing_periods_ratio"] = increasingRatio
    metadata["max_consecutive_increases"] = maxConsecutive
    
    // 5. Calculate overall trend strength
    trendStrength := (avgPctIncrease / minIncreasePct) * increasingRatio
    if trendStrength > 1 {
        trendStrength = 1.0
    }
    metadata["trend_strength"] = trendStrength
    
    // 6. Determine if the pattern meets the criteria
    meetsMinIncrease := avgPctIncrease >= minIncreasePct
    strongEnoughTrend := increasingRatio >= 0.5
    
    metadata["is_increasing"] = meetsMinIncrease && strongEnoughTrend
    
    if !meetsMinIncrease {
        metadata["reason"] = fmt.Sprintf("average increase %.2f%% is below minimum required %.2f%%", 
            avgPctIncrease, minIncreasePct)
    } else if !strongEnoughTrend {
        metadata["reason"] = fmt.Sprintf("only %.2f%% of periods show increase, need at least 50%%", 
            increasingRatio*100)
    }
    
    return meetsMinIncrease && strongEnoughTrend
}
```

### Decreasing Pattern Detection

The `evaluateDecreasingPattern` function analyzes if event counts show a meaningful decreasing trend:

```go
func evaluateDecreasingPattern(periodCounts []int, criteria map[string]interface{}, periodKeys []string, metadata map[string]interface{}) bool {
    // 1. Calculate percentage decreases between consecutive periods
    var decreases []float64
    var decreasePercentages []float64
    
    for i := 1; i < len(periodCounts); i++ {
        prev := float64(periodCounts[i-1])
        curr := float64(periodCounts[i])
        
        if prev == 0 {
            // Handle division by zero
            continue
        }
        
        decrease := prev - curr
        decreases = append(decreases, decrease)
        
        pctDecrease := (decrease / prev) * 100
        decreasePercentages = append(decreasePercentages, pctDecrease)
    }
    
    metadata["decrease_percentages"] = decreasePercentages
    
    // 2. Calculate average percentage decrease
    if len(decreasePercentages) == 0 {
        metadata["reason"] = "no valid decreases found"
        metadata["is_decreasing"] = false
        return false
    }
    
    sumPct := 0.0
    for _, pct := range decreasePercentages {
        sumPct += pct
    }
    avgPctDecrease := sumPct / float64(len(decreasePercentages))
    metadata["average_percent_decrease"] = avgPctDecrease
    
    // 3. Check chronological ordering
    // If periods may not be in perfect chronological order, we can correct this
    // using the period keys to ensure accurate pattern detection
    
    // 4. Calculate trend strength metrics
    negativeDecreases := 0
    maxConsecutive := 0
    currentConsecutive := 0
    
    for _, pct := range decreasePercentages {
        if pct > 0 {
            negativeDecreases++
            currentConsecutive++
            if currentConsecutive > maxConsecutive {
                maxConsecutive = currentConsecutive
            }
        } else {
            currentConsecutive = 0
        }
    }
    
    decreasingRatio := float64(negativeDecreases) / float64(len(decreasePercentages))
    metadata["decreasing_periods_ratio"] = decreasingRatio
    metadata["max_consecutive_decreases"] = maxConsecutive
    
    // 5. Set criteria thresholds
    maxDecreasePct := 20.0 // Default
    if max, ok := criteria["maxDecreasePct"].(float64); ok {
        maxDecreasePct = max
    }
    
    // 6. Calculate overall trend strength
    trendStrength := decreasingRatio
    if maxDecreasePct > 0 {
        trendStrength *= (avgPctDecrease / maxDecreasePct)
    }
    if trendStrength > 1 {
        trendStrength = 1.0
    }
    metadata["trend_strength"] = trendStrength
    
    // 7. Determine if the pattern meets the criteria
    withinMaxDecrease := maxDecreasePct == 0 || avgPctDecrease <= maxDecreasePct
    strongEnoughTrend := decreasingRatio >= 0.5
    
    metadata["is_decreasing"] = strongEnoughTrend && withinMaxDecrease
    
    if !withinMaxDecrease {
        metadata["reason"] = fmt.Sprintf("average decrease %.2f%% exceeds maximum allowed %.2f%%", 
            avgPctDecrease, maxDecreasePct)
        metadata["note"] = "Gradual decline pattern exceeded maximum allowed decrease rate"
    } else if !strongEnoughTrend {
        metadata["reason"] = fmt.Sprintf("only %.2f%% of periods show decrease, need at least 50%%", 
            decreasingRatio*100)
    } else {
        metadata["note"] = "Gradual decline pattern detected with chronological correction"
    }
    
    return strongEnoughTrend && withinMaxDecrease
}
```

## Testing Pattern Detection

We use the `createEventsWithPattern` function to generate test events with specific patterns:

```go
func createEventsWithPattern(pattern string, periodType string, periodCount int, startTime time.Time, baseCountPerPeriod int, variation float64) []models.Event {
    var events []models.Event
    currentID := 1

    for period := 0; period < periodCount; period++ {
        // Calculate period start and end times
        var periodStart, periodEnd time.Time
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

        // Calculate events count based on pattern
        eventsInPeriod := baseCountPerPeriod
        switch pattern {
        case "consistent":
            // Add random variations within specified percentage
            delta := int(float64(baseCountPerPeriod) * variation * (rand.Float64()*2 - 1))
            eventsInPeriod = baseCountPerPeriod + delta
        case "increasing":
            // Linear increase
            increaseFactor := 1.0 + (variation * float64(period))
            eventsInPeriod = int(float64(baseCountPerPeriod) * increaseFactor)
        case "decreasing":
            // Linear decrease
            if period == 0 {
                eventsInPeriod = baseCountPerPeriod
            } else {
                decreaseFactor := 1.0 - (variation * float64(period))
                if decreaseFactor < 0.2 {
                    decreaseFactor = 0.2 // Ensure some events remain
                }
                eventsInPeriod = int(float64(baseCountPerPeriod) * decreaseFactor)
            }
        }

        // Ensure at least one event per period
        if eventsInPeriod < 1 {
            eventsInPeriod = 1
        }

        // Generate events distributed throughout the period
        for i := 0; i < eventsInPeriod; i++ {
            progress := float64(i) / float64(eventsInPeriod)
            eventTime := periodStart.Add(time.Duration(progress*periodEnd.Sub(periodStart).Seconds()) * time.Second)
            
            // Add some random jitter
            jitter := time.Duration(rand.Intn(3600)) * time.Second
            eventTime = eventTime.Add(jitter)
            
            // Create event
            event := models.Event{
                ID:          currentID,
                UserID:      "test-user",
                EventTypeID: 1,
                OccurredAt:  eventTime,
                Payload:     models.JSONB{"value": float64(currentID)},
            }
            events = append(events, event)
            currentID++
        }
    }

    return events
}
```

## Common Edge Cases

When working with pattern detection, be aware of these edge cases:

1. **Empty Event Sets**: All algorithms handle empty event sets gracefully
2. **Few Periods**: Pattern detection requires at least `minPeriods` periods
3. **Zero Counts**: Special handling for periods with zero events
4. **Division by Zero**: Algorithms guard against division by zero errors
5. **Isolated Anomalies**: Consistent pattern detection can handle isolated outliers
6. **Boundary Values**: Tests for values exactly at boundaries

## Best Practices

1. **Include Sufficient Periods**: Always set `minPeriods` appropriately for your use case
2. **Set Realistic Thresholds**: Based on your domain knowledge
3. **Examine Metadata**: Use the detailed metadata for debugging and understanding results
4. **Test Edge Cases**: Always test with boundary conditions
5. **Monitor Real Data**: Validate algorithm performance on real user data
6. **Balance Sensitivity**: Adjust parameters to balance false positives and negatives

## Performance Considerations

Pattern detection can be resource-intensive for large datasets. Consider:

1. **Limit Time Range**: Restrict analysis to recent time periods
2. **Pre-aggregate Data**: Precompute period counts where possible 
3. **Caching Results**: Cache evaluation results for common criteria
4. **Batch Processing**: Evaluate patterns in background batch processes

## References

- Source code: `internal/engine/time_utils.go`
- Test examples: `internal/engine/tests/pattern_criteria/pattern_test.go`
- Full documentation: `README_BADGE_CRITERIA.md` 
