package cmd

import (
	"fmt"

	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/models"
	"github.com/michaelenger/innbundet/parser"
	"github.com/spf13/cobra"
)

// Run the sync command
func runSyncCommand(cmd *cobra.Command, args []string) error {
	conf, err := config.FromFile(configFile)
	if err != nil {
		return err
	}

	db, err := db.Init(conf.DatabaseFile)
	if err != nil {
		return err
	}

	feeds := []models.Feed{}
	result := db.Find(&feeds)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("Found %d feeds\n", len(feeds))

	for _, feed := range feeds {
		fmt.Printf("Syncing %s (%d)\n", feed.Url, feed.ID)

		feed, items, err := parser.ParseFeed(feed.Url)
		if err != nil {
			fmt.Printf(" ERROR! Unable to parse feed: %s\n", err)
			continue // don't stop us from syncing the other feeds
		}

		feed, _, err = models.CreateOrUpdateFeed(db, feed)
		if err != nil {
			return err
		}

		fmt.Println("..updated metadata")

		createdCount := 0
		updatedCount := 0
		for _, item := range items {
			item.Feed = *feed
			_, created, err := models.CreateOrUpdateFeedItem(db, item)
			if err != nil {
				return err
			}

			if created {
				createdCount += 1
			} else {
				updatedCount += 1
			}
		}

		if createdCount != 0 {
			fmt.Printf("..created %d feed items\n", createdCount)
		}
		if updatedCount != 0 {
			fmt.Printf("..updated %d feed items\n", updatedCount)
		}
	}

	return nil
}

// Sync command - download
var syncCommand = &cobra.Command{
	Use:   "sync",
	Short: "Syncronise the feeds",
	Long:  "Syncronise the feeds, downloading any new feed items and removing old ones",
	Run:   wrapRunFn(runSyncCommand),
}

// Initialise the sync command
func init() {
	rootCmd.AddCommand(syncCommand)
}
