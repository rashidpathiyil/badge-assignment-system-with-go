package engine

import (
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/logging"
)

func TestIsDynamicTimeVariable(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"$NOW", true},
		{"$NOW(-30d)", true},
		{"$NOW(-1y)", true},
		{"$NOW(-1y-3M)", true},
		{"2023-12-01T00:00:00Z", false},
		{"not-a-time-var", false},
		{"$NOTAVALIDVAR", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := IsDynamicTimeVariable(test.input)
			if result != test.expected {
				t.Errorf("IsDynamicTimeVariable(%q) = %v, expected %v", test.input, result, test.expected)
			}
		})
	}
}

func TestParseDynamicTimeVariable(t *testing.T) {
	// Create a fixed time cache for testing
	fixedTime := time.Date(2023, 12, 15, 12, 0, 0, 0, time.UTC)
	cache := &TimeVariableCache{
		now: fixedTime,
	}

	tests := []struct {
		input    string
		expected time.Time
		hasError bool
	}{
		{"$NOW", fixedTime, false},
		{"$NOW(-30d)", fixedTime.AddDate(0, 0, -30), false},
		{"$NOW(-1y)", fixedTime.AddDate(-1, 0, 0), false},
		{"$NOW(-1y-3M)", fixedTime.AddDate(-1, -3, 0), false},
		{"$NOW(-1h-30m)", fixedTime.Add(-1*time.Hour - 30*time.Minute), false},
		{"$NOW(not-valid)", time.Time{}, true},
		{"invalid", time.Time{}, true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := ParseDynamicTimeVariable(test.input, cache)

			if test.hasError {
				if err == nil {
					t.Errorf("ParseDynamicTimeVariable(%q) expected error, got nil", test.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseDynamicTimeVariable(%q) unexpected error: %v", test.input, err)
				}

				if !result.Equal(test.expected) {
					t.Errorf("ParseDynamicTimeVariable(%q) = %v, expected %v", test.input, result, test.expected)
				}
			}
		})
	}
}

func TestTimeValueInBadgeCriteria(t *testing.T) {
	// Create a mock badge criteria with dynamic time variables
	criteria := map[string]interface{}{
		"timestamp": map[string]interface{}{
			"$gte": "$NOW(-30d)",
		},
	}

	// Create a rule engine with a fixed time for testing
	re := &RuleEngine{
		TimeVarCache: &TimeVariableCache{
			now: time.Date(2023, 12, 15, 12, 0, 0, 0, time.UTC),
		},
		Logger: logging.NewLogger("TEST", logging.LogLevelInfo),
	}

	// Create a test event with a timestamp 15 days ago (should match)
	eventTime := re.TimeVarCache.now.AddDate(0, 0, -15)

	// Manually test the timestamp condition
	result, err := re.evaluateTimestampCondition(eventTime, criteria["timestamp"].(map[string]interface{}))

	if err != nil {
		t.Errorf("evaluateTimestampCondition unexpected error: %v", err)
	}

	if !result {
		t.Errorf("evaluateTimestampCondition should return true for event 15 days ago with 30 day window")
	}

	// Create a test event with a timestamp 45 days ago (should not match)
	oldEventTime := re.TimeVarCache.now.AddDate(0, 0, -45)

	result, err = re.evaluateTimestampCondition(oldEventTime, criteria["timestamp"].(map[string]interface{}))

	if err != nil {
		t.Errorf("evaluateTimestampCondition unexpected error: %v", err)
	}

	if result {
		t.Errorf("evaluateTimestampCondition should return false for event 45 days ago with 30 day window")
	}
}

func TestComplexTimeVariableCriteria(t *testing.T) {
	// Create a rule engine with a fixed time for testing
	fixedTime := time.Date(2023, 12, 15, 12, 0, 0, 0, time.UTC)
	re := &RuleEngine{
		TimeVarCache: &TimeVariableCache{
			now: fixedTime,
		},
		Logger: logging.NewLogger("TEST", logging.LogLevelInfo),
	}

	// Create a mock complex criteria with multiple time variables
	// Similar to the "Loyal Active Customer" example in the proposal
	criteria := map[string]interface{}{
		"user": map[string]interface{}{
			"created_at": map[string]interface{}{
				"$lte": "$NOW(-1y)", // Account at least 1 year old
			},
			"subscription": map[string]interface{}{
				"expires_at": map[string]interface{}{
					"$gte": "$NOW", // Subscription not expired
				},
			},
		},
		"last_purchase": map[string]interface{}{
			"$gte": "$NOW(-90d)", // Purchase in last 90 days
		},
		"recent_activity": map[string]interface{}{
			"$gte": "$NOW(-30d)", // Activity in last 30 days
		},
	}

	// Test user.created_at condition with account 2 years old (should match)
	userCreatedAt := fixedTime.AddDate(-2, 0, 0)
	result, err := re.parseTimeValueWithCache(criteria["user"].(map[string]interface{})["created_at"].(map[string]interface{})["$lte"])
	if err != nil {
		t.Errorf("parseTimeValueWithCache unexpected error: %v", err)
	}
	if !userCreatedAt.Before(result) {
		t.Errorf("Account 2 years old should be before $NOW(-1y)")
	}

	// Test user.subscription.expires_at condition with unexpired subscription (should match)
	subExpiresAt := fixedTime.AddDate(0, 1, 0) // Expires in 1 month
	result, err = re.parseTimeValueWithCache(criteria["user"].(map[string]interface{})["subscription"].(map[string]interface{})["expires_at"].(map[string]interface{})["$gte"])
	if err != nil {
		t.Errorf("parseTimeValueWithCache unexpected error: %v", err)
	}
	if !subExpiresAt.After(result) {
		t.Errorf("Subscription expiring in 1 month should be after $NOW")
	}

	// Test last_purchase condition with purchase 60 days ago (should match)
	lastPurchase := fixedTime.AddDate(0, 0, -60)
	result, err = re.parseTimeValueWithCache(criteria["last_purchase"].(map[string]interface{})["$gte"])
	if err != nil {
		t.Errorf("parseTimeValueWithCache unexpected error: %v", err)
	}
	if lastPurchase.Before(result) {
		t.Errorf("Purchase 60 days ago should be after or equal to $NOW(-90d)")
	}

	// Test recent_activity condition with activity 45 days ago (should not match)
	recentActivity := fixedTime.AddDate(0, 0, -45)
	result, err = re.parseTimeValueWithCache(criteria["recent_activity"].(map[string]interface{})["$gte"])
	if err != nil {
		t.Errorf("parseTimeValueWithCache unexpected error: %v", err)
	}
	if recentActivity.After(result) || recentActivity.Equal(result) {
		t.Errorf("Activity 45 days ago should be before $NOW(-30d)")
	}
}
