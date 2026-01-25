package cmd

import (
	"github.com/bithostio/bh/internal/tui"
	"github.com/spf13/cobra"
)

var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Launch the interactive TUI dashboard",
	Long: `Launch the interactive terminal user interface for managing servers.

The TUI provides a visual dashboard for:
  - Viewing and managing servers
  - Creating new servers with a step-by-step wizard
  - Managing SSH keys`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run()
	},
}

func init() {
	rootCmd.AddCommand(manageCmd)
}
