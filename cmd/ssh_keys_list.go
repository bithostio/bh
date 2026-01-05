package cmd

import (
	"fmt"
	"strconv"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/ui"
	"github.com/spf13/cobra"
)

var sshKeysListPage int

var sshKeysListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all SSH keys",
	RunE:    runSSHKeysList,
}

func init() {
	sshKeysCmd.AddCommand(sshKeysListCmd)
	sshKeysListCmd.Flags().IntVar(&sshKeysListPage, "page", 0, "Page number (optional)")
}

func runSSHKeysList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	resp, err := client.ListSSHKeys(sshKeysListPage)
	if err != nil {
		return err
	}

	if len(resp.Keys) == 0 {
		fmt.Println("No SSH keys found.")
		fmt.Println("\nAdd a key with: bh ssh-keys add")
		return nil
	}

	rows := [][]string{
		{"ID", "Label", "Fingerprint"},
	}

	for _, key := range resp.Keys {
		fingerprint := key.Key
		if len(fingerprint) > 50 {
			fingerprint = fingerprint[:47] + "..."
		}

		rows = append(rows, []string{
			strconv.Itoa(key.ID),
			ui.Truncate(key.Label, 30),
			fingerprint,
		})
	}

	ui.PrintTable(rows)
	ui.DisplayPagination(resp.Meta.Pagination.Page, resp.Meta.Pagination.HasPreviousPage, resp.Meta.Pagination.HasNextPage)

	return nil
}
