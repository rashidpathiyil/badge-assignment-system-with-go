package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationSetup(t *testing.T) {
	// This test verifies that the integration test setup is working
	assert.NotZero(t, EventTypeID, "EventTypeID should be initialized")
	assert.NotZero(t, BadgeID, "BadgeID should be initialized")
	assert.Equal(t, "test-user-integration", TestUserID, "TestUserID should be initialized")
}
