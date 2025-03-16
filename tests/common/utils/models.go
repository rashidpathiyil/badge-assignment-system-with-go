package utils

import "time"

// EventType represents an event type in the API
type EventType struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ConditionType represents a condition type in the API
// WARNING: This feature is not fully implemented in the current system.
// The database schema and API endpoints exist, but there is no JavaScript engine
// to execute the evaluation_logic field. Badge criteria use the direct JSON-based
// approach with operators instead.
type ConditionType struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	EvaluationLogic string    `json:"evaluation_logic"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Badge represents a badge in the API
type Badge struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BadgeCriteria represents badge criteria in the API
type BadgeCriteria struct {
	ID             int                    `json:"id"`
	BadgeID        int                    `json:"badge_id"`
	FlowDefinition map[string]interface{} `json:"flow_definition"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// BadgeWithCriteria combines Badge and BadgeCriteria
type BadgeWithCriteria struct {
	Badge    Badge         `json:"badge"`
	Criteria BadgeCriteria `json:"criteria"`
}

// UserBadge represents a badge awarded to a user
type UserBadge struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	AwardedAt   time.Time `json:"awarded_at"`
	Metadata    string    `json:"metadata"`
}

// EventTypeRequest is used to create a new event type
type EventTypeRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
}

// ConditionTypeRequest is used to create a new condition type
// WARNING: This feature is not fully implemented in the current system.
type ConditionTypeRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	EvaluationLogic string `json:"evaluation_logic"`
}

// BadgeRequest is used to create a new badge
type BadgeRequest struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	ImageURL       string                 `json:"image_url"`
	FlowDefinition map[string]interface{} `json:"flow_definition"`
}

// EventRequest is used to submit an event
type EventRequest struct {
	EventType string                 `json:"event_type"`
	UserID    string                 `json:"user_id"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp string                 `json:"timestamp,omitempty"`
}
