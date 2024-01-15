package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// Run the sync command
func runSyncCommand(cmd *cobra.Command, args []string) {
	logger := log.Default()

	logger.Print("TODO")
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
