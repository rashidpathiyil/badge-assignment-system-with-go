package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	serverURL string
	badgesDir string
	outputDir string
)

// Embed all JSON files from the badges directory
//
//go:embed badges/*.json badges/event_types/*.json
var exampleFS embed.FS

// Command variables - Import/Export Commands
var listExamplesCmd = &cobra.Command{
	Use:   "list-examples",
	Short: "List all example badges and event types included with the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ListExampleBadgesAndEventTypes(); err != nil {
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

var exportExamplesCmd = &cobra.Command{
	Use:   "export-examples",
	Short: "Export all example badges and event types to a directory",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ExportExampleBadgesAndEventTypes(outputDir); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("All example badges and event types exported to %s\n", outputDir)
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

// New commands for listing and exporting all items
var listBadgesCmd = &cobra.Command{
	Use:   "list-badges",
	Short: "List all badges in the system",
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := ListBadges(client); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

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

var exportAllBadgesCmd = &cobra.Command{
	Use:   "export-all-badges [output_directory]",
	Short: "Export all badges from the system to a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := ExportAllBadges(client, args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var exportAllEventTypesCmd = &cobra.Command{
	Use:   "export-all-event-types [output_directory]",
	Short: "Export all event types from the system to a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := ExportAllEventTypes(client, args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// New commands for importing all items
var importAllBadgesCmd = &cobra.Command{
	Use:   "import-all-badges [directory_path]",
	Short: "Import all badge definitions from a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := ImportAllBadges(client, args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var importAllEventTypesCmd = &cobra.Command{
	Use:   "import-all-event-types [directory_path]",
	Short: "Import all event type definitions from a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := NewAPIClient(serverURL)
		if err := ImportAllEventTypes(client, args[0]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "badgecli",
	Short: "Badge CLI provides import and export functionality for badges and event types",
	Long: `Badge CLI is a tool to manage badges and event types in the badge assignment system.
It helps with importing, exporting, and providing examples of badge and event type definitions.`,
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.badgecli.yaml)")
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8080", "server URL")

	// Export flags
	exportExamplesCmd.Flags().StringVar(&outputDir, "output-dir", "./examples", "directory to export examples to")

	// Add commands
	rootCmd.AddCommand(listExamplesCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(exportExamplesCmd)
	rootCmd.AddCommand(importEventTypeCmd)
	rootCmd.AddCommand(exportEventTypeCmd)

	// Add new commands
	rootCmd.AddCommand(listBadgesCmd)
	rootCmd.AddCommand(listEventTypesCmd)
	rootCmd.AddCommand(exportAllBadgesCmd)
	rootCmd.AddCommand(exportAllEventTypesCmd)
	rootCmd.AddCommand(importAllBadgesCmd)
	rootCmd.AddCommand(importAllEventTypesCmd)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".badgecli" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(".badgecli")
	}

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Load .env file if present
	_ = godotenv.Load()

	// Read server URL from environment variable if set
	if serverEnv := os.Getenv("BADGE_SERVER_URL"); serverEnv != "" {
		serverURL = serverEnv
	}
}

func main() {
	Execute()
}
