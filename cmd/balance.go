package cmd

import (
	"fmt"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Show your account balance",
	RunE:  runBalance,
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}

func runBalance(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client := api.NewClient(cfg.BaseURL, cfg.APIKey)

	balance, err := client.GetBalance()
	if err != nil {
		return fmt.Errorf(api.HandleError(err))
	}

	displayBalance(balance)

	return nil
}

func displayBalance(balance float64) {
	threshold := 5.0 // Low balance threshold

	bold := color.New(color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("\nCurrent Balance: %s\n", bold(fmt.Sprintf("$%.2f", balance)))

	if balance < threshold {
		fmt.Printf("\n%s Your balance is low!\n", yellow("Warning:"))
		fmt.Printf("Top up at: %s\n\n", cyan("https://dashboard.bithost.io/billing"))
	} else {
		fmt.Println()
	}
}
