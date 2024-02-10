package cmd

import (
	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/server"
	"github.com/spf13/cobra"
)

// Port to serve the site on
var port int32

// Run the server command
func runServerCommand(cmd *cobra.Command, args []string) error {
	conf, err := config.FromFile(configFile)
	if err != nil {
		return err
	}

	db, err := db.Init(conf.DatabaseFile)
	if err != nil {
		return err
	}

	serv, err := server.Init(db, conf)
	if err != nil {
		return err
	}

	return serv.Start(":8080")
}

// Server command - run the web app
var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Run the web app",
	Long:  "Run the web app",
	RunE:  runServerCommand,
}

// Initialise the server command
func init() {
	serverCommand.Flags().Int32VarP(&port, "port", "p", 8080, "Port to serve the app on")

	rootCmd.AddCommand(serverCommand)
}
