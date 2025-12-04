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
)

var regionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "List available regions for a provider",
	RunE:  runRegions,
}

func init() {
	rootCmd.AddCommand(regionsCmd)
	regionsCmd.Flags().IntVarP(&regionsProviderID, "provider", "p", 0, "Provider ID (required)")
	regionsCmd.MarkFlagRequired("provider")
}

func runRegions(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	regions, err := client.ListRegions(regionsProviderID)
	if err != nil {
		return err
	}

	if len(regions) == 0 {
		fmt.Println("No regions found.")
		return nil
	}

	rows := [][]string{
		{"ID", "Name", "Slug", "Location"},
	}

	for _, region := range regions {
		rows = append(rows, []string{
			strconv.Itoa(region.ID),
			region.Name,
			region.Slug,
			region.Origin,
		})
	}

	ui.PrintTable(rows)

	return nil
}
