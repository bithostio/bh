package cmd

import (
	"fmt"
	"strconv"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/bithostio/bh/internal/ui"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <server-id>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a server",
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

var (
	forceDelete bool
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Skip confirmation prompt")
}

func runDelete(cmd *cobra.Command, args []string) error {
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
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Delete server %d", serverID),
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrAbort {
				fmt.Println("Cancelled.")
				return nil
			}
			return err
		}
	}

	stop := ui.ShowProgress(fmt.Sprintf("Deleting server %d...", serverID))
	err = client.DeleteServer(serverID)
	stop()

	if err != nil {
		return fmt.Errorf(api.HandleError(err))
	}

	fmt.Printf("%s Server %d deleted successfully\n", ui.Green("✓"), serverID)

	return nil
}
