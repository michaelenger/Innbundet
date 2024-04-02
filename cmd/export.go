package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/log"
	"github.com/michaelenger/innbundet/models"
	"github.com/spf13/cobra"
)

// Run the command
func runExportCommand(cmd *cobra.Command, args []string) error {
	conf, err := config.FromFile(configFile)
	if err != nil {
		return err
	}

	db, err := db.Init(conf.DatabaseFile)
	if err != nil {
		return err
	}

	// Get the feed and items
	log.Debug("Retrieving feeds...")
	feeds := []models.Feed{}
	result := db.Order("title asc").Find(&feeds)
	if result.Error != nil {
		return result.Error
	}

	bytes, err := json.MarshalIndent(feeds, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	return nil
}

// Export command - Spit out the feeds as JSON
var exportCommand = &cobra.Command{
	Use:   "export",
	Short: "Export the feeds as JSON",
	Long:  "Export the feeds as JSON, ordered by feed title.",
	Run:   wrapRunFn(runExportCommand),
}

// Initialise the command
func init() {
	rootCmd.AddCommand(exportCommand)
}
