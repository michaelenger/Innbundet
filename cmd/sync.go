package cmd

import (
	"fmt"

	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/models"
	"github.com/michaelenger/innbundet/parser"
	"github.com/rs/zerolog/log"
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

	log.Debug().Msg("Fetching feeds...")
	feeds := []models.Feed{}
	result := db.Find(&feeds)
	if result.Error != nil {
		return result.Error
	}

	log.Info().Msg(fmt.Sprintf("Found %d feeds", len(feeds)))
	for _, feed := range feeds {
		log.Info().
			Uint("id", feed.ID).
			Msg("Syncing feed...")

		feed, items, err := parser.ParseFeed(feed.Url)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Unable to parse feed")
			continue // don't stop us from syncing the other feeds
		}

		feed, _, err = models.CreateOrUpdateFeed(db, feed)
		if err != nil {
			return err
		}

		log.Debug().
			Uint("id", feed.ID).
			Msg("..updated metadata")

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
			log.Info().
				Uint("id", feed.ID).
				Msg(fmt.Sprintf("..added %d feed items", createdCount))
		}
		if updatedCount != 0 {
			log.Debug().
				Uint("id", feed.ID).
				Msg(fmt.Sprintf("..updated %d feed items", updatedCount))
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
