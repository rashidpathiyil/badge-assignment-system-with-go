package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
)

// APIClientInterface abstracts the APIClient for testability
type APIClientInterface interface {
	// Badge operations
	GetBadges() ([]Badge, error)
	GetBadgeWithCriteria(id string) (*BadgeWithCriteria, error)
	CreateBadge(badge *NewBadgeRequest) (*BadgeWithCriteria, error)

	// Event type operations
	GetEventTypes() ([]EventType, error)
	GetEventTypeByID(id string) (*EventType, error)
	CreateEventType(eventType *NewEventTypeRequest) (*EventType, error)
}

// ImportBadge imports a badge from a JSON file
func ImportBadge(client APIClientInterface, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var badge NewBadgeRequest
	if err := json.Unmarshal(data, &badge); err != nil {
		return fmt.Errorf("failed to parse badge JSON: %w", err)
	}

	createdBadge, err := client.CreateBadge(&badge)
	if err != nil {
		return fmt.Errorf("failed to create badge: %w", err)
	}

	fmt.Printf("Badge '%s' imported successfully!\n", badge.Name)
	fmt.Printf("ID: %d\n", createdBadge.Badge.ID)

	return nil
}

// ExportBadge exports a badge to a JSON file
func ExportBadge(client APIClientInterface, id, outputPath string) error {
	badgeWithCriteria, err := client.GetBadgeWithCriteria(id)
	if err != nil {
		return fmt.Errorf("failed to get badge with criteria: %w", err)
	}

	// Convert to NewBadgeRequest format
	exportBadge := NewBadgeRequest{
		Name:           badgeWithCriteria.Badge.Name,
		Description:    badgeWithCriteria.Badge.Description,
		ImageURL:       badgeWithCriteria.Badge.ImageURL,
		FlowDefinition: badgeWithCriteria.Criteria.FlowDefinition,
	}

	jsonData, err := json.MarshalIndent(exportBadge, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing badge: %w", err)
	}

	if outputPath == "" {
		outputPath = fmt.Sprintf("badge_%s.json", id)
	}

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Badge exported to %s\n", outputPath)
	return nil
}

// ExportAllBadges exports all badges from the system to a directory
func ExportAllBadges(client APIClientInterface, outputDir string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	// Get all badges
	badges, err := client.GetBadges()
	if err != nil {
		return fmt.Errorf("failed to get badges: %w", err)
	}

	if len(badges) == 0 {
		fmt.Println("No badges found in the system.")
		return nil
	}

	// Export each badge
	for _, badge := range badges {
		// Get badge with criteria
		badgeWithCriteria, err := client.GetBadgeWithCriteria(strconv.Itoa(badge.ID))
		if err != nil {
			fmt.Printf("Warning: couldn't get criteria for badge ID %d: %v\n", badge.ID, err)
			continue
		}

		// Convert to NewBadgeRequest format
		exportBadge := NewBadgeRequest{
			Name:           badge.Name,
			Description:    badge.Description,
			ImageURL:       badge.ImageURL,
			FlowDefinition: badgeWithCriteria.Criteria.FlowDefinition,
		}

		// Create filename
		filename := fmt.Sprintf("%s_%d.json", sanitizeFilename(badge.Name), badge.ID)
		outputPath := filepath.Join(outputDir, filename)

		// Write to file
		jsonData, err := json.MarshalIndent(exportBadge, "", "  ")
		if err != nil {
			fmt.Printf("Warning: couldn't serialize badge %d: %v\n", badge.ID, err)
			continue
		}

		if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
			fmt.Printf("Warning: couldn't write badge %d to file: %v\n", badge.ID, err)
			continue
		}

		fmt.Printf("Exported badge ID %d (%s) to %s\n", badge.ID, badge.Name, outputPath)
	}

	fmt.Printf("Exported %d badges to %s\n", len(badges), outputDir)
	return nil
}

// ListBadges lists all badges in the system
func ListBadges(client APIClientInterface) error {
	badges, err := client.GetBadges()
	if err != nil {
		return fmt.Errorf("failed to get badges: %w", err)
	}

	if len(badges) == 0 {
		fmt.Println("No badges found in the system.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tDescription\tActive")
	fmt.Fprintln(w, "--\t----\t-----------\t------")

	for _, badge := range badges {
		fmt.Fprintf(w, "%d\t%s\t%s\t%t\n",
			badge.ID,
			badge.Name,
			badge.Description,
			badge.Active)
	}
	w.Flush()

	fmt.Printf("\nTotal badges: %d\n", len(badges))
	return nil
}

// ImportEventType imports an event type from a JSON file
func ImportEventType(client APIClientInterface, filePath string) error {
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
func ExportEventType(client APIClientInterface, id, outputPath string) error {
	eventType, err := client.GetEventTypeByID(id)
	if err != nil {
		return fmt.Errorf("failed to get event type: %w", err)
	}

	// Convert to NewEventTypeRequest format
	exportEventType := NewEventTypeRequest{
		Name:        eventType.Name,
		Description: eventType.Description,
		Schema:      eventType.Schema,
	}

	jsonData, err := json.MarshalIndent(exportEventType, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing event type: %w", err)
	}

	if outputPath == "" {
		outputPath = fmt.Sprintf("event_type_%s.json", id)
	}

	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Event type exported to %s\n", outputPath)
	return nil
}

// ExportAllEventTypes exports all event types from the system to a directory
func ExportAllEventTypes(client APIClientInterface, outputDir string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	// Get all event types
	eventTypes, err := client.GetEventTypes()
	if err != nil {
		return fmt.Errorf("failed to get event types: %w", err)
	}

	if len(eventTypes) == 0 {
		fmt.Println("No event types found in the system.")
		return nil
	}

	// Export each event type
	for _, eventType := range eventTypes {
		// Convert to NewEventTypeRequest format
		exportEventType := NewEventTypeRequest{
			Name:        eventType.Name,
			Description: eventType.Description,
			Schema:      eventType.Schema,
		}

		// Create filename
		filename := fmt.Sprintf("%s_%d.json", sanitizeFilename(eventType.Name), eventType.ID)
		outputPath := filepath.Join(outputDir, filename)

		// Write to file
		jsonData, err := json.MarshalIndent(exportEventType, "", "  ")
		if err != nil {
			fmt.Printf("Warning: couldn't serialize event type %d: %v\n", eventType.ID, err)
			continue
		}

		if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
			fmt.Printf("Warning: couldn't write event type %d to file: %v\n", eventType.ID, err)
			continue
		}

		fmt.Printf("Exported event type ID %d (%s) to %s\n", eventType.ID, eventType.Name, outputPath)
	}

	fmt.Printf("Exported %d event types to %s\n", len(eventTypes), outputDir)
	return nil
}

// ListEventTypes lists all event types in the system
func ListEventTypes(client APIClientInterface) error {
	eventTypes, err := client.GetEventTypes()
	if err != nil {
		return fmt.Errorf("failed to get event types: %w", err)
	}

	if len(eventTypes) == 0 {
		fmt.Println("No event types found in the system.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tDescription")
	fmt.Fprintln(w, "--\t----\t-----------")

	for _, eventType := range eventTypes {
		fmt.Fprintf(w, "%d\t%s\t%s\n",
			eventType.ID,
			eventType.Name,
			eventType.Description)
	}
	w.Flush()

	fmt.Printf("\nTotal event types: %d\n", len(eventTypes))
	return nil
}

// ListExampleBadgesAndEventTypes lists all example badges and event types included with the CLI
func ListExampleBadgesAndEventTypes() error {
	// List badge examples
	fmt.Println("📛 Example Badges:")
	fmt.Println("=================")

	badgeEntries, err := fs.Glob(exampleFS, "badges/*.json")
	if err != nil {
		return fmt.Errorf("failed to read badge examples: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Name\tDescription\tFile")
	fmt.Fprintln(w, "----\t-----------\t----")

	for _, entry := range badgeEntries {
		if strings.HasSuffix(entry, ".json") && !strings.Contains(entry, "/event_types/") {
			data, err := exampleFS.ReadFile(entry)
			if err != nil {
				fmt.Printf("Warning: couldn't read %s: %v\n", entry, err)
				continue
			}

			var badge struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}

			if err := json.Unmarshal(data, &badge); err != nil {
				fmt.Printf("Warning: couldn't parse %s: %v\n", entry, err)
				continue
			}

			fmt.Fprintf(w, "%s\t%s\t%s\n", badge.Name, badge.Description, entry)
		}
	}
	w.Flush()

	// List event type examples
	fmt.Println("\n⚡ Example Event Types:")
	fmt.Println("=====================")

	eventTypeEntries, err := fs.Glob(exampleFS, "badges/event_types/*.json")
	if err != nil {
		return fmt.Errorf("failed to read event type examples: %w", err)
	}

	w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Name\tDescription\tFile")
	fmt.Fprintln(w, "----\t-----------\t----")

	for _, entry := range eventTypeEntries {
		if strings.HasSuffix(entry, ".json") {
			data, err := exampleFS.ReadFile(entry)
			if err != nil {
				fmt.Printf("Warning: couldn't read %s: %v\n", entry, err)
				continue
			}

			var eventType struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			}

			if err := json.Unmarshal(data, &eventType); err != nil {
				fmt.Printf("Warning: couldn't parse %s: %v\n", entry, err)
				continue
			}

			fmt.Fprintf(w, "%s\t%s\t%s\n", eventType.Name, eventType.Description, entry)
		}
	}
	w.Flush()

	return nil
}

// ExportExampleBadgesAndEventTypes exports all example badges and event types to a directory
func ExportExampleBadgesAndEventTypes(outputDir string) error {
	// Create output directories
	badgesOutputDir := filepath.Join(outputDir, "badges")
	eventTypesOutputDir := filepath.Join(outputDir, "event_types")

	dirs := []string{outputDir, badgesOutputDir, eventTypesOutputDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Export badges
	badgeEntries, err := fs.Glob(exampleFS, "badges/*.json")
	if err != nil {
		return fmt.Errorf("failed to read badge examples: %w", err)
	}

	for _, entry := range badgeEntries {
		if strings.HasSuffix(entry, ".json") && !strings.Contains(entry, "/event_types/") {
			data, err := exampleFS.ReadFile(entry)
			if err != nil {
				fmt.Printf("Warning: couldn't read %s: %v\n", entry, err)
				continue
			}

			outputFile := filepath.Join(badgesOutputDir, filepath.Base(entry))
			if err := os.WriteFile(outputFile, data, 0644); err != nil {
				fmt.Printf("Warning: couldn't write %s: %v\n", outputFile, err)
				continue
			}

			fmt.Printf("Exported badge to %s\n", outputFile)
		}
	}

	// Export event types
	eventTypeEntries, err := fs.Glob(exampleFS, "badges/event_types/*.json")
	if err != nil {
		return fmt.Errorf("failed to read event type examples: %w", err)
	}

	for _, entry := range eventTypeEntries {
		if strings.HasSuffix(entry, ".json") {
			data, err := exampleFS.ReadFile(entry)
			if err != nil {
				fmt.Printf("Warning: couldn't read %s: %v\n", entry, err)
				continue
			}

			outputFile := filepath.Join(eventTypesOutputDir, filepath.Base(entry))
			if err := os.WriteFile(outputFile, data, 0644); err != nil {
				fmt.Printf("Warning: couldn't write %s: %v\n", outputFile, err)
				continue
			}

			fmt.Printf("Exported event type to %s\n", outputFile)
		}
	}

	return nil
}

// Helper function to sanitize filenames
func sanitizeFilename(name string) string {
	// Replace spaces and special characters with underscores
	sanitized := strings.ToLower(name)
	sanitized = strings.ReplaceAll(sanitized, " ", "_")
	sanitized = strings.ReplaceAll(sanitized, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, "\\", "_")
	sanitized = strings.ReplaceAll(sanitized, ":", "_")
	sanitized = strings.ReplaceAll(sanitized, "*", "_")
	sanitized = strings.ReplaceAll(sanitized, "?", "_")
	sanitized = strings.ReplaceAll(sanitized, "\"", "_")
	sanitized = strings.ReplaceAll(sanitized, "<", "_")
	sanitized = strings.ReplaceAll(sanitized, ">", "_")
	sanitized = strings.ReplaceAll(sanitized, "|", "_")

	return sanitized
}

// ImportAllBadges imports all badge definitions from a directory
func ImportAllBadges(client APIClientInterface, dirPath string) error {
	// Check if directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("failed to access directory %s: %w", dirPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dirPath)
	}

	// Get list of JSON files in the directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var jsonFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			jsonFiles = append(jsonFiles, filepath.Join(dirPath, entry.Name()))
		}
	}

	if len(jsonFiles) == 0 {
		return fmt.Errorf("no JSON files found in directory %s", dirPath)
	}

	// Import each file
	importCount := 0
	failCount := 0

	for _, filePath := range jsonFiles {
		// Read the file
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: couldn't read %s: %v\n", filePath, err)
			failCount++
			continue
		}

		// Try to parse as badge
		var badge NewBadgeRequest
		if err := json.Unmarshal(data, &badge); err != nil {
			fmt.Printf("Warning: couldn't parse %s as a badge: %v\n", filePath, err)
			failCount++
			continue
		}

		// Skip if missing required fields
		if badge.Name == "" || badge.Description == "" || badge.FlowDefinition == nil {
			fmt.Printf("Warning: skipping %s - doesn't appear to be a valid badge definition\n", filePath)
			failCount++
			continue
		}

		// Import the badge
		createdBadge, err := client.CreateBadge(&badge)
		if err != nil {
			fmt.Printf("Error importing badge %s: %v\n", filePath, err)
			failCount++
			continue
		}

		fmt.Printf("Imported badge '%s' (ID: %d) from %s\n", badge.Name, createdBadge.Badge.ID, filePath)
		importCount++
	}

	fmt.Printf("\nImport summary: %d badges imported, %d failed\n", importCount, failCount)
	return nil
}

// ImportAllEventTypes imports all event type definitions from a directory
func ImportAllEventTypes(client APIClientInterface, dirPath string) error {
	// Check if directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("failed to access directory %s: %w", dirPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dirPath)
	}

	// Get list of JSON files in the directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var jsonFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			jsonFiles = append(jsonFiles, filepath.Join(dirPath, entry.Name()))
		}
	}

	if len(jsonFiles) == 0 {
		return fmt.Errorf("no JSON files found in directory %s", dirPath)
	}

	// Import each file
	importCount := 0
	failCount := 0

	for _, filePath := range jsonFiles {
		// Read the file
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Warning: couldn't read %s: %v\n", filePath, err)
			failCount++
			continue
		}

		// Try to parse as event type
		var eventType NewEventTypeRequest
		if err := json.Unmarshal(data, &eventType); err != nil {
			fmt.Printf("Warning: couldn't parse %s as an event type: %v\n", filePath, err)
			failCount++
			continue
		}

		// Skip if missing required fields
		if eventType.Name == "" || eventType.Description == "" || eventType.Schema == nil {
			fmt.Printf("Warning: skipping %s - doesn't appear to be a valid event type definition\n", filePath)
			failCount++
			continue
		}

		// Import the event type
		createdEventType, err := client.CreateEventType(&eventType)
		if err != nil {
			fmt.Printf("Error importing event type %s: %v\n", filePath, err)
			failCount++
			continue
		}

		fmt.Printf("Imported event type '%s' (ID: %d) from %s\n", eventType.Name, createdEventType.ID, filePath)
		importCount++
	}

	fmt.Printf("\nImport summary: %d event types imported, %d failed\n", importCount, failCount)
	return nil
}
