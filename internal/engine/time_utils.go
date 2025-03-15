package engine

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/badge-assignment-system/internal/models"
)

// isWeekend checks if a given time is a weekend (Saturday or Sunday)
func isWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// isHoliday checks if a given date is in the list of holidays
func isHoliday(t time.Time, holidays []string) bool {
	dateStr := t.Format("2006-01-02")
	for _, holiday := range holidays {
		if holiday == dateStr {
			return true
		}
	}
	return false
}

// getPeriodKey returns a string key for grouping events by time period
func getPeriodKey(t time.Time, periodType string) (string, error) {
	switch periodType {
	case "day":
		return t.Format("2006-01-02"), nil
	case "week":
		year, week := t.ISOWeek()
		return fmt.Sprintf("%d-W%d", year, week), nil
	case "month":
		return t.Format("2006-01"), nil
	case "quarter":
		quarter := (int(t.Month()) + 2) / 3
		return fmt.Sprintf("%d-Q%d", t.Year(), quarter), nil
	case "year":
		return t.Format("2006"), nil
	default:
		return "", fmt.Errorf("unsupported period type: %s", periodType)
	}
}

// groupEventsByPeriod groups events by time period and returns a map of period keys to event counts
func groupEventsByPeriod(events []models.Event, periodType string) (map[string]int, []string, error) {
	periodCounts := make(map[string]int)
	var periods []string

	for _, event := range events {
		periodKey, err := getPeriodKey(event.OccurredAt, periodType)
		if err != nil {
			return nil, nil, err
		}

		if _, exists := periodCounts[periodKey]; !exists {
			periods = append(periods, periodKey)
		}
		periodCounts[periodKey]++
	}

	// Sort periods chronologically
	sort.Strings(periods)

	return periodCounts, periods, nil
}

// evaluateConsistentPattern checks if event counts follow a consistent pattern
func evaluateConsistentPattern(counts []int, maxDeviation float64) (bool, map[string]interface{}) {
	if len(counts) == 0 {
		return false, nil
	}

	// Calculate average
	sum := 0
	for _, count := range counts {
		sum += count
	}
	avg := float64(sum) / float64(len(counts))

	// Calculate maximum deviation from average
	maxDev := 0.0
	for _, count := range counts {
		dev := math.Abs(float64(count)-avg) / avg
		if dev > maxDev {
			maxDev = dev
		}
	}

	result := maxDev <= maxDeviation
	metadata := map[string]interface{}{
		"average":       avg,
		"max_deviation": maxDev,
		"is_consistent": result,
	}

	return result, metadata
}

// evaluateIncreasingPattern checks if event counts follow an increasing pattern
func evaluateIncreasingPattern(counts []int, minIncreasePct float64) (bool, map[string]interface{}) {
	if len(counts) < 2 {
		return false, nil
	}

	// Calculate percentage increase between consecutive periods
	increases := 0
	totalPctIncrease := 0.0

	for i := 1; i < len(counts); i++ {
		if counts[i] > counts[i-1] {
			increases++
			pctIncrease := (float64(counts[i]) - float64(counts[i-1])) / float64(counts[i-1]) * 100
			totalPctIncrease += pctIncrease
		}
	}

	avgPctIncrease := 0.0
	if increases > 0 {
		avgPctIncrease = totalPctIncrease / float64(increases)
	}

	increasingRatio := float64(increases) / float64(len(counts)-1)
	result := increasingRatio >= 0.75 && avgPctIncrease >= minIncreasePct

	metadata := map[string]interface{}{
		"increasing_periods_ratio": increasingRatio,
		"average_percent_increase": avgPctIncrease,
		"is_increasing":            result,
	}

	return result, metadata
}

// evaluateDecreasingPattern checks if event counts follow a decreasing pattern
func evaluateDecreasingPattern(counts []int, maxDecreasePct float64) (bool, map[string]interface{}) {
	if len(counts) < 2 {
		return false, nil
	}

	// Calculate percentage decrease between consecutive periods
	decreases := 0
	totalPctDecrease := 0.0

	for i := 1; i < len(counts); i++ {
		if counts[i] < counts[i-1] {
			decreases++
			pctDecrease := (float64(counts[i-1]) - float64(counts[i])) / float64(counts[i-1]) * 100
			totalPctDecrease += pctDecrease
		}
	}

	avgPctDecrease := 0.0
	if decreases > 0 {
		avgPctDecrease = totalPctDecrease / float64(decreases)
	}

	decreasingRatio := float64(decreases) / float64(len(counts)-1)
	result := decreasingRatio >= 0.75 && avgPctDecrease <= maxDecreasePct

	metadata := map[string]interface{}{
		"decreasing_periods_ratio": decreasingRatio,
		"average_percent_decrease": avgPctDecrease,
		"is_decreasing":            result,
	}

	return result, metadata
}

// findValidSequence checks if there is a valid sequence of events matching the criteria
func findValidSequence(sequenceEvents [][]models.Event, maxGapSeconds int, requireStrict bool) (bool, map[string]interface{}) {
	if len(sequenceEvents) == 0 {
		return false, nil
	}

	// For each event in the first group, try to find a valid sequence
	for _, firstEvent := range sequenceEvents[0] {
		metadata := map[string]interface{}{
			"start_event_id": firstEvent.ID,
			"sequence_found": false,
		}

		current := firstEvent
		valid := true
		var sequenceIds []int
		sequenceIds = append(sequenceIds, current.ID)

		// Try to find matching events for each subsequent group
		for i := 1; i < len(sequenceEvents); i++ {
			found := false

			for _, nextEvent := range sequenceEvents[i] {
				// Check if this event happened after the current one
				if nextEvent.OccurredAt.After(current.OccurredAt) {
					// Check maximum gap if specified
					if maxGapSeconds > 0 {
						gapSeconds := int(nextEvent.OccurredAt.Sub(current.OccurredAt).Seconds())
						if gapSeconds > maxGapSeconds {
							continue // Gap too large
						}
					}

					// This event is a match
					current = nextEvent
					sequenceIds = append(sequenceIds, current.ID)
					found = true
					break
				}
			}

			if !found {
				valid = false
				break
			}
		}

		if valid {
			metadata["sequence_found"] = true
			metadata["sequence_event_ids"] = sequenceIds
			metadata["end_event_id"] = current.ID
			return true, metadata
		}
	}

	return false, nil
}

// parseTimeWindow parses time window criteria and returns start and end time
func parseTimeWindow(criteria map[string]interface{}) (time.Time, time.Time, error) {
	var startTime, endTime time.Time
	var err error

	// Check for explicit start/end times
	if startStr, ok := criteria["start"].(string); ok {
		startTime, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return startTime, endTime, fmt.Errorf("invalid start time format: %w", err)
		}
	}

	if endStr, ok := criteria["end"].(string); ok {
		endTime, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return startTime, endTime, fmt.Errorf("invalid end time format: %w", err)
		}
	}

	// Check for relative time window (e.g., "last 30d")
	if lastStr, ok := criteria["last"].(string); ok {
		endTime = time.Now()

		// Parse duration (e.g., "30d", "2w", "1m")
		re := regexp.MustCompile(`^(\d+)([dwmqy])$`)
		matches := re.FindStringSubmatch(lastStr)

		if len(matches) != 3 {
			return startTime, endTime, fmt.Errorf("invalid duration format: %s", lastStr)
		}

		value, _ := strconv.Atoi(matches[1])
		unit := matches[2]

		switch unit {
		case "d": // days
			startTime = endTime.AddDate(0, 0, -value)
		case "w": // weeks
			startTime = endTime.AddDate(0, 0, -value*7)
		case "m": // months
			startTime = endTime.AddDate(0, -value, 0)
		case "q": // quarters
			startTime = endTime.AddDate(0, -value*3, 0)
		case "y": // years
			startTime = endTime.AddDate(-value, 0, 0)
		default:
			return startTime, endTime, fmt.Errorf("unsupported duration unit: %s", unit)
		}
	}

	// If no time window specified, assume all time
	if startTime.IsZero() {
		startTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	if endTime.IsZero() {
		endTime = time.Now()
	}

	return startTime, endTime, nil
}

// filterEventsByTimeWindow filters events based on a time window
func filterEventsByTimeWindow(events []models.Event, timeWindow map[string]interface{}) ([]models.Event, error) {
	startTime, endTime, err := parseTimeWindow(timeWindow)
	if err != nil {
		return nil, err
	}

	var filteredEvents []models.Event

	// Filter events within the time window
	for _, event := range events {
		if (event.OccurredAt.Equal(startTime) || event.OccurredAt.After(startTime)) &&
			(event.OccurredAt.Equal(endTime) || event.OccurredAt.Before(endTime)) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	// Handle business days only
	if businessDaysOnly, ok := timeWindow["businessDaysOnly"].(bool); ok && businessDaysOnly {
		var businessDayEvents []models.Event

		for _, event := range filteredEvents {
			if !isWeekend(event.OccurredAt) {
				businessDayEvents = append(businessDayEvents, event)
			}
		}

		return businessDayEvents, nil
	}

	return filteredEvents, nil
}

// calculateDuration calculates the duration between two events in the specified unit
func calculateDuration(start, end time.Time, unit string) float64 {
	duration := end.Sub(start)

	switch unit {
	case "minutes":
		return duration.Minutes()
	case "hours":
		return duration.Hours()
	case "days":
		return duration.Hours() / 24
	case "seconds":
		return duration.Seconds()
	default:
		return duration.Hours() // default to hours
	}
}
