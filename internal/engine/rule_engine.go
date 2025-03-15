package engine

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/badge-assignment-system/internal/models"
)

// RuleEngine handles the dynamic evaluation of badge criteria against events
type RuleEngine struct {
	DB *models.DB
}

// NewRuleEngine creates a new rule engine
func NewRuleEngine(db *models.DB) *RuleEngine {
	return &RuleEngine{
		DB: db,
	}
}

// EvaluateBadgeCriteria checks if a user meets the criteria for a badge
func (re *RuleEngine) EvaluateBadgeCriteria(badgeID int, userID string) (bool, map[string]interface{}, error) {
	// Get badge with criteria
	badgeWithCriteria, err := re.DB.GetBadgeWithCriteria(badgeID)
	if err != nil {
		return false, nil, fmt.Errorf("failed to get badge criteria: %w", err)
	}

	// Extract the criteria flow definition
	flowDefinition := badgeWithCriteria.Criteria.FlowDefinition

	// Evaluate the criteria
	metadata := make(map[string]interface{})
	result, err := re.evaluateFlow(flowDefinition, userID, metadata)
	if err != nil {
		return false, nil, fmt.Errorf("criteria evaluation failed: %w", err)
	}

	return result, metadata, nil
}

// evaluateFlow recursively evaluates a badge criteria flow definition
func (re *RuleEngine) evaluateFlow(flow models.JSONB, userID string, metadata map[string]interface{}) (bool, error) {
	// Check if this is an event-based criterion
	if eventType, hasEventType := flow["event"].(string); hasEventType {
		// Get the event type ID
		eventTypeObj, err := re.DB.GetEventTypeByName(eventType)
		if err != nil {
			return false, fmt.Errorf("event type '%s' not found: %w", eventType, err)
		}

		// Get criteria for this event
		criteria, hasCriteria := flow["criteria"].(map[string]interface{})
		if !hasCriteria {
			return false, errors.New("invalid criteria format: missing 'criteria' field")
		}

		// Get events for this user and event type
		events, err := re.DB.GetUserEventsByType(userID, eventTypeObj.ID)
		if err != nil {
			return false, fmt.Errorf("failed to get user events: %w", err)
		}

		return re.evaluateEventCriteria(criteria, events, metadata)
	}

	// Handle logical operators
	for operator, value := range flow {
		switch operator {
		case "$and":
			return re.evaluateAndOperator(value, userID, metadata)
		case "$or":
			return re.evaluateOrOperator(value, userID, metadata)
		case "$not":
			return re.evaluateNotOperator(value, userID, metadata)
		// Time-based operators
		case "$timePeriod":
			criteria, ok := value.(map[string]interface{})
			if !ok {
				return false, fmt.Errorf("$timePeriod requires a criteria object")
			}
			// Get all events for the user (across all event types)
			events, err := re.DB.GetUserEvents(userID)
			if err != nil {
				return false, fmt.Errorf("failed to get user events: %w", err)
			}
			return re.evaluateTimePeriodCriteria(criteria, events, metadata)
		case "$pattern":
			criteria, ok := value.(map[string]interface{})
			if !ok {
				return false, fmt.Errorf("$pattern requires a criteria object")
			}
			events, err := re.DB.GetUserEvents(userID)
			if err != nil {
				return false, fmt.Errorf("failed to get user events: %w", err)
			}
			return re.evaluatePatternCriteria(criteria, events, metadata)
		case "$sequence":
			criteria, ok := value.(map[string]interface{})
			if !ok {
				return false, fmt.Errorf("$sequence requires a criteria object")
			}
			return re.evaluateSequenceCriteria(criteria, userID, metadata)
		case "$gap":
			criteria, ok := value.(map[string]interface{})
			if !ok {
				return false, fmt.Errorf("$gap requires a criteria object")
			}
			events, err := re.DB.GetUserEvents(userID)
			if err != nil {
				return false, fmt.Errorf("failed to get user events: %w", err)
			}
			return re.evaluateGapCriteria(criteria, events, metadata)
		case "$duration":
			criteria, ok := value.(map[string]interface{})
			if !ok {
				return false, fmt.Errorf("$duration requires a criteria object")
			}
			events, err := re.DB.GetUserEvents(userID)
			if err != nil {
				return false, fmt.Errorf("failed to get user events: %w", err)
			}
			return re.evaluateDurationCriteria(criteria, events, metadata)
		case "$aggregate":
			criteria, ok := value.(map[string]interface{})
			if !ok {
				return false, fmt.Errorf("$aggregate requires a criteria object")
			}
			events, err := re.DB.GetUserEvents(userID)
			if err != nil {
				return false, fmt.Errorf("failed to get user events: %w", err)
			}
			return re.evaluateAggregationCriteria(criteria, events, metadata)
		case "$timeWindow":
			criteria, ok := value.(map[string]interface{})
			if !ok {
				return false, fmt.Errorf("$timeWindow requires a criteria object")
			}
			subFlow, ok := criteria["flow"].(map[string]interface{})
			if !ok {
				return false, fmt.Errorf("$timeWindow requires a 'flow' object")
			}

			// Parse the time window
			windowStart, windowEnd, err := parseTimeWindow(criteria)
			if err != nil {
				return false, err
			}

			// Create a sub-metadata map to capture results within the time window
			windowMetadata := make(map[string]interface{})

			// Evaluate the sub-flow with time window constraints
			// We'll need to modify any DB queries to include the time window filter
			// This would require additional implementation in the DB functions
			tempUserID := fmt.Sprintf("%s|%s|%s", userID, windowStart.Format(time.RFC3339), windowEnd.Format(time.RFC3339))
			result, err := re.evaluateFlow(models.JSONB(subFlow), tempUserID, windowMetadata)
			if err != nil {
				return false, err
			}

			// Merge the window metadata with the parent metadata
			for k, v := range windowMetadata {
				metadata[fmt.Sprintf("window_%s", k)] = v
			}

			return result, nil
		}
	}

	return false, errors.New("unsupported flow definition format")
}

// evaluateAndOperator handles the $and operator
func (re *RuleEngine) evaluateAndOperator(conditions interface{}, userID string, metadata map[string]interface{}) (bool, error) {
	conditionsArray, ok := conditions.([]interface{})
	if !ok {
		return false, errors.New("$and operator requires an array of conditions")
	}

	for _, condition := range conditionsArray {
		conditionMap, ok := condition.(map[string]interface{})
		if !ok {
			return false, errors.New("each condition in $and must be an object")
		}

		result, err := re.evaluateFlow(models.JSONB(conditionMap), userID, metadata)
		if err != nil {
			return false, err
		}

		if !result {
			return false, nil // Short-circuit: if any condition is false, the whole AND is false
		}
	}

	return true, nil // All conditions passed
}

// evaluateOrOperator handles the $or operator
func (re *RuleEngine) evaluateOrOperator(conditions interface{}, userID string, metadata map[string]interface{}) (bool, error) {
	conditionsArray, ok := conditions.([]interface{})
	if !ok {
		return false, errors.New("$or operator requires an array of conditions")
	}

	for _, condition := range conditionsArray {
		conditionMap, ok := condition.(map[string]interface{})
		if !ok {
			return false, errors.New("each condition in $or must be an object")
		}

		result, err := re.evaluateFlow(models.JSONB(conditionMap), userID, metadata)
		if err != nil {
			return false, err
		}

		if result {
			return true, nil // Short-circuit: if any condition is true, the whole OR is true
		}
	}

	return false, nil // No conditions passed
}

// evaluateNotOperator handles the $not operator
func (re *RuleEngine) evaluateNotOperator(condition interface{}, userID string, metadata map[string]interface{}) (bool, error) {
	conditionMap, ok := condition.(map[string]interface{})
	if !ok {
		return false, errors.New("$not operator requires a condition object")
	}

	result, err := re.evaluateFlow(models.JSONB(conditionMap), userID, metadata)
	if err != nil {
		return false, err
	}

	return !result, nil // Negate the result
}

// evaluateEventCriteria evaluates criteria against a set of events
func (re *RuleEngine) evaluateEventCriteria(criteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	// Handle count criteria
	if countCriteria, hasCount := criteria["count"].(map[string]interface{}); hasCount {
		return re.evaluateCountCriteria(countCriteria, events, metadata)
	}

	// Filter events based on criteria
	filteredEvents, err := re.filterEvents(criteria, events)
	if err != nil {
		return false, err
	}

	// Store filtered events count in metadata
	metadata["filtered_event_count"] = len(filteredEvents)
	if len(filteredEvents) > 0 {
		metadata["first_event_id"] = filteredEvents[0].ID
		metadata["last_event_id"] = filteredEvents[len(filteredEvents)-1].ID
	}

	// If we get here, the criteria is considered met if there are any events that match
	return len(filteredEvents) > 0, nil
}

// filterEvents filters events based on criteria
func (re *RuleEngine) filterEvents(criteria map[string]interface{}, events []models.Event) ([]models.Event, error) {
	var filteredEvents []models.Event

	for _, event := range events {
		passes, err := re.eventMatchesCriteria(event, criteria)
		if err != nil {
			return nil, err
		}
		if passes {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents, nil
}

// eventMatchesCriteria checks if an event matches the given criteria
func (re *RuleEngine) eventMatchesCriteria(event models.Event, criteria map[string]interface{}) (bool, error) {
	for field, conditionValue := range criteria {
		// Skip the "count" field as it's handled separately
		if field == "count" {
			continue
		}

		// Handle timestamp field specially
		if field == "timestamp" {
			conditionMap, ok := conditionValue.(map[string]interface{})
			if !ok {
				return false, errors.New("timestamp condition must be an object")
			}
			matches, err := re.evaluateTimestampCondition(event.OccurredAt, conditionMap)
			if err != nil {
				return false, err
			}
			if !matches {
				return false, nil
			}
			continue
		}

		// For other fields, check in the payload
		if event.Payload == nil {
			return false, nil
		}

		// Get the field value from the event payload
		fieldValue, exists := event.Payload[field]
		if !exists {
			return false, nil
		}

		// If conditionValue is a comparison operator object
		if conditionMap, ok := conditionValue.(map[string]interface{}); ok {
			matches, err := re.evaluateComparison(fieldValue, conditionMap)
			if err != nil {
				return false, err
			}
			if !matches {
				return false, nil
			}
		} else {
			// Direct equality comparison
			if !reflect.DeepEqual(fieldValue, conditionValue) {
				return false, nil
			}
		}
	}

	return true, nil
}

// evaluateTimestampCondition evaluates timestamp-specific conditions
func (re *RuleEngine) evaluateTimestampCondition(timestamp time.Time, conditions map[string]interface{}) (bool, error) {
	for operator, value := range conditions {
		switch operator {
		case "$gte":
			compareTime, err := parseTimeValue(value)
			if err != nil {
				return false, err
			}
			if !timestamp.Equal(compareTime) && !timestamp.After(compareTime) {
				return false, nil
			}
		case "$gt":
			compareTime, err := parseTimeValue(value)
			if err != nil {
				return false, err
			}
			if !timestamp.After(compareTime) {
				return false, nil
			}
		case "$lte":
			compareTime, err := parseTimeValue(value)
			if err != nil {
				return false, err
			}
			if !timestamp.Equal(compareTime) && !timestamp.Before(compareTime) {
				return false, nil
			}
		case "$lt":
			compareTime, err := parseTimeValue(value)
			if err != nil {
				return false, err
			}
			if !timestamp.Before(compareTime) {
				return false, nil
			}
		case "$eq":
			compareTime, err := parseTimeValue(value)
			if err != nil {
				return false, err
			}
			if !timestamp.Equal(compareTime) {
				return false, nil
			}
		case "$ne":
			compareTime, err := parseTimeValue(value)
			if err != nil {
				return false, err
			}
			if timestamp.Equal(compareTime) {
				return false, nil
			}
		default:
			return false, fmt.Errorf("unsupported timestamp operator: %s", operator)
		}
	}
	return true, nil
}

// parseTimeValue converts a string or RFC3339 time value to time.Time
func parseTimeValue(value interface{}) (time.Time, error) {
	if timeStr, ok := value.(string); ok {
		return time.Parse(time.RFC3339, timeStr)
	}
	return time.Time{}, errors.New("timestamp value must be a string in RFC3339 format")
}

// evaluateComparison evaluates comparison operators on values
func (re *RuleEngine) evaluateComparison(fieldValue interface{}, conditions map[string]interface{}) (bool, error) {
	for operator, compareValue := range conditions {
		switch operator {
		case "$eq":
			if !reflect.DeepEqual(fieldValue, compareValue) {
				return false, nil
			}
		case "$ne":
			if reflect.DeepEqual(fieldValue, compareValue) {
				return false, nil
			}
		case "$gt":
			result, err := compareValues(fieldValue, compareValue, func(a, b float64) bool { return a > b })
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		case "$gte":
			result, err := compareValues(fieldValue, compareValue, func(a, b float64) bool { return a >= b })
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		case "$lt":
			result, err := compareValues(fieldValue, compareValue, func(a, b float64) bool { return a < b })
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		case "$lte":
			result, err := compareValues(fieldValue, compareValue, func(a, b float64) bool { return a <= b })
			if err != nil {
				return false, err
			}
			if !result {
				return false, nil
			}
		case "$in":
			if !isInArray(fieldValue, compareValue) {
				return false, nil
			}
		case "$nin":
			if isInArray(fieldValue, compareValue) {
				return false, nil
			}
		case "$regex":
			return false, errors.New("$regex operator not implemented yet")
		default:
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

// evaluateCountCriteria checks if the number of events meets the count criteria
func (re *RuleEngine) evaluateCountCriteria(countCriteria map[string]interface{}, events []models.Event, metadata map[string]interface{}) (bool, error) {
	// First, filter events based on other criteria in the parent object
	filteredEvents, err := re.filterEvents(map[string]interface{}{}, events)
	if err != nil {
		return false, err
	}

	count := len(filteredEvents)
	metadata["event_count"] = count

	// Evaluate the count criteria
	for operator, value := range countCriteria {
		compareValue, err := toFloat64(value)
		if err != nil {
			return false, fmt.Errorf("invalid count comparison value: %w", err)
		}

		floatCount := float64(count)

		switch operator {
		case "$eq":
			if floatCount != compareValue {
				return false, nil
			}
		case "$ne":
			if floatCount == compareValue {
				return false, nil
			}
		case "$gt":
			if floatCount <= compareValue {
				return false, nil
			}
		case "$gte":
			if floatCount < compareValue {
				return false, nil
			}
		case "$lt":
			if floatCount >= compareValue {
				return false, nil
			}
		case "$lte":
			if floatCount > compareValue {
				return false, nil
			}
		default:
			return false, fmt.Errorf("unsupported count operator: %s", operator)
		}
	}

	return true, nil
}

// ProcessEvents processes a batch of events and awards badges if criteria are met
func (re *RuleEngine) ProcessEvents(userID string) error {
	// Get all active badges
	badges, err := re.DB.GetActiveBadges()
	if err != nil {
		return fmt.Errorf("failed to get active badges: %w", err)
	}

	for _, badge := range badges {
		// Check if user already has this badge
		userBadges, err := re.DB.GetUserBadges(userID)
		if err != nil {
			return fmt.Errorf("failed to get user badges: %w", err)
		}

		hasThisBadge := false
		for _, ub := range userBadges {
			if ub.BadgeID == badge.ID {
				hasThisBadge = true
				break
			}
		}

		if hasThisBadge {
			continue
		}

		// Evaluate badge criteria
		meets, metadata, err := re.EvaluateBadgeCriteria(badge.ID, userID)
		if err != nil {
			// Log the error but continue with other badges
			fmt.Printf("Error evaluating badge %d for user %s: %v\n", badge.ID, userID, err)
			continue
		}

		// If criteria met, award the badge
		if meets {
			userBadge := models.UserBadge{
				UserID:   userID,
				BadgeID:  badge.ID,
				Metadata: models.JSONB(metadata),
			}

			if err := re.DB.AwardBadgeToUser(&userBadge); err != nil {
				return fmt.Errorf("failed to award badge %d to user %s: %w", badge.ID, userID, err)
			}

			fmt.Printf("Badge '%s' (ID: %d) awarded to user %s\n", badge.Name, badge.ID, userID)
		}
	}

	return nil
}

// ProcessEvent processes a single event and checks if it triggers any badge awards
func (re *RuleEngine) ProcessEvent(event *models.Event) error {
	// Process badges for the user who triggered the event
	return re.ProcessEvents(event.UserID)
}
