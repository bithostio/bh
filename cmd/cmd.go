package cmd

import (
	"fmt"
	"os"

	"github.com/bithostio/bh/internal/cli"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bh",
	Short: "Bithost.io CLI - create and manage cloud servers",
	Long: `Bithost.io CLI - create and manage cloud servers.

Run 'bh tui' for interactive mode, or 'bh auth' to get started.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, cli.Red(err.Error()))
		os.Exit(1)
	}
}
