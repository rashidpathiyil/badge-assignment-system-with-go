package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

// APIClientInterface abstracts the APIClient for testability
type APIClientInterface interface {
	GetBadges() ([]Badge, error)
	GetBadgeByID(id string) (*Badge, error)
	GetBadgeWithCriteria(id string) (*BadgeWithCriteria, error)
	CreateBadge(badge *NewBadgeRequest) (*BadgeWithCriteria, error)
	DeleteBadge(id string) error
}

// List displays all badges in the system
func List(client APIClientInterface) error {
	badges, err := client.GetBadges()
	if err != nil {
		return fmt.Errorf("failed to get badges: %w", err)
	}

	if len(badges) == 0 {
		fmt.Println("No badges found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tDescription\tActive\tCreated At")
	fmt.Fprintln(w, "--\t----\t-----------\t------\t----------")
	for _, badge := range badges {
		fmt.Fprintf(w, "%d\t%s\t%s\t%t\t%s\n", badge.ID, badge.Name, badge.Description, badge.Active, badge.CreatedAt.Format("2006-01-02"))
	}
	w.Flush()

	return nil
}

// GetBadge displays a badge by ID
func GetBadge(client APIClientInterface, id string) error {
	badge, err := client.GetBadgeByID(id)
	if err != nil {
		return fmt.Errorf("failed to get badge: %w", err)
	}

	badgeWithCriteria, err := client.GetBadgeWithCriteria(id)
	if err != nil {
		fmt.Printf("Warning: Could not fetch badge criteria: %v\n", err)
	}

	fmt.Printf("ID: %d\n", badge.ID)
	fmt.Printf("Name: %s\n", badge.Name)
	fmt.Printf("Description: %s\n", badge.Description)
	fmt.Printf("Image URL: %s\n", badge.ImageURL)
	fmt.Printf("Active: %t\n", badge.Active)
	fmt.Printf("Created At: %s\n", badge.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated At: %s\n", badge.UpdatedAt.Format("2006-01-02 15:04:05"))

	if badgeWithCriteria != nil {
		fmt.Println("\nCriteria:")
		prettyJSON, err := json.MarshalIndent(badgeWithCriteria.Criteria.FlowDefinition, "", "  ")
		if err != nil {
			fmt.Printf("Error formatting criteria: %v\n", err)
		} else {
			fmt.Println(string(prettyJSON))
		}
	}

	return nil
}

// AddBadge creates a new badge with the given name
func AddBadge(client APIClientInterface, name, description, imageURL, flowDefFile string) error {
	if name == "" {
		return fmt.Errorf("badge name is required")
	}

	var flowDef map[string]interface{}

	// Read flow definition from file
	if flowDefFile != "" {
		data, err := os.ReadFile(flowDefFile)
		if err != nil {
			return fmt.Errorf("failed to read flow definition file: %w", err)
		}

		if err := json.Unmarshal(data, &flowDef); err != nil {
			return fmt.Errorf("failed to parse flow definition: %w", err)
		}
	} else {
		// Default simple flow definition
		flowDef = map[string]interface{}{
			"event": "generic-event",
			"criteria": map[string]interface{}{
				"count": map[string]interface{}{
					"$gte": 1,
				},
			},
		}
	}

	// Check for required event types
	requiredEventTypes := ExtractRequiredEventTypes(flowDef)
	if len(requiredEventTypes) > 0 {
		fmt.Printf("Badge requires the following event types: %v\n", requiredEventTypes)

		// Check if the required event types exist
		eventTypeClient, ok := client.(EventTypeClientInterface)
		if !ok {
			fmt.Println("Warning: Cannot check for required event types. Please ensure they are installed.")
		} else {
			eventTypeStatus, err := CheckEventTypesExist(eventTypeClient, requiredEventTypes)
			if err != nil {
				fmt.Printf("Warning: Failed to check for required event types: %v\n", err)
			} else {
				// Collect missing event types
				missingEventTypes := make([]string, 0)
				for et, exists := range eventTypeStatus {
					if !exists {
						missingEventTypes = append(missingEventTypes, et)
					}
				}

				// Prompt to install missing event types
				if len(missingEventTypes) > 0 {
					if err := PromptForEventTypeInstall(eventTypeClient, missingEventTypes); err != nil {
						return fmt.Errorf("failed to install required event types: %w", err)
					}
				}
			}
		}
	}

	badge := &NewBadgeRequest{
		Name:           name,
		Description:    description,
		ImageURL:       imageURL,
		FlowDefinition: flowDef,
	}

	createdBadge, err := client.CreateBadge(badge)
	if err != nil {
		return fmt.Errorf("failed to create badge: %w", err)
	}

	fmt.Printf("Badge created successfully!\n")
	fmt.Printf("ID: %d\n", createdBadge.Badge.ID)
	fmt.Printf("Name: %s\n", createdBadge.Badge.Name)

	return nil
}

// ExtractRequiredEventTypes extracts the event types required by a badge's flow definition
func ExtractRequiredEventTypes(flowDef map[string]interface{}) []string {
	eventTypes := make(map[string]bool)

	// Recursively search for event types in the flow definition
	var extractEvents func(def map[string]interface{})
	extractEvents = func(def map[string]interface{}) {
		// Check for direct event property
		if event, ok := def["event"].(string); ok {
			// Store the event type as kebab-case
			kebabEvent := strings.ToLower(strings.ReplaceAll(event, " ", "-"))
			eventTypes[kebabEvent] = true
		}

		// Check for nested AND/OR conditions
		for _, key := range []string{"$and", "$or"} {
			if conditions, ok := def[key].([]map[string]interface{}); ok {
				for _, condition := range conditions {
					extractEvents(condition)
				}
			}
		}

		// Check for other nested maps
		for k, v := range def {
			if nestedMap, ok := v.(map[string]interface{}); ok && k != "criteria" && k != "payload" {
				extractEvents(nestedMap)
			}
		}
	}

	extractEvents(flowDef)

	// Convert map to slice
	result := make([]string, 0, len(eventTypes))
	for event := range eventTypes {
		result = append(result, event)
	}

	return result
}

// CheckEventTypesExist checks if the required event types exist in the system
func CheckEventTypesExist(client EventTypeClientInterface, eventTypeKeys []string) (map[string]bool, error) {
	// Get all event types from the system
	eventTypes, err := client.GetEventTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to get event types: %w", err)
	}

	// Create maps for kebab-case versions
	existingEventTypes := make(map[string]bool)

	for _, et := range eventTypes {
		// Create the kebab-case version
		kebabKey := strings.ToLower(strings.ReplaceAll(et.Name, " ", "-"))
		existingEventTypes[kebabKey] = true
	}

	// Check if all required event types exist
	result := make(map[string]bool)
	for _, key := range eventTypeKeys {
		// Convert to kebab-case just in case
		kebabKey := strings.ToLower(strings.ReplaceAll(key, " ", "-"))
		result[kebabKey] = existingEventTypes[kebabKey]
	}

	return result, nil
}

// PromptForEventTypeInstall prompts the user to install missing event types
func PromptForEventTypeInstall(client EventTypeClientInterface, missingEventTypes []string) error {
	fmt.Println("The following required event types are missing:")
	for _, key := range missingEventTypes {
		fmt.Printf("- %s\n", key)
	}

	for _, key := range missingEventTypes {
		fmt.Printf("Would you like to install the '%s' event type? (y/n): ", key)
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			return fmt.Errorf("error reading response: %w", err)
		}

		if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
			if err := SavePredefinedEventType(client, key); err != nil {
				return fmt.Errorf("failed to save event type '%s': %w", key, err)
			}
			fmt.Printf("Event type '%s' installed successfully!\n", key)
		} else {
			fmt.Printf("Skipping installation of '%s' event type.\n", key)
		}
	}

	return nil
}

// SavePredefinedBadge adds a predefined badge to the system
func SavePredefinedBadge(client APIClientInterface, key string) error {
	badge, ok := predefinedBadges[key]
	if !ok {
		// List available predefined badges
		fmt.Println("Predefined badge not found. Available badges:")
		for k, v := range predefinedBadges {
			fmt.Printf("- %s: %s\n", k, v.Name)
		}
		return fmt.Errorf("predefined badge '%s' not found", key)
	}

	// Check for required event types
	requiredEventTypes := ExtractRequiredEventTypes(badge.FlowDefinition)
	if len(requiredEventTypes) > 0 {
		fmt.Printf("Badge '%s' requires the following event types: %v\n", badge.Name, requiredEventTypes)

		// Check if the required event types exist
		eventTypeClient, ok := client.(EventTypeClientInterface)
		if !ok {
			fmt.Println("Warning: Cannot check for required event types. Please ensure they are installed.")
		} else {
			eventTypeStatus, err := CheckEventTypesExist(eventTypeClient, requiredEventTypes)
			if err != nil {
				fmt.Printf("Warning: Failed to check for required event types: %v\n", err)
			} else {
				// Collect missing event types
				missingEventTypes := make([]string, 0)
				for et, exists := range eventTypeStatus {
					if !exists {
						missingEventTypes = append(missingEventTypes, et)
					}
				}

				// Prompt to install missing event types
				if len(missingEventTypes) > 0 {
					if err := PromptForEventTypeInstall(eventTypeClient, missingEventTypes); err != nil {
						return fmt.Errorf("failed to install required event types: %w", err)
					}
				}
			}
		}
	}

	createdBadge, err := client.CreateBadge(&badge)
	if err != nil {
		return fmt.Errorf("failed to create badge: %w", err)
	}

	fmt.Printf("Predefined badge '%s' created successfully!\n", badge.Name)
	fmt.Printf("ID: %d\n", createdBadge.Badge.ID)

	return nil
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

	// Check for required event types
	requiredEventTypes := ExtractRequiredEventTypes(badge.FlowDefinition)
	if len(requiredEventTypes) > 0 {
		fmt.Printf("Badge '%s' requires the following event types: %v\n", badge.Name, requiredEventTypes)

		// Check if the required event types exist
		eventTypeClient, ok := client.(EventTypeClientInterface)
		if !ok {
			fmt.Println("Warning: Cannot check for required event types. Please ensure they are installed.")
		} else {
			eventTypeStatus, err := CheckEventTypesExist(eventTypeClient, requiredEventTypes)
			if err != nil {
				fmt.Printf("Warning: Failed to check for required event types: %v\n", err)
			} else {
				// Collect missing event types
				missingEventTypes := make([]string, 0)
				for et, exists := range eventTypeStatus {
					if !exists {
						missingEventTypes = append(missingEventTypes, et)
					}
				}

				// Prompt to install missing event types
				if len(missingEventTypes) > 0 {
					if err := PromptForEventTypeInstall(eventTypeClient, missingEventTypes); err != nil {
						return fmt.Errorf("failed to install required event types: %w", err)
					}
				}
			}
		}
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

// SaveAllPredefinedBadges saves all predefined badges to individual JSON files
func SaveAllPredefinedBadges(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	for key, badge := range predefinedBadges {
		fileName := filepath.Join(directory, key+".json")
		jsonData, err := json.MarshalIndent(badge, "", "  ")
		if err != nil {
			return fmt.Errorf("error serializing badge %s: %w", key, err)
		}

		// Add newline at the end to ensure proper JSON formatting
		jsonData = append(jsonData, '\n')

		if err := os.WriteFile(fileName, jsonData, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fileName, err)
		}

		fmt.Printf("Saved %s to %s\n", badge.Name, fileName)
	}

	return nil
}

// ListPredefinedBadges lists all available predefined badges
func ListPredefinedBadges() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Key\tName\tDescription")
	fmt.Fprintln(w, "---\t----\t-----------")
	for key, badge := range predefinedBadges {
		fmt.Fprintf(w, "%s\t%s\t%s\n", key, badge.Name, badge.Description)
	}
	w.Flush()
}
