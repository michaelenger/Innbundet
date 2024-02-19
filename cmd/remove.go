package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/models"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Run the command
func runRemoveCommand(cmd *cobra.Command, args []string) error {
	feedId := args[0]

	conf, err := config.FromFile(configFile)
	if err != nil {
		return err
	}

	db, err := db.Init(conf.DatabaseFile)
	if err != nil {
		return err
	}

	// Get the feed and items
	log.Debug().
		Str("id", feedId).
		Msg("Retrieving the feed")
	feed := models.Feed{}
	result := db.First(&feed, feedId)
	if result.Error != nil {
		return result.Error
	}

	log.Debug().
		Str("id", feedId).
		Msg("Retrieving the feed items")
	feedItems := []models.FeedItem{}
	result = db.Where("feed_id = ?", feed.ID).Find(&feedItems)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("Remove feed \"%s\" and its %d feed items? (y/N): ", feed.Title, len(feedItems))
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	input = strings.TrimSuffix(input, "\n")
	if input == "" && strings.ToLower(input) != "y" {
		return nil
	}

	if len(feedItems) != 0 {
		log.Debug().Msg("Deleting feed items...")
		result = db.Delete(&feedItems)
		if result.Error != nil {
			return result.Error
		}
		log.Debug().Msg(fmt.Sprintf("%d rows deleted", result.RowsAffected))
	}

	log.Debug().
		Uint("id", feed.ID).
		Msg("Deleting feed...")
	result = db.Delete(&feed)
	if result.Error != nil {
		return result.Error
	}

	log.Debug().Msg(fmt.Sprintf("%d rows deleted", result.RowsAffected))

	return nil
}

// Remove command - Remove a feed
var removeCommand = &cobra.Command{
	Use:   "remove [id]",
	Short: "Remove a feed",
	Long:  "Remove a feed to the list of feeds, including all its feed items",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   wrapRunFn(runRemoveCommand),
}

// Initialise the command
func init() {
	rootCmd.AddCommand(removeCommand)
}
