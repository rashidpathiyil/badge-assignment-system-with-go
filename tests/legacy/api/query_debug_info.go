package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// TestQueryDebugInfo queries the database directly to examine events and badge data
func TestQueryDebugInfo(t *testing.T) {
	// Skip in normal test runs unless explicitly requested
	if os.Getenv("RUN_DEBUG_TESTS") != "true" {
		t.Skip("Skipping debug test. Set RUN_DEBUG_TESTS=true to run this test.")
	}

	// Read DB connection parameters from environment
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "badge_system"
	}

	// Create log file
	logFile, err := os.Create("db_debug.log")
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	log := func(format string, args ...interface{}) {
		msg := fmt.Sprintf(format, args...)
		t.Log(msg)
		fmt.Fprintln(logFile, msg)
	}

	log("Starting database debug query test")

	// Connect to database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log("Error connecting to database: %v", err)
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log("Failed to ping database: %v", err)
		t.Fatalf("Failed to ping database: %v", err)
	}
	log("Successfully connected to database")

	// Query recent event types
	log("\n--- Recent Event Types ---")
	rows, err := db.Query("SELECT id, name, schema FROM event_types ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		log("Error querying event types: %v", err)
		t.Fatalf("Error querying event types: %v", err)
	}

	for rows.Next() {
		var id int
		var name string
		var schema string
		err = rows.Scan(&id, &name, &schema)
		if err != nil {
			log("Error scanning row: %v", err)
			continue
		}
		log("Event Type ID: %d, Name: %s, Schema: %s", id, name, schema)
	}
	rows.Close()

	// Query recent condition types
	log("\n--- Recent Condition Types ---")
	rows, err = db.Query("SELECT id, name, evaluation_logic FROM condition_types ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		log("Error querying condition types: %v", err)
	} else {
		for rows.Next() {
			var id int
			var name string
			var logic string
			err = rows.Scan(&id, &name, &logic)
			if err != nil {
				log("Error scanning row: %v", err)
				continue
			}
			log("Condition Type ID: %d, Name: %s", id, name)
			log("Logic: %s", logic)
		}
		rows.Close()
	}

	// Query recent badges
	log("\n--- Recent Badges ---")
	rows, err = db.Query("SELECT id, name, description FROM badges ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		log("Error querying badges: %v", err)
	} else {
		for rows.Next() {
			var id int
			var name, description string
			err = rows.Scan(&id, &name, &description)
			if err != nil {
				log("Error scanning row: %v", err)
				continue
			}
			log("Badge ID: %d, Name: %s, Description: %s", id, name, description)
		}
		rows.Close()
	}

	// Query recent badge criteria
	log("\n--- Recent Badge Criteria ---")
	rows, err = db.Query("SELECT id, badge_id, flow_definition FROM badge_criteria ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		log("Error querying badge criteria: %v", err)
	} else {
		for rows.Next() {
			var id, badgeID int
			var flowDefinition string
			err = rows.Scan(&id, &badgeID, &flowDefinition)
			if err != nil {
				log("Error scanning row: %v", err)
				continue
			}
			log("Criteria ID: %d, Badge ID: %d", id, badgeID)
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, []byte(flowDefinition), "", "  ")
			if err != nil {
				log("Raw flow definition: %s", flowDefinition)
			} else {
				log("Flow definition: %s", prettyJSON.String())
			}
		}
		rows.Close()
	}

	// Query recent events
	log("\n--- Recent Events ---")
	rows, err = db.Query(`
		SELECT e.id, e.user_id, et.name as event_type, e.payload, e.created_at 
		FROM events e
		JOIN event_types et ON e.event_type_id = et.id
		ORDER BY e.created_at DESC LIMIT 10
	`)
	if err != nil {
		log("Error querying events: %v", err)
	} else {
		for rows.Next() {
			var id int
			var userID, eventType, payload string
			var createdAt time.Time
			err = rows.Scan(&id, &userID, &eventType, &payload, &createdAt)
			if err != nil {
				log("Error scanning row: %v", err)
				continue
			}
			log("Event ID: %d, User ID: %s, Type: %s, Created: %v", id, userID, eventType, createdAt)
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, []byte(payload), "", "  ")
			if err != nil {
				log("Raw payload: %s", payload)
			} else {
				log("Payload: %s", prettyJSON.String())
			}
		}
		rows.Close()
	}

	// Query recent user badges
	log("\n--- Recent User Badges ---")
	rows, err = db.Query(`
		SELECT ub.id, ub.user_id, b.name as badge_name, ub.awarded_at 
		FROM user_badges ub
		JOIN badges b ON ub.badge_id = b.id
		ORDER BY ub.awarded_at DESC LIMIT 10
	`)
	if err != nil {
		log("Error querying user badges: %v", err)
	} else {
		for rows.Next() {
			var id int
			var userID, badgeName string
			var awardedAt time.Time
			err = rows.Scan(&id, &userID, &badgeName, &awardedAt)
			if err != nil {
				log("Error scanning row: %v", err)
				continue
			}
			log("User Badge ID: %d, User ID: %s, Badge: %s, Awarded: %v", id, userID, badgeName, awardedAt)
		}
		rows.Close()
	}

	// Look for debug_user specifically from recent tests
	log("\n--- Debug User Events ---")
	rows, err = db.Query(`
		SELECT e.id, e.user_id, et.name as event_type, e.payload, e.created_at 
		FROM events e
		JOIN event_types et ON e.event_type_id = et.id
		WHERE e.user_id LIKE 'debug-user-%'
		ORDER BY e.created_at DESC LIMIT 10
	`)
	if err != nil {
		log("Error querying debug user events: %v", err)
	} else {
		for rows.Next() {
			var id int
			var userID, eventType, payload string
			var createdAt time.Time
			err = rows.Scan(&id, &userID, &eventType, &payload, &createdAt)
			if err != nil {
				log("Error scanning row: %v", err)
				continue
			}
			log("Event ID: %d, User ID: %s, Type: %s, Created: %v", id, userID, eventType, createdAt)
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, []byte(payload), "", "  ")
			if err != nil {
				log("Raw payload: %s", payload)
			} else {
				log("Payload: %s", prettyJSON.String())
			}
		}
		rows.Close()
	}

	log("Database debug query test completed")
}
