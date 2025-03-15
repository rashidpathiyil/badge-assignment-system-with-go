package engine

import (
	"fmt"
	"sort"

	"github.com/badge-assignment-system/internal/models"
)

// evaluateTimePeriodCriteria handles time period counts (days, weeks, months)
func (re *RuleEngine) evaluateTimePeriodCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	// Parse and validate criteria
	var timePeriodCriteria models.TimePeriodCriteria

	if periodType, ok := criteria["periodType"].(string); ok {
		timePeriodCriteria.PeriodType = periodType
	} else {
		return false, fmt.Errorf("missing or invalid periodType in timePeriod criteria")
	}

	if count, ok := criteria["count"].(map[string]interface{}); ok {
		timePeriodCriteria.Count = count
	}

	if excludeWeekends, ok := criteria["excludeWeekends"].(bool); ok {
		timePeriodCriteria.ExcludeWeekends = excludeWeekends
	}

	if excludeHolidays, ok := criteria["excludeHolidays"].(bool); ok {
		timePeriodCriteria.ExcludeHolidays = excludeHolidays
	}

	if holidays, ok := criteria["holidays"].([]interface{}); ok {
		for _, holiday := range holidays {
			if holidayStr, ok := holiday.(string); ok {
				timePeriodCriteria.Holidays = append(timePeriodCriteria.Holidays, holidayStr)
			}
		}
	}

	// Filter events based on exclusions
	var filteredEvents []models.Event
	for _, event := range events {
		// Check weekend exclusion
		if timePeriodCriteria.ExcludeWeekends && isWeekend(event.OccurredAt) {
			continue
		}

		// Check holiday exclusion
		if timePeriodCriteria.ExcludeHolidays && isHoliday(event.OccurredAt, timePeriodCriteria.Holidays) {
			continue
		}

		filteredEvents = append(filteredEvents, event)
	}

	// Group events by the specified time period
	uniquePeriods := make(map[string]bool)

	for _, event := range filteredEvents {
		periodKey, err := getPeriodKey(event.OccurredAt, timePeriodCriteria.PeriodType)
		if err != nil {
			return false, err
		}
		uniquePeriods[periodKey] = true
	}

	// Count unique periods
	periodCount := len(uniquePeriods)
	metadata["unique_period_count"] = periodCount

	// If there's a count criteria, evaluate it
	if len(timePeriodCriteria.Count) > 0 {
		return re.evaluateNumericCriteria(float64(periodCount), timePeriodCriteria.Count)
	}

	// Default behavior: criterion is met if there's at least one period
	return periodCount > 0, nil
}

// evaluatePatternCriteria evaluates patterns in event frequency over time
func (re *RuleEngine) evaluatePatternCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	// Parse and validate criteria
	var patternCriteria models.PatternCriteria

	if pattern, ok := criteria["pattern"].(string); ok {
		patternCriteria.Pattern = pattern
	} else {
		return false, fmt.Errorf("missing or invalid pattern in pattern criteria")
	}

	if periodType, ok := criteria["periodType"].(string); ok {
		patternCriteria.PeriodType = periodType
	} else {
		return false, fmt.Errorf("missing or invalid periodType in pattern criteria")
	}

	if minPeriods, ok := criteria["minPeriods"].(float64); ok {
		patternCriteria.MinPeriods = int(minPeriods)
	} else {
		patternCriteria.MinPeriods = 3 // Default to requiring at least 3 periods
	}

	if minIncreasePct, ok := criteria["minIncreasePct"].(float64); ok {
		patternCriteria.MinIncreasePct = minIncreasePct
	} else {
		patternCriteria.MinIncreasePct = 5.0 // Default to 5% increase
	}

	if maxDecreasePct, ok := criteria["maxDecreasePct"].(float64); ok {
		patternCriteria.MaxDecreasePct = maxDecreasePct
	} else {
		patternCriteria.MaxDecreasePct = 5.0 // Default to 5% decrease
	}

	if maxDeviation, ok := criteria["maxDeviation"].(float64); ok {
		patternCriteria.MaxDeviation = maxDeviation
	} else {
		patternCriteria.MaxDeviation = 0.2 // Default to 20% deviation
	}

	// Group events by period
	periodCounts, periods, err := groupEventsByPeriod(events, patternCriteria.PeriodType)
	if err != nil {
		return false, err
	}

	// Check if we have enough periods
	if len(periods) < patternCriteria.MinPeriods {
		metadata["reason"] = fmt.Sprintf("not enough periods: %d < %d required", len(periods), patternCriteria.MinPeriods)
		return false, nil
	}

	// Convert period counts to a slice in chronological order
	counts := make([]int, len(periods))
	for i, period := range periods {
		counts[i] = periodCounts[period]
	}

	metadata["period_keys"] = periods
	metadata["period_counts"] = counts

	// Evaluate the pattern
	var result bool
	var patternMetadata map[string]interface{}

	switch patternCriteria.Pattern {
	case "consistent":
		result, patternMetadata = evaluateConsistentPattern(counts, patternCriteria.MaxDeviation)
	case "increasing":
		result, patternMetadata = evaluateIncreasingPattern(counts, patternCriteria.MinIncreasePct)
	case "decreasing":
		result, patternMetadata = evaluateDecreasingPattern(counts, patternCriteria.MaxDecreasePct)
	default:
		return false, fmt.Errorf("unsupported pattern type: %s", patternCriteria.Pattern)
	}

	// Merge pattern-specific metadata
	for k, v := range patternMetadata {
		metadata[k] = v
	}

	return result, nil
}

// evaluateSequenceCriteria checks if events occur in a specific sequence
func (re *RuleEngine) evaluateSequenceCriteria(criteria map[string]interface{}, userID string, metadata map[string]interface{}) (bool, error) {
	// Parse and validate criteria
	var sequenceCriteria models.SequenceCriteria

	if sequence, ok := criteria["sequence"].([]interface{}); ok {
		for _, item := range sequence {
			if eventType, ok := item.(string); ok {
				sequenceCriteria.Sequence = append(sequenceCriteria.Sequence, eventType)
			} else {
				return false, fmt.Errorf("invalid event type in sequence")
			}
		}
	} else {
		return false, fmt.Errorf("missing or invalid sequence in sequence criteria")
	}

	if len(sequenceCriteria.Sequence) == 0 {
		return false, fmt.Errorf("sequence cannot be empty")
	}

	if maxGapSeconds, ok := criteria["maxGapSeconds"].(float64); ok {
		sequenceCriteria.MaxGapSeconds = int(maxGapSeconds)
	}

	if requireStrict, ok := criteria["requireStrict"].(bool); ok {
		sequenceCriteria.RequireStrict = requireStrict
	}

	// For each event type in sequence, fetch the events
	sequenceEvents := make([][]models.Event, len(sequenceCriteria.Sequence))

	for i, eventType := range sequenceCriteria.Sequence {
		eventTypeObj, err := re.DB.GetEventTypeByName(eventType)
		if err != nil {
			return false, fmt.Errorf("event type '%s' not found: %w", eventType, err)
		}

		events, err := re.DB.GetUserEventsByType(userID, eventTypeObj.ID)
		if err != nil {
			return false, fmt.Errorf("failed to get events for type '%s': %w", eventType, err)
		}

		// Sort events by timestamp
		sort.Slice(events, func(i, j int) bool {
			return events[i].OccurredAt.Before(events[j].OccurredAt)
		})

		sequenceEvents[i] = events

		// Early termination: if any event type has no events, sequence is impossible
		if len(events) == 0 {
			metadata["missing_event_type"] = eventType
			return false, nil
		}
	}

	// Check for a valid sequence
	result, sequenceMetadata := findValidSequence(sequenceEvents, sequenceCriteria.MaxGapSeconds, sequenceCriteria.RequireStrict)

	// Merge sequence-specific metadata
	for k, v := range sequenceMetadata {
		metadata[k] = v
	}

	return result, nil
}

// evaluateGapCriteria checks for gaps in event occurrence
func (re *RuleEngine) evaluateGapCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	// Parse and validate criteria
	var gapCriteria models.GapCriteria

	if maxGapHours, ok := criteria["maxGapHours"].(float64); ok {
		gapCriteria.MaxGapHours = maxGapHours
	} else {
		return false, fmt.Errorf("missing or invalid maxGapHours in gap criteria")
	}

	if minGapHours, ok := criteria["minGapHours"].(float64); ok {
		gapCriteria.MinGapHours = minGapHours
	}

	if periodType, ok := criteria["periodType"].(string); ok {
		gapCriteria.PeriodType = periodType
	} else {
		gapCriteria.PeriodType = "all" // Default to checking all time
	}

	if excludeConditions, ok := criteria["excludeConditions"].(map[string]interface{}); ok {
		gapCriteria.ExcludeConditions = excludeConditions
	}

	// Filter events with exclusion conditions if provided
	filteredEvents := events
	if len(gapCriteria.ExcludeConditions) > 0 {
		var err error
		filteredEvents, err = re.filterEvents(gapCriteria.ExcludeConditions, events)
		if err != nil {
			return false, err
		}
	}

	// Need at least 2 events to have a gap
	if len(filteredEvents) < 2 {
		metadata["reason"] = "not enough events to check for gaps"
		return false, nil
	}

	// Sort events by time
	sort.Slice(filteredEvents, func(i, j int) bool {
		return filteredEvents[i].OccurredAt.Before(filteredEvents[j].OccurredAt)
	})

	// Check gaps between consecutive events
	maxGapFound := 0.0
	minGapFound := float64(100 * 365 * 24) // Initialize to a large value (100 years in hours)

	for i := 1; i < len(filteredEvents); i++ {
		gapHours := filteredEvents[i].OccurredAt.Sub(filteredEvents[i-1].OccurredAt).Hours()

		if gapHours > maxGapFound {
			maxGapFound = gapHours
			metadata["max_gap_hours"] = maxGapFound
			metadata["max_gap_event_ids"] = []int{filteredEvents[i-1].ID, filteredEvents[i].ID}
		}

		if gapHours < minGapFound {
			minGapFound = gapHours
			metadata["min_gap_hours"] = minGapFound
			metadata["min_gap_event_ids"] = []int{filteredEvents[i-1].ID, filteredEvents[i].ID}
		}
	}

	// Check if gaps meet the criteria
	if gapCriteria.MaxGapHours > 0 && maxGapFound > gapCriteria.MaxGapHours {
		metadata["gap_criterion_failed"] = "max_gap"
		return false, nil
	}

	if gapCriteria.MinGapHours > 0 && minGapFound < gapCriteria.MinGapHours {
		metadata["gap_criterion_failed"] = "min_gap"
		return false, nil
	}

	return true, nil
}

// evaluateDurationCriteria assesses time duration between related events
func (re *RuleEngine) evaluateDurationCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	// Parse and validate criteria
	var durationCriteria models.DurationCriteria

	if startEvent, ok := criteria["startEvent"].(map[string]interface{}); ok {
		durationCriteria.StartEvent = startEvent
	} else {
		return false, fmt.Errorf("missing or invalid startEvent in duration criteria")
	}

	if endEvent, ok := criteria["endEvent"].(map[string]interface{}); ok {
		durationCriteria.EndEvent = endEvent
	} else {
		return false, fmt.Errorf("missing or invalid endEvent in duration criteria")
	}

	if duration, ok := criteria["duration"].(map[string]interface{}); ok {
		durationCriteria.Duration = duration
	}

	if unit, ok := criteria["unit"].(string); ok {
		durationCriteria.Unit = unit
	} else {
		durationCriteria.Unit = "hours" // Default to hours
	}

	// Filter events to find start and end events
	startEvents, err := re.filterEvents(durationCriteria.StartEvent, events)
	if err != nil {
		return false, err
	}

	if len(startEvents) == 0 {
		metadata["reason"] = "no matching start events"
		return false, nil
	}

	endEvents, err := re.filterEvents(durationCriteria.EndEvent, events)
	if err != nil {
		return false, err
	}

	if len(endEvents) == 0 {
		metadata["reason"] = "no matching end events"
		return false, nil
	}

	// Sort events by time
	sort.Slice(startEvents, func(i, j int) bool {
		return startEvents[i].OccurredAt.Before(startEvents[j].OccurredAt)
	})

	sort.Slice(endEvents, func(i, j int) bool {
		return endEvents[i].OccurredAt.Before(endEvents[j].OccurredAt)
	})

	// Find valid start-end pairs
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
			}
		}
	}

	if len(validPairs) == 0 {
		metadata["reason"] = "no valid start-end event pairs found"
		return false, nil
	}

	// Calculate durations for all valid pairs
	durations := make([]float64, len(validPairs))
	for i, pair := range validPairs {
		durations[i] = calculateDuration(pair.start.OccurredAt, pair.end.OccurredAt, durationCriteria.Unit)
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

	metadata["longest_duration"] = longestDuration
	metadata["longest_duration_start_event_id"] = validPairs[longestPairIndex].start.ID
	metadata["longest_duration_end_event_id"] = validPairs[longestPairIndex].end.ID

	// If there's a duration criteria, evaluate it (using the shortest duration as this is typically the goal)
	if len(durationCriteria.Duration) > 0 {
		return re.evaluateNumericCriteria(shortestDuration, durationCriteria.Duration)
	}

	// Default behavior: criterion is met if there's at least one valid pair
	return len(validPairs) > 0, nil
}

// evaluateAggregationCriteria handles min, max, avg calculations
func (re *RuleEngine) evaluateAggregationCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	// Parse and validate criteria
	var aggregationCriteria models.AggregationCriteria

	if aggType, ok := criteria["type"].(string); ok {
		aggregationCriteria.Type = aggType
	} else {
		return false, fmt.Errorf("missing or invalid type in aggregation criteria")
	}

	if field, ok := criteria["field"].(string); ok {
		aggregationCriteria.Field = field
	} else {
		return false, fmt.Errorf("missing or invalid field in aggregation criteria")
	}

	if value, ok := criteria["value"].(map[string]interface{}); ok {
		aggregationCriteria.Value = value
	}

	if timeWindow, ok := criteria["timeWindow"].(map[string]interface{}); ok {
		aggregationCriteria.TimeWindow = timeWindow

		// Filter events by time window if specified
		filteredEvents, err := filterEventsByTimeWindow(events, timeWindow)
		if err != nil {
			return false, err
		}
		events = filteredEvents
	}

	// Extract values for the specified field from events
	var values []float64

	for _, event := range events {
		if event.Payload == nil {
			continue
		}

		fieldValue, exists := event.Payload[aggregationCriteria.Field]
		if !exists {
			continue
		}

		numValue, err := toFloat64(fieldValue)
		if err != nil {
			continue
		}

		values = append(values, numValue)
	}

	if len(values) == 0 {
		metadata["reason"] = fmt.Sprintf("no valid values found for field '%s'", aggregationCriteria.Field)
		return false, nil
	}

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
	case "max":
		result = values[0]
		for _, v := range values {
			if v > result {
				result = v
			}
		}
	case "avg":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		result = sum / float64(len(values))
	case "sum":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		result = sum
	case "count":
		result = float64(len(values))
	default:
		return false, fmt.Errorf("unsupported aggregation type: %s", aggregationCriteria.Type)
	}

	metadata[fmt.Sprintf("%s_%s", aggregationCriteria.Type, aggregationCriteria.Field)] = result

	// If there's a value criteria, evaluate it
	if len(aggregationCriteria.Value) > 0 {
		return re.evaluateNumericCriteria(result, aggregationCriteria.Value)
	}

	// Default behavior: criterion is met if there's at least one value
	return len(values) > 0, nil
}

// evaluateNumericCriteria evaluates a numeric value against criteria operators
func (re *RuleEngine) evaluateNumericCriteria(value float64, criteria map[string]interface{}) (bool, error) {
	for operator, compareValue := range criteria {
		compareFloat, err := toFloat64(compareValue)
		if err != nil {
			return false, fmt.Errorf("invalid comparison value for numeric criteria: %w", err)
		}

		switch operator {
		case "$eq":
			if value != compareFloat {
				return false, nil
			}
		case "$ne":
			if value == compareFloat {
				return false, nil
			}
		case "$gt":
			if value <= compareFloat {
				return false, nil
			}
		case "$gte":
			if value < compareFloat {
				return false, nil
			}
		case "$lt":
			if value >= compareFloat {
				return false, nil
			}
		case "$lte":
			if value > compareFloat {
				return false, nil
			}
		default:
			return false, fmt.Errorf("unsupported numeric operator: %s", operator)
		}
	}

	return true, nil
}
