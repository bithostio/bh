package cmd

import (
	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/ui"
	"github.com/spf13/cobra"
)

var serversCmd = &cobra.Command{
	Use:     "servers",
	Aliases: []string{"ls", "list"},
	Short:   "List all your servers",
	RunE:    runServers,
}

func init() {
	rootCmd.AddCommand(serversCmd)
}

func runServers(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	servers, err := client.ListServers()
	if err != nil {
		return err
	}

	ui.DisplayServers(servers)

	return nil
}
