-- Script to delete all data from all tables in the badge_system database
-- This script preserves the schema but removes all rows

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
