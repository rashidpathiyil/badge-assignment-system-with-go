package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONB is a type for handling PostgreSQL JSONB data
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, &j)
}

// EventType represents the event_types table
type EventType struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Schema      JSONB     `db:"schema" json:"schema"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// ConditionType represents the condition_types table
// WARNING: This feature is not fully implemented in the current system.
// The database schema and API endpoints exist, but there is no JavaScript engine
// to execute the evaluation_logic field. Badge criteria use the direct JSON-based
// approach with operators instead.
type ConditionType struct {
	ID              int       `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	Description     string    `db:"description" json:"description"`
	EvaluationLogic string    `db:"evaluation_logic" json:"evaluation_logic"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// Badge represents the badges table
type Badge struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	ImageURL    string    `db:"image_url" json:"image_url"`
	Active      bool      `db:"active" json:"active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// BadgeCriteria represents the badge_criteria table
type BadgeCriteria struct {
	ID             int       `db:"id" json:"id"`
	BadgeID        int       `db:"badge_id" json:"badge_id"`
	FlowDefinition JSONB     `db:"flow_definition" json:"flow_definition"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// UserBadge represents the user_badges table
type UserBadge struct {
	ID        int       `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	BadgeID   int       `db:"badge_id" json:"badge_id"`
	AwardedAt time.Time `db:"awarded_at" json:"awarded_at"`
	Metadata  JSONB     `db:"metadata" json:"metadata"`
}

// Event represents the events table
type Event struct {
	ID          int       `db:"id" json:"id"`
	EventTypeID int       `db:"event_type_id" json:"event_type_id"`
	UserID      string    `db:"user_id" json:"user_id"`
	Payload     JSONB     `db:"payload" json:"payload"`
	OccurredAt  time.Time `db:"occurred_at" json:"occurred_at"`
}

// BadgeWithCriteria combines Badge and BadgeCriteria for easier handling
type BadgeWithCriteria struct {
	Badge    Badge         `json:"badge"`
	Criteria BadgeCriteria `json:"criteria"`
}

// EventTypeWithSchema represents EventType with accessible schema
type EventTypeWithSchema struct {
	EventType
	SchemaFields map[string]interface{} `json:"schema_fields"`
}

// NewEventRequest is used for creating a new event
type NewEventRequest struct {
	EventType string                 `json:"event_type"`
	UserID    string                 `json:"user_id"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp string                 `json:"timestamp,omitempty"`
}

// NewBadgeRequest is used for creating a new badge
type NewBadgeRequest struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	ImageURL       string                 `json:"image_url"`
	FlowDefinition map[string]interface{} `json:"flow_definition"`
}

// UpdateBadgeRequest is used for updating an existing badge
type UpdateBadgeRequest struct {
	Name           string                 `json:"name,omitempty"`
	Description    string                 `json:"description,omitempty"`
	ImageURL       string                 `json:"image_url,omitempty"`
	Active         *bool                  `json:"active,omitempty"`
	FlowDefinition map[string]interface{} `json:"flow_definition,omitempty"`
}

// NewEventTypeRequest is used for creating a new event type
type NewEventTypeRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
}

// UpdateEventTypeRequest is used for updating an existing event type
type UpdateEventTypeRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
}

// NewConditionTypeRequest is used for creating a new condition type
// WARNING: This feature is not fully implemented in the current system.
type NewConditionTypeRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	EvaluationLogic string `json:"evaluation_logic"`
}

// UpdateConditionTypeRequest is used for updating an existing condition type
// WARNING: This feature is not fully implemented in the current system.
type UpdateConditionTypeRequest struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	EvaluationLogic string `json:"evaluation_logic,omitempty"`
}
