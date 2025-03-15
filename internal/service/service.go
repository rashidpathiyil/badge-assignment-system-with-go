package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/badge-assignment-system/internal/engine"
	"github.com/badge-assignment-system/internal/models"
)

// Service handles business logic for the badge system
type Service struct {
	DB         *models.DB
	RuleEngine *engine.RuleEngine
}

// NewService creates a new service
func NewService(db *models.DB) *Service {
	return &Service{
		DB:         db,
		RuleEngine: engine.NewRuleEngine(db),
	}
}

// CreateEventType creates a new event type
func (s *Service) CreateEventType(req *models.NewEventTypeRequest) (*models.EventType, error) {
	// Validate request
	if req.Name == "" {
		return nil, errors.New("event type name is required")
	}

	// Check if event type with same name already exists
	_, err := s.DB.GetEventTypeByName(req.Name)
	if err == nil {
		return nil, fmt.Errorf("event type with name '%s' already exists", req.Name)
	}

	// Create event type
	eventType := &models.EventType{
		Name:        req.Name,
		Description: req.Description,
		Schema:      models.JSONB(req.Schema),
	}

	if err := s.DB.CreateEventType(eventType); err != nil {
		return nil, fmt.Errorf("failed to create event type: %w", err)
	}

	return eventType, nil
}

// GetEventTypes gets all event types
func (s *Service) GetEventTypes() ([]models.EventType, error) {
	return s.DB.GetEventTypes()
}

// GetEventTypeByID gets an event type by ID
func (s *Service) GetEventTypeByID(id int) (*models.EventType, error) {
	eventType, err := s.DB.GetEventTypeByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get event type: %w", err)
	}
	return &eventType, nil
}

// UpdateEventType updates an existing event type
func (s *Service) UpdateEventType(id int, req *models.UpdateEventTypeRequest) (*models.EventType, error) {
	// Get existing event type
	eventType, err := s.DB.GetEventTypeByID(id)
	if err != nil {
		return nil, fmt.Errorf("event type not found: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		// Check if name already exists (for different ID)
		existingET, err := s.DB.GetEventTypeByName(req.Name)
		if err == nil && existingET.ID != id {
			return nil, fmt.Errorf("event type with name '%s' already exists", req.Name)
		}
		eventType.Name = req.Name
	}

	if req.Description != "" {
		eventType.Description = req.Description
	}

	if req.Schema != nil {
		eventType.Schema = models.JSONB(req.Schema)
	}

	// Update in database
	if err := s.DB.UpdateEventType(&eventType); err != nil {
		return nil, fmt.Errorf("failed to update event type: %w", err)
	}

	return &eventType, nil
}

// DeleteEventType deletes an event type
func (s *Service) DeleteEventType(id int) error {
	return s.DB.DeleteEventType(id)
}

// CreateBadge creates a new badge with criteria
func (s *Service) CreateBadge(req *models.NewBadgeRequest) (*models.BadgeWithCriteria, error) {
	// Validate request
	if req.Name == "" {
		return nil, errors.New("badge name is required")
	}

	if req.FlowDefinition == nil {
		return nil, errors.New("flow definition is required")
	}

	// Create badge
	badge := &models.Badge{
		Name:        req.Name,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Active:      true,
	}

	// Create criteria
	criteria := &models.BadgeCriteria{
		FlowDefinition: models.JSONB(req.FlowDefinition),
	}

	// Save to database
	if err := s.DB.CreateBadge(badge, criteria); err != nil {
		return nil, fmt.Errorf("failed to create badge: %w", err)
	}

	// Return the created badge with criteria
	return &models.BadgeWithCriteria{
		Badge:    *badge,
		Criteria: *criteria,
	}, nil
}

// GetBadges gets all badges
func (s *Service) GetBadges() ([]models.Badge, error) {
	return s.DB.GetBadges()
}

// GetActiveBadges gets all active badges
func (s *Service) GetActiveBadges() ([]models.Badge, error) {
	return s.DB.GetActiveBadges()
}

// GetBadgeByID gets a badge by ID
func (s *Service) GetBadgeByID(id int) (*models.Badge, error) {
	badge, err := s.DB.GetBadgeByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get badge: %w", err)
	}
	return &badge, nil
}

// GetBadgeWithCriteria gets a badge with its criteria
func (s *Service) GetBadgeWithCriteria(id int) (*models.BadgeWithCriteria, error) {
	badgeWithCriteria, err := s.DB.GetBadgeWithCriteria(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get badge with criteria: %w", err)
	}
	return &badgeWithCriteria, nil
}

// UpdateBadge updates an existing badge
func (s *Service) UpdateBadge(id int, req *models.UpdateBadgeRequest) (*models.Badge, error) {
	// Get existing badge
	badge, err := s.DB.GetBadgeByID(id)
	if err != nil {
		return nil, fmt.Errorf("badge not found: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		badge.Name = req.Name
	}

	if req.Description != "" {
		badge.Description = req.Description
	}

	if req.ImageURL != "" {
		badge.ImageURL = req.ImageURL
	}

	if req.Active != nil {
		badge.Active = *req.Active
	}

	// Prepare criteria if flow definition is provided
	var criteria *models.BadgeCriteria
	if req.FlowDefinition != nil {
		criteria = &models.BadgeCriteria{
			BadgeID:        id,
			FlowDefinition: models.JSONB(req.FlowDefinition),
		}
	}

	// Update in database
	if err := s.DB.UpdateBadge(&badge, criteria); err != nil {
		return nil, fmt.Errorf("failed to update badge: %w", err)
	}

	return &badge, nil
}

// DeleteBadge deletes a badge
func (s *Service) DeleteBadge(id int) error {
	return s.DB.DeleteBadge(id)
}

// ProcessEvent processes an event and potentially awards badges
func (s *Service) ProcessEvent(req *models.NewEventRequest) error {
	// Validate request
	if req.EventType == "" {
		return errors.New("event type is required")
	}

	if req.UserID == "" {
		return errors.New("user ID is required")
	}

	// Get event type
	eventType, err := s.DB.GetEventTypeByName(req.EventType)
	if err != nil {
		return fmt.Errorf("event type '%s' not found: %w", req.EventType, err)
	}

	// Determine the timestamp
	var occurredAt time.Time
	if req.Timestamp != "" {
		occurredAt, err = time.Parse(time.RFC3339, req.Timestamp)
		if err != nil {
			return fmt.Errorf("invalid timestamp format: %w", err)
		}
	} else {
		occurredAt = time.Now()
	}

	// Create and save the event
	event := &models.Event{
		EventTypeID: eventType.ID,
		UserID:      req.UserID,
		Payload:     models.JSONB(req.Payload),
		OccurredAt:  occurredAt,
	}

	if err := s.DB.CreateEvent(event); err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	// Process the event to check if it triggers any badges
	if err := s.RuleEngine.ProcessEvent(event); err != nil {
		return fmt.Errorf("failed to process event for badge evaluation: %w", err)
	}

	return nil
}

// GetUserBadges gets all badges awarded to a user
func (s *Service) GetUserBadges(userID string) ([]map[string]interface{}, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	return s.DB.GetUserBadgeDetails(userID)
}

// CreateConditionType creates a new condition type
func (s *Service) CreateConditionType(req *models.NewConditionTypeRequest) (*models.ConditionType, error) {
	// Validate request
	if req.Name == "" {
		return nil, errors.New("condition type name is required")
	}

	// Create condition type
	conditionType := &models.ConditionType{
		Name:            req.Name,
		Description:     req.Description,
		EvaluationLogic: req.EvaluationLogic,
	}

	if err := s.DB.CreateConditionType(conditionType); err != nil {
		return nil, fmt.Errorf("failed to create condition type: %w", err)
	}

	return conditionType, nil
}

// GetConditionTypes gets all condition types
func (s *Service) GetConditionTypes() ([]models.ConditionType, error) {
	return s.DB.GetAllConditionTypes()
}

// GetConditionTypeByID gets a condition type by ID
func (s *Service) GetConditionTypeByID(id int) (*models.ConditionType, error) {
	conditionType, err := s.DB.GetConditionTypeByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get condition type: %w", err)
	}
	return &conditionType, nil
}

// UpdateConditionType updates an existing condition type
func (s *Service) UpdateConditionType(id int, req *models.UpdateConditionTypeRequest) (*models.ConditionType, error) {
	// Get existing condition type
	conditionType, err := s.DB.GetConditionTypeByID(id)
	if err != nil {
		return nil, fmt.Errorf("condition type not found: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		conditionType.Name = req.Name
	}

	if req.Description != "" {
		conditionType.Description = req.Description
	}

	if req.EvaluationLogic != "" {
		conditionType.EvaluationLogic = req.EvaluationLogic
	}

	// Update in database
	if err := s.DB.UpdateConditionType(&conditionType); err != nil {
		return nil, fmt.Errorf("failed to update condition type: %w", err)
	}

	return &conditionType, nil
}

// DeleteConditionType deletes a condition type
func (s *Service) DeleteConditionType(id int) error {
	return s.DB.DeleteConditionType(id)
}
