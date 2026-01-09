package cmd

import (
	"fmt"
	"os"

	"github.com/bithostio/bh/internal/cli"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bh",
	Short: "bithost.io CLI - Manage your servers from the command line",
	Long: `bh is a command-line interface for Bithost.io that allows you to:
  - Create and manage servers
  - Check your account balance
  - List and delete servers

Get started by running: bh auth`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, cli.Red(err.Error()))
		os.Exit(1)
	}
}
