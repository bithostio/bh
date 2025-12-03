package cmd

import (
	"fmt"
	"strings"

	"github.com/bithostio/bh/internal/config"
	"github.com/manifoldco/promptui"
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
	prompt := promptui.Prompt{
		Label: "Enter your Bithost.io API key",
		Mask:  '*',
		Validate: func(input string) error {
			input = strings.TrimSpace(input)
			if len(input) < 10 {
				return fmt.Errorf("API key seems too short")
			}
			return nil
		},
	}

	apiKey, err := prompt.Run()
	if err != nil {
		return err
	}

	apiKey = strings.TrimSpace(apiKey)

	cfg := &config.Config{
		APIKey:  apiKey,
		BaseURL: "https://dashboard.bithost.io/api/v1",
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	configPath, _ := config.DefaultConfigPath()
	fmt.Printf("\n✓ API key saved to %s\n", configPath)
	fmt.Println("\nYou're all set! Try running: bh balance")

	return nil
}
