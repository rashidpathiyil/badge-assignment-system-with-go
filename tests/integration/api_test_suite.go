package integration

import (
	"testing"

	"github.com/badge-assignment-system/internal/testutil"
)

// TestAPIIntegrationSuite runs the API integration tests in the correct order
// to ensure proper dependencies between tests
func TestAPIIntegrationSuite(t *testing.T) {
	// Skip if API is not ready
	if !IsAPIReady() {
		SkipIfNotReady(t, "API is not available")
		return
	}

	// Run the test suite in order - only including tests that we've implemented
	tests := []testing.InternalTest{
		{Name: "LogicalOperators", F: TestLogicalOperators},
		{Name: "TimeBasedCriteria", F: TestTimeBasedCriteria},
		{Name: "NegativeScenarios", F: TestNegativeScenarios},
		{Name: "InactiveBadge", F: TestInactiveBadge},
	}

	RunTestSuite(t, tests)
}

// IsAPIReady checks if the API is ready for testing
func IsAPIReady() bool {
	resp := testutil.MakeRequest("GET", "/health", nil)
	return resp.Error == nil && resp.StatusCode == 200
}
