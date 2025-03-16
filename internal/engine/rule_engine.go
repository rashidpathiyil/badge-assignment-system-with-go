package engine

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/badge-assignment-system/internal/logging"
	"github.com/badge-assignment-system/internal/models"
)

// Ensure that *models.DB implements DBInterface
var _ DBInterface = (*models.DB)(nil)

// DBInterface defines the database operations needed by the rule engine
type DBInterface interface {
	GetBadgeWithCriteria(id int) (models.BadgeWithCriteria, error)
	GetEventTypeByName(name string) (models.EventType, error)
	GetUserEventsByType(userID string, eventTypeID int) ([]models.Event, error)
	GetUserEvents(userID string) ([]models.Event, error)
	GetActiveBadges() ([]models.Badge, error)
	GetUserBadges(userID string) ([]models.UserBadge, error)
	AwardBadgeToUser(userBadge *models.UserBadge) error
}

// RuleEngine handles the dynamic evaluation of badge criteria against events
type RuleEngine struct {
	DB           DBInterface
	Logger       *logging.Logger
	TimeVarCache *TimeVariableCache
}

// NewRuleEngine creates a new rule engine
func NewRuleEngine(db DBInterface) *RuleEngine {
	return &RuleEngine{
		DB:           db,
		Logger:       logging.NewLogger("RULE-ENGINE", logging.LogLevelInfo),
		TimeVarCache: NewTimeVariableCache(),
	}
}

// SetLogLevel sets the logging level for the rule engine
func (re *RuleEngine) SetLogLevel(level logging.LogLevel) {
	re.Logger.SetLevel(level)
}

// EvaluateBadgeCriteria checks if a user meets the criteria for a badge
func (re *RuleEngine) EvaluateBadgeCriteria(badgeID int, userID string) (bool, map[string]interface{}, error) {
	re.Logger.Debug("Evaluating badge criteria for badge ID %d and user %s", badgeID, userID)

	// Reset time variable cache for new evaluation
	re.TimeVarCache = NewTimeVariableCache()

	// Get badge with criteria
	badgeWithCriteria, err := re.DB.GetBadgeWithCriteria(badgeID)
	if err != nil {
		re.Logger.Error("Failed to get badge criteria: %v", err)
		return false, nil, fmt.Errorf("failed to get badge criteria: %w", err)
	}

	re.Logger.Debug("Retrieved badge criteria for badge ID: %d", badgeID)

	// Extract the criteria flow definition
	flowDefinition := badgeWithCriteria.Criteria.FlowDefinition

	// Evaluate the criteria
	metadata := make(map[string]interface{})
	re.Logger.Debug("Starting flow evaluation for badge %d", badgeID)
	result, err := re.evaluateFlow(flowDefinition, userID, metadata)
	if err != nil {
		re.Logger.Error("Criteria evaluation failed: %v", err)
		return false, nil, fmt.Errorf("criteria evaluation failed: %w", err)
	}

	re.Logger.Debug("Badge %d criteria evaluation result: %v with %d metadata items",
		badgeID, result, len(metadata))

	return result, metadata, nil
}

// evaluateFlow recursively evaluates a badge criteria flow definition
func (re *RuleEngine) evaluateFlow(flow models.JSONB, userID string, metadata map[string]interface{}) (bool, error) {
	// Check if this is an event-based criterion
	if eventType, hasEventType := flow["event"].(string); hasEventType {
		re.Logger.Debug("Evaluating event-based criterion for event type: %s", eventType)

		// Get the event type ID
		eventTypeObj, err := re.DB.GetEventTypeByName(eventType)
		if err != nil {
			re.Logger.Error("Event type '%s' not found: %v", eventType, err)
			return false, fmt.Errorf("event type '%s' not found: %w", eventType, err)
		}

		// Get criteria for this event
		criteria, hasCriteria := flow["criteria"].(map[string]interface{})
		if !hasCriteria {
			re.Logger.Error("Invalid criteria format: missing 'criteria' field")
			return false, errors.New("invalid criteria format: missing 'criteria' field")
		}

		// Get events for this user and event type
		re.Logger.Debug("Retrieving events for user %s and event type %s (ID: %d)",
			userID, eventType, eventTypeObj.ID)
		events, err := re.DB.GetUserEventsByType(userID, eventTypeObj.ID)
		if err != nil {
			re.Logger.Error("Failed to get user events: %v", err)
			return false, fmt.Errorf("failed to get user events: %w", err)
		}

		re.Logger.Debug("Found %d events of type '%s' for user %s", len(events), eventType, userID)
		return re.evaluateEventCriteria(criteria, events, metadata)
	}

	// Handle logical operators
	for operator, value := range flow {
		re.Logger.Debug("Processing operator: %s", operator)

		switch operator {
		case "$and":
			re.Logger.Debug("Evaluating $and operator")
			return re.evaluateAndOperator(value, userID, metadata)
		case "$or":
			re.Logger.Debug("Evaluating $or operator")
			return re.evaluateOrOperator(value, userID, metadata)
		case "$not":
			re.Logger.Debug("Evaluating $not operator")
			return re.evaluateNotOperator(value, userID, metadata)
		// Time-based operators
		case "$timePeriod":
			re.Logger.Debug("Evaluating $timePeriod operator")
			criteria, ok := value.(map[string]interface{})
			if !ok {
				re.Logger.Error("$timePeriod requires a criteria object")
				return false, fmt.Errorf("$timePeriod requires a criteria object")
			}
			// Get all events for the user (across all event types)
			re.Logger.Debug("Retrieving all events for user %s for time period evaluation", userID)
			events, err := re.DB.GetUserEvents(userID)
			if err != nil {
				re.Logger.Error("Failed to get user events: %v", err)
				return false, fmt.Errorf("failed to get user events: %w", err)
			}
			re.Logger.Debug("Found %d total events for user %s", len(events), userID)
			return re.evaluateTimePeriodCriteria(criteria, events, metadata)
		case "$pattern":
			re.Logger.Debug("Evaluating $pattern operator")
			criteria, ok := value.(map[string]interface{})
			if !ok {
				re.Logger.Error("$pattern requires a criteria object")
				return false, fmt.Errorf("$pattern requires a criteria object")
			}
			events, err := re.DB.GetUserEvents(userID)
			if err != nil {
				re.Logger.Error("Failed to get user events: %v", err)
				return false, fmt.Errorf("failed to get user events: %w", err)
			}
			re.Logger.Debug("Found %d total events for user %s", len(events), userID)
			return re.evaluatePatternCriteria(criteria, events, metadata)
		case "$sequence":
			re.Logger.Debug("Evaluating $sequence operator")
			criteria, ok := value.(map[string]interface{})
			if !ok {
				re.Logger.Error("$sequence requires a criteria object")
				return false, fmt.Errorf("$sequence requires a criteria object")
			}
			return re.evaluateSequenceCriteria(criteria, userID, metadata)
		case "$gap":
			re.Logger.Debug("Evaluating $gap operator")
			criteria, ok := value.(map[string]interface{})
			if !ok {
				re.Logger.Error("$gap requires a criteria object")
				return false, fmt.Errorf("$gap requires a criteria object")
			}
			events, err := re.DB.GetUserEvents(userID)
			if err != nil {
				re.Logger.Error("Failed to get user events: %v", err)
				return false, fmt.Errorf("failed to get user events: %w", err)
			}
			re.Logger.Debug("Found %d total events for user %s", len(events), userID)
			return re.evaluateGapCriteria(criteria, events, metadata)
		case "$duration":
			re.Logger.Debug("Evaluating $duration operator")
			criteria, ok := value.(map[string]interface{})
			if !ok {
				re.Logger.Error("$duration requires a criteria object")
				return false, fmt.Errorf("$duration requires a criteria object")
			}
			events, err := re.DB.GetUserEvents(userID)
			if err != nil {
				re.Logger.Error("Failed to get user events: %v", err)
				return false, fmt.Errorf("failed to get user events: %w", err)
			}
			re.Logger.Debug("Found %d total events for user %s", len(events), userID)
			return re.evaluateDurationCriteria(criteria, events, metadata)
		case "$aggregate":
			re.Logger.Debug("Evaluating $aggregate operator")
			criteria, ok := value.(map[string]interface{})
			if !ok {
				re.Logger.Error("$aggregate requires a criteria object")
				return false, fmt.Errorf("$aggregate requires a criteria object")
			}
			events, err := re.DB.GetUserEvents(userID)
			if err != nil {
				re.Logger.Error("Failed to get user events: %v", err)
				return false, fmt.Errorf("failed to get user events: %w", err)
			}
			re.Logger.Debug("Found %d total events for user %s", len(events), userID)
			return re.evaluateAggregationCriteria(criteria, events, metadata)
		case "$timeWindow":
			re.Logger.Debug("Evaluating $timeWindow operator")
			criteria, ok := value.(map[string]interface{})
			if !ok {
				re.Logger.Error("$timeWindow requires a criteria object")
				return false, fmt.Errorf("$timeWindow requires a criteria object")
			}
			subFlow, ok := criteria["flow"].(map[string]interface{})
			if !ok {
				re.Logger.Error("$timeWindow requires a 'flow' object")
				return false, fmt.Errorf("$timeWindow requires a 'flow' object")
			}

			// Parse the time window
			windowStart, windowEnd, err := parseTimeWindow(criteria, re.TimeVarCache)
			if err != nil {
				re.Logger.Error("Failed to parse time window: %v", err)
				return false, err
			}

			re.Logger.Debug("Time window from %s to %s",
				windowStart.Format(time.RFC3339), windowEnd.Format(time.RFC3339))

			// Create a sub-metadata map to capture results within the time window
			windowMetadata := make(map[string]interface{})

			// Evaluate the sub-flow with time window constraints
			tempUserID := fmt.Sprintf("%s|%s|%s", userID, windowStart.Format(time.RFC3339), windowEnd.Format(time.RFC3339))
			re.Logger.Debug("Evaluating subflow with time window constraints using temp user ID: %s", tempUserID)
			result, err := re.evaluateFlow(models.JSONB(subFlow), tempUserID, windowMetadata)
			if err != nil {
				re.Logger.Error("Error evaluating time window subflow: %v", err)
				return false, err
			}

			// Merge the window metadata with the parent metadata
			for k, v := range windowMetadata {
				metadata[fmt.Sprintf("window_%s", k)] = v
			}

			re.Logger.Debug("Time window evaluation result: %v", result)
			return result, nil
		default:
			re.Logger.Warning("Unsupported operator: %s", operator)
		}
	}

	re.Logger.Error("Unsupported flow definition format")
	return false, errors.New("unsupported flow definition format")
}

// evaluateAndOperator handles the $and operator
func (re *RuleEngine) evaluateAndOperator(conditions interface{}, userID string, metadata map[string]interface{}) (bool, error) {
	conditionsArray, ok := conditions.([]interface{})
	if !ok {
		re.Logger.Error("$and operator requires an array of conditions")
		return false, errors.New("$and operator requires an array of conditions")
	}

	re.Logger.Debug("Evaluating $and operator with %d conditions", len(conditionsArray))

	for i, condition := range conditionsArray {
		re.Logger.Debug("Evaluating condition %d/%d in $and", i+1, len(conditionsArray))

		conditionMap, ok := condition.(map[string]interface{})
		if !ok {
			re.Logger.Error("Each condition in $and must be an object")
			return false, errors.New("each condition in $and must be an object")
		}

		result, err := re.evaluateFlow(models.JSONB(conditionMap), userID, metadata)
		if err != nil {
			re.Logger.Error("Error evaluating condition %d in $and: %v", i+1, err)
			return false, err
		}

		if !result {
			re.Logger.Debug("Condition %d in $and evaluated to false, short-circuiting", i+1)
			return false, nil // Short-circuit: if any condition is false, the whole AND is false
		}

		re.Logger.Debug("Condition %d in $and evaluated to true", i+1)
	}

	re.Logger.Debug("All conditions in $and evaluated to true")
	return true, nil // All conditions passed
}

// evaluateOrOperator handles the $or operator
func (re *RuleEngine) evaluateOrOperator(conditions interface{}, userID string, metadata map[string]interface{}) (bool, error) {
	conditionsArray, ok := conditions.([]interface{})
	if !ok {
		re.Logger.Error("$or operator requires an array of conditions")
		return false, errors.New("$or operator requires an array of conditions")
	}

	re.Logger.Debug("Evaluating $or operator with %d conditions", len(conditionsArray))

	for i, condition := range conditionsArray {
		re.Logger.Debug("Evaluating condition %d/%d in $or", i+1, len(conditionsArray))

		conditionMap, ok := condition.(map[string]interface{})
		if !ok {
			re.Logger.Error("Each condition in $or must be an object")
			return false, errors.New("each condition in $or must be an object")
		}

		result, err := re.evaluateFlow(models.JSONB(conditionMap), userID, metadata)
		if err != nil {
			re.Logger.Error("Error evaluating condition %d in $or: %v", i+1, err)
			return false, err
		}

		if result {
			re.Logger.Debug("Condition %d in $or evaluated to true, short-circuiting", i+1)
			return true, nil // Short-circuit: if any condition is true, the whole OR is true
		}

		re.Logger.Debug("Condition %d in $or evaluated to false", i+1)
	}

	re.Logger.Debug("All conditions in $or evaluated to false")
	return false, nil // No conditions passed
}

// evaluateNotOperator handles the $not operator
func (re *RuleEngine) evaluateNotOperator(condition interface{}, userID string, metadata map[string]interface{}) (bool, error) {
	conditionMap, ok := condition.(map[string]interface{})
	if !ok {
		re.Logger.Error("$not operator requires a condition object")
		return false, errors.New("$not operator requires a condition object")
	}

	re.Logger.Debug("Evaluating $not operator")

	result, err := re.evaluateFlow(models.JSONB(conditionMap), userID, metadata)
	if err != nil {
		re.Logger.Error("Error evaluating condition in $not: %v", err)
		return false, err
	}

	re.Logger.Debug("Condition in $not evaluated to %v, negating to %v", result, !result)
	return !result, nil // Negate the result
}

// evaluateEventCriteria evaluates criteria against a set of events
func (re *RuleEngine) evaluateEventCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	re.Logger.Debug("Evaluating event criteria against %d events", len(events))

	// Handle event count criteria
	if eventCountCriteria, hasEventCount := criteria["$eventCount"].(map[string]interface{}); hasEventCount {
		re.Logger.Debug("Detected $eventCount criteria, evaluating")
		return re.evaluateEventCountCriteria(eventCountCriteria, events, metadata)
	}

	// Filter events based on criteria
	re.Logger.Debug("Filtering %d events based on criteria", len(events))
	filteredEvents, err := re.filterEvents(criteria, events)
	if err != nil {
		re.Logger.Error("Error filtering events: %v", err)
		return false, err
	}

	// Store filtered events count in metadata
	metadata["filtered_event_count"] = len(filteredEvents)
	if len(filteredEvents) > 0 {
		metadata["first_event_id"] = filteredEvents[0].ID
		metadata["last_event_id"] = filteredEvents[len(filteredEvents)-1].ID
		re.Logger.Debug("Filtered to %d events, first ID: %d, last ID: %d",
			len(filteredEvents), filteredEvents[0].ID, filteredEvents[len(filteredEvents)-1].ID)
	} else {
		re.Logger.Debug("No events matched the criteria")
	}

	// If we get here, the criteria is considered met if there are any events that match
	result := len(filteredEvents) > 0
	re.Logger.Debug("Event criteria evaluation result: %v", result)
	return result, nil
}

// filterEvents filters events based on criteria
func (re *RuleEngine) filterEvents(criteria map[string]interface{}, events []models.Event) ([]models.Event, error) {
	var filteredEvents []models.Event

	re.Logger.Trace("Starting to filter %d events", len(events))
	for _, event := range events {
		passes, err := re.eventMatchesCriteria(event, criteria)
		if err != nil {
			re.Logger.Error("Error matching event %d against criteria: %v", event.ID, err)
			return nil, err
		}
		if passes {
			filteredEvents = append(filteredEvents, event)
			re.Logger.Trace("Event ID %d matched criteria", event.ID)
		}
	}

	re.Logger.Trace("Filtered to %d/%d events", len(filteredEvents), len(events))
	return filteredEvents, nil
}

// eventMatchesCriteria checks if an event matches the given criteria
func (re *RuleEngine) eventMatchesCriteria(event models.Event, criteria map[string]interface{}) (bool, error) {
	for field, conditionValue := range criteria {
		// Skip the event count field as it's handled separately
		if field == "$eventCount" {
			continue
		}

		// Handle timestamp field specially
		if field == "timestamp" {
			conditionMap, ok := conditionValue.(map[string]interface{})
			if !ok {
				re.Logger.Error("Timestamp condition must be an object")
				return false, errors.New("timestamp condition must be an object")
			}
			matches, err := re.evaluateTimestampCondition(event.OccurredAt, conditionMap)
			if err != nil {
				re.Logger.Error("Error evaluating timestamp condition: %v", err)
				return false, err
			}
			if !matches {
				re.Logger.Trace("Event ID %d timestamp did not match condition", event.ID)
				return false, nil
			}
			continue
		}

		// For other fields, check in the payload
		if event.Payload == nil {
			re.Logger.Trace("Event ID %d has no payload, criteria not met", event.ID)
			return false, nil
		}

		// Get the field value from the event payload
		fieldValue, exists := event.Payload[field]
		if !exists {
			re.Logger.Trace("Event ID %d payload does not contain field '%s'", event.ID, field)
			return false, nil
		}

		// If conditionValue is a comparison operator object
		if conditionMap, ok := conditionValue.(map[string]interface{}); ok {
			matches, err := re.evaluateComparison(fieldValue, conditionMap)
			if err != nil {
				re.Logger.Error("Error evaluating comparison for field '%s': %v", field, err)
				return false, err
			}
			if !matches {
				re.Logger.Trace("Event ID %d field '%s' did not match comparison", event.ID, field)
				return false, nil
			}
		} else {
			// Direct equality comparison
			if !reflect.DeepEqual(fieldValue, conditionValue) {
				re.Logger.Trace("Event ID %d field '%s' did not match direct equality", event.ID, field)
				return false, nil
			}
		}
	}

	return true, nil
}

// evaluateTimestampCondition evaluates timestamp-specific conditions
func (re *RuleEngine) evaluateTimestampCondition(timestamp time.Time, conditions map[string]interface{}) (bool, error) {
	for operator, value := range conditions {
		re.Logger.Trace("Evaluating timestamp operator %s against %v", operator, timestamp)

		switch operator {
		case "$gte":
			compareTime, err := re.parseTimeValueWithCache(value)
			if err != nil {
				re.Logger.Error("Error parsing time value for $gte: %v", err)
				return false, err
			}
			if !timestamp.Equal(compareTime) && !timestamp.After(compareTime) {
				re.Logger.Trace("Timestamp %v is not >= %v", timestamp, compareTime)
				return false, nil
			}
		case "$gt":
			compareTime, err := re.parseTimeValueWithCache(value)
			if err != nil {
				re.Logger.Error("Error parsing time value for $gt: %v", err)
				return false, err
			}
			if !timestamp.After(compareTime) {
				re.Logger.Trace("Timestamp %v is not > %v", timestamp, compareTime)
				return false, nil
			}
		case "$lte":
			compareTime, err := re.parseTimeValueWithCache(value)
			if err != nil {
				re.Logger.Error("Error parsing time value for $lte: %v", err)
				return false, err
			}
			if !timestamp.Equal(compareTime) && !timestamp.Before(compareTime) {
				re.Logger.Trace("Timestamp %v is not <= %v", timestamp, compareTime)
				return false, nil
			}
		case "$lt":
			compareTime, err := re.parseTimeValueWithCache(value)
			if err != nil {
				re.Logger.Error("Error parsing time value for $lt: %v", err)
				return false, err
			}
			if !timestamp.Before(compareTime) {
				re.Logger.Trace("Timestamp %v is not < %v", timestamp, compareTime)
				return false, nil
			}
		case "$eq":
			compareTime, err := re.parseTimeValueWithCache(value)
			if err != nil {
				re.Logger.Error("Error parsing time value for $eq: %v", err)
				return false, err
			}
			if !timestamp.Equal(compareTime) {
				re.Logger.Trace("Timestamp %v is not == %v", timestamp, compareTime)
				return false, nil
			}
		case "$ne":
			compareTime, err := re.parseTimeValueWithCache(value)
			if err != nil {
				re.Logger.Error("Error parsing time value for $ne: %v", err)
				return false, err
			}
			if timestamp.Equal(compareTime) {
				re.Logger.Trace("Timestamp %v is not != %v", timestamp, compareTime)
				return false, nil
			}
		default:
			re.Logger.Error("Unsupported timestamp operator: %s", operator)
			return false, fmt.Errorf("unsupported timestamp operator: %s", operator)
		}
	}
	return true, nil
}

// parseTimeValueWithCache converts a string or RFC3339 time value to time.Time using the TimeVariableCache
func (re *RuleEngine) parseTimeValueWithCache(value interface{}) (time.Time, error) {
	if timeStr, ok := value.(string); ok {
		// Check if this is a dynamic time variable
		if IsDynamicTimeVariable(timeStr) {
			return ParseDynamicTimeVariable(timeStr, re.TimeVarCache)
		}

		// Otherwise parse as normal RFC3339 time
		return time.Parse(time.RFC3339, timeStr)
	}
	return time.Time{}, errors.New("timestamp value must be a string in RFC3339 format or dynamic time variable")
}

// parseTimeValue converts a string or RFC3339 time value to time.Time
// Global helper function for backward compatibility
func parseTimeValue(value interface{}) (time.Time, error) {
	if timeStr, ok := value.(string); ok {
		// Check if this is a dynamic time variable
		if IsDynamicTimeVariable(timeStr) {
			cache := NewTimeVariableCache()
			return ParseDynamicTimeVariable(timeStr, cache)
		}

		// Otherwise parse as normal RFC3339 time
		return time.Parse(time.RFC3339, timeStr)
	}
	return time.Time{}, errors.New("timestamp value must be a string in RFC3339 format or dynamic time variable")
}

// evaluateComparison evaluates comparison operators on values
func (re *RuleEngine) evaluateComparison(fieldValue interface{}, conditions map[string]interface{}) (bool, error) {
	for operator, compareValue := range conditions {
		re.Logger.Trace("Evaluating comparison operator %s", operator)

		switch operator {
		case "$eq":
			if !reflect.DeepEqual(fieldValue, compareValue) {
				re.Logger.Trace("Value %v is not == %v", fieldValue, compareValue)
				return false, nil
			}
		case "$ne":
			if reflect.DeepEqual(fieldValue, compareValue) {
				re.Logger.Trace("Value %v is not != %v", fieldValue, compareValue)
				return false, nil
			}
		case "$gt":
			result, err := compareValues(fieldValue, compareValue, func(a, b float64) bool { return a > b })
			if err != nil {
				re.Logger.Error("Error in $gt comparison: %v", err)
				return false, err
			}
			if !result {
				re.Logger.Trace("Value %v is not > %v", fieldValue, compareValue)
				return false, nil
			}
		case "$gte":
			result, err := compareValues(fieldValue, compareValue, func(a, b float64) bool { return a >= b })
			if err != nil {
				re.Logger.Error("Error in $gte comparison: %v", err)
				return false, err
			}
			if !result {
				re.Logger.Trace("Value %v is not >= %v", fieldValue, compareValue)
				return false, nil
			}
		case "$lt":
			result, err := compareValues(fieldValue, compareValue, func(a, b float64) bool { return a < b })
			if err != nil {
				re.Logger.Error("Error in $lt comparison: %v", err)
				return false, err
			}
			if !result {
				re.Logger.Trace("Value %v is not < %v", fieldValue, compareValue)
				return false, nil
			}
		case "$lte":
			result, err := compareValues(fieldValue, compareValue, func(a, b float64) bool { return a <= b })
			if err != nil {
				re.Logger.Error("Error in $lte comparison: %v", err)
				return false, err
			}
			if !result {
				re.Logger.Trace("Value %v is not <= %v", fieldValue, compareValue)
				return false, nil
			}
		case "$in":
			if !isInArray(fieldValue, compareValue) {
				re.Logger.Trace("Value %v is not in array %v", fieldValue, compareValue)
				return false, nil
			}
		case "$nin":
			if isInArray(fieldValue, compareValue) {
				re.Logger.Trace("Value %v is in array %v (should not be)", fieldValue, compareValue)
				return false, nil
			}
		case "$regex":
			re.Logger.Warning("$regex operator not implemented yet")
			return false, errors.New("$regex operator not implemented yet")
		default:
			re.Logger.Error("Unsupported operator: %s", operator)
			return false, fmt.Errorf("unsupported operator: %s", operator)
		}
	}
	return true, nil
}

// compareValues compares two values using the provided comparison function
func compareValues(a, b interface{}, compare func(float64, float64) bool) (bool, error) {
	// Convert values to float64 for comparison
	aFloat, err := toFloat64(a)
	if err != nil {
		return false, err
	}

	bFloat, err := toFloat64(b)
	if err != nil {
		return false, err
	}

	return compare(aFloat, bFloat), nil
}

// toFloat64 converts an interface{} to float64
func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case string:
		// Try to convert string to float64
		var result float64
		_, err := fmt.Sscanf(v, "%f", &result)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string '%s' to float64: %w", v, err)
		}
		return result, nil
	default:
		return 0, fmt.Errorf("cannot convert type %T to float64", value)
	}
}

// isInArray checks if a value is in an array
func isInArray(value interface{}, array interface{}) bool {
	arrayValue := reflect.ValueOf(array)
	if arrayValue.Kind() != reflect.Slice && arrayValue.Kind() != reflect.Array {
		return false
	}

	for i := 0; i < arrayValue.Len(); i++ {
		if reflect.DeepEqual(value, arrayValue.Index(i).Interface()) {
			return true
		}
	}
	return false
}

// evaluateEventCountCriteria checks if the number of events meets the count criteria
func (re *RuleEngine) evaluateEventCountCriteria(eventCountCriteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	re.Logger.Debug("Evaluating event count criteria against %d events", len(events))

	// First, filter events based on other criteria in the parent object
	filteredEvents, err := re.filterEvents(map[string]interface{}{}, events)
	if err != nil {
		re.Logger.Error("Error filtering events for event count criteria: %v", err)
		return false, err
	}

	count := len(filteredEvents)
	metadata["event_count"] = count
	re.Logger.Debug("Event count criteria: working with %d filtered events", count)

	// Evaluate the event count criteria
	for operator, value := range eventCountCriteria {
		compareValue, err := toFloat64(value)
		if err != nil {
			re.Logger.Error("Invalid event count comparison value: %v", err)
			return false, fmt.Errorf("invalid event count comparison value: %w", err)
		}

		floatCount := float64(count)
		re.Logger.Debug("Evaluating event count operator %s: %v %s %v",
			operator, floatCount, operator, compareValue)

		switch operator {
		case "$eq":
			if floatCount != compareValue {
				re.Logger.Debug("Event count criterion not met: %v != %v", floatCount, compareValue)
				return false, nil
			}
		case "$ne":
			if floatCount == compareValue {
				re.Logger.Debug("Event count criterion not met: %v == %v", floatCount, compareValue)
				return false, nil
			}
		case "$gt":
			if floatCount <= compareValue {
				re.Logger.Debug("Event count criterion not met: %v <= %v", floatCount, compareValue)
				return false, nil
			}
		case "$gte":
			if floatCount < compareValue {
				re.Logger.Debug("Event count criterion not met: %v < %v", floatCount, compareValue)
				return false, nil
			}
		case "$lt":
			if floatCount >= compareValue {
				re.Logger.Debug("Event count criterion not met: %v >= %v", floatCount, compareValue)
				return false, nil
			}
		case "$lte":
			if floatCount > compareValue {
				re.Logger.Debug("Event count criterion not met: %v > %v", floatCount, compareValue)
				return false, nil
			}
		default:
			re.Logger.Error("Unsupported event count operator: %s", operator)
			return false, fmt.Errorf("unsupported event count operator: %s", operator)
		}
	}

	re.Logger.Debug("Event count criteria met")
	return true, nil
}

// ProcessEvents processes a batch of events and awards badges if criteria are met
func (re *RuleEngine) ProcessEvents(userID string) error {
	re.Logger.Info("Processing events for user %s", userID)

	// Get all active badges
	badges, err := re.DB.GetActiveBadges()
	if err != nil {
		re.Logger.Error("Failed to retrieve active badges: %v", err)
		return fmt.Errorf("failed to retrieve active badges: %w", err)
	}
	re.Logger.Debug("Retrieved %d active badges", len(badges))

	// Get user's existing badges
	userBadges, err := re.DB.GetUserBadges(userID)
	if err != nil {
		re.Logger.Error("Failed to retrieve user badges: %v", err)
		return fmt.Errorf("failed to retrieve user badges: %w", err)
	}
	re.Logger.Debug("User %s already has %d badges", userID, len(userBadges))

	// Create a map for fast lookup of existing badges
	userBadgeMap := make(map[int]bool)
	for _, badge := range userBadges {
		userBadgeMap[badge.BadgeID] = true
		re.Logger.Trace("User already has badge ID %d", badge.BadgeID)
	}

	// Process each badge
	var awarded int = 0
	for _, badge := range badges {
		re.Logger.Debug("Evaluating badge ID %d: %s", badge.ID, badge.Name)

		// Skip badges the user already has
		if userBadgeMap[badge.ID] {
			re.Logger.Debug("Badge ID %d already awarded to user %s, skipping", badge.ID, userID)
			continue
		}

		// Evaluate badge criteria
		re.Logger.Debug("Evaluating criteria for badge ID %d", badge.ID)
		result, metadata, err := re.EvaluateBadgeCriteria(badge.ID, userID)
		if err != nil {
			re.Logger.Error("Error evaluating criteria for badge ID %d: %v", badge.ID, err)
			continue
		}

		// Award badge if criteria are met
		if result {
			re.Logger.Info("Badge criteria met for badge ID %d (%s) for user %s",
				badge.ID, badge.Name, userID)

			userBadge := &models.UserBadge{
				UserID:   userID,
				BadgeID:  badge.ID,
				Metadata: models.JSONB(metadata),
			}
			err = re.DB.AwardBadgeToUser(userBadge)
			if err != nil {
				re.Logger.Error("Failed to award badge ID %d to user %s: %v",
					badge.ID, userID, err)
				continue
			}
			awarded++
			re.Logger.Info("Badge ID %d (%s) awarded to user %s", badge.ID, badge.Name, userID)
		} else {
			re.Logger.Debug("Badge criteria not met for badge ID %d for user %s", badge.ID, userID)
		}
	}

	re.Logger.Info("Badge processing complete for user %s - %d new badges awarded", userID, awarded)
	return nil
}

// ProcessEvent processes a single event and checks if it triggers any badge awards
func (re *RuleEngine) ProcessEvent(event *models.Event) error {
	re.Logger.Debug("Processing event ID %d of type %d for user %s",
		event.ID, event.EventTypeID, event.UserID)

	// Process badges for the user who triggered the event
	return re.ProcessEvents(event.UserID)
}
