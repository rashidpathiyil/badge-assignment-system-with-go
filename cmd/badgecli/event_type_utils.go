package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
)

// ListPredefinedEventTypes lists all available predefined event types
func ListPredefinedEventTypes() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Key\tName\tDescription")
	fmt.Fprintln(w, "---\t----\t-----------")
	for key, et := range predefinedEventTypes {
		fmt.Fprintf(w, "%s\t%s\t%s\n", key, et.Name, et.Description)
	}
	w.Flush()
}

// SavePredefinedEventType adds a predefined event type to the system
func SavePredefinedEventType(client EventTypeClientInterface, key string) error {
	eventType, ok := predefinedEventTypes[key]
	if !ok {
		// List available predefined event types
		fmt.Println("Predefined event type not found. Available event types:")
		for k, v := range predefinedEventTypes {
			fmt.Printf("- %s: %s\n", k, v.Name)
		}
		return fmt.Errorf("predefined event type '%s' not found", key)
	}

	createdEventType, err := client.CreateEventType(&eventType)
	if err != nil {
		return fmt.Errorf("failed to create event type: %w", err)
	}

	fmt.Printf("Predefined event type '%s' created successfully!\n", eventType.Name)
	fmt.Printf("ID: %d\n", createdEventType.ID)

	return nil
}

// SaveAllPredefinedEventTypes saves all predefined event types to individual JSON files
func SaveAllPredefinedEventTypes(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	for key, eventType := range predefinedEventTypes {
		fileName := filepath.Join(directory, key+".json")
		jsonData, err := json.MarshalIndent(eventType, "", "  ")
		if err != nil {
			return fmt.Errorf("error serializing event type %s: %w", key, err)
		}

		// Add newline at the end to ensure proper JSON formatting
		jsonData = append(jsonData, '\n')

		if err := os.WriteFile(fileName, jsonData, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fileName, err)
		}

		fmt.Printf("Saved %s to %s\n", eventType.Name, fileName)
	}

	return nil
}
