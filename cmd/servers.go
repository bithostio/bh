package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "servers",
	Short: "Manage servers",
	Long:  `Create, list, and delete servers.`,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
