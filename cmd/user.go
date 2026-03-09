package cmd

import (
	"fmt"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/cli"
	"github.com/bithostio/bh/internal/config"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Show your user information",
	RunE:  runUser,
}

func init() {
	rootCmd.AddCommand(userCmd)
}

func runUser(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	user, err := client.GetUser()
	if err != nil {
		return err
	}

	displayUser(user)

	return nil
}

func displayUser(user *api.User) {
	threshold := 5.0

	fmt.Printf("\nAccount: %s (%s)\n", cli.Bold(user.FullName), user.Email)
	fmt.Printf("Balance: %s\n", cli.Bold(fmt.Sprintf("$%.2f", user.Balance)))
	if user.ServerLimit > 0 {
		fmt.Printf("Server Limit: %s\n", cli.Bold(fmt.Sprintf("%d", user.ServerLimit)))
	}

	if user.Balance < threshold {
		fmt.Printf("\n%s Your balance is low!\n", cli.Yellow("Warning:"))
		fmt.Printf("Top up at: %s\n\n", cli.Cyan("https://dashboard.bithost.io"))
	} else {
		fmt.Println()
	}
}
