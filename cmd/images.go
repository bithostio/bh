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
	imagesProviderID   int
	imagesArchitecture string
)

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List available OS images",
	RunE:  runImages,
}

func init() {
	rootCmd.AddCommand(imagesCmd)
	imagesCmd.Flags().IntVarP(&imagesProviderID, "provider", "p", 0, "Provider ID (required)")
	imagesCmd.Flags().StringVarP(&imagesArchitecture, "arch", "a", "x86", "Architecture (x86 or arm)")
	imagesCmd.MarkFlagRequired("provider")
}

func runImages(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	images, err := client.ListImages(imagesProviderID, imagesArchitecture)
	if err != nil {
		return err
	}

	if len(images) == 0 {
		fmt.Println("No images found.")
		return nil
	}

	rows := [][]string{
		{"ID", "Name", "Distribution", "Architecture"},
	}

	for _, image := range images {
		rows = append(rows, []string{
			strconv.Itoa(image.ID),
			image.Name,
			image.Distribution,
			image.Architecture,
		})
	}

	ui.PrintTable(rows)

	return nil
}
