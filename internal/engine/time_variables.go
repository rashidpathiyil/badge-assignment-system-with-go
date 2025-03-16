package engine

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// TimeVariableCache caches time variables within a single badge evaluation
type TimeVariableCache struct {
	now time.Time
}

// NewTimeVariableCache creates a new time variable cache with current time
func NewTimeVariableCache() *TimeVariableCache {
	return &TimeVariableCache{
		now: time.Now().UTC(),
	}
}

// IsDynamicTimeVariable checks if a string value is a dynamic time variable
func IsDynamicTimeVariable(value string) bool {
	return len(value) >= 4 && value[0:4] == "$NOW"
}

// ParseDynamicTimeVariable parses a dynamic time variable and returns the corresponding time
func ParseDynamicTimeVariable(value string, cache *TimeVariableCache) (time.Time, error) {
	if !IsDynamicTimeVariable(value) {
		return time.Time{}, fmt.Errorf("not a dynamic time variable: %s", value)
	}

	// Basic $NOW with no adjustments
	if value == "$NOW" {
		return cache.now, nil
	}

	// $NOW with adjustments: $NOW(-30d), $NOW(-1y-3M), etc.
	adjustmentRegex := regexp.MustCompile(`^\$NOW\(([-+][0-9]+[smhdwMy](?:[-+][0-9]+[smhdwMy])*)\)$`)
	matches := adjustmentRegex.FindStringSubmatch(value)
	if len(matches) != 2 {
		return time.Time{}, fmt.Errorf("invalid $NOW syntax: %s", value)
	}

	adjustmentStr := matches[1]
	adjusted := cache.now

	// Parse individual adjustments (e.g., -30d, +1y)
	unitRegex := regexp.MustCompile(`([-+][0-9]+)([smhdwMy])`)
	adjustments := unitRegex.FindAllStringSubmatch(adjustmentStr, -1)

	for _, adj := range adjustments {
		if len(adj) != 3 {
			continue
		}

		// Parse the number (with sign)
		amount, err := strconv.Atoi(adj[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid adjustment amount: %s", adj[1])
		}

		// Apply the adjustment based on the unit
		unit := adj[2]
		adjusted = applyTimeAdjustment(adjusted, amount, unit)
	}

	return adjusted, nil
}

// applyTimeAdjustment applies a time adjustment based on amount and unit
func applyTimeAdjustment(t time.Time, amount int, unit string) time.Time {
	switch unit {
	case "s": // seconds
		return t.Add(time.Duration(amount) * time.Second)
	case "m": // minutes
		return t.Add(time.Duration(amount) * time.Minute)
	case "h": // hours
		return t.Add(time.Duration(amount) * time.Hour)
	case "d": // days
		return t.AddDate(0, 0, amount)
	case "w": // weeks
		return t.AddDate(0, 0, amount*7)
	case "M": // months
		return t.AddDate(0, amount, 0)
	case "y": // years
		return t.AddDate(amount, 0, 0)
	default:
		return t
	}
}
