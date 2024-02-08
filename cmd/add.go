package cmd

import (
	"bufio"
	"fmt"
	"log"
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
func runAddCommand(cmd *cobra.Command, args []string) {
	logger := log.Default()
	url := args[0]

	// Get feed URLs
	feedUrls, err := parser.FindFeedUrls(url)
	if err != nil {
		logger.Fatal(err)
	}

	if len(feedUrls) == 0 {
		logger.Fatal(fmt.Sprintf("Unable to find a feed in %s", url))
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
			logger.Fatal(err)
		}

		input = strings.TrimSuffix(input, "\n")
		index, err = strconv.Atoi(input)
		if err != nil {
			logger.Fatal(err)
		}
	}

	url = feedUrls[index]

	// Parse the feed
	feed, items, err := parser.ParseFeed(url)
	if err != nil {
		logger.Fatal(err)
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

		return
	}

	// Read config file
	conf, err := config.FromFile(configFile)
	if err != nil {
		logger.Fatal(err)
	}

	db, err := db.Init(conf.DatabaseFile)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("Adding/updating feed from %s", feed.Url)
	feed, created, err := models.CreateOrUpdateFeed(db, feed)
	if err != nil {
		logger.Fatal(err)
	}

	if created {
		logger.Printf("..created (id=%d)", feed.ID)
	} else {
		logger.Printf("..updated (id=%d)", feed.ID)
	}

	for _, item := range items {
		logger.Printf("Adding/updating feed item: %s", item.Link)
		item.Feed = *feed
		item, created, err = models.CreateOrUpdateFeedItem(db, item)
		if err != nil {
			logger.Fatal(err)
		}

		if created {
			logger.Printf("..created (id=%d)", item.ID)
		} else {
			logger.Printf("..updated (id=%d)", item.ID)
		}
	}
}

// Add command - Add a new feed
var addCommand = &cobra.Command{
	Use:   "add [url]",
	Short: "Add a feed",
	Long:  "Add a feed to the list of feeds, syncing it in the process",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run:   runAddCommand,
}

// Initialise the sync command
func init() {
	addCommand.Flags().BoolVarP(&verifyFeed, "verify", "v", false, "Verify feed instead of adding it")

	rootCmd.AddCommand(addCommand)
}
