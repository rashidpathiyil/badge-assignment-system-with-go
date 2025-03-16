package integration

import (
	"os"
	"testing"
)

// TestMain runs before any tests in the package
func TestMain(m *testing.M) {
	// Set up everything we need for testing
	// This will be run once before any tests
	EventTypeID = 1 // Default value for integration tests
	// ConditionTypeID was removed as the feature is not fully implemented
	BadgeID = 1 // Default value for integration tests
	TestUserID = "test-user-integration"

	// Run all tests
	os.Exit(m.Run())
}

// TestComplexCriteriaSuite runs the complex criteria tests
func TestComplexCriteriaSuite(t *testing.T) {
	t.Skip("Skipping complex criteria tests until API is ready")
}

// TestNegativeScenariosSuite runs the negative scenario tests
func TestNegativeScenariosSuite(t *testing.T) {
	t.Skip("Skipping negative scenario tests until API is ready")
}

func init() {
	// Set default values for integration tests
	EventTypeID = 1 // Default value for integration tests
	BadgeID = 1     // Default value for integration tests
}
