package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage servers",
	Long:  `Create, list, and delete servers.`,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
