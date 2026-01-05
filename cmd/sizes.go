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
	sizesProviderID int
	sizesRegionID   int
)

var sizesCmd = &cobra.Command{
	Use:   "sizes",
	Short: "List available server sizes/plans",
	RunE:  runSizes,
}

func init() {
	rootCmd.AddCommand(sizesCmd)
	sizesCmd.Flags().IntVarP(&sizesProviderID, "provider", "p", 0, "Provider ID (required)")
	sizesCmd.Flags().IntVarP(&sizesRegionID, "region", "r", 0, "Region ID (required)")
	sizesCmd.MarkFlagRequired("provider")
	sizesCmd.MarkFlagRequired("region")
}

func runSizes(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	sizes, err := client.ListSizes(sizesRegionID, sizesProviderID)
	if err != nil {
		return err
	}

	if len(sizes) == 0 {
		fmt.Println("No sizes found.")
		return nil
	}

	rows := [][]string{
		{"ID", "Name", "Memory", "CPU", "Disk", "Price/mo"},
	}

	for _, size := range sizes {
		rows = append(rows, []string{
			strconv.Itoa(size.ID),
			size.Name,
			fmt.Sprintf("%d MB", size.Memory),
			size.Processor,
			fmt.Sprintf("%d GB", size.Disk),
			ui.FormatMoney(size.Price),
		})
	}

	ui.PrintTable(rows)

	return nil
}
