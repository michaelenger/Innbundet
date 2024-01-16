package cmd

import (
	"log"

	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/models"
	"github.com/michaelenger/innbundet/sync"
	"github.com/spf13/cobra"
)

// Run the sync command
func runSyncCommand(cmd *cobra.Command, args []string) {
	logger := log.Default()

	manager, err := db.Init()
	if err != nil {
		logger.Fatal(err)
	}

	feeds := []models.Feed{}
	result := manager.Find(&feeds)
	if result.Error != nil {
		logger.Fatal(result.Error)
	}

	logger.Printf("Found %d feeds", len(feeds))

	for _, feed := range feeds {
		err = sync.SyncFeed(manager, &feed)
		if err != nil {
			logger.Printf("ERROR: %v", err)
		}
	}
}

// Sync command - download
var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "Syncronise the feeds",
	Long:  "Syncronise the feeds, downloading any new feed items and removing old ones",
	Run:   runSyncCommand,
}

// Initialise the sync command
func init() {
	rootCmd.AddCommand(syncCommand)
}
