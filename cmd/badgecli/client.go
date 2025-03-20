package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIClient provides methods to interact with the badge API
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Badge represents a badge in the system
type Badge struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BadgeCriteria represents the criteria for a badge
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

// NewBadgeRequest is used for creating a new badge
type NewBadgeRequest struct {
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	ImageURL       string                 `json:"image_url"`
	FlowDefinition map[string]interface{} `json:"flow_definition"`
}

// EventType represents an event type in the system
type EventType struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// NewEventTypeRequest is used for creating a new event type
type NewEventTypeRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetBadges gets all badges from the system
func (c *APIClient) GetBadges() ([]Badge, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/api/v1/badges", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var badges []Badge
	if err := json.NewDecoder(resp.Body).Decode(&badges); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return badges, nil
}

// GetBadgeWithCriteria gets a badge with its criteria by ID
func (c *APIClient) GetBadgeWithCriteria(id string) (*BadgeWithCriteria, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/api/v1/admin/badges/%s/criteria", c.BaseURL, id))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var badgeWithCriteria BadgeWithCriteria
	if err := json.NewDecoder(resp.Body).Decode(&badgeWithCriteria); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &badgeWithCriteria, nil
}

// CreateBadge creates a new badge
func (c *APIClient) CreateBadge(badge *NewBadgeRequest) (*BadgeWithCriteria, error) {
	jsonData, err := json.Marshal(badge)
	if err != nil {
		return nil, fmt.Errorf("error serializing badge: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/api/v1/admin/badges", c.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var badgeWithCriteria BadgeWithCriteria
	if err := json.NewDecoder(resp.Body).Decode(&badgeWithCriteria); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &badgeWithCriteria, nil
}

// GetEventTypes gets all event types from the system
func (c *APIClient) GetEventTypes() ([]EventType, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/api/v1/admin/event-types", c.BaseURL))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var eventTypes []EventType
	if err := json.NewDecoder(resp.Body).Decode(&eventTypes); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return eventTypes, nil
}

// GetEventTypeByID gets an event type by ID
func (c *APIClient) GetEventTypeByID(id string) (*EventType, error) {
	resp, err := c.HTTPClient.Get(fmt.Sprintf("%s/api/v1/admin/event-types/%s", c.BaseURL, id))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var eventType EventType
	if err := json.NewDecoder(resp.Body).Decode(&eventType); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &eventType, nil
}

// CreateEventType creates a new event type
func (c *APIClient) CreateEventType(eventType *NewEventTypeRequest) (*EventType, error) {
	jsonData, err := json.Marshal(eventType)
	if err != nil {
		return nil, fmt.Errorf("error serializing event type: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		fmt.Sprintf("%s/api/v1/admin/event-types", c.BaseURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var createdEventType EventType
	if err := json.NewDecoder(resp.Body).Decode(&createdEventType); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &createdEventType, nil
}
