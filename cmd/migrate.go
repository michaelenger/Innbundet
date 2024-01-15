package cmd

import (
	"log"

	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/models"
	"github.com/spf13/cobra"
)

// Run the migrate command
func runMigrateCommand(cmd *cobra.Command, args []string) {
	logger := log.Default()

	// Database
	err := db.Init()
	if err != nil {
		logger.Fatal(err)
	}

	manager := db.DbManager()

	logger.Print("Migrating database models")
	manager.AutoMigrate(&models.Feed{})
	manager.AutoMigrate(&models.FeedItem{})
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
	rootCmd.AddCommand(migrateCommand)
}
