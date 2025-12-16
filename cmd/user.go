package cmd

import (
	"fmt"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/fatih/color"
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

func displayUser(user *api.UserResponse) {
	threshold := 5.0

	bold := color.New(color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("\nAccount: %s (%s)\n", bold(user.FullName), user.Email)
	fmt.Printf("Balance: %s\n", bold(fmt.Sprintf("$%.2f", user.Balance)))
	fmt.Printf("Server Limit: %s\n", bold(fmt.Sprintf("%d", user.ServerLimit)))

	if user.Balance < threshold {
		fmt.Printf("\n%s Your balance is low!\n", yellow("Warning:"))
		fmt.Printf("Top up at: %s\n\n", cyan("https://dashboard.bithost.io/billing"))
	} else {
		fmt.Println()
	}
}
