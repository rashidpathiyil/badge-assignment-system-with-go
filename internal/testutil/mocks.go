package testutil

import (
	"time"

	"github.com/badge-assignment-system/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the database interface that can be used
// in place of models.DB in tests
type MockDB struct {
	mock.Mock
}

// NewMockDB creates a new MockDB
func NewMockDB() *MockDB {
	return &MockDB{}
}

// GetBadgeWithCriteria mocks retrieving a badge with its criteria
func (m *MockDB) GetBadgeWithCriteria(badgeID int) (models.BadgeWithCriteria, error) {
	args := m.Called(badgeID)
	if args.Get(0) == nil {
		return models.BadgeWithCriteria{}, args.Error(1)
	}
	// Convert from pointer to value
	return *args.Get(0).(*models.BadgeWithCriteria), args.Error(1)
}

// GetActiveBadges mocks retrieving all active badges
func (m *MockDB) GetActiveBadges() ([]models.Badge, error) {
	args := m.Called()
	return args.Get(0).([]models.Badge), args.Error(1)
}

// GetUserBadges mocks retrieving badges for a user
func (m *MockDB) GetUserBadges(userID string) ([]models.UserBadge, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.UserBadge), args.Error(1)
}

// AwardBadgeToUser mocks awarding a badge to a user
func (m *MockDB) AwardBadgeToUser(userBadge *models.UserBadge) error {
	args := m.Called(userBadge)
	return args.Error(0)
}

// GetEventsByType mocks retrieving events by type
func (m *MockDB) GetEventsByType(userID string, eventTypeIDs []int, startTime, endTime interface{}) ([]models.Event, error) {
	args := m.Called(userID, eventTypeIDs, startTime, endTime)
	return args.Get(0).([]models.Event), args.Error(1)
}

// CreateEvent mocks creating a new event
func (m *MockDB) CreateEvent(event *models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

// GetUserEventsByType mocks retrieving events of a specific type for a user
func (m *MockDB) GetUserEventsByType(userID string, eventTypeID int) ([]models.Event, error) {
	args := m.Called(userID, eventTypeID)
	return args.Get(0).([]models.Event), args.Error(1)
}

// GetUserEvents mocks retrieving all events for a specific user
func (m *MockDB) GetUserEvents(userID string) ([]models.Event, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Event), args.Error(1)
}

// GetEventType mocks retrieving an event type
func (m *MockDB) GetEventType(eventTypeID int) (*models.EventType, error) {
	args := m.Called(eventTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EventType), args.Error(1)
}

// GetEventTypeByName mocks retrieving an event type by name
func (m *MockDB) GetEventTypeByName(name string) (models.EventType, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return models.EventType{}, args.Error(1)
	}
	// If a pointer was passed, convert it to a value
	if et, ok := args.Get(0).(*models.EventType); ok {
		return *et, args.Error(1)
	}
	return args.Get(0).(models.EventType), args.Error(1)
}

// GetConditionType mocks retrieving a condition type
func (m *MockDB) GetConditionType(conditionTypeID int) (*models.ConditionType, error) {
	args := m.Called(conditionTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ConditionType), args.Error(1)
}

// CreateTestBadgeWithCriteria creates a test badge with criteria for testing
func CreateTestBadgeWithCriteria(badgeID int, name string, criteria map[string]interface{}) *models.BadgeWithCriteria {
	return &models.BadgeWithCriteria{
		Badge: models.Badge{
			ID:          badgeID,
			Name:        name,
			Description: "Test Badge Description",
			ImageURL:    "https://example.com/badge.png",
			Active:      true,
		},
		Criteria: models.BadgeCriteria{
			ID:             badgeID,
			BadgeID:        badgeID,
			FlowDefinition: models.JSONB(criteria),
		},
	}
}

// CreateTestEvent creates a test event for testing
func CreateTestEvent(id int, userID string, eventTypeID int, payload map[string]interface{}) models.Event {
	return models.Event{
		ID:          id,
		UserID:      userID,
		EventTypeID: eventTypeID,
		OccurredAt:  time.Now(),
		Payload:     models.JSONB(payload),
	}
}
