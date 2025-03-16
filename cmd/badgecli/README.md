# Badge CLI

A command-line interface for managing badges in the Badge Assignment System.

## Installation

Before using the Badge CLI, ensure you have Go installed on your system.

1. Install the required dependencies:
   ```bash
   go get github.com/spf13/cobra
   go get github.com/spf13/viper
   go get github.com/joho/godotenv
   ```

2. Build the CLI tool:
   ```bash
   go build -o badgecli
   ```

## Usage

The Badge CLI provides various commands to manage badges and event types in the system. By default, the CLI connects to a server running at http://localhost:8080, but you can specify a different server URL using the `--server` flag.

The CLI interacts with both public and admin API endpoints:
- Public endpoints: Used for listing badges and viewing badge details
- Admin endpoints: Used for creating, updating, and deleting badges and event types

### Global Flags

- `--server`: Server URL (default "http://localhost:8080")
- `--badges-dir`: Directory for badge and event type definitions (default "./badges")
- `--config`: Config file (default is $HOME/.badgecli.yaml)

### Badge Commands

#### List Badges

List all badges in the system:

```bash
./badgecli list
```

#### Get Badge

Get details of a badge by ID:

```bash
./badgecli get 1
```

#### Add Badge

Add a new badge to the system:

```bash
./badgecli add "New Badge" --description "Description of the badge" --image "https://example.com/badge.png" --flow flow_definition.json
```

When adding a badge, the CLI automatically checks if the required event types exist in the system. If any required event types are missing, the CLI will prompt you to install them before creating the badge.

#### Save Predefined Badge

Save a predefined badge to the system:

```bash
./badgecli save-predefined early-bird
```

To see a list of available predefined badges:

```bash
./badgecli save-predefined
```

Just like with the `add` command, the CLI will check for required event types and prompt you to install any that are missing.

#### List Predefined Badges

List all available predefined badges:

```bash
./badgecli list-predefined
```

This will display a table showing the key, name, and description of all predefined badges.

#### Save All Predefined Badges

Save all predefined badges to JSON files:

```bash
./badgecli save-all-predefined
```

This will create JSON files for each predefined badge in the badges directory.

#### Import Badge

Import a badge from a JSON file:

```bash
./badgecli import badge_file.json
```

When importing a badge, the CLI will check for required event types and prompt you to install any that are missing, ensuring all dependencies are properly set up.

#### Export Badge

Export a badge to a JSON file:

```bash
./badgecli export 1 exported_badge.json
```

### Event Type Commands

#### List Event Types

List all event types in the system:

```bash
./badgecli list-event-types
```

#### Get Event Type

Get details of an event type by ID:

```bash
./badgecli get-event-type 1
```

#### Add Event Type

Add a new event type to the system:

```bash
./badgecli add-event-type "New Event Type" --description "Description of the event type" --schema schema.json
```

#### Update Event Type

Update an existing event type:

```bash
./badgecli update-event-type 1 --name "Updated Name" --description "Updated description" --schema updated_schema.json
```

#### Delete Event Type

Delete an event type:

```bash
./badgecli delete-event-type 1
```

#### Import Event Type

Import an event type from a JSON file:

```bash
./badgecli import-event-type event_type_file.json
```

#### Export Event Type

Export an event type to a JSON file:

```bash
./badgecli export-event-type 1 exported_event_type.json
```

#### List Predefined Event Types

List all predefined event types:

```bash
./badgecli list-predefined-event-types
```

#### Save Predefined Event Type

Save a predefined event type to the system:

```bash
./badgecli save-predefined-event-type check-in
```

#### Save All Predefined Event Types

Save all predefined event types to JSON files:

```bash
./badgecli save-all-event-types
```

This will create JSON files for each predefined event type in the badges/event_types directory.

## Predefined Event Types

The CLI includes the following predefined event types that are required for the badge system:

1. **check-in**: User check-in event, used for attendance tracking
2. **check-out**: User check-out event, used for attendance tracking
3. **work-log**: User work log entry for tracking hours worked
4. **meeting-attendance**: User attendance at a meeting
5. **bug-report**: User submitted a bug report
6. **task-completion**: User completed a task
7. **project-collaboration**: User collaborated on a team project

## Predefined Badges

The CLI includes the following predefined badges that can be added to the system:

1. **early-bird**: Checked in before 9 AM for 5 consecutive days.
2. **workaholic**: Logged 40+ hours in a week.
3. **consistency-king**: Checked in & out without missing a day for a month.
4. **team-player**: Collaborated on 5+ team projects in a month.
5. **overtime-warrior**: Logged extra 10 hours in a week.
6. **meeting-maestro**: Attended 10+ meetings in a month.
7. **bug-hunter**: Reported 5+ issues that got fixed.
8. **task-master**: Completed 50+ tasks in a month.

## Environment Variables

The CLI can be configured using environment variables:

- `BADGE_SERVER_URL`: The URL of the badge server
- `BADGE_DEFINITIONS_DIR`: Directory to store badge definitions

You can set these variables in a `.env` file in the current directory.

## Example Workflow

1. Set up event types for the badge system:
   ```bash
   ./badgecli save-all-event-types
   ```

2. Add badges to the system (the CLI will automatically prompt for any missing event types):
   ```bash
   ./badgecli save-predefined early-bird
   ./badgecli save-predefined workaholic
   ```

   The CLI will detect that these badges require specific event types and prompt you to install them if they don't exist.

3. List all badges in the system:
   ```bash
   ./badgecli list
   ``` 
