# Database Utilities

This folder contains utilities for database management and operations for the Badge Assignment System.

## Clear Database Tool

The utilities in this folder provide an easy way to clear all data from the database tables while preserving the schema (tables, functions, triggers, etc.). This is particularly useful during development and testing when you need to reset the database to a clean state.

### Files

- `clear_db.go` - Go script to connect to the database and execute the cleanup command
- `clear_db.sql` - SQL script with the commands to truncate all tables

### Usage

To clear all data from the database:

```bash
# From the project root directory
cd tools/db-utils
go run clear_db.go
```

Or directly from the project root:

```bash
go run tools/db-utils/clear_db.go
```

### How It Works

The tool:

1. Connects to the database using the credentials in your `.env` file
2. Temporarily disables triggers to avoid constraint violations
3. Uses PostgreSQL's dynamic SQL capabilities to loop through all tables
4. Executes a `TRUNCATE TABLE CASCADE` command on each table
5. Re-enables triggers afterward

### Prerequisites

- Access to the PostgreSQL database specified in your `.env` file
- The `.env` file must be properly configured with database connection details

### Caution

**WARNING**: This tool deletes ALL DATA from ALL TABLES in the database. It should only be used in development and testing environments, never in production. 
