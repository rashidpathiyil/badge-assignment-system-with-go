# Badge CLI

Badge CLI is a command line tool for importing and exporting badge and event type definitions to and from the Badge Assignment System.

## Overview

The Badge CLI provides the following functionality:

- Import badge and event type definitions from JSON files (individually or in bulk)
- Export badge and event type definitions to JSON files (individually or in bulk)
- View and export example badge and event type definitions
- List all badges and event types in the system

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
   go build -o badgecli cmd/badgecli/*.go
   ```

## Usage

### Example Commands

#### List examples included with the CLI

```bash
./badgecli list-examples
```

#### Export examples to a directory

```bash
./badgecli export-examples --output-dir ./my-examples
```

#### Import a badge definition

```bash
./badgecli import /path/to/my-badge.json
```

#### Import all badge definitions from a directory

```bash
./badgecli import-all-badges ./my-badges
```

#### Export a badge definition by ID

```bash
./badgecli export 123 /path/to/exported-badge.json
```

#### Import an event type definition

```bash
./badgecli import-event-type /path/to/my-event-type.json
```

#### Import all event type definitions from a directory

```bash
./badgecli import-all-event-types ./my-event-types
```

#### Export an event type definition by ID

```bash
./badgecli export-event-type 456 /path/to/exported-event-type.json
```

#### List all badges in the system

```bash
./badgecli list-badges
```

#### List all event types in the system

```bash
./badgecli list-event-types
```

#### Export all badges from the system to a directory

```bash
./badgecli export-all-badges ./my-badges
```

#### Export all event types from the system to a directory

```bash
./badgecli export-all-event-types ./my-event-types
```

## Badge and Event Type Definitions

The CLI includes example badge and event type definitions in the `badges` directory. These examples showcase recommended patterns and best practices for defining badges and event types.

### Badge Definition Examples

Badges are defined with the following format:

```json
{
  "name": "Early Bird",
  "description": "Checked in before 9 AM for 5 consecutive days",
  "image_url": "https://example.com/badges/early-bird.png",
  "flow_definition": {
    "operator": "$count",
    "event_type": "check_in",
    "conditions": {
      "timestamp": {
        "operator": "$lt",
        "value": "09:00:00"
      }
    },
    "min_count": 5,
    "time_window": "5d",
    "consecutive": true
  }
}
```

Alternative format with $NOW dynamic variable support:

```json
{
  "name": "Early Bird",
  "description": "Checked in before 9 AM for 5 consecutive days.",
  "image_url": "https://example.com/badges/early-bird.png",
  "flow_definition": {
    "criteria": {
      "$eventCount": {
        "$gte": 5
      },
      "timestamp": {
        "$gte": "$NOW(-5d)",
        "$lte": "$NOW()"
      },
      "payload": {
        "time": {
          "$lt": "09:00:00"
        }
      }
    },
    "event": "check-in"
  }
}
```

### Event Type Definition Example

Event types are defined with the following format:

```json
{
  "name": "check_in",
  "description": "User check-in event, used for attendance tracking",
  "schema": {
    "type": "object",
    "properties": {
      "user_id": {
        "type": "string"
      },
      "timestamp": {
        "type": "string",
        "format": "date-time"
      },
      "location": {
        "type": "string"
      },
      "time": {
        "description": "The time of check-in in HH:MM:SS format",
        "format": "time",
        "type": "string"
      },
      "date": {
        "format": "date",
        "type": "string"
      }
    },
    "required": ["user_id", "timestamp"]
  }
}
```

## Configuration

The CLI can be configured using:

1. Command line flags:
   ```bash
   ./badgecli --server http://api.example.com
   ```

2. Environment variables:
   ```bash
   export BADGE_SERVER_URL=http://api.example.com
   ```

3. Configuration file (default: `$HOME/.badgecli.yaml`):
   ```yaml
   server: http://api.example.com
   ```

To specify a custom configuration file:
```bash
./badgecli --config /path/to/config.yaml
``` 
