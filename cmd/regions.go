package cmd

import (
	"fmt"
	"strconv"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/ui"
	"github.com/spf13/cobra"
)

var (
	regionsProviderID int
	regionsPage       int
)

var regionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "List available regions for a provider",
	RunE:  runRegions,
}

func init() {
	rootCmd.AddCommand(regionsCmd)
	regionsCmd.Flags().IntVarP(&regionsProviderID, "provider", "p", 0, "Provider ID (required)")
	regionsCmd.Flags().IntVar(&regionsPage, "page", 0, "Page number (optional)")
	regionsCmd.MarkFlagRequired("provider")
}

func runRegions(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	resp, err := client.ListRegions(regionsProviderID, regionsPage)
	if err != nil {
		return err
	}

	if len(resp.Regions) == 0 {
		fmt.Println("No regions found.")
		return nil
	}

	rows := [][]string{
		{"ID", "Name", "Slug", "Location"},
	}

	for _, region := range resp.Regions {
		rows = append(rows, []string{
			strconv.Itoa(region.ID),
			region.Name,
			region.Slug,
			region.Origin,
		})
	}

	ui.PrintTable(rows)
	ui.DisplayPagination(resp.Meta.Pagination.Page, resp.Meta.Pagination.HasPreviousPage, resp.Meta.Pagination.HasNextPage)

	return nil
}
