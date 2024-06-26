package cmd

import (
	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/log"
	"github.com/michaelenger/innbundet/models"
	"github.com/spf13/cobra"
)

// Whether we want to include example data
var includeExampleData bool

// Run the migrate command
func runMigrateCommand(cmd *cobra.Command, args []string) error {
	conf, err := config.FromFile(configFile)
	if err != nil {
		return err
	}

	db, err := db.Init(conf.DatabaseFile)
	if err != nil {
		return err
	}

	log.Info("Migrating database models...")
	db.AutoMigrate(&models.Feed{})
	db.AutoMigrate(&models.FeedItem{})

	if includeExampleData {
		log.Info("Adding example data...")
		image := "https://michaelenger.com/assets/happybass.png"
		db.Create(&models.Feed{
			Url:         "https://michaelenger.com/feed.rss",
			Title:       "Michael Enger",
			Link:        "https://michaelenger.com",
			Description: "I am a professional. This is my website.",
			Image:       &image,
		})
		image = "https://i0.wp.com/pluralistic.net/wp-content/uploads/2020/02/cropped-guillotine-French-Revolution.jpg?fit=32%2C32&#038;ssl=1"
		db.Create(&models.Feed{
			Url:         "https://pluralistic.net/feed/",
			Title:       "Pluralistic: Daily links from Cory Doctorow",
			Link:        "https://pluralistic.net/",
			Description: "No trackers, no ads. Black type, white background. Privacy policy: we don't collect or retain any data at all ever period.",
			Image:       &image,
		})
		db.Create(&models.Feed{
			Url:         "https://boilingsteam.com/feed/rss-feed.xml",
			Title:       "Boiling Steam",
			Link:        "https://boilingsteam.com",
			Description: "PC Gaming on Linux is so Tomorrow!",
			Image:       nil,
		})
	}

	return nil
}

// Migrate command - migrate database models
var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Initialise/migrate the database models",
	Long:  "Initialise/migrate the database models",
	Run:   wrapRunFn(runMigrateCommand),
}

// Initialise the migrate command
func init() {
	migrateCommand.Flags().BoolVarP(&includeExampleData, "include-example-data", "e", false, "Include example data")

	rootCmd.AddCommand(migrateCommand)
}
