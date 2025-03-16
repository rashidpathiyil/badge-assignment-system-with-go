package integration

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/badge-assignment-system/internal/testutil"
)

// Setup constants
const (
	DefaultTimeout  = 30 * time.Second
	DefaultInterval = 1 * time.Second
)

// Global variables for shared test data
var (
	EventTypeID int = 1 // Default value for testing
	// ConditionTypeID int = 1 // Default value for testing - removed as feature not implemented
	BadgeID    int = 1 // Default value for testing
	TestUserID     = "test-user-integration"
)

// SkipIfNotReady skips a test if the required resources aren't ready
func SkipIfNotReady(t *testing.T, reason string) {
	t.Skipf("Skipping test: %s", reason)
}

// SetupTest sets up the test environment
func SetupTest() {
	// Check if the API is available
	waitForAPI()
}

// waitForAPI waits for the API to be available
func waitForAPI() {
	log.Println("Waiting for API to be available...")

	timeout := DefaultTimeout
	interval := DefaultInterval

	// Override with environment variables if set
	if val := os.Getenv("TEST_API_TIMEOUT"); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			timeout = duration
		}
	}

	if val := os.Getenv("TEST_API_INTERVAL"); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			interval = duration
		}
	}

	start := time.Now()
	for {
		if time.Since(start) > timeout {
			log.Fatalf("Timed out waiting for API to be available after %v", timeout)
		}

		resp := testutil.MakeRequest("GET", "/health", nil)
		if resp.Error == nil && resp.StatusCode == 200 {
			log.Println("API is available!")
			return
		}

		log.Printf("API not available yet, waiting %v...", interval)
		time.Sleep(interval)
	}
}

// RunTestSuite runs a test suite in the correct order
func RunTestSuite(t *testing.T, tests []testing.InternalTest) {
	for i, test := range tests {
		testName := fmt.Sprintf("%d-%s", i+1, test.Name)
		t.Run(testName, test.F)
	}
}
