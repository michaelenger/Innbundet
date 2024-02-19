package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/models"
	"github.com/michaelenger/innbundet/parser"
	"github.com/spf13/cobra"
)

// Whether we want to verify the feed instead of adding it
var verifyFeed bool

// Run the add command
func runAddCommand(cmd *cobra.Command, args []string) error {
	url := args[0]

	fmt.Printf("Getting feed URL from: %s\n", url)
	feedUrls, err := parser.FindFeedUrls(url)
	if err != nil {
		return err
	}

	if len(feedUrls) == 0 {
		return errors.New(fmt.Sprintf("Unable to find a feed in %s", url))
	}

	index := 0
	if len(feedUrls) > 1 {
		fmt.Printf("Found %d feed URLs\n", len(feedUrls))
		for i, feedUrl := range feedUrls {
			fmt.Printf(" [%d] %s\n", i+1, feedUrl)
		}

		fmt.Print("Select feed: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		input = strings.TrimSuffix(input, "\n")
		index, err = strconv.Atoi(input)
		if err != nil {
			return err
		}
	}

	url = feedUrls[index]
	fmt.Printf("Parsing feed from URL: %s\n", url)

	// Parse the feed
	feed, items, err := parser.ParseFeed(url)
	if err != nil {
		return err
	}

	if verifyFeed {
		fmt.Println("Feed")
		fmt.Printf("  Url:         %s\n", feed.Url)
		fmt.Printf("  Title:       %s\n", feed.Title)
		fmt.Printf("  Link:        %s\n", feed.Link)
		fmt.Printf("  Description: %s\n", feed.Description)
		image := "<none>"
		if feed.Image != nil {
			image = *feed.Image
		}
		fmt.Printf("  Image:       %s\n", image)
		fmt.Println("")

		for _, item := range items {
			fmt.Println("FeedItem")
			fmt.Printf("  Title       %s\n", item.Title)
			fmt.Printf("  Link        %s\n", item.Link)
			fmt.Printf("  Description %s\n", item.Description)
			fmt.Printf("  Author      %s\n", item.Author)
			image = "<none>"
			if item.Image != nil {
				image = *item.Image
			}
			fmt.Printf("  Image       %s\n", image)
			fmt.Printf("  Published   %s\n", item.Published)
			fmt.Println("")
		}

		return nil
	}

	// Read config file
	conf, err := config.FromFile(configFile)
	if err != nil {
		return err
	}

	db, err := db.Init(conf.DatabaseFile)
	if err != nil {
		return err
	}

	fmt.Printf("Processing feed...")
	feed, created, err := models.CreateOrUpdateFeed(db, feed)
	if err != nil {
		return err
	}

	if created {
		fmt.Printf("created (id=%d)\n", feed.ID)
	} else {
		fmt.Printf("updated (id=%d)\n", feed.ID)
	}

	for _, item := range items {
		fmt.Printf("Adding/updating feed item (%s)...", item.Link)
		item.Feed = *feed
		item, created, err = models.CreateOrUpdateFeedItem(db, item)
		if err != nil {
			return err
		}

		if created {
			fmt.Printf("created (id=%d)\n", item.ID)
		} else {
			fmt.Printf("updated (id=%d)\n", item.ID)
		}
	}

	return nil
}

// Add command - Add a new feed
var addCommand = &cobra.Command{
	Use:   "add [url]",
	Short: "Add a feed",
	Long:  "Add a feed to the list of feeds, syncing it in the process",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   wrapRunFn(runAddCommand),
}

// Initialise the sync command
func init() {
	addCommand.Flags().BoolVarP(&verifyFeed, "verify", "v", false, "Verify feed instead of adding it")

	rootCmd.AddCommand(addCommand)
}
