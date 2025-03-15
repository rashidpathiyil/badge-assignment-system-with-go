package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB is the database connection
type DB struct {
	*sqlx.DB
}

// NewDB creates a new database connection
func NewDB() (*DB, error) {
	// Get database connection details from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "badge_system")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Open database connection
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to the database successfully")
	return &DB{db}, nil
}

// Helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEventTypes retrieves all event types
func (db *DB) GetEventTypes() ([]EventType, error) {
	var eventTypes []EventType
	err := db.Select(&eventTypes, "SELECT * FROM event_types ORDER BY name")
	return eventTypes, err
}

// GetEventTypeByID retrieves an event type by ID
func (db *DB) GetEventTypeByID(id int) (EventType, error) {
	var eventType EventType
	err := db.Get(&eventType, "SELECT * FROM event_types WHERE id = $1", id)
	return eventType, err
}

// GetEventTypeByName retrieves an event type by name
func (db *DB) GetEventTypeByName(name string) (EventType, error) {
	var eventType EventType
	err := db.Get(&eventType, "SELECT * FROM event_types WHERE name = $1", name)
	return eventType, err
}

// CreateEventType creates a new event type
func (db *DB) CreateEventType(et *EventType) error {
	query := `
		INSERT INTO event_types (name, description, schema)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	return db.QueryRow(query, et.Name, et.Description, et.Schema).Scan(&et.ID, &et.CreatedAt, &et.UpdatedAt)
}

// UpdateEventType updates an existing event type
func (db *DB) UpdateEventType(et *EventType) error {
	query := `
		UPDATE event_types
		SET name = $1, description = $2, schema = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`
	return db.QueryRow(query, et.Name, et.Description, et.Schema, et.ID).Scan(&et.UpdatedAt)
}

// DeleteEventType deletes an event type
func (db *DB) DeleteEventType(id int) error {
	_, err := db.Exec("DELETE FROM event_types WHERE id = $1", id)
	return err
}

// GetBadges retrieves all badges
func (db *DB) GetBadges() ([]Badge, error) {
	var badges []Badge
	err := db.Select(&badges, "SELECT * FROM badges ORDER BY name")
	return badges, err
}

// GetActiveBadges retrieves all active badges
func (db *DB) GetActiveBadges() ([]Badge, error) {
	var badges []Badge
	err := db.Select(&badges, "SELECT * FROM badges WHERE active = true ORDER BY name")
	return badges, err
}

// GetBadgeByID retrieves a badge by ID
func (db *DB) GetBadgeByID(id int) (Badge, error) {
	var badge Badge
	err := db.Get(&badge, "SELECT * FROM badges WHERE id = $1", id)
	return badge, err
}

// GetBadgeWithCriteria retrieves a badge with its criteria
func (db *DB) GetBadgeWithCriteria(id int) (BadgeWithCriteria, error) {
	var result BadgeWithCriteria

	// Get the badge
	badge, err := db.GetBadgeByID(id)
	if err != nil {
		return result, err
	}
	result.Badge = badge

	// Get the criteria
	var criteria BadgeCriteria
	err = db.Get(&criteria, "SELECT * FROM badge_criteria WHERE badge_id = $1", id)
	if err != nil {
		return result, err
	}
	result.Criteria = criteria

	return result, nil
}

// CreateBadge creates a new badge and its criteria
func (db *DB) CreateBadge(badge *Badge, criteria *BadgeCriteria) error {
	// Start a transaction
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert badge
	query := `
		INSERT INTO badges (name, description, image_url, active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	err = tx.QueryRow(query, badge.Name, badge.Description, badge.ImageURL, badge.Active).
		Scan(&badge.ID, &badge.CreatedAt, &badge.UpdatedAt)
	if err != nil {
		return err
	}

	// Set badge ID in criteria
	criteria.BadgeID = badge.ID

	// Insert criteria
	query = `
		INSERT INTO badge_criteria (badge_id, flow_definition)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`
	err = tx.QueryRow(query, criteria.BadgeID, criteria.FlowDefinition).
		Scan(&criteria.ID, &criteria.CreatedAt, &criteria.UpdatedAt)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// UpdateBadge updates an existing badge and its criteria
func (db *DB) UpdateBadge(badge *Badge, criteria *BadgeCriteria) error {
	// Start a transaction
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update badge
	badgeQuery := `
		UPDATE badges
		SET name = $1, description = $2, image_url = $3, active = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at`
	err = tx.QueryRow(badgeQuery, badge.Name, badge.Description, badge.ImageURL, badge.Active, badge.ID).
		Scan(&badge.UpdatedAt)
	if err != nil {
		return err
	}

	// If criteria update is requested
	if criteria != nil && criteria.FlowDefinition != nil {
		// Check if criteria exists
		var count int
		err = tx.Get(&count, "SELECT COUNT(*) FROM badge_criteria WHERE badge_id = $1", badge.ID)
		if err != nil {
			return err
		}

		if count > 0 {
			// Update existing criteria
			criteriaQuery := `
				UPDATE badge_criteria
				SET flow_definition = $1, updated_at = NOW()
				WHERE badge_id = $2
				RETURNING id, updated_at`
			err = tx.QueryRow(criteriaQuery, criteria.FlowDefinition, badge.ID).
				Scan(&criteria.ID, &criteria.UpdatedAt)
		} else {
			// Insert new criteria
			criteriaQuery := `
				INSERT INTO badge_criteria (badge_id, flow_definition)
				VALUES ($1, $2)
				RETURNING id, created_at, updated_at`
			err = tx.QueryRow(criteriaQuery, badge.ID, criteria.FlowDefinition).
				Scan(&criteria.ID, &criteria.CreatedAt, &criteria.UpdatedAt)
		}
		if err != nil {
			return err
		}
	}

	// Commit transaction
	return tx.Commit()
}

// DeleteBadge deletes a badge and its criteria
func (db *DB) DeleteBadge(id int) error {
	_, err := db.Exec("DELETE FROM badges WHERE id = $1", id)
	return err
}

// CreateEvent creates a new event
func (db *DB) CreateEvent(event *Event) error {
	query := `
		INSERT INTO events (event_type_id, user_id, payload, occurred_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	return db.QueryRow(query, event.EventTypeID, event.UserID, event.Payload, event.OccurredAt).
		Scan(&event.ID)
}

// GetUserEvents retrieves all events for a specific user
func (db *DB) GetUserEvents(userID string) ([]Event, error) {
	var events []Event
	err := db.Select(&events, "SELECT * FROM events WHERE user_id = $1 ORDER BY occurred_at DESC", userID)
	return events, err
}

// GetUserEventsByType retrieves events of a specific type for a user
func (db *DB) GetUserEventsByType(userID string, eventTypeID int) ([]Event, error) {
	var events []Event
	err := db.Select(&events, "SELECT * FROM events WHERE user_id = $1 AND event_type_id = $2 ORDER BY occurred_at DESC",
		userID, eventTypeID)
	return events, err
}

// GetUserBadges retrieves all badges awarded to a user
func (db *DB) GetUserBadges(userID string) ([]UserBadge, error) {
	var userBadges []UserBadge
	err := db.Select(&userBadges, "SELECT * FROM user_badges WHERE user_id = $1 ORDER BY awarded_at DESC", userID)
	return userBadges, err
}

// AwardBadgeToUser awards a badge to a user
func (db *DB) AwardBadgeToUser(userBadge *UserBadge) error {
	// Check if user already has this badge
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM user_badges WHERE user_id = $1 AND badge_id = $2",
		userBadge.UserID, userBadge.BadgeID)
	if err != nil {
		return err
	}

	// If user already has the badge, don't award it again
	if count > 0 {
		return nil
	}

	// Award the badge
	query := `
		INSERT INTO user_badges (user_id, badge_id, metadata)
		VALUES ($1, $2, $3)
		RETURNING id, awarded_at`
	return db.QueryRow(query, userBadge.UserID, userBadge.BadgeID, userBadge.Metadata).
		Scan(&userBadge.ID, &userBadge.AwardedAt)
}

// GetUserBadgeDetails retrieves details of badges awarded to a user
func (db *DB) GetUserBadgeDetails(userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT b.id, b.name, b.description, b.image_url, ub.awarded_at, ub.metadata
		FROM user_badges ub
		JOIN badges b ON ub.badge_id = b.id
		WHERE ub.user_id = $1
		ORDER BY ub.awarded_at DESC`

	rows, err := db.Queryx(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		err := rows.MapScan(row)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}

// GetAllConditionTypes retrieves all condition types
func (db *DB) GetAllConditionTypes() ([]ConditionType, error) {
	var conditionTypes []ConditionType
	err := db.Select(&conditionTypes, "SELECT * FROM condition_types ORDER BY name")
	return conditionTypes, err
}

// GetConditionTypeByID retrieves a condition type by ID
func (db *DB) GetConditionTypeByID(id int) (ConditionType, error) {
	var conditionType ConditionType
	err := db.Get(&conditionType, "SELECT * FROM condition_types WHERE id = $1", id)
	return conditionType, err
}

// CreateConditionType creates a new condition type
func (db *DB) CreateConditionType(ct *ConditionType) error {
	query := `
		INSERT INTO condition_types (name, description, evaluation_logic)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	return db.QueryRow(query, ct.Name, ct.Description, ct.EvaluationLogic).
		Scan(&ct.ID, &ct.CreatedAt, &ct.UpdatedAt)
}

// UpdateConditionType updates an existing condition type
func (db *DB) UpdateConditionType(ct *ConditionType) error {
	query := `
		UPDATE condition_types
		SET name = $1, description = $2, evaluation_logic = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`
	return db.QueryRow(query, ct.Name, ct.Description, ct.EvaluationLogic, ct.ID).
		Scan(&ct.UpdatedAt)
}

// DeleteConditionType deletes a condition type
func (db *DB) DeleteConditionType(id int) error {
	_, err := db.Exec("DELETE FROM condition_types WHERE id = $1", id)
	return err
}
