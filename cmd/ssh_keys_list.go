package cmd

import (
	"fmt"
	"strconv"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/ui"
	"github.com/spf13/cobra"
)

var sshKeysListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all SSH keys",
	RunE:    runSSHKeysList,
}

func init() {
	sshKeysCmd.AddCommand(sshKeysListCmd)
}

func runSSHKeysList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	keys, err := client.ListSSHKeys()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		fmt.Println("No SSH keys found.")
		fmt.Println("\nAdd a key with: bh ssh-keys add")
		return nil
	}

	rows := [][]string{
		{"ID", "Label", "Fingerprint"},
	}

	for _, key := range keys {
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

	return nil
}
