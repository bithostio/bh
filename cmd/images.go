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
	imagesProvider     string
	imagesArchitecture string
	imagesPage         int
)

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List available OS images",
	RunE:  runImages,
}

func init() {
	rootCmd.AddCommand(imagesCmd)
	imagesCmd.Flags().StringVarP(&imagesProvider, "provider", "p", "", "Provider slug (required)")
	imagesCmd.Flags().StringVarP(&imagesArchitecture, "arch", "a", "x86", "Architecture (x86 or arm)")
	imagesCmd.Flags().IntVar(&imagesPage, "page", 0, "Page number (optional)")
	_ = imagesCmd.MarkFlagRequired("provider")
}

func runImages(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	resp, err := client.ListImages(imagesProvider, imagesArchitecture, imagesPage)
	if err != nil {
		return err
	}

	if len(resp.Images) == 0 {
		fmt.Println("No images found.")
		return nil
	}

	rows := [][]string{
		{"ID", "Name", "Distribution", "Architecture"},
	}

	for _, image := range resp.Images {
		rows = append(rows, []string{
			strconv.Itoa(image.ID),
			image.Name,
			image.Distribution,
			image.Architecture,
		})
	}

	cli.PrintTable(rows)
	cli.DisplayPagination(resp.Meta.Pagination)

	return nil
}
