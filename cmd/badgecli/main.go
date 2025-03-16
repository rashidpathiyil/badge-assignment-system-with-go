package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	serverURL   string
	badgesDir   string
	description string
	imageURL    string
	flowDefFile string
	schemaFile  string
)

// Command variables - Badge Commands
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all badges in the system",
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := List(client); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var getCmd = &cobra.Command{
	Use:   "get [badge_id]",
	Short: "Get details of a badge by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := GetBadge(client, args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add [badge_name]",
	Short: "Add a new badge to the system",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := AddBadge(client, args[0], description, imageURL, flowDefFile); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var savePredefCmd = &cobra.Command{
	Use:   "save-predefined [badge_key]",
	Short: "Save a predefined badge to the system",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// List predefined badges instead
			ListPredefinedBadges()
			return
		}

		client := NewAPIClient(serverURL)
		if err := SavePredefinedBadge(client, args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var importCmd = &cobra.Command{
	Use:   "import [file_path]",
	Short: "Import a badge definition from a JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := ImportBadge(client, args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var exportCmd = &cobra.Command{
	Use:   "export [badge_id] [file_path]",
	Short: "Export a badge definition to a JSON file",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		var outputPath string
		if len(args) > 1 {
			outputPath = args[1]
		}
		if err := ExportBadge(client, args[0], outputPath); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var saveAllPredefCmd = &cobra.Command{
	Use:   "save-all-predefined",
	Short: "Save all predefined badges to JSON files in the badges directory",
	Run: func(cmd *cobra.Command, args []string) {
		if err := SaveAllPredefinedBadges(badgesDir); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("All predefined badges saved to %s\n", badgesDir)
	},
}

var listPredefinedCmd = &cobra.Command{
	Use:   "list-predefined",
	Short: "List all predefined badges",
	Run: func(cmd *cobra.Command, args []string) {
		ListPredefinedBadges()
	},
}

// Command variables - Event Type Commands
var listEventTypesCmd = &cobra.Command{
	Use:   "list-event-types",
	Short: "List all event types in the system",
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := ListEventTypes(client); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var getEventTypeCmd = &cobra.Command{
	Use:   "get-event-type [event_type_id]",
	Short: "Get details of an event type by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := GetEventType(client, args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var addEventTypeCmd = &cobra.Command{
	Use:   "add-event-type [event_type_name]",
	Short: "Add a new event type to the system",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := AddEventType(client, args[0], description, schemaFile); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var updateEventTypeCmd = &cobra.Command{
	Use:   "update-event-type [event_type_id]",
	Short: "Update an existing event type",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := UpdateEventType(client, args[0], cmd.Flag("name").Value.String(), description, schemaFile); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var deleteEventTypeCmd = &cobra.Command{
	Use:   "delete-event-type [event_type_id]",
	Short: "Delete an event type",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := DeleteEventType(client, args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var importEventTypeCmd = &cobra.Command{
	Use:   "import-event-type [file_path]",
	Short: "Import an event type definition from a JSON file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := ImportEventType(client, args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var exportEventTypeCmd = &cobra.Command{
	Use:   "export-event-type [event_type_id] [file_path]",
	Short: "Export an event type definition to a JSON file",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		var outputPath string
		if len(args) > 1 {
			outputPath = args[1]
		}
		if err := ExportEventType(client, args[0], outputPath); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var listPredefinedEventTypesCmd = &cobra.Command{
	Use:   "list-predefined-event-types",
	Short: "List all predefined event types",
	Run: func(cmd *cobra.Command, args []string) {
		ListPredefinedEventTypes()
	},
}

var savePredefEventTypeCmd = &cobra.Command{
	Use:   "save-predefined-event-type [event_type_key]",
	Short: "Save a predefined event type to the system",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// List predefined event types instead
			ListPredefinedEventTypes()
			return
		}

		client := NewAPIClient(serverURL)
		if err := SavePredefinedEventType(client, args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var saveAllEventTypesCmd = &cobra.Command{
	Use:   "save-all-event-types",
	Short: "Save all predefined event types to JSON files in a directory",
	Run: func(cmd *cobra.Command, args []string) {
		eventTypesDir := badgesDir + "/event_types"
		if err := SaveAllPredefinedEventTypes(eventTypesDir); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("All predefined event types saved to %s\n", eventTypesDir)
	},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	Execute()
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "badgecli",
	Short: "A CLI for managing badges in the Badge Assignment System",
	Long: `A command line interface for managing badges in the Badge Assignment System. 
This tool allows you to create, list, view, update, and delete badges and event types, as well as
import and export badge and event type definitions to and from JSON files.`,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.badgecli.yaml)")
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8080", "server URL")
	rootCmd.PersistentFlags().StringVar(&badgesDir, "badges-dir", "./badges", "directory for badge definitions")

	// Add specific flags for badge commands
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Badge description")
	addCmd.Flags().StringVarP(&imageURL, "image", "i", "", "Badge image URL")
	addCmd.Flags().StringVarP(&flowDefFile, "flow", "f", "", "Path to flow definition JSON file")

	// Add specific flags for event type commands
	addEventTypeCmd.Flags().StringVarP(&description, "description", "d", "", "Event type description")
	addEventTypeCmd.Flags().StringVarP(&schemaFile, "schema", "s", "", "Path to schema JSON file")

	updateEventTypeCmd.Flags().String("name", "", "New name for the event type")
	updateEventTypeCmd.Flags().StringVarP(&description, "description", "d", "", "Event type description")
	updateEventTypeCmd.Flags().StringVarP(&schemaFile, "schema", "s", "", "Path to schema JSON file")

	// Add badge commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(savePredefCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(saveAllPredefCmd)
	rootCmd.AddCommand(listPredefinedCmd)

	// Add event type commands
	rootCmd.AddCommand(listEventTypesCmd)
	rootCmd.AddCommand(getEventTypeCmd)
	rootCmd.AddCommand(addEventTypeCmd)
	rootCmd.AddCommand(updateEventTypeCmd)
	rootCmd.AddCommand(deleteEventTypeCmd)
	rootCmd.AddCommand(importEventTypeCmd)
	rootCmd.AddCommand(exportEventTypeCmd)
	rootCmd.AddCommand(listPredefinedEventTypesCmd)
	rootCmd.AddCommand(savePredefEventTypeCmd)
	rootCmd.AddCommand(saveAllEventTypesCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".badgecli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".badgecli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// If environment variables are set, use them
	if os.Getenv("BADGE_SERVER_URL") != "" {
		serverURL = os.Getenv("BADGE_SERVER_URL")
	}

	if os.Getenv("BADGE_DEFINITIONS_DIR") != "" {
		badgesDir = os.Getenv("BADGE_DEFINITIONS_DIR")
	}

	// Ensure the badges directory exists
	if _, err := os.Stat(badgesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(badgesDir, 0755); err != nil {
			log.Fatalf("Error creating badges directory: %v", err)
		}
	}
}
