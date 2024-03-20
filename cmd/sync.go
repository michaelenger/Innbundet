package cmd

import (
	"fmt"
	"sync"

	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/models"
	"github.com/michaelenger/innbundet/parser"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func syncFeed(wg *sync.WaitGroup, db *gorm.DB, feedID uint) {
	defer wg.Done()

	var feed *models.Feed
	db.Where("id = ?", feedID).Limit(1).Find(&feed)

	outputText := fmt.Sprintf("Syncing feed %s (%d) ...", feed.Title, feed.ID)
	feed, items, err := parser.ParseFeed(feed.Url)
	if err != nil {
		fmt.Printf("%s ERROR! Failed to parse feed: %s", outputText, err)
		return
	}

	feed, _, err = models.CreateOrUpdateFeed(db, feed)
	if err != nil {
		fmt.Printf("%s ERROR! Failed to create/update feed: %s", outputText, err)
		return
	}

	createdCount := 0
	updatedCount := 0
	for _, item := range items {
		item.Feed = *feed
		_, created, err := models.CreateOrUpdateFeedItem(db, item)
		if err != nil {
			fmt.Printf("%s ERROR! Failed to create/update feed item: %s", outputText, err)
			return
		}

		if created {
			createdCount += 1
		} else {
			updatedCount += 1
		}
	}

	outputText += " SUCCESS! "

	if createdCount != 0 {
		outputText += fmt.Sprintf("%d items created", createdCount)
	}
	if updatedCount != 0 {
		if createdCount != 0 {
			outputText += ", "
		}
		outputText += fmt.Sprintf("%d items updated", updatedCount)
	}

	fmt.Println(outputText)
}

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

	fmt.Printf("Fetching feeds... ")
	feeds := []models.Feed{}
	result := db.Find(&feeds)
	if result.Error != nil {
		return result.Error
	}

	var wg sync.WaitGroup

	fmt.Printf("%d found\n", len(feeds))
	for _, feed := range feeds {
		wg.Add(1)
		go syncFeed(&wg, db, feed.ID)
	}

	wg.Wait()

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
