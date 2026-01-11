package cmd

import (
	"fmt"
	"strconv"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/cli"
	"github.com/spf13/cobra"
)

var (
	sizesProviderID int
	sizesRegionID   int
	sizesPage       int
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
	sizesCmd.Flags().IntVar(&sizesPage, "page", 0, "Page number (optional)")
	_ = sizesCmd.MarkFlagRequired("provider")
	_ = sizesCmd.MarkFlagRequired("region")
}

func runSizes(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	resp, err := client.ListSizes(sizesRegionID, sizesProviderID, sizesPage)
	if err != nil {
		return err
	}

	if len(resp.Sizes) == 0 {
		fmt.Println("No sizes found.")
		return nil
	}

	rows := [][]string{
		{"ID", "Name", "Memory", "CPU", "Disk", "Price/mo"},
	}

	for _, size := range resp.Sizes {
		rows = append(rows, []string{
			strconv.Itoa(size.ID),
			size.Name,
			cli.FormatMemory(size.Memory),
			size.Processor,
			cli.FormatStorage(size.Disk),
			cli.FormatMoney(size.Price),
		})
	}

	cli.PrintTable(rows)
	cli.DisplayPagination(resp.Meta.Pagination)

	return nil
}
