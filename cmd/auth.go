package cmd

import (
	"fmt"
	"strings"

	"github.com/bithostio/bh/internal/cli"
	"github.com/bithostio/bh/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Configure API authentication",
	Long: `Configure your bithost.io API key for authentication.

Get your API key from: https://dashboard.bithost.io/api_keys`,
	RunE: runAuth,
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func runAuth(cmd *cobra.Command, args []string) error {
	apiKey, err := cli.Input("Enter your Bithost.io API key",
		cli.WithMask('*'),
		cli.WithValidation(func(input string) error {
			if len(strings.TrimSpace(input)) < 10 {
				return fmt.Errorf("API key seems too short")
			}
			return nil
		}),
	)
	if err != nil {
		return err
	}

	apiKey = strings.TrimSpace(apiKey)

	cfg, err := config.Load()
	if err != nil {
		cfg = config.New()
	}

	cfg.APIKey = apiKey

	if err := config.Save(cfg); err != nil {
		return err
	}

	configPath, _ := config.DefaultConfigPath()
	fmt.Printf("\n✓ API key saved to %s\n", configPath)
	fmt.Println("\nYou're all set! Try running: bh balance")

	return nil
}
