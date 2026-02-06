package cmd

import (
	"fmt"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/cli"
	"github.com/bithostio/bh/internal/config"
	"github.com/spf13/cobra"
)

var providersPage int

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "List available cloud providers",
	RunE:  runProviders,
}

func init() {
	rootCmd.AddCommand(providersCmd)
	providersCmd.Flags().IntVar(&providersPage, "page", 0, "Page number (optional)")
}

func runProviders(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	resp, err := client.ListProviders(providersPage)
	if err != nil {
		return err
	}

	if len(resp.Providers) == 0 {
		fmt.Println("No providers found.")
		return nil
	}

	rows := [][]string{
		{"Slug", "Name"},
	}

	for _, provider := range resp.Providers {
		rows = append(rows, []string{
			provider.Slug,
			provider.Name,
		})
	}

	cli.PrintTable(rows)

	return nil
}
