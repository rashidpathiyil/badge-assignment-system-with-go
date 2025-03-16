package engine

import (
	"fmt"
	"sort"

	"github.com/badge-assignment-system/internal/models"
)

// evaluateTimePeriodCriteria handles time period counts (days, weeks, months)
func (re *RuleEngine) evaluateTimePeriodCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	re.Logger.Debug("Evaluating time period criteria with %d events", len(events))

	// Parse and validate criteria
	var timePeriodCriteria models.TimePeriodCriteria

	if periodType, ok := criteria["periodType"].(string); ok {
		timePeriodCriteria.PeriodType = periodType
		re.Logger.Debug("Time period type: %s", periodType)
	} else {
		re.Logger.Error("Missing or invalid periodType in timePeriod criteria")
		return false, fmt.Errorf("missing or invalid periodType in timePeriod criteria")
	}

	if periodCount, ok := criteria["periodCount"].(map[string]interface{}); ok {
		timePeriodCriteria.PeriodCount = periodCount
		re.Logger.Debug("Period count criteria: %v", periodCount)
	}

	if excludeWeekends, ok := criteria["excludeWeekends"].(bool); ok {
		timePeriodCriteria.ExcludeWeekends = excludeWeekends
		re.Logger.Debug("Exclude weekends: %v", excludeWeekends)
	}

	if excludeHolidays, ok := criteria["excludeHolidays"].(bool); ok {
		timePeriodCriteria.ExcludeHolidays = excludeHolidays
		re.Logger.Debug("Exclude holidays: %v", excludeHolidays)
	}

	if holidays, ok := criteria["holidays"].([]interface{}); ok {
		for _, holiday := range holidays {
			if holidayStr, ok := holiday.(string); ok {
				timePeriodCriteria.Holidays = append(timePeriodCriteria.Holidays, holidayStr)
			}
		}
		re.Logger.Debug("Holidays to exclude: %v", timePeriodCriteria.Holidays)
	}

	// Filter events based on exclusions
	var filteredEvents []models.Event
	for _, event := range events {
		// Check weekend exclusion
		if timePeriodCriteria.ExcludeWeekends && isWeekend(event.OccurredAt) {
			re.Logger.Trace("Excluding event ID %d (weekend: %s)", event.ID, event.OccurredAt.Weekday().String())
			continue
		}

		// Check holiday exclusion
		if timePeriodCriteria.ExcludeHolidays && isHoliday(event.OccurredAt, timePeriodCriteria.Holidays) {
			re.Logger.Trace("Excluding event ID %d (holiday: %s)", event.ID, event.OccurredAt.Format("2006-01-02"))
			continue
		}

		filteredEvents = append(filteredEvents, event)
	}

	re.Logger.Debug("After exclusions: %d/%d events remain", len(filteredEvents), len(events))

	// Group events by the specified time period
	uniquePeriods := make(map[string]bool)

	for _, event := range filteredEvents {
		periodKey, err := getPeriodKey(event.OccurredAt, timePeriodCriteria.PeriodType)
		if err != nil {
			re.Logger.Error("Failed to get period key: %v", err)
			return false, err
		}
		uniquePeriods[periodKey] = true
		re.Logger.Trace("Event ID %d belongs to period: %s", event.ID, periodKey)
	}

	// Count unique periods
	periodCount := len(uniquePeriods)
	metadata["unique_period_count"] = periodCount
	re.Logger.Debug("Unique period count: %d", periodCount)

	// If there's a period count criteria, evaluate it
	if len(timePeriodCriteria.PeriodCount) > 0 {
		re.Logger.Debug("Evaluating period count criteria against period count: %d", periodCount)
		result, err := re.evaluateNumericCriteria(float64(periodCount), timePeriodCriteria.PeriodCount)
		if err != nil {
			re.Logger.Error("Error evaluating numeric criteria: %v", err)
			return false, err
		}
		re.Logger.Debug("Time period criteria evaluation result: %v", result)
		return result, nil
	}

	// Default behavior: criterion is met if there's at least one period
	result := periodCount > 0
	re.Logger.Debug("Time period criteria evaluation result: %v", result)
	return result, nil
}

// evaluatePatternCriteria evaluates patterns in event frequency over time
func (re *RuleEngine) evaluatePatternCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	re.Logger.Debug("Evaluating pattern criteria with %d events", len(events))

	// Parse and validate criteria
	var patternCriteria models.PatternCriteria

	if pattern, ok := criteria["pattern"].(string); ok {
		patternCriteria.Pattern = pattern
		re.Logger.Debug("Pattern type: %s", pattern)
	} else {
		re.Logger.Error("Missing or invalid pattern in pattern criteria")
		return false, fmt.Errorf("missing or invalid pattern in pattern criteria")
	}

	if periodType, ok := criteria["periodType"].(string); ok {
		patternCriteria.PeriodType = periodType
		re.Logger.Debug("Period type: %s", periodType)
	} else {
		re.Logger.Error("Missing or invalid periodType in pattern criteria")
		return false, fmt.Errorf("missing or invalid periodType in pattern criteria")
	}

	if minPeriods, ok := criteria["minPeriods"].(float64); ok {
		patternCriteria.MinPeriods = int(minPeriods)
		re.Logger.Debug("Minimum periods: %d", patternCriteria.MinPeriods)
	} else {
		patternCriteria.MinPeriods = 3 // Default to requiring at least 3 periods
		re.Logger.Debug("Using default minimum periods: %d", patternCriteria.MinPeriods)
	}

	if minIncreasePct, ok := criteria["minIncreasePct"].(float64); ok {
		patternCriteria.MinIncreasePct = minIncreasePct
		re.Logger.Debug("Minimum increase percentage: %.2f%%", patternCriteria.MinIncreasePct)
	} else {
		patternCriteria.MinIncreasePct = 5.0 // Default to 5% increase
		re.Logger.Debug("Using default minimum increase percentage: %.2f%%", patternCriteria.MinIncreasePct)
	}

	if maxDecreasePct, ok := criteria["maxDecreasePct"].(float64); ok {
		patternCriteria.MaxDecreasePct = maxDecreasePct
		re.Logger.Debug("Maximum decrease percentage: %.2f%%", patternCriteria.MaxDecreasePct)
	} else {
		patternCriteria.MaxDecreasePct = 5.0 // Default to 5% decrease
		re.Logger.Debug("Using default maximum decrease percentage: %.2f%%", patternCriteria.MaxDecreasePct)
	}

	if maxDeviation, ok := criteria["maxDeviation"].(float64); ok {
		patternCriteria.MaxDeviation = maxDeviation
		re.Logger.Debug("Maximum deviation: %.2f", patternCriteria.MaxDeviation)
	} else {
		patternCriteria.MaxDeviation = 0.2 // Default to 20% deviation
		re.Logger.Debug("Using default maximum deviation: %.2f", patternCriteria.MaxDeviation)
	}

	// Group events by period
	periodCounts, periods, err := groupEventsByPeriod(events, patternCriteria.PeriodType)
	if err != nil {
		re.Logger.Error("Failed to group events by period: %v", err)
		return false, err
	}

	re.Logger.Debug("Period grouping resulted in %d unique periods", len(periods))

	// Check if we have enough periods
	if len(periods) < patternCriteria.MinPeriods {
		reason := fmt.Sprintf("not enough periods: %d < %d required", len(periods), patternCriteria.MinPeriods)
		metadata["reason"] = reason
		re.Logger.Debug("Pattern criteria not met: %s", reason)
		return false, nil
	}

	// Convert period counts to a slice in chronological order
	counts := make([]int, len(periods))
	for i, period := range periods {
		counts[i] = periodCounts[period]
		re.Logger.Trace("Period %s has %d events", period, periodCounts[period])
	}

	metadata["period_keys"] = periods
	metadata["period_counts"] = counts

	re.Logger.Debug("Period counts: %v", counts)

	// Evaluate the pattern
	var result bool

	switch patternCriteria.Pattern {
	case "consistent":
		re.Logger.Debug("Evaluating consistent pattern")
		critMap := map[string]interface{}{
			"maxDeviation": patternCriteria.MaxDeviation,
		}
		result = evaluateConsistentPattern(counts, critMap, metadata)
	case "increasing":
		re.Logger.Debug("Evaluating increasing pattern")
		critMap := map[string]interface{}{
			"minIncreasePct": patternCriteria.MinIncreasePct,
			"minPeriods":     float64(len(counts)),
		}
		result = evaluateIncreasingPattern(counts, critMap, periods, metadata)
	case "decreasing":
		re.Logger.Debug("Evaluating decreasing pattern")
		critMap := map[string]interface{}{
			"maxDecreasePct": patternCriteria.MaxDecreasePct,
			"minPeriods":     float64(len(counts)),
		}
		result = evaluateDecreasingPattern(counts, critMap, periods, metadata)
	default:
		re.Logger.Error("Unsupported pattern type: %s", patternCriteria.Pattern)
		return false, fmt.Errorf("unsupported pattern type: %s", patternCriteria.Pattern)
	}

	re.Logger.Debug("Pattern criteria evaluation result: %v", result)
	return result, nil
}

// evaluateSequenceCriteria checks if events occur in a specific sequence
func (re *RuleEngine) evaluateSequenceCriteria(criteria map[string]interface{}, userID string, metadata map[string]interface{}) (bool, error) {
	re.Logger.Debug("Evaluating sequence criteria for user %s", userID)

	// Parse and validate criteria
	var sequenceCriteria models.SequenceCriteria

	if sequence, ok := criteria["sequence"].([]interface{}); ok {
		for _, item := range sequence {
			if eventType, ok := item.(string); ok {
				sequenceCriteria.Sequence = append(sequenceCriteria.Sequence, eventType)
			} else {
				re.Logger.Error("Invalid event type in sequence")
				return false, fmt.Errorf("invalid event type in sequence")
			}
		}
		re.Logger.Debug("Sequence: %v", sequenceCriteria.Sequence)
	} else {
		re.Logger.Error("Missing or invalid sequence in sequence criteria")
		return false, fmt.Errorf("missing or invalid sequence in sequence criteria")
	}

	if len(sequenceCriteria.Sequence) == 0 {
		re.Logger.Error("Sequence cannot be empty")
		return false, fmt.Errorf("sequence cannot be empty")
	}

	if maxGapSeconds, ok := criteria["maxGapSeconds"].(float64); ok {
		sequenceCriteria.MaxGapSeconds = int(maxGapSeconds)
		re.Logger.Debug("Maximum gap between sequence events: %d seconds", sequenceCriteria.MaxGapSeconds)
	}

	if requireStrict, ok := criteria["requireStrict"].(bool); ok {
		sequenceCriteria.RequireStrict = requireStrict
		re.Logger.Debug("Require strict sequence: %v", sequenceCriteria.RequireStrict)
	}

	// For each event type in sequence, fetch the events
	re.Logger.Debug("Fetching events for %d event types in sequence", len(sequenceCriteria.Sequence))
	sequenceEvents := make([][]models.Event, len(sequenceCriteria.Sequence))

	for i, eventType := range sequenceCriteria.Sequence {
		eventTypeObj, err := re.DB.GetEventTypeByName(eventType)
		if err != nil {
			re.Logger.Error("Event type '%s' not found: %v", eventType, err)
			return false, fmt.Errorf("event type '%s' not found: %w", eventType, err)
		}

		events, err := re.DB.GetUserEventsByType(userID, eventTypeObj.ID)
		if err != nil {
			re.Logger.Error("Failed to get events for type '%s': %v", eventType, err)
			return false, fmt.Errorf("failed to get events for type '%s': %w", eventType, err)
		}

		// Sort events by timestamp
		sort.Slice(events, func(i, j int) bool {
			return events[i].OccurredAt.Before(events[j].OccurredAt)
		})

		sequenceEvents[i] = events
		re.Logger.Debug("Found %d events of type '%s'", len(events), eventType)

		// Early termination: if any event type has no events, sequence is impossible
		if len(events) == 0 {
			metadata["missing_event_type"] = eventType
			re.Logger.Debug("Sequence criteria not met: missing event type '%s'", eventType)
			return false, nil
		}
	}

	// Check for a valid sequence
	re.Logger.Debug("Checking for valid sequence")
	result, sequenceMetadata := findValidSequence(sequenceEvents, sequenceCriteria.MaxGapSeconds, sequenceCriteria.RequireStrict)

	// Merge sequence-specific metadata
	for k, v := range sequenceMetadata {
		metadata[k] = v
		re.Logger.Debug("Sequence metadata: %s = %v", k, v)
	}

	re.Logger.Debug("Sequence criteria evaluation result: %v", result)
	return result, nil
}

// evaluateGapCriteria checks for gaps in event occurrence
func (re *RuleEngine) evaluateGapCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	re.Logger.Debug("Evaluating gap criteria with %d events", len(events))

	// Parse and validate criteria
	var gapCriteria models.GapCriteria

	if maxGapHours, ok := criteria["maxGapHours"].(float64); ok {
		gapCriteria.MaxGapHours = maxGapHours
		re.Logger.Debug("Maximum gap hours: %.2f", maxGapHours)
	} else {
		re.Logger.Error("Missing or invalid maxGapHours in gap criteria")
		return false, fmt.Errorf("missing or invalid maxGapHours in gap criteria")
	}

	if minGapHours, ok := criteria["minGapHours"].(float64); ok {
		gapCriteria.MinGapHours = minGapHours
		re.Logger.Debug("Minimum gap hours: %.2f", minGapHours)
	}

	if periodType, ok := criteria["periodType"].(string); ok {
		gapCriteria.PeriodType = periodType
		re.Logger.Debug("Period type: %s", periodType)
	} else {
		gapCriteria.PeriodType = "all" // Default to checking all time
		re.Logger.Debug("Using default period type: all")
	}

	if excludeConditions, ok := criteria["excludeConditions"].(map[string]interface{}); ok {
		gapCriteria.ExcludeConditions = excludeConditions
		re.Logger.Debug("Exclude conditions: %v", excludeConditions)
	}

	// Filter events with exclusion conditions if provided
	filteredEvents := events
	if len(gapCriteria.ExcludeConditions) > 0 {
		re.Logger.Debug("Filtering events based on exclusion conditions")
		var err error
		filteredEvents, err = re.filterEvents(gapCriteria.ExcludeConditions, events)
		if err != nil {
			re.Logger.Error("Error filtering events: %v", err)
			return false, err
		}
		re.Logger.Debug("After filtering: %d/%d events remain", len(filteredEvents), len(events))
	}

	// Need at least 2 events to have a gap
	if len(filteredEvents) < 2 {
		metadata["reason"] = "not enough events to check for gaps"
		re.Logger.Debug("Gap criteria not met: not enough events to check for gaps (minimum 2 required)")
		return false, nil
	}

	// Sort events by time
	re.Logger.Debug("Sorting events by timestamp")
	sort.Slice(filteredEvents, func(i, j int) bool {
		return filteredEvents[i].OccurredAt.Before(filteredEvents[j].OccurredAt)
	})

	// Check gaps between consecutive events
	maxGapFound := 0.0
	minGapFound := float64(100 * 365 * 24) // Initialize to a large value (100 years in hours)

	re.Logger.Debug("Analyzing gaps between events")
	for i := 1; i < len(filteredEvents); i++ {
		gapHours := filteredEvents[i].OccurredAt.Sub(filteredEvents[i-1].OccurredAt).Hours()
		re.Logger.Trace("Gap between events %d and %d: %.2f hours",
			filteredEvents[i-1].ID, filteredEvents[i].ID, gapHours)

		if gapHours > maxGapFound {
			maxGapFound = gapHours
			metadata["max_gap_hours"] = maxGapFound
			metadata["max_gap_event_ids"] = []int{filteredEvents[i-1].ID, filteredEvents[i].ID}
			re.Logger.Debug("New maximum gap: %.2f hours between events %d and %d",
				maxGapFound, filteredEvents[i-1].ID, filteredEvents[i].ID)
		}

		if gapHours < minGapFound {
			minGapFound = gapHours
			metadata["min_gap_hours"] = minGapFound
			metadata["min_gap_event_ids"] = []int{filteredEvents[i-1].ID, filteredEvents[i].ID}
			re.Logger.Debug("New minimum gap: %.2f hours between events %d and %d",
				minGapFound, filteredEvents[i-1].ID, filteredEvents[i].ID)
		}
	}

	// Check if gaps meet the criteria
	if gapCriteria.MaxGapHours > 0 && maxGapFound > gapCriteria.MaxGapHours {
		metadata["gap_criterion_failed"] = "max_gap"
		re.Logger.Debug("Gap criterion not met: maximum gap %.2f hours exceeds limit of %.2f hours",
			maxGapFound, gapCriteria.MaxGapHours)
		return false, nil
	}

	if gapCriteria.MinGapHours > 0 && minGapFound < gapCriteria.MinGapHours {
		metadata["gap_criterion_failed"] = "min_gap"
		re.Logger.Debug("Gap criterion not met: minimum gap %.2f hours is below requirement of %.2f hours",
			minGapFound, gapCriteria.MinGapHours)
		return false, nil
	}

	re.Logger.Debug("Gap criteria met")
	return true, nil
}

// evaluateDurationCriteria assesses time duration between related events
func (re *RuleEngine) evaluateDurationCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	re.Logger.Debug("Evaluating duration criteria with %d events", len(events))

	// Parse and validate criteria
	var durationCriteria models.DurationCriteria

	if startEvent, ok := criteria["startEvent"].(map[string]interface{}); ok {
		durationCriteria.StartEvent = startEvent
		re.Logger.Debug("Start event criteria: %v", startEvent)
	} else {
		re.Logger.Error("Missing or invalid startEvent in duration criteria")
		return false, fmt.Errorf("missing or invalid startEvent in duration criteria")
	}

	if endEvent, ok := criteria["endEvent"].(map[string]interface{}); ok {
		durationCriteria.EndEvent = endEvent
		re.Logger.Debug("End event criteria: %v", endEvent)
	} else {
		re.Logger.Error("Missing or invalid endEvent in duration criteria")
		return false, fmt.Errorf("missing or invalid endEvent in duration criteria")
	}

	if duration, ok := criteria["duration"].(map[string]interface{}); ok {
		durationCriteria.Duration = duration
		re.Logger.Debug("Duration criteria: %v", duration)
	}

	if unit, ok := criteria["unit"].(string); ok {
		durationCriteria.Unit = unit
		re.Logger.Debug("Time unit: %s", unit)
	} else {
		durationCriteria.Unit = "hours" // Default to hours
		re.Logger.Debug("Using default time unit: hours")
	}

	// Filter events to find start and end events
	re.Logger.Debug("Filtering events to find start events")
	startEvents, err := re.filterEvents(durationCriteria.StartEvent, events)
	if err != nil {
		re.Logger.Error("Error filtering start events: %v", err)
		return false, err
	}

	if len(startEvents) == 0 {
		metadata["reason"] = "no matching start events"
		re.Logger.Debug("Duration criteria not met: no matching start events")
		return false, nil
	}
	re.Logger.Debug("Found %d matching start events", len(startEvents))

	re.Logger.Debug("Filtering events to find end events")
	endEvents, err := re.filterEvents(durationCriteria.EndEvent, events)
	if err != nil {
		re.Logger.Error("Error filtering end events: %v", err)
		return false, err
	}

	if len(endEvents) == 0 {
		metadata["reason"] = "no matching end events"
		re.Logger.Debug("Duration criteria not met: no matching end events")
		return false, nil
	}
	re.Logger.Debug("Found %d matching end events", len(endEvents))

	// Sort events by time
	sort.Slice(startEvents, func(i, j int) bool {
		return startEvents[i].OccurredAt.Before(startEvents[j].OccurredAt)
	})

	sort.Slice(endEvents, func(i, j int) bool {
		return endEvents[i].OccurredAt.Before(endEvents[j].OccurredAt)
	})

	// Find valid start-end pairs
	re.Logger.Debug("Finding valid start-end pairs")
	var validPairs []struct {
		start models.Event
		end   models.Event
	}

	for _, startEvent := range startEvents {
		for _, endEvent := range endEvents {
			// End event must occur after start event
			if endEvent.OccurredAt.After(startEvent.OccurredAt) {
				validPairs = append(validPairs, struct {
					start models.Event
					end   models.Event
				}{startEvent, endEvent})
				re.Logger.Trace("Found valid pair: start event %d -> end event %d",
					startEvent.ID, endEvent.ID)
			}
		}
	}

	if len(validPairs) == 0 {
		metadata["reason"] = "no valid start-end event pairs found"
		re.Logger.Debug("Duration criteria not met: no valid start-end event pairs found")
		return false, nil
	}
	re.Logger.Debug("Found %d valid start-end pairs", len(validPairs))

	// Calculate durations for all valid pairs
	durations := make([]float64, len(validPairs))
	for i, pair := range validPairs {
		durations[i] = calculateDuration(pair.start.OccurredAt, pair.end.OccurredAt, durationCriteria.Unit)
		re.Logger.Trace("Pair %d duration: %.2f %s", i, durations[i], durationCriteria.Unit)
	}

	// Find shortest and longest durations
	shortestDuration := durations[0]
	longestDuration := durations[0]
	shortestPairIndex := 0
	longestPairIndex := 0

	for i, duration := range durations {
		if duration < shortestDuration {
			shortestDuration = duration
			shortestPairIndex = i
		}
		if duration > longestDuration {
			longestDuration = duration
			longestPairIndex = i
		}
	}

	metadata["shortest_duration"] = shortestDuration
	metadata["shortest_duration_start_event_id"] = validPairs[shortestPairIndex].start.ID
	metadata["shortest_duration_end_event_id"] = validPairs[shortestPairIndex].end.ID
	re.Logger.Debug("Shortest duration: %.2f %s between events %d and %d",
		shortestDuration, durationCriteria.Unit,
		validPairs[shortestPairIndex].start.ID, validPairs[shortestPairIndex].end.ID)

	metadata["longest_duration"] = longestDuration
	metadata["longest_duration_start_event_id"] = validPairs[longestPairIndex].start.ID
	metadata["longest_duration_end_event_id"] = validPairs[longestPairIndex].end.ID
	re.Logger.Debug("Longest duration: %.2f %s between events %d and %d",
		longestDuration, durationCriteria.Unit,
		validPairs[longestPairIndex].start.ID, validPairs[longestPairIndex].end.ID)

	// If there's a duration criteria, evaluate it (using the shortest duration as this is typically the goal)
	if len(durationCriteria.Duration) > 0 {
		re.Logger.Debug("Evaluating duration criteria against shortest duration: %.2f %s",
			shortestDuration, durationCriteria.Unit)
		result, err := re.evaluateNumericCriteria(shortestDuration, durationCriteria.Duration)
		if err != nil {
			re.Logger.Error("Error evaluating numeric criteria: %v", err)
			return false, err
		}
		re.Logger.Debug("Duration criteria evaluation result: %v", result)
		return result, nil
	}

	// Default behavior: criterion is met if there's at least one valid pair
	re.Logger.Debug("Duration criteria met (default behavior with at least one valid pair)")
	return len(validPairs) > 0, nil
}

// evaluateAggregationCriteria handles min, max, avg calculations
func (re *RuleEngine) evaluateAggregationCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	re.Logger.Debug("Evaluating aggregation criteria with %d events", len(events))

	// Parse and validate criteria
	var aggregationCriteria models.AggregationCriteria

	if aggType, ok := criteria["type"].(string); ok {
		aggregationCriteria.Type = aggType
		re.Logger.Debug("Aggregation type: %s", aggType)
	} else {
		re.Logger.Error("Missing or invalid type in aggregation criteria")
		return false, fmt.Errorf("missing or invalid type in aggregation criteria")
	}

	if field, ok := criteria["field"].(string); ok {
		aggregationCriteria.Field = field
		re.Logger.Debug("Field to aggregate: %s", field)
	} else {
		re.Logger.Error("Missing or invalid field in aggregation criteria")
		return false, fmt.Errorf("missing or invalid field in aggregation criteria")
	}

	if value, ok := criteria["value"].(map[string]interface{}); ok {
		aggregationCriteria.Value = value
		re.Logger.Debug("Value criteria: %v", value)
	}

	if timeWindow, ok := criteria["timeWindow"].(map[string]interface{}); ok {
		aggregationCriteria.TimeWindow = timeWindow
		re.Logger.Debug("Time window: %v", timeWindow)

		// Filter events by time window if specified
		re.Logger.Debug("Filtering events by time window")
		filteredEvents, err := filterEventsByTimeWindowWithCache(events, timeWindow, re.TimeVarCache)
		if err != nil {
			re.Logger.Error("Error filtering events by time window: %v", err)
			return false, err
		}
		events = filteredEvents
		re.Logger.Debug("After time window filtering: %d events remain", len(events))
	}

	// Extract values for the specified field from events
	var values []float64

	for _, event := range events {
		if event.Payload == nil {
			re.Logger.Trace("Event ID %d has no payload, skipping", event.ID)
			continue
		}

		fieldValue, exists := event.Payload[aggregationCriteria.Field]
		if !exists {
			re.Logger.Trace("Event ID %d does not have field '%s', skipping",
				event.ID, aggregationCriteria.Field)
			continue
		}

		numValue, err := toFloat64(fieldValue)
		if err != nil {
			re.Logger.Trace("Event ID %d field '%s' value '%v' could not be converted to float: %v",
				event.ID, aggregationCriteria.Field, fieldValue, err)
			continue
		}

		values = append(values, numValue)
		re.Logger.Trace("Event ID %d field '%s' value: %v", event.ID, aggregationCriteria.Field, numValue)
	}

	if len(values) == 0 {
		reason := fmt.Sprintf("no valid values found for field '%s'", aggregationCriteria.Field)
		metadata["reason"] = reason
		re.Logger.Debug("Aggregation criteria not met: %s", reason)
		return false, nil
	}
	re.Logger.Debug("Found %d valid values for field '%s'", len(values), aggregationCriteria.Field)

	// Calculate the aggregation
	var result float64

	switch aggregationCriteria.Type {
	case "min":
		result = values[0]
		for _, v := range values {
			if v < result {
				result = v
			}
		}
		re.Logger.Debug("Minimum value: %.2f", result)
	case "max":
		result = values[0]
		for _, v := range values {
			if v > result {
				result = v
			}
		}
		re.Logger.Debug("Maximum value: %.2f", result)
	case "avg":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		result = sum / float64(len(values))
		re.Logger.Debug("Average value: %.2f", result)
	case "sum":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		result = sum
		re.Logger.Debug("Sum value: %.2f", result)
	case "count":
		result = float64(len(values))
		re.Logger.Debug("Count value: %.0f", result)
	default:
		re.Logger.Error("Unsupported aggregation type: %s", aggregationCriteria.Type)
		return false, fmt.Errorf("unsupported aggregation type: %s", aggregationCriteria.Type)
	}

	metadata[fmt.Sprintf("%s_%s", aggregationCriteria.Type, aggregationCriteria.Field)] = result

	// If there's a value criteria, evaluate it
	if len(aggregationCriteria.Value) > 0 {
		re.Logger.Debug("Evaluating value criteria against aggregation result: %.2f", result)
		evalResult, err := re.evaluateNumericCriteria(result, aggregationCriteria.Value)
		if err != nil {
			re.Logger.Error("Error evaluating numeric criteria: %v", err)
			return false, err
		}
		re.Logger.Debug("Aggregation criteria evaluation result: %v", evalResult)
		return evalResult, nil
	}

	// Default behavior: criterion is met if there's at least one value
	re.Logger.Debug("Aggregation criteria met (default behavior with at least one value)")
	return len(values) > 0, nil
}

// evaluateNumericCriteria evaluates a numeric value against criteria operators
func (re *RuleEngine) evaluateNumericCriteria(value float64, criteria map[string]interface{}) (bool, error) {
	re.Logger.Debug("Evaluating numeric criteria against value: %.2f", value)

	for operator, compareValue := range criteria {
		compareFloat, err := toFloat64(compareValue)
		if err != nil {
			re.Logger.Error("Invalid comparison value for numeric criteria: %v", err)
			return false, fmt.Errorf("invalid comparison value for numeric criteria: %w", err)
		}

		re.Logger.Debug("Checking %.2f %s %.2f", value, operator, compareFloat)

		switch operator {
		case "$eq":
			if value != compareFloat {
				re.Logger.Debug("Numeric criterion not met: %.2f != %.2f", value, compareFloat)
				return false, nil
			}
		case "$ne":
			if value == compareFloat {
				re.Logger.Debug("Numeric criterion not met: %.2f == %.2f", value, compareFloat)
				return false, nil
			}
		case "$gt":
			if value <= compareFloat {
				re.Logger.Debug("Numeric criterion not met: %.2f <= %.2f", value, compareFloat)
				return false, nil
			}
		case "$gte":
			if value < compareFloat {
				re.Logger.Debug("Numeric criterion not met: %.2f < %.2f", value, compareFloat)
				return false, nil
			}
		case "$lt":
			if value >= compareFloat {
				re.Logger.Debug("Numeric criterion not met: %.2f >= %.2f", value, compareFloat)
				return false, nil
			}
		case "$lte":
			if value > compareFloat {
				re.Logger.Debug("Numeric criterion not met: %.2f > %.2f", value, compareFloat)
				return false, nil
			}
		default:
			re.Logger.Error("Unsupported numeric operator: %s", operator)
			return false, fmt.Errorf("unsupported numeric operator: %s", operator)
		}
	}

	re.Logger.Debug("Numeric criteria met")
	return true, nil
}
