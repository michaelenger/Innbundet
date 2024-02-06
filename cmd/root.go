package cmd

import "github.com/spf13/cobra"

// Path to the config file (used in all commands)
var configFile string

// Root command
var rootCmd = &cobra.Command{
	Use:   "innbundet",
	Short: "Personal RSS/ATOM reader",
}

func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml", "Config file to read")

	return rootCmd.Execute()
}
