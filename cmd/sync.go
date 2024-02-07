package cmd

import (
	"log"

	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/models"
	"github.com/michaelenger/innbundet/parser"
	"github.com/spf13/cobra"
)

// Run the sync command
func runSyncCommand(cmd *cobra.Command, args []string) {
	logger := log.Default()

	// Read config file
	conf, err := config.FromFile(configFile)
	if err != nil {
		logger.Fatal(err)
	}

	db, err := db.Init(conf.DatabaseFile)
	if err != nil {
		logger.Fatal(err)
	}

	feeds := []models.Feed{}
	result := db.Find(&feeds)
	if result.Error != nil {
		logger.Fatal(result.Error)
	}

	logger.Printf("Found %d feeds", len(feeds))

	for _, feed := range feeds {
		logger.Printf("Syncing %s (%d)", feed.Url, feed.ID)

		feed, items, err := parser.ParseFeed(feed.Url)
		if err != nil {
			logger.Fatal(err)
		}

		feed, _, err = models.CreateOrUpdateFeed(db, feed)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Print("..updated metadata")

		createdCount := 0
		updatedCount := 0
		for _, item := range items {
			item.Feed = *feed
			_, created, err := models.CreateOrUpdateFeedItem(db, item)
			if err != nil {
				logger.Fatal(err)
			}

			if created {
				createdCount += 1
			} else {
				updatedCount += 1
			}
		}

		if createdCount != 0 {
			logger.Printf("..added %d feed items", createdCount)
		}
		if updatedCount != 0 {
			logger.Printf("..updated %d feed items", updatedCount)
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
