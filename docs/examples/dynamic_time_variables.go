package examples

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Example of how dynamic time variables could be implemented in the rule engine

// TimeVariableResolver resolves dynamic time variables in criteria
type TimeVariableResolver struct {
	now time.Time // Cached current time for consistent evaluation
}

// NewTimeVariableResolver creates a new resolver with the current time
func NewTimeVariableResolver() *TimeVariableResolver {
	return &TimeVariableResolver{
		now: time.Now().UTC(),
	}
}

// WithFixedTime creates a resolver with a fixed time (for testing)
func NewTimeVariableResolverWithFixedTime(fixedTime time.Time) *TimeVariableResolver {
	return &TimeVariableResolver{
		now: fixedTime,
	}
}

// ResolveCriteria resolves any time variables in the criteria map
func (r *TimeVariableResolver) ResolveCriteria(criteria map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	for key, value := range criteria {
		// If value is a nested map, recursively resolve it
		if nestedMap, ok := value.(map[string]interface{}); ok {
			result[key] = r.ResolveCriteria(nestedMap)
			continue
		}
		
		// If value is a string, check for time variables
		if strValue, ok := value.(string); ok {
			if strings.HasPrefix(strValue, "$NOW") {
				// Resolve the $NOW variable
				resolvedTime, err := r.resolveNowVariable(strValue)
				if err == nil {
					// Format as RFC3339 for timestamp comparison
					result[key] = resolvedTime.Format(time.RFC3339)
				} else {
					// If there's an error, keep the original value
					result[key] = strValue
					fmt.Printf("Error resolving time variable '%s': %v\n", strValue, err)
				}
				continue
			}
		}
		
		// Default: keep the original value
		result[key] = value
	}
	
	return result
}

// Regular expression to parse $NOW syntax
// Matches $NOW, $NOW(+1d), $NOW(-30d), etc.
var nowRegex = regexp.MustCompile(`^\$NOW(?:\(([+-])(\d+)([smhdwMy])\))?$`)

// resolveNowVariable parses and resolves a $NOW time variable
func (r *TimeVariableResolver) resolveNowVariable(variable string) (time.Time, error) {
	matches := nowRegex.FindStringSubmatch(variable)
	
	// If no match or just $NOW without adjustment
	if len(matches) == 0 || (len(matches) == 1 && matches[0] == "$NOW") {
		return r.now, nil
	}
	
	// Extract adjustment parts
	sign := matches[1]      // + or -
	amount, err := strconv.Atoi(matches[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid amount in time variable: %v", err)
	}
	unit := matches[3]      // s, m, h, d, w, M, y
	
	// Apply negative adjustment if sign is -
	if sign == "-" {
		amount = -amount
	}
	
	// Calculate adjusted time
	return r.adjustTime(amount, unit)
}

// adjustTime applies the time adjustment based on amount and unit
func (r *TimeVariableResolver) adjustTime(amount int, unit string) (time.Time, error) {
	switch unit {
	case "s":
		return r.now.Add(time.Duration(amount) * time.Second), nil
	case "m":
		return r.now.Add(time.Duration(amount) * time.Minute), nil
	case "h":
		return r.now.Add(time.Duration(amount) * time.Hour), nil
	case "d":
		return r.now.AddDate(0, 0, amount), nil
	case "w":
		return r.now.AddDate(0, 0, 7*amount), nil
	case "M":
		return r.now.AddDate(0, amount, 0), nil
	case "y":
		return r.now.AddDate(amount, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unknown time unit: %s", unit)
	}
}

// ExampleUsage demonstrates how the TimeVariableResolver would be used
func ExampleUsage() {
	// Example badge criteria with time variables
	criteria := map[string]interface{}{
		"$eventCount": map[string]interface{}{
			"$gte": float64(5),
		},
		"timestamp": map[string]interface{}{
			"$gte": "$NOW(-30d)", // Events from the last 30 days
		},
	}
	
	// Create a resolver
	resolver := NewTimeVariableResolver()
	
	// Resolve time variables in the criteria
	resolvedCriteria := resolver.ResolveCriteria(criteria)
	
	// The resolved criteria would now have concrete timestamps
	fmt.Printf("Original criteria: %+v\n", criteria)
	fmt.Printf("Resolved criteria: %+v\n", resolvedCriteria)
	
	// For testing with a fixed time
	fixedTime := time.Date(2023, 12, 15, 12, 0, 0, 0, time.UTC)
	testResolver := NewTimeVariableResolverWithFixedTime(fixedTime)
	
	// Test a more complex criteria
	complexCriteria := map[string]interface{}{
		"$eventCount": map[string]interface{}{
			"$gte": float64(3),
		},
		"timestamp": map[string]interface{}{
			"$gte": "$NOW(-1y)",  // From 1 year ago
			"$lte": "$NOW",       // To now
		},
		"nested": map[string]interface{}{
			"created_at": "$NOW(-1M)", // 1 month ago
		},
	}
	
	resolvedComplex := testResolver.ResolveCriteria(complexCriteria)
	fmt.Printf("Resolved complex criteria with fixed time: %+v\n", resolvedComplex)
}

// How to integrate with the existing rule engine
func integrationExample() {
	// 1. During badge criteria evaluation, before evaluating the criteria:
	// criteria := GetBadgeCriteria()
	
	// 2. Create a resolver
	// resolver := NewTimeVariableResolver()
	
	// 3. Resolve time variables
	// resolvedCriteria := resolver.ResolveCriteria(criteria)
	
	// 4. Evaluate the resolved criteria against events
	// EvaluateBadgeCriteria(resolvedCriteria, events)
} 
