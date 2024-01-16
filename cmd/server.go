package cmd

import (
	"log"

	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/server"
	"github.com/spf13/cobra"
)

// Port to serve the site on
var port int32

// Run the server command
func runServerCommand(cmd *cobra.Command, args []string) {
	logger := log.Default()

	// Database
	_, err := db.Init()
	if err != nil {
		logger.Fatal(err)
	}

	// Server
	serv, err := server.Init()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Fatal(serv.Start(":8080"))
}

// Server command - run the web app
var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Run the web app",
	Long:  "Run the web app",
	Run:   runServerCommand,
}

// Initialise the server command
func init() {
	serverCommand.Flags().Int32VarP(&port, "port", "p", 8080, "Port to serve the app on")

	rootCmd.AddCommand(serverCommand)
}
