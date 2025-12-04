package cmd

import (
	"fmt"
	"strconv"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/ui"
	"github.com/spf13/cobra"
)

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List available cloud providers",
	RunE:  runProviders,
}

func init() {
	rootCmd.AddCommand(providersCmd)
}

func runProviders(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	providers, err := client.ListProviders()
	if err != nil {
		return err
	}

	if len(providers) == 0 {
		fmt.Println("No providers found.")
		return nil
	}

	rows := [][]string{
		{"ID", "Name"},
	}

	for _, provider := range providers {
		rows = append(rows, []string{
			strconv.Itoa(provider.ID),
			provider.Name,
		})
	}

	ui.PrintTable(rows)

	return nil
}
