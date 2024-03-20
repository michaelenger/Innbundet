package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// Path to the config file (used in all commands)
var configFile string
var debug bool

// Root command
var rootCmd = &cobra.Command{
	Use:   "innbundet",
	Short: "Personal RSS/ATOM reader",
}

type runFunc func(*cobra.Command, []string) error

// Wrap a run function in one which handles any errors it returns.
func wrapRunFn(fn runFunc) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		err := fn(cmd, args)
		if err != nil {
			fmt.Printf("Error! %s\n", err)
			os.Exit(1)
		}
	}
}

func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "Config file to read")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Whether to show debug output")

	return rootCmd.Execute()
}
