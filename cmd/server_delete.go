package cmd

import (
	"fmt"
	"strconv"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/cli"
	"github.com/bithostio/bh/internal/config"
	"github.com/spf13/cobra"
)

var (
	forceDelete bool
)

var serverDeleteCmd = &cobra.Command{
	Use:     "delete <server-id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a server",
	Args:    cobra.ExactArgs(1),
	RunE:    runServerDelete,
}

func init() {
	serverCmd.AddCommand(serverDeleteCmd)
	serverDeleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Skip confirmation prompt")
}

func runServerDelete(cmd *cobra.Command, args []string) error {
	serverID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid server ID: %s", args[0])
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	if !forceDelete {
		confirmed, err := cli.Confirm(fmt.Sprintf("Delete server %d?", serverID))
		if err != nil {
			if err == cli.ErrAborted {
				fmt.Println("Cancelled.")
				return nil
			}
			return err
		}
		if !confirmed {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	stop, cancel := cli.ShowProgress(fmt.Sprintf("Deleting server %d...", serverID))
	err = client.DeleteServer(serverID)

	if err != nil {
		cancel()
		return err
	}

	stop()

	fmt.Printf("%s Server %d deleted successfully\n", cli.Green("✓"), serverID)

	return nil
}
