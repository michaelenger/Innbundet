package cmd

import (
	"fmt"
	"os"

	"github.com/michaelenger/innbundet/config"
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("Starting server")
	return serv.Start(fmt.Sprintf(":%s", port))
}

// Server command - run the web app
var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Run the web app",
	Long:  "Run the web app",
	Run:   wrapRunFn(runServerCommand),
}

// Initialise the server command
func init() {
	rootCmd.AddCommand(serverCommand)
}
