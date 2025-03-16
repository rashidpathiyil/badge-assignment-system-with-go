 package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create DB connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	// Check connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	// SQL to clear all tables
	clearSQL := `
	DO $$
	DECLARE
		row record;
	BEGIN
		-- Disable triggers temporarily
		SET session_replication_role = 'replica';

		-- Loop through all tables in the current schema, excluding PostgreSQL system tables
		FOR row IN 
			SELECT tablename FROM pg_tables WHERE schemaname = 'public'
		LOOP
			EXECUTE 'TRUNCATE TABLE "' || row.tablename || '" CASCADE';
		END LOOP;

		-- Re-enable triggers
		SET session_replication_role = 'origin';
	END $$;
	`

	// Execute the SQL
	fmt.Println("Clearing all tables in the database...")
	_, err = db.Exec(clearSQL)
	if err != nil {
		log.Fatalf("Error clearing tables: %v", err)
	}

	fmt.Println("All tables have been successfully cleared while preserving the schema.")
}
