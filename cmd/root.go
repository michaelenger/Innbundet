package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "innbundet",
	Short: "Personal RSS/ATOM reader",
}

func Execute() error {
	return rootCmd.Execute()
}
