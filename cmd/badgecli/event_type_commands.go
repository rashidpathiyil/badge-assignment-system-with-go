package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

// EventTypeClientInterface abstracts the APIClient for event type operations
type EventTypeClientInterface interface {
	GetEventTypes() ([]EventType, error)
	GetEventTypeByID(id string) (*EventType, error)
	CreateEventType(eventType *NewEventTypeRequest) (*EventType, error)
	UpdateEventType(id string, eventType *UpdateEventTypeRequest) (*EventType, error)
	DeleteEventType(id string) error
}

// ListEventTypes displays all event types in the system
func ListEventTypes(client EventTypeClientInterface) error {
	eventTypes, err := client.GetEventTypes()
	if err != nil {
		return fmt.Errorf("failed to get event types: %w", err)
	}

	if len(eventTypes) == 0 {
		fmt.Println("No event types found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tDescription\tCreated At")
	fmt.Fprintln(w, "--\t----\t-----------\t----------")
	for _, et := range eventTypes {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", et.ID, et.Name, et.Description, et.CreatedAt.Format("2006-01-02"))
	}
	w.Flush()

	return nil
}

// GetEventType displays an event type by ID
func GetEventType(client EventTypeClientInterface, id string) error {
	eventType, err := client.GetEventTypeByID(id)
	if err != nil {
		return fmt.Errorf("failed to get event type: %w", err)
	}

	fmt.Printf("ID: %d\n", eventType.ID)
	fmt.Printf("Name: %s\n", eventType.Name)
	fmt.Printf("Description: %s\n", eventType.Description)
	fmt.Printf("Created At: %s\n", eventType.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated At: %s\n", eventType.UpdatedAt.Format("2006-01-02 15:04:05"))

	fmt.Println("\nSchema:")
	prettyJSON, err := json.MarshalIndent(eventType.Schema, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting schema: %v\n", err)
	} else {
		fmt.Println(string(prettyJSON))
	}

	return nil
}

// AddEventType creates a new event type
func AddEventType(client EventTypeClientInterface, name, description string, schemaFile string) error {
	if name == "" {
		return fmt.Errorf("event type name is required")
	}

	var schema map[string]interface{}

	// Read schema from file
	if schemaFile != "" {
		data, err := os.ReadFile(schemaFile)
		if err != nil {
			return fmt.Errorf("failed to read schema file: %w", err)
		}

		if err := json.Unmarshal(data, &schema); err != nil {
			return fmt.Errorf("failed to parse schema: %w", err)
		}
	} else {
		// Default simple schema
		schema = map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"timestamp": map[string]interface{}{
					"type":   "string",
					"format": "date-time",
				},
			},
			"required": []string{"timestamp"},
		}
	}

	eventType := &NewEventTypeRequest{
		Name:        name,
		Description: description,
		Schema:      schema,
	}

	createdEventType, err := client.CreateEventType(eventType)
	if err != nil {
		return fmt.Errorf("failed to create event type: %w", err)
	}

	fmt.Printf("Event type created successfully!\n")
	fmt.Printf("ID: %d\n", createdEventType.ID)
	fmt.Printf("Name: %s\n", createdEventType.Name)

	return nil
}

// UpdateEventType updates an existing event type
func UpdateEventType(client EventTypeClientInterface, id, name, description string, schemaFile string) error {
	updateReq := &UpdateEventTypeRequest{}

	if name != "" {
		updateReq.Name = name
	}

	if description != "" {
		updateReq.Description = description
	}

	// Read schema from file if provided
	if schemaFile != "" {
		data, err := os.ReadFile(schemaFile)
		if err != nil {
			return fmt.Errorf("failed to read schema file: %w", err)
		}

		var schema map[string]interface{}
		if err := json.Unmarshal(data, &schema); err != nil {
			return fmt.Errorf("failed to parse schema: %w", err)
		}

		updateReq.Schema = schema
	}

	// Don't update if no fields provided
	if name == "" && description == "" && schemaFile == "" {
		return fmt.Errorf("at least one field to update must be provided")
	}

	updatedEventType, err := client.UpdateEventType(id, updateReq)
	if err != nil {
		return fmt.Errorf("failed to update event type: %w", err)
	}

	fmt.Printf("Event type updated successfully!\n")
	fmt.Printf("ID: %d\n", updatedEventType.ID)
	fmt.Printf("Name: %s\n", updatedEventType.Name)

	return nil
}

// DeleteEventType deletes an event type
func DeleteEventType(client EventTypeClientInterface, id string) error {
	if err := client.DeleteEventType(id); err != nil {
		return fmt.Errorf("failed to delete event type: %w", err)
	}

	fmt.Printf("Event type with ID %s deleted successfully!\n", id)
	return nil
}

// ImportEventType imports an event type from a JSON file
func ImportEventType(client EventTypeClientInterface, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var eventType NewEventTypeRequest
	if err := json.Unmarshal(data, &eventType); err != nil {
		return fmt.Errorf("failed to parse event type JSON: %w", err)
	}

	createdEventType, err := client.CreateEventType(&eventType)
	if err != nil {
		return fmt.Errorf("failed to create event type: %w", err)
	}

	fmt.Printf("Event type '%s' imported successfully!\n", eventType.Name)
	fmt.Printf("ID: %d\n", createdEventType.ID)

	return nil
}

// ExportEventType exports an event type to a JSON file
func ExportEventType(client EventTypeClientInterface, id, outputPath string) error {
	eventType, err := client.GetEventTypeByID(id)
	if err != nil {
		return fmt.Errorf("failed to get event type: %w", err)
	}

	// Convert to NewEventTypeRequest format for clean export
	exportEventType := NewEventTypeRequest{
		Name:        eventType.Name,
		Description: eventType.Description,
		Schema:      eventType.Schema,
	}

	jsonData, err := json.MarshalIndent(exportEventType, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing event type: %w", err)
	}

	// Add newline at the end to ensure proper JSON formatting
	jsonData = append(jsonData, '\n')

	if outputPath == "" {
		outputPath = fmt.Sprintf("event_type_%s.json", id)
	}

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Event type exported to %s\n", outputPath)
	return nil
}
