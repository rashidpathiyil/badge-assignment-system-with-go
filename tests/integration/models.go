package integration

import "time"

// EventTypeRequest represents a request to create an event type
type EventTypeRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
}

// EventTypeResponse represents a response for an event type
type EventTypeResponse struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// NOTE: ConditionTypeRequest and ConditionTypeResponse structs have been removed
// as the condition types feature is not fully implemented in the current system.

// BadgeRequest represents a request to create a badge
type BadgeRequest struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	ImageURL       string                 `json:"image_url"`
	FlowDefinition map[string]interface{} `json:"flow_definition"`
	Active         bool                   `json:"active"`
}

// BadgeResponse represents a response for a badge
type BadgeResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BadgeCriteriaResponse represents a response for a badge with criteria
type BadgeCriteriaResponse struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	ImageURL    string                 `json:"image_url"`
	Active      bool                   `json:"active"`
	Criteria    map[string]interface{} `json:"criteria"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// EventRequest represents a request to process an event
type EventRequest struct {
	EventType string                 `json:"event_type"`
	UserID    string                 `json:"user_id"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp,omitempty"`
}

// EventResponse represents a response for an event
type EventResponse struct {
	ID          int                    `json:"id"`
	EventTypeID int                    `json:"event_type_id"`
	UserID      string                 `json:"user_id"`
	Data        map[string]interface{} `json:"data"`
	ProcessedAt time.Time              `json:"processed_at"`
	CreatedAt   time.Time              `json:"created_at"`
}

// UserBadgeResponse represents a badge assigned to a user
type UserBadgeResponse struct {
	BadgeID     int       `json:"badge_id"`
	UserID      string    `json:"user_id"`
	BadgeName   string    `json:"badge_name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	AwardedAt   time.Time `json:"awarded_at"`
}

// UserBadgesResponse represents a list of badges assigned to a user
type UserBadgesResponse struct {
	UserID string              `json:"user_id"`
	Badges []UserBadgeResponse `json:"badges"`
}
