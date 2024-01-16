package cmd

import (
	"database/sql"
	"log"

	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/models"
	"github.com/spf13/cobra"
)

// Whether we want to include example data
var includeExampleData bool

// Run the migrate command
func runMigrateCommand(cmd *cobra.Command, args []string) {
	logger := log.Default()

	// Database
	manager, err := db.Init()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Print("Migrating database models")
	manager.AutoMigrate(&models.Feed{})
	manager.AutoMigrate(&models.FeedItem{})

	if includeExampleData {
		logger.Print("Adding example data")
		manager.Create(&models.Feed{
			Url:         "https://michaelenger.com/feed.rss",
			Title:       "Michael Enger",
			Link:        sql.NullString{"https://michaelenger.com", true},
			Description: sql.NullString{"I am a professional. This is my website.", true},
			Image:       sql.NullString{"https://michaelenger.com/assets/happybass.png", true},
		})
		manager.Create(&models.Feed{
			Url:         "https://pluralistic.net/feed/",
			Title:       "Pluralistic: Daily links from Cory Doctorow",
			Link:        sql.NullString{"https://pluralistic.net/", true},
			Description: sql.NullString{"No trackers, no ads. Black type, white background. Privacy policy: we don't collect or retain any data at all ever period.", true},
			Image:       sql.NullString{"https://i0.wp.com/pluralistic.net/wp-content/uploads/2020/02/cropped-guillotine-French-Revolution.jpg?fit=32%2C32&#038;ssl=1", true},
		})
		manager.Create(&models.Feed{
			Url:         "https://boilingsteam.com/feed/rss-feed.xml",
			Title:       "Boiling Steam",
			Link:        sql.NullString{"https://boilingsteam.com", true},
			Description: sql.NullString{"PC Gaming on Linux is so Tomorrow!", true},
			Image:       sql.NullString{"", false},
		})
	}
}

// Migrate command - migrate database models
var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Initialise/migrate the database models",
	Long:  "Initialise/migrate the database models",
	Run:   runMigrateCommand,
}

// Initialise the migrate command
func init() {
	migrateCommand.Flags().BoolVarP(&includeExampleData, "include-example-data", "d", false, "Include example data")

	rootCmd.AddCommand(migrateCommand)
}
